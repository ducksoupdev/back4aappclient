package object

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	back4appBaseUrl     = "https://parseapi.back4app.com"
	contentTypeHeader   = "Content-type"
	contentTypeValue    = "application/json"
	applicationIdHeader = "X-Parse-Application-Id"
	restApiKeyHeader    = "X-Parse-REST-API-Key"
	sessionTokenHeader  = "X-Parse-Session-Token"
)

type Object struct {
	httpClient    *http.Client
	baseUrl       *url.URL
	applicationId string
	restApiKey    string
	sessionToken  string
}

func NewObject(applicationId string, restApiKey string, sessionToken string, httpClient *http.Client, baseUrl *url.URL) *Object {
	c := &Object{
		httpClient:    httpClient,
		baseUrl:       baseUrl,
		applicationId: applicationId,
		restApiKey:    restApiKey,
		sessionToken:  sessionToken,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}
	if c.baseUrl == nil {
		c.baseUrl, _ = url.Parse(back4appBaseUrl)
	}
	return c
}

func (c *Object) Create(className string, data map[string]interface{}) (map[string]interface{}, error) {
	// create the URL
	createUrl, _ := url.Parse(fmt.Sprintf("/classes/%s", className))
	createClassUrl := c.baseUrl.ResolveReference(createUrl)

	// create the body
	marshalled, _ := json.Marshal(data)

	// create the request
	req, _ := http.NewRequest("POST", createClassUrl.String(), bytes.NewReader(marshalled))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusCreated {
		return nil, errors.New(fmt.Sprintf("unable to create object: %d", resp.StatusCode))
	}

	// parse the result
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Object) Delete(className string, id string) (bool, error) {
	// create the URL
	deleteUrl, _ := url.Parse(fmt.Sprintf("/classes/%s/%s", className, id))
	deleteClassUrl := c.baseUrl.ResolveReference(deleteUrl)

	// create the request
	req, _ := http.NewRequest("DELETE", deleteClassUrl.String(), nil)
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		return false, errors.New(fmt.Sprintf("unable to delete object: %d", resp.StatusCode))
	}

	return true, nil
}

func (c *Object) Read(className string, id string) (map[string]interface{}, error) {
	// create the URL
	readUrl, _ := url.Parse(fmt.Sprintf("/classes/%s/%s", className, id))
	readClassUrl := c.baseUrl.ResolveReference(readUrl)

	// create the request
	req, _ := http.NewRequest("GET", readClassUrl.String(), nil)
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unable to read object: %d", resp.StatusCode))
	}

	// parse the result
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Object) List(className string) (map[string][]map[string]interface{}, error) {
	// create the URL
	listUrl, _ := url.Parse(fmt.Sprintf("/classes/%s", className))
	listClassUrl := c.baseUrl.ResolveReference(listUrl)

	// create the request
	req, _ := http.NewRequest("GET", listClassUrl.String(), nil)
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	// check the status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unable to list object: %d", resp.StatusCode))
	}

	// parse the result
	var result map[string][]map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Object) Update(className string, id string, data map[string]interface{}) (bool, error) {
	// create the URL
	updateUrl, _ := url.Parse(fmt.Sprintf("/classes/%s/%s", className, id))
	updateClassUrl := c.baseUrl.ResolveReference(updateUrl)

	// create the body
	marshalled, _ := json.Marshal(data)

	// create the request
	req, _ := http.NewRequest("PUT", updateClassUrl.String(), bytes.NewReader(marshalled))
	req.Header.Add(contentTypeHeader, contentTypeValue)
	req.Header.Add(applicationIdHeader, c.applicationId)
	req.Header.Add(restApiKeyHeader, c.restApiKey)
	req.Header.Add(sessionTokenHeader, c.sessionToken)

	// make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// check the status code
	if resp.StatusCode != http.StatusOK {
		return false, errors.New(fmt.Sprintf("unable to update object: %d", resp.StatusCode))
	}

	return true, nil
}
