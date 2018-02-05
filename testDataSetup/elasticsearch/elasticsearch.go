package elasticsearch

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/rchttp"
)

// ErrorUnexpectedStatusCode represents the error message to be returned when
// the status received from elastic is not as expected
var ErrorUnexpectedStatusCode = errors.New("unexpected status code from elasticsearch api")

var client = rchttp.DefaultClient

// DeleteIndex removes a specified index from elasticsearch
func DeleteIndex(path string) error {
	logData := log.Data{"url": path}

	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorC("failed to create url for elasticsearch call", err, logData)
		return err
	}
	path = URL.String()

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		log.ErrorC("failed to create request for call to elasticsearch", err, logData)
		return err
	}

	resp, err := client.Do(context.Background(), req)
	if err != nil {
		log.ErrorC("failed to call elasticsearch", err, logData)
		return err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		log.Error(ErrorUnexpectedStatusCode, logData)
		return ErrorUnexpectedStatusCode
	}

	return nil
}
