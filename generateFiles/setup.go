package generateFiles

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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

func sendV4FileToAWS(region, bucket, filename string) error {

	config := aws.NewConfig().WithRegion(region)

	store := &Store{
		config: config,
		bucket: bucket,
	}

	session, err := session.NewSession(store.config)
	if err != nil {
		return err
	}

	v4File, err := os.Open(filename)
	if err != nil {
		return err
	}

	uploader := s3manager.NewUploader(session)

	uploader.Upload(&s3manager.UploadInput{})

	// the AWS uploader automatically handles large files breaking them into parts and using the multi part API.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   v4File,
		Bucket: &store.bucket,
		Key:    &filename,
	})
	if err != nil {
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
