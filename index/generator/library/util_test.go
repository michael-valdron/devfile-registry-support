package library

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devfile/registry-support/index/generator/schema"
	"github.com/stretchr/testify/assert"
)

const (
	zipType string = "application/zip"
)

func TestDownloadRemoteStack(t *testing.T) {
	tests := []struct {
		name string
		git  *schema.Git
		path string
	}{
		{
			"Case 1: Maven Java (Without subDir)",
			&schema.Git{
				Url:        "https://github.com/odo-devfiles/springboot-ex.git",
				RemoteName: "origin",
			},
			filepath.Join(os.TempDir(), "springboot-ex"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hiddenGitPath := filepath.Join(tt.path, ".git")

			if err := DownloadRemoteStack(tt.git, tt.path, false); err != nil {
				t.Errorf("Git download to bytes failed: %v", err)
			}

			if _, err := os.Stat(tt.path); os.IsNotExist(err) {
				t.Errorf("%s does not exist but is suppose to", tt.path)
			} else if _, err := os.Stat(hiddenGitPath); os.IsExist(err) {
				t.Errorf(".git exist but isn't suppose to within %s", tt.path)
			}

			if err := os.RemoveAll(tt.path); err != nil {
				t.Logf("Deleting %s failed.", tt.path)
			}
		})
	}
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

func TestDownloadStackFromGit(t *testing.T) {
	tests := []struct {
		name string
		git  *schema.Git
		path string
	}{
		{
			"Case 1: Maven Java (Without subDir)",
			&schema.Git{
				Url:        "https://github.com/odo-devfiles/springboot-ex.git",
				RemoteName: "origin",
			},
			filepath.Join(os.TempDir(), "springboot-ex"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hiddenGitPath := filepath.Join(tt.path, ".git")
			bytes, err := DownloadStackFromGit(tt.git, tt.path, false)

			if err != nil {
				t.Errorf("Git download to bytes failed: %v", err)
			} else if _, err := os.Stat(hiddenGitPath); os.IsExist(err) {
				t.Errorf(".git exist but isn't suppose to within %s", hiddenGitPath)
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
