package library

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	zipType string = "application/zip"
)

func TestDownloadRemoteStack(t *testing.T) {
	assert.Fail(t, "Not Implemented.")
}

func TestDownloadStackFromZipUrl(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
	}{
		{
			"Case 1: Java Quarkus (Without subDir)",
			map[string]string{
				"ZipUrl": "https://code.quarkus.io/d?e=io.quarkus%3Aquarkus-resteasy&e=io.quarkus%3Aquarkus-micrometer&e=io.quarkus%3Aquarkus-smallrye-health&e=io.quarkus%3Aquarkus-openshift&cn=devfile",
				"SubDir": "",
			},
		},
		// {
		// 	"Case 2: With subDir",
		// 	map[string]string{
		// 		"ZipUrl": "",
		// 		"SubDir": "",
		// 	},
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zipUrlBasename := filepath.Base(tt.params["ZipUrl"])
			zipUrlBasename = strings.ReplaceAll(zipUrlBasename, filepath.Ext(zipUrlBasename), "")
			zipPath := filepath.Join(os.TempDir(), zipUrlBasename)
			bytes, err := DownloadStackFromZipUrl(tt.params["ZipUrl"], tt.params["SubDir"], zipPath)

			if err != nil {
				t.Errorf("Zip download to bytes failed: %v", err)
			}

			resultantType := http.DetectContentType(bytes)

			if resultantType != zipType {
				t.Errorf("Content type of download not matching expected. Expected: %s, Actual: %s",
					zipType, resultantType)
			}
		})
	}
}

func TestGetSubDir(t *testing.T) {
	assert.Fail(t, "Not Implemented.")
}

func TestCopyFileWithFs(t *testing.T) {
	assert.Fail(t, "Not Implemented.")
}

func TestCopyDirWithFS(t *testing.T) {
	assert.Fail(t, "Not Implemented.")
}

func TestCleanDir(t *testing.T) {
	assert.Fail(t, "Not Implemented.")
}

func TestZipDir(t *testing.T) {
	assert.Fail(t, "Not Implemented.")
}
