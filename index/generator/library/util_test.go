package library

import (
	"net/http"
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
	bytes, err := DownloadStackFromZipUrl("")

	if err != nil {
		t.Errorf("Zip download to bytes failed: %v", err)
	}

	resultantType := http.DetectContentType(bytes)

	if resultantType != zipType {
		t.Errorf("Content type of download not matching expected. Expected: %s, Actual: %s",
			zipType, resultantType)
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
