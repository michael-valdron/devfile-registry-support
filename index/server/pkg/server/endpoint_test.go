package server_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/devfile/registry-support/index/server/pkg/server"
	"github.com/gin-gonic/gin"
)

func TestServeHealthCheck(t *testing.T) {
	var got gin.H

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	server.ServeHealthCheck(c)

	wantStatusCode := 200
	if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, wantStatusCode) {
		t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, wantStatusCode)
		return
	}

	wantContentType := "application/json"
	header := w.Header()
	if gotContentType := strings.Split(header.Get("Content-Type"), ";")[0]; !reflect.DeepEqual(gotContentType, wantContentType) {
		t.Errorf("Did not get expected content type, Got: %v, Expected: %v", gotContentType, wantContentType)
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

	wantMessage := "the server is up and running"
	gotMessage, found := got["message"]
	if !found {
		t.Error("Did not get any body or message.")
		return
	} else if !reflect.DeepEqual(gotMessage, wantMessage) {
		t.Errorf("Did not get expected body or message, Got: %v, Expected: %v", gotMessage, wantMessage)
		return
	}
}

func TestServeDevfileIndexV1(t *testing.T) {
	// TODO: Create testing data for ServeDevfileIndexV1 mock testing
	tests := []struct {
		name     string
		params   gin.Params
		wantCode int
	}{
		{
			name: "Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "name", Value: "nodejs"},
				gin.Param{Key: "starterProjectName", Value: "nodejs-starter"},
			},
			wantCode: 200,
		},
		{
			name: "Not Found Response Test",
			params: gin.Params{
				gin.Param{Key: "name", Value: "node"},
			},
			wantCode: 404,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = test.params

			server.ServeDevfileIndexV1(c)

			// TODO: Insert checks
		})
	}
}
