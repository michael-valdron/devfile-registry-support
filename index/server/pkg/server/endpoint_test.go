package server_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/devfile/registry-support/index/server/pkg/server"
	"github.com/gin-gonic/gin"
)

func TestServeHealthCheck(t *testing.T) {
	var got gin.H
	var expected interface{}

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	server.ServeHealthCheck(c)

	expected = 200
	if w.Code != expected {
		t.Errorf("Did not get expected status code, Got: %v, Expected: %v", w.Code, expected)
		return
	}

	expected = "application/json"
	if header := w.Header(); header.Get("Content-Type") != expected {
		t.Errorf("Did not get expected content type, Got: %v, Expected: %v", header.Get("Content-Type"), expected)
		return
	}

	bytes, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
		return
	}

	if err = json.Unmarshal(bytes, &got); err != nil {
		t.Fatalf("Did not expect error: %v", err)
		return
	}

	expected = "the server is up and running"
	if got["message"] != expected {
		t.Errorf("Did not get expected body or message, Got: %v, Expected: %v", got["message"], expected)
		return
	}
}

func TestServeDevfileIndexV1(t *testing.T) {
	// TODO: Create testing data for ServeDevfileIndexV1 mock testing

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// TODO: Insert params

	server.ServeDevfileIndexV1(c)

	// TODO: Insert checks
}
