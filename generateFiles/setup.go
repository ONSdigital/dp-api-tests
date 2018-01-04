package generateFiles

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/ONSdigital/go-ns/log"
	ons3 "github.com/ONSdigital/go-ns/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func sendV4FileToAWS(region, bucket, filename string) (string, error) {

	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	session, err := session.NewSession(store.config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return "", err
	}

	v4File, err := os.Open(filename)
	if err != nil {
		log.ErrorC("failed to open file", err, nil)
		return "", err
	}

	log.Info("successfully retrieved file", nil)

	uploader := s3manager.NewUploader(session)

	// the AWS uploader automatically handles large files breaking them into
	//  parts and using the multi part API.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Body:   v4File,
		Bucket: &store.bucket,
		Key:    &filename,
	})
	if err != nil {
		log.ErrorC("failed to upload file", err, nil)
		return "", err
	}

	return result.Location, nil
}

func getS3File(region, s3URL string) (io.ReadCloser, error) {
	s3, err := ons3.New(region)
	if err != nil {
		log.Error(err, nil)
		return nil, err
	}

	file, err := s3.Get(s3URL)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func getS3FileSize(region, bucket, filename string) (*int64, error) {
	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	session, err := session.NewSession(store.config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return nil, err
	}

	svc := s3.New(session)

	input := &s3.GetObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(bucket),
	}

	ctx := context.Background()
	result, err := svc.GetObjectWithContext(ctx, input)
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
