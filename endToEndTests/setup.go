package generateFiles

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	reqst "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var client = http.Client{}

type Client interface {
	SDKClient
	CryptoClient
}

type ClientImpl struct {
	SDKClient
	CryptoClient
}

type SDKClient interface {
	GetObject(*s3.GetObjectInput) (*s3.GetObjectOutput, error)
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
	GetObjectWithContext(aws.Context, *s3.GetObjectInput, ...reqst.Option) (*s3.GetObjectOutput, error)
}

type CryptoClient interface {
	GetObjectWithPSK(*s3.GetObjectInput, []byte) (*s3.GetObjectOutput, error)
	PutObjectWithPSK(*s3.PutObjectInput, []byte) (*s3.PutObjectOutput, error)
}

type request struct {
	Body    []byte
	Context string
	Header  string
	Method  string
	URL     string
}

// Store provides file storage via S3.
type Store struct {
	config *aws.Config
	bucket string
}

// TODO Once export services have been updated with encryption and decryption
// remove decrypt boolean flag
func sendV4FileToAWS(region, bucket, filename string, encrypt bool) error {
	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	sess, err := session.NewSession(store.config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return err
	}

	v4File, err := os.Open(filename)
	if err != nil {
		log.ErrorC("failed to open file", err, nil)
		return err
	}
	log.Info("successfully retrieved file", nil)

	client, err := getClient(sess, encrypt)
	if err != nil {
		log.ErrorC("failed to create client", err, nil)
		return err
	}

	putObject := &s3.PutObjectInput{
		Key:    &filename,
		Bucket: &bucket,
		Body:   v4File,
	}

	if encrypt {
		psk := createPSK()
		pskStr := hex.EncodeToString(psk)

		err := vaultClient.WriteKey(cfg.VaultPath, filename, pskStr)
		if err != nil {
			log.ErrorC("failed to write to vault", err, nil)
			return err
		}

		_, err = client.PutObjectWithPSK(putObject, psk)
		if err != nil {
			log.ErrorC("failed to upload file", err, nil)
			return err
		}

	} else {
		_, err = client.PutObject(putObject)
		if err != nil {
			log.ErrorC("failed to upload file", err, nil)
			return err
		}
	}

	return nil
}

func getS3File(region, bucket, filename string, decrypt bool) (io.ReadCloser, error) {
	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	sess, err := session.NewSession(store.config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return nil, err
	}

	client, err := getClient(sess, decrypt)
	if err != nil {
		log.ErrorC("failed to create client", err, nil)
		return nil, err
	}

	input := &s3.GetObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(bucket),
	}

	var output *s3.GetObjectOutput
	if decrypt {
		pskStr, err := vaultClient.ReadKey(cfg.VaultPath, filename)
		if err != nil {
			return nil, err
		}
		psk, err := hex.DecodeString(pskStr)
		if err != nil {
			return nil, err
		}

		output, err = client.GetObjectWithPSK(input, psk)
		if err != nil {
			log.ErrorC("encountered error retrieving csv file", err, nil)
			return nil, err
		}
	} else {
		output, err = client.GetObject(input)
		if err != nil {
			log.ErrorC("encountered error retrieving csv file", err, nil)
			return nil, err
		}
	}

	return output.Body, nil
}

func getS3FileSize(region, bucket, filename string, decrypt bool) (*int64, error) {
	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	sess, err := session.NewSession(store.config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return nil, err
	}

	client, err := getClient(sess, decrypt)
	if err != nil {
		log.ErrorC("failed to create client", err, nil)
		return nil, err
	}

	input := &s3.GetObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(bucket),
	}

	ctx := context.Background()
	result, err := client.GetObjectWithContext(ctx, input)
	if err != nil {
		log.ErrorC("failed to find file", err, nil)
		return nil, err
	}
	defer result.Body.Close()

	size := result.ContentLength

	return size, nil
}

func deleteS3File(region, bucket, filename string) error {

	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	session, err := session.NewSession(store.config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return err
	}

	svc := s3.New(session)

	input := &s3.DeleteObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(bucket),
	}

	if _, err = svc.DeleteObject(input); err != nil {
		log.ErrorC("failed to remove file", err, nil)
		return err
	}

	return nil
}

func newRequest(context, data, header, method, url string) *request {
	body := []byte(data)

	return &request{
		Body:    body,
		Context: context,
		Header:  header,
		Method:  method,
		URL:     url,
	}
}

func (req request) sendRequest() (*http.Response, error) {
	response, err := makeRequest(req.createRequest())
	if err != nil {
		return nil, err
	}

	return response, nil
}

func makeRequest(req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (req request) createRequest() *http.Request {
	var reader io.Reader
	if req.Body != nil {
		reader = bytes.NewBuffer(req.Body)
	}

	r, _ := http.NewRequest(req.Method, req.URL, reader)

	if req.Header != "" {
		r.Header.Set("Internal-Token", req.Header)
	}

	return r
}

func getClient(sess *session.Session, decrypt bool) (client Client, err error) {
	logData := log.Data{"encryption_disabled_flag": cfg.EncryptionDisabled, "decrpyt_flag": decrypt}

	cryptoConfig := &s3crypto.Config{HasUserDefinedPSK: true}

	log.Debug("setting up s3 client", logData)

	c := ClientImpl{
		SDKClient: s3.New(sess),
	}

	if decrypt {
		c.CryptoClient = s3crypto.New(sess, cryptoConfig)
	}

	return c, nil
}

func createPSK() []byte {
	key := make([]byte, 16)
	rand.Read(key)

	return key
}
