package generateFiles

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var client = http.Client{}

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
func sendV4FileToAWS(region, bucket, filename string, decrypt bool) error {
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

	client, err := getClient(sess, decrypt)
	if err != nil {
		log.ErrorC("failed to create client", err, nil)
		return err
	}

	putObject := &s3.PutObjectInput{
		Key:    &filename,
		Bucket: &bucket,
		Body:   v4File,
	}

	_, err = client.PutObject(putObject)
	if err != nil {
		log.ErrorC("failed to upload file", err, nil)
		return err
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

	output, err := client.GetObject(input)
	if err != nil {
		log.ErrorC("encountered error retrieving csv file", err, nil)
		return nil, err
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

func getClient(sess *session.Session, decrypt bool) (client s3iface.S3API, err error) {
	logData := log.Data{"encryption_disabled_flag": cfg.EncryptionDisabled, "decrpyt_flag": decrypt}

	// Encrypt v4 file with PrivateKey
	if !cfg.EncryptionDisabled && decrypt {
		log.Debug("accessing encryption/decryption client", logData)
		privateKey, err := getPrivateKey([]byte(cfg.PrivateKey))
		if err != nil {
			return nil, err
		}
		cryptoConfig := &s3crypto.Config{PrivateKey: privateKey}
		client = s3crypto.New(sess, cryptoConfig)
	} else {
		log.Debug("regular client", logData)
		client = s3.New(sess)
	}

	return client, nil
}

func getPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("invalid RSA PRIVATE KEY provided")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
