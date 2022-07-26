package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devfile/registry-support/index/server/pkg/ocitest"
	"github.com/gin-gonic/gin"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

const (
	ociServerIP = "127.0.0.1:5000"
)

var (
	manifests = map[string]map[string]ocispec.Manifest{
		"java-maven": {
			"1.1.0": {
				Versioned: specs.Versioned{SchemaVersion: 2},
				Config: ocispec.Descriptor{
					MediaType: devfileConfigMediaType,
				},
				Layers: []ocispec.Descriptor{
					{
						MediaType: devfileMediaType,
						Digest:    "b81a4a857ebbd6b7093c38703e3b7c6d7a2652abfd55898f82fdea45634fd549",
						Size:      1251,
						Annotations: map[string]string{
							"org.opencontainers.image.title": devfileName,
						},
					},
				},
			},
		},
	}
)

func serveManifest(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		getManifest(c)
	}
}

func getManifest(c *gin.Context) {
	name, ref := c.Param("name"), c.Param("ref")
	var (
		stackManifest ocispec.Manifest
		found         bool
		bytes         []byte
		err           error
	)

	if strings.HasPrefix(ref, "sha256:") {
		stackManifests, found := manifests[name]
		if !found {
			notFoundManifest(c, ref)
			return
		}

		found = false
		for _, manifest := range stackManifests {
			dgst, err := digestEntity(manifest)
			if err != nil {
				log.Fatal(err)
			} else if reflect.DeepEqual(ref, dgst) {
				stackManifest = manifest
				found = true
				break
			}
		}

		if !found {
			notFoundManifest(c, ref)
			return
		}
	} else {
		stackManifest, found = manifests[name][ref]

		if !found {
			notFoundManifest(c, ref)
			return
		}
	}

	bytes, err = json.Marshal(stackManifest)
	if err != nil {
		log.Fatal(err)
	}

	c.Data(http.StatusOK, ocispec.MediaTypeImageManifest, bytes)
}

func notFoundManifest(c *gin.Context, tag string) {
	c.JSON(http.StatusNotFound, ocitest.WriteErrors([]ocitest.ResponseError{
		{
			Code:    "MANIFEST_UNKNOWN",
			Message: "manifest unknown",
			Detail: ocitest.ResponseErrorDetails{
				Tag: tag,
			},
		},
	}))
}

func digestEntity(e interface{}) (string, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}

	return digest.FromBytes(bytes).String(), nil
}

func digestFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	dgst, err := digest.FromReader(file)
	if err != nil {
		return "", err
	}

	return dgst.String(), nil
}

func serveBlob(c *gin.Context) {
	name, dgst := c.Param("name"), c.Param("digest")
	stackRoot := filepath.Join(stacksPath, name)
	stackRootList, err := ioutil.ReadDir(stackRoot)
	if err != nil {
		log.Fatal(err)
	}
	var (
		blobPath string
		found    bool
	)

	found = false
	for _, stackFile := range stackRootList {
		fpath := filepath.Join(stackRoot, stackFile.Name())
		fdgst, err := digestFile(fpath)
		if err != nil {
			log.Fatal(err)
		} else if reflect.DeepEqual(dgst, strings.TrimLeft(fdgst, "sha256:")) {
			blobPath = fpath
			found = true
			break
		}
	}

	if !found {
		notFoundBlob(c)
		return
	}

	file, err := os.Open(blobPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	c.Data(http.StatusOK, http.DetectContentType(bytes), bytes)
}

func notFoundBlob(c *gin.Context) {
	c.Data(http.StatusNotFound, "plain/text", []byte{})
}

func setupMockOCIServer() (func(), error) {
	mockOCIServer := ocitest.NewMockOCIServer()

	// Pull Routes
	mockOCIServer.ServeManifest = serveManifest
	mockOCIServer.ServeBlob = serveBlob

	if err := mockOCIServer.Start(ociServerIP); err != nil {
		return nil, err
	}

	return mockOCIServer.Close, nil
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

func TestMockOCIServer(t *testing.T) {
	mockOCIServer := ocitest.NewMockOCIServer()
	if err := mockOCIServer.Start(ociServerIP); err != nil {
		t.Errorf("Failed to setup mock OCI server: %v", err)
		return
	}
	defer mockOCIServer.Close()
	setupVars()

	resp, err := http.Get(fmt.Sprintf("http://%s", filepath.Join(ociServerIP, "/v2/ping")))
	if err != nil {
		t.Errorf("Error in request: %v", err)
		return
	}

	if !reflect.DeepEqual(resp.StatusCode, 200) {
		t.Errorf("Did not get expected status code, Got: %v, Expected: %v", resp.StatusCode, 200)
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

			c.Params = append(c.Params, test.params...)

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

			c.Params = append(c.Params, test.params...)

			serveDevfileIndexV2WithType(c)

			if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, test.wantCode) {
				t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, test.wantCode)
				return
			}
		})
	}
}

func TestServeDevfile(t *testing.T) {
	tests := []struct {
		name     string
		params   gin.Params
		wantCode int
	}{
		{
			name: "Fetch Devfile",
			params: gin.Params{
				gin.Param{Key: "name", Value: "java-maven"},
			},
			wantCode: 200,
		},
	}

	closeServer, err := setupMockOCIServer()
	if err != nil {
		t.Errorf("Did not setup mock OCI server properly: %v", err)
		return
	}
	defer closeServer()
	setupVars()

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = append(c.Params, test.params...)

			serveDevfile(c)

			if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, test.wantCode) {
				t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, test.wantCode)
				return
			}
		})
	}
}
