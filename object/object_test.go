package object

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewObject(t *testing.T) {
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, nil)
	assert.NotNil(t, c.httpClient)
	assert.NotNil(t, c.baseUrl)
}

func TestInitialize(t *testing.T) {
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, nil)
	assert.NotNil(t, c.httpClient)
	assert.NotNil(t, c.baseUrl)
	assert.Equal(t, c.applicationId, "applicationId")
	assert.Equal(t, c.restApiKey, "restApiKey")
	assert.Equal(t, c.sessionToken, "sessionToken")
	assert.Equal(t, c.baseUrl.String(), "https://parseapi.back4app.com")
}

func TestCreate(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"username":"username", "objectId":"objectId"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	obj, _ := c.create("className", make(map[string]interface{}))
	assert.NotNil(t, obj)
	assert.Equal(t, "username", obj["username"])
}

func TestCreateError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	obj, err := c.create("className", make(map[string]interface{}))
	assert.Nil(t, obj)
	assert.Error(t, err)
	assert.Equal(t, "unable to create object: 400", err.Error())
}

func TestDelete(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isDeleted, _ := c.delete("className", "id")
	assert.True(t, isDeleted)
}

func TestDeleteError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isDeleted, err := c.delete("className", "id")
	assert.False(t, isDeleted)
	assert.Error(t, err)
	assert.Equal(t, "unable to delete object: 400", err.Error())
}

func TestRead(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"item":"item"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	item, _ := c.read("className", "id")
	assert.NotNil(t, item)
	assert.Equal(t, "item", item["item"])
}

func TestReadError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	item, err := c.read("className", "id")
	assert.Nil(t, item)
	assert.Error(t, err)
	assert.Equal(t, "unable to read object: 400", err.Error())
}

func TestList(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"results":[{"item":"item"},{"item2":"item2"}]}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	list, _ := c.list("className")
	assert.NotNil(t, list)
	assert.Len(t, list["results"], 2)
}

func TestListError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	list, err := c.list("className")
	assert.Nil(t, list)
	assert.Error(t, err)
	assert.Equal(t, "unable to list object: 400", err.Error())
}

func TestUpdate(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"updatedAt":"updatedAt"}`))
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isUpdated, _ := c.update("className", "id", make(map[string]interface{}))
	assert.True(t, isUpdated)
}

func TestUpdateError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer svr.Close()
	b, _ := url.Parse(svr.URL)
	c := NewObject("applicationId", "restApiKey", "sessionToken", nil, b)
	isUpdated, err := c.update("className", "id", make(map[string]interface{}))
	assert.False(t, isUpdated)
	assert.Error(t, err)
	assert.Equal(t, "unable to update object: 400", err.Error())
}
