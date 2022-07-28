package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/devfile/library/pkg/devfile/parser"
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
						Digest:    "sha256:b81a4a857ebbd6b7093c38703e3b7c6d7a2652abfd55898f82fdea45634fd549",
						Size:      1251,
						Annotations: map[string]string{
							"org.opencontainers.image.title": devfileName,
						},
					},
				},
			},
		},
		"go": {
			"1.1.0": {
				Versioned: specs.Versioned{SchemaVersion: 2},
				Config: ocispec.Descriptor{
					MediaType: devfileConfigMediaType,
				},
				Layers: []ocispec.Descriptor{
					{
						MediaType: devfileMediaType,
						Digest:    "sha256:4cad7c1629ba848245205a08107b296adba307f77ca3635b16920473589cb12e",
						Size:      1086,
						Annotations: map[string]string{
							"org.opencontainers.image.title": devfileName,
						},
					},
				},
			},
			"1.2.0": {
				Versioned: specs.Versioned{SchemaVersion: 2},
				Config: ocispec.Descriptor{
					MediaType: devfileConfigMediaType,
				},
				Layers: []ocispec.Descriptor{
					{
						MediaType: devfileMediaType,
						Digest:    "sha256:bb4c6b96292bbcd48f445436f7945399a4d314b111ee976d6235199e854bfb68",
						Size:      1091,
						Annotations: map[string]string{
							"org.opencontainers.image.title": devfileName,
						},
					},
				},
			},
		},
		"java-quarkus": {
			"1.1.0": {
				Versioned: specs.Versioned{SchemaVersion: 2},
				Config: ocispec.Descriptor{
					MediaType: devfileConfigMediaType,
				},
				Layers: []ocispec.Descriptor{
					{
						MediaType: devfileMediaType,
						Digest:    "sha256:6143ffeab064d0b9a6991a6e1025a1b63d01be87e197fabb8a34d5a9fc3fcba0",
						Size:      1935,
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
	var (
		blobPath string
		found    bool
		err      error
	)

	found = false
	err = filepath.WalkDir(stackRoot, func(path string, d fs.DirEntry, err error) error {
		var fdgst string

		if err != nil {
			return err
		}

		if found || d.IsDir() {
			return nil
		}

		fdgst, err = digestFile(path)
		if err != nil {
			return err
		} else if reflect.DeepEqual(dgst, fdgst) {
			blobPath = path
			found = true
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	} else if !found {
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

	if !reflect.DeepEqual(resp.StatusCode, http.StatusOK) {
		t.Errorf("Did not get expected status code, Got: %v, Expected: %v", resp.StatusCode, http.StatusOK)
	}
}

func TestServeHealthCheck(t *testing.T) {
	var got gin.H

	gin.SetMode(gin.TestMode)

	setupVars()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	serveHealthCheck(c)

	wantStatusCode := http.StatusOK
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
	const wantStatusCode = http.StatusOK

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
			wantCode: http.StatusOK,
		},
		{
			name: "GET /index/sample - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "sample"},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "GET /index/all - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "all"},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "GET /index/notatype - Type Not Found Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "notatype"},
			},
			wantCode: http.StatusNotFound,
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
	const wantStatusCode = http.StatusOK

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
			wantCode: http.StatusOK,
		},
		{
			name: "GET /v2index/sample - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "sample"},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "GET /v2index/all - Successful Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "all"},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "GET /v2index/notatype - Type Not Found Response Test",
			params: gin.Params{
				gin.Param{Key: "type", Value: "notatype"},
			},
			wantCode: http.StatusNotFound,
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
		name              string
		params            gin.Params
		wantCode          int
		wantSchemaVersion string
		wantError         bool
	}{
		{
			name: "GET /devfiles/java-maven - Fetch Java Maven Devfile",
			params: gin.Params{
				gin.Param{Key: "name", Value: "java-maven"},
			},
			wantCode:          http.StatusOK,
			wantSchemaVersion: "2.2.0",
		},
		{
			name: "GET /devfiles/go - Fetch Go Devfile",
			params: gin.Params{
				gin.Param{Key: "name", Value: "go"},
			},
			wantCode:          http.StatusOK,
			wantSchemaVersion: "2.0.0",
		},
		{
			name: "GET /devfiles/not-exist - Fetch Non-Existent Devfile",
			params: gin.Params{
				gin.Param{Key: "name", Value: "not-exist"},
			},
			wantCode:  http.StatusNotFound,
			wantError: true,
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
			} else if !test.wantError {
				bytes := w.Body.Bytes()
				content, err := parser.ParseFromData(bytes)
				if err != nil {
					t.Fatalf("Did not expect error: %v", err)
				}

				if gotSchemaVersion := content.Data.GetSchemaVersion(); !reflect.DeepEqual(gotSchemaVersion, test.wantSchemaVersion) {
					t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotSchemaVersion, test.wantSchemaVersion)
				}
			}
		})
	}
}

func TestServeDevfileWithVersion(t *testing.T) {
	tests := []struct {
		name              string
		params            gin.Params
		wantCode          int
		wantSchemaVersion string
		wantError         bool
	}{
		{
			name: "GET /devfiles/go/default - Fetch Go Devfile With Default Version",
			params: gin.Params{
				gin.Param{Key: "name", Value: "go"},
				gin.Param{Key: "version", Value: "default"},
			},
			wantCode:          http.StatusOK,
			wantSchemaVersion: "2.0.0",
		},
		{
			name: "GET /devfiles/go/latest - Fetch Go Devfile With Latest Version",
			params: gin.Params{
				gin.Param{Key: "name", Value: "go"},
				gin.Param{Key: "version", Value: "latest"},
			},
			wantCode:          http.StatusOK,
			wantSchemaVersion: "2.1.0",
		},
		{
			name: "GET /devfiles/go/1.2.0 - Fetch Go Devfile With Specific Version",
			params: gin.Params{
				gin.Param{Key: "name", Value: "go"},
				gin.Param{Key: "version", Value: "1.2.0"},
			},
			wantCode:          http.StatusOK,
			wantSchemaVersion: "2.1.0",
		},
		{
			name: "GET /devfiles/not-exist/latest - Fetch Non-Existent Devfile With Latest Version",
			params: gin.Params{
				gin.Param{Key: "name", Value: "not-exist"},
				gin.Param{Key: "version", Value: "latest"},
			},
			wantCode:  http.StatusNotFound,
			wantError: true,
		},
		{
			name: "GET /devfiles/java-maven/not-exist - Fetch Java Maven Devfile With Non-Existent Version",
			params: gin.Params{
				gin.Param{Key: "name", Value: "java-maven"},
				gin.Param{Key: "version", Value: "non-exist"},
			},
			wantCode:  http.StatusNotFound,
			wantError: true,
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

			serveDevfileWithVersion(c)

			if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, test.wantCode) {
				t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, test.wantCode)
			} else if !test.wantError {
				bytes := w.Body.Bytes()
				content, err := parser.ParseFromData(bytes)
				if err != nil {
					t.Fatalf("Did not expect error: %v", err)
				}

				if gotSchemaVersion := content.Data.GetSchemaVersion(); !reflect.DeepEqual(gotSchemaVersion, test.wantSchemaVersion) {
					t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotSchemaVersion, test.wantSchemaVersion)
				}
			}
		})
	}
}

func TestServeDevfileStarterProject(t *testing.T) {
	const wantContentType = "application/zip"
	tests := []struct {
		name      string
		params    gin.Params
		wantCode  int
		wantError bool
	}{
		{
			name: "GET /devfiles/java-maven/starter-projects/springbootproject - Fetch Java Maven 'springbootproject' Starter Project",
			params: gin.Params{
				gin.Param{Key: "name", Value: "java-maven"},
				gin.Param{Key: "starterProjectName", Value: "springbootproject"},
			},
			wantCode: http.StatusAccepted,
		},
		{
			name: "GET /devfiles/go/starter-projects/go-starter - Fetch Go 'go-starter' Starter Project",
			params: gin.Params{
				gin.Param{Key: "name", Value: "go"},
				gin.Param{Key: "starterProjectName", Value: "go-starter"},
			},
			wantCode: http.StatusAccepted,
		},
		{
			name: "GET /devfiles/java-quarkus/starter-projects/community - Fetch Go 'community' Starter Project",
			params: gin.Params{
				gin.Param{Key: "name", Value: "java-quarkus"},
				gin.Param{Key: "starterProjectName", Value: "community"},
			},
			wantCode: http.StatusAccepted,
		},
		{
			name: "GET /devfiles/not-exist/starter-projects/some - Fetch 'some' starter project from Non-Existent stack",
			params: gin.Params{
				gin.Param{Key: "name", Value: "not-exist"},
				gin.Param{Key: "starterProjectName", Value: "some"},
			},
			wantCode:  http.StatusNotFound,
			wantError: true,
		},
		{
			name: "GET /devfiles/java-maven/starter-projects/not-exist - Fetch Non-Existent starter project from Java Maven stack",
			params: gin.Params{
				gin.Param{Key: "name", Value: "java-maven"},
				gin.Param{Key: "starterProjectName", Value: "not-exist"},
			},
			wantCode:  http.StatusNotFound,
			wantError: true,
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

			serveDevfileStarterProject(c)

			if gotStatusCode := w.Code; !reflect.DeepEqual(gotStatusCode, test.wantCode) {
				t.Errorf("Did not get expected status code, Got: %v, Expected: %v", gotStatusCode, test.wantCode)
			} else if !test.wantError {
				gotContentType := http.DetectContentType(w.Body.Bytes())
				if !reflect.DeepEqual(gotContentType, wantContentType) {
					t.Errorf("Did not get expected content-type, Got: %v, Expected: %v", gotContentType, wantContentType)
				}
			}
		})
	}
}
