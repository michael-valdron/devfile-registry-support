package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devfile/registry-support/index/server/pkg/util"
	"github.com/gin-gonic/gin"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	ociServerIP = "127.0.0.1:5000"
)

type responseError struct {
	code    string `json:"code"`
	message string `json:"message"`
	detail  string `json:"detail"`
}

func writeErrors(errors []responseError) ([]byte, error) {
	return json.Marshal(struct {
		errors []responseError `json:"errors"`
	}{
		errors: errors,
	})
}

func validateMethod(handle http.HandlerFunc, allowedMethods ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if arrayList := util.ConvertStringArrayToArrayList(allowedMethods); arrayList.Contains(r.Method) {
			handle(w, r)
		} else {
			bytes, err := writeErrors([]responseError{
				{
					code:    fmt.Sprintf("%d", http.StatusBadRequest),
					message: fmt.Sprintf("%s method not supported for route %s", r.Method, r.URL.Path),
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			if _, err = w.Write(bytes); err != nil {
				log.Fatal(err)
			}
		}
	})
}

func getManifest(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(ocispec.Manifest{})
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", ocispec.MediaTypeImageManifest)
	if _, err = w.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

func getBlob(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(ocispec.Descriptor{})
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", devfileMediaType)
	if _, err = w.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

func setupMockOCIServer() (func(), error) {
	router := http.NewServeMux()
	router.Handle("/v2/:name/manifests/:ref", validateMethod(getManifest, http.MethodGet, http.MethodHead))
	router.Handle("/v2/:name/blob/:digest", validateMethod(getBlob, http.MethodGet, http.MethodHead))
	testOCIServer := httptest.NewUnstartedServer(router)

	l, err := net.Listen("tcp", ociServerIP)
	if err != nil {
		return testOCIServer.Close, fmt.Errorf("Unexpected error while creating listener: %v", err)
	}

	testOCIServer.Listener.Close()
	testOCIServer.Listener = l

	testOCIServer.Start()

	return testOCIServer.Close, nil
}

func setupVars() {
	var registryPath string

	if _, found := os.LookupEnv("DEVFILE_REGISTRY"); found {
		registryPath = os.Getenv("DEVFILE_REGISTRY")
	} else {
		registryPath = "../../tests/registry"
	}

	if stacksPath == "" {
		stacksPath = filepath.Join(registryPath, "stacks")
	}
	if samplesPath == "" {
		samplesPath = filepath.Join(registryPath, "samples")
	}
	if indexPath == "" {
		indexPath = filepath.Join(registryPath, "index_main.json")
	}
	if sampleIndexPath == "" {
		sampleIndexPath = filepath.Join(registryPath, "index_extra.json")
	}
	if stackIndexPath == "" {
		stackIndexPath = filepath.Join(registryPath, "index_registry.json")
	}
}

func TestServeHealthCheck(t *testing.T) {
	var got gin.H

	gin.SetMode(gin.TestMode)

	setupVars()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serveHealthCheck(c)

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
	const wantStatusCode = 200

	setupVars()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serveDevfileIndexV1(c)

	if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, wantStatusCode) {
		t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, wantStatusCode)
		return
	}
}

func TestServeDevfileIndexV1WithType(t *testing.T) {
	setupVars()
	tests := []struct {
		name     string
		params   gin.Params
		wantCode int
	}{
		{
			name: "GET /index/stack - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "stack"},
			},
			wantCode: 200,
		},
		{
			name: "GET /index/sample - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "sample"},
			},
			wantCode: 200,
		},
		{
			name: "GET /index/all - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "all"},
			},
			wantCode: 200,
		},
		{
			name: "GET /index/notatype - Type Not Found Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "notatype"},
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

			serveDevfileIndexV1WithType(c)

			if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, test.wantCode) {
				t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, test.wantCode)
				return
			}
		})
	}
}

func TestServeDevfileIndexV2(t *testing.T) {
	const wantStatusCode = 200

	setupVars()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serveDevfileIndexV2(c)

	if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, wantStatusCode) {
		t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, wantStatusCode)
		return
	}
}

func TestServeDevfileIndexV2WithType(t *testing.T) {
	setupVars()
	tests := []struct {
		name     string
		params   gin.Params
		wantCode int
	}{
		{
			name: "GET /v2index/stack - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "stack"},
			},
			wantCode: 200,
		},
		{
			name: "GET /v2index/sample - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "sample"},
			},
			wantCode: 200,
		},
		{
			name: "GET /v2index/all - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "all"},
			},
			wantCode: 200,
		},
		{
			name: "GET /v2index/notatype - Type Not Found Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "notatype"},
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

			serveDevfileIndexV2WithType(c)

			if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, test.wantCode) {
				t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, test.wantCode)
				return
			}
		})
	}
}
