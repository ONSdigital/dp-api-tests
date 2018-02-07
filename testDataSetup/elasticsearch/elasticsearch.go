package elasticsearch

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/rchttp"
)

type Index struct {
	InstanceID   string
	Dimension    string
	TestDataFile string
	URL          string
	MappingsFile string
}

var client = rchttp.DefaultClient

// ErrorUnexpectedStatusCode represents the error message to be returned when
// the status received from elastic is not as expected
var ErrorUnexpectedStatusCode = errors.New("unexpected status code from elasticsearch api")

// DeleteIndex removes a specified index from elasticsearch
func DeleteIndex(path string) (int, error) {
	logData := log.Data{"url": path}

	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorC("failed to create url for elasticsearch call", err, logData)
		return 0, err
	}
	path = URL.String()

	req, err := http.NewRequest("DELETE", path, nil)
	if err != nil {
		log.ErrorC("failed to create request for call to elasticsearch", err, logData)
		return 0, err
	}

	resp, err := client.Do(context.Background(), req)
	if err != nil {
		log.ErrorC("failed to call elasticsearch", err, logData)
		return 0, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode
	if resp.StatusCode != http.StatusOK {
		log.Error(ErrorUnexpectedStatusCode, logData)
		return resp.StatusCode, ErrorUnexpectedStatusCode
	}

	return http.StatusOK, nil
}

// Dimension represents the data stored in an elasticsearch document
type Dimension struct {
	Code             string `json:"code"`
	HasData          bool   `json:"has_data"`
	Label            string `json:"label"`
	NumberOfChildren int    `json:"number_of_children"`
	URL              string `json:"url"`
}

// CreateSearchIndex represents the creation and loading of test data into an index for testing
func (i *Index) CreateSearchIndex() error {
	index := i.URL + "/" + i.InstanceID + "_" + i.Dimension
	// Remove index
	statusCode, err := DeleteIndex(index)
	if err != nil {
		if statusCode != http.StatusNotFound {
			log.ErrorC("failed to delete index", err, log.Data{"path": index})
			return err
		}
	}

	ctx := context.Background()

	// Create index
	indexMappings, err := ioutil.ReadFile(i.MappingsFile)
	if err != nil {
		return err
	}

	_, _, err = CallElastic(ctx, index, "PUT", indexMappings)
	if err != nil {
		log.ErrorC("fail to create index", err, log.Data{"path": index})
		return err
	}

	// Add docs to index
	file, err := os.Open(i.TestDataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slice := []byte(line)
		var dimension *Dimension
		if err = json.Unmarshal(slice, &dimension); err != nil {
			log.ErrorC("unable to unmarshal bytes to json", err, log.Data{"line": line})
			return err
		}
		path := index + "/hierarchy/" + dimension.Code
		_, _, err = CallElastic(ctx, path, "PUT", slice)
		if err != nil {
			log.ErrorC("encountered error writing to elasticsearch index", err, log.Data{"json_file": i.TestDataFile, "payload": scanner.Text()})
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	log.Info("successfully loaded data into elasticsearch", log.Data{"json_file": i.TestDataFile})
	return nil
}

// CallElastic builds a request to elastic search based on the method, path and payload
func CallElastic(ctx context.Context, path, method string, payload interface{}) ([]byte, int, error) {

	URL, err := url.Parse(path)
	if err != nil {
		log.ErrorC("failed to create url for elastic call", err, nil)
		return nil, 0, err
	}
	path = URL.String()
	logData := log.Data{"url": path}

	var req *http.Request

	if payload != nil {
		req, err = http.NewRequest(method, path, bytes.NewReader(payload.([]byte)))
		req.Header.Add("Content-type", "application/json")
		logData["payload"] = string(payload.([]byte))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	// check req, above, didn't error
	if err != nil {
		log.ErrorC("failed to create request for call to elastic", err, logData)
		return nil, 0, err
	}

	resp, err := client.Do(ctx, req)
	if err != nil {
		log.ErrorC("failed to call elastic", err, logData)
		return nil, 0, err
	}
	defer resp.Body.Close()

	logData["http_code"] = resp.StatusCode

	jsonBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.ErrorC("failed to read response body from call to elastic", err, logData)
		return nil, resp.StatusCode, err
	}

	logData["json_body"] = string(jsonBody)
	logData["status_code"] = resp.StatusCode

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		log.ErrorC("failed", ErrorUnexpectedStatusCode, logData)
		return nil, resp.StatusCode, ErrorUnexpectedStatusCode
	}

	return jsonBody, resp.StatusCode, nil
}
