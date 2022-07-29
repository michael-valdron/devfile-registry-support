package ocitest

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type ResponseErrorDetails struct {
	Tag string `json:"Tag"`
}

type ResponseError struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Detail  ResponseErrorDetails `json:"detail"`
}

type MockOCIServer struct {
	httpserver    *httptest.Server
	router        *gin.Engine
	ServeManifest func(c *gin.Context)
	ServeBlob     func(c *gin.Context)
}

func servePing(c *gin.Context) {
	data, err := json.Marshal(gin.H{
		"message": "ok",
	})
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, data)
}

func WriteErrors(errors []ResponseError) gin.H {
	return gin.H{
		"errors": errors,
	}
}

func NewMockOCIServer() *MockOCIServer {
	gin.SetMode(gin.TestMode)

	mockOCIServer := &MockOCIServer{
		// Create router engine of mock OCI server
		router: gin.Default(),
	}

	// Create mock OCI server using the router engine
	mockOCIServer.httpserver = httptest.NewUnstartedServer(mockOCIServer.router)

	return mockOCIServer
}

func (server *MockOCIServer) Start(listenAddr string) error {
	// Testing Route for checking mock OCI server
	server.router.GET("/v2/ping", servePing)

	// Pull Routes
	// Fetch manifest routes
	if server.ServeManifest != nil {
		server.router.GET("/v2/devfile-catalog/:name/manifests/:ref", server.ServeManifest)
		server.router.HEAD("/v2/devfile-catalog/:name/manifests/:ref", server.ServeManifest)
	}

	// Fetch blob routes
	if server.ServeBlob != nil {
		server.router.GET("/v2/devfile-catalog/:name/blobs/:digest", server.ServeBlob)
		server.router.HEAD("/v2/devfile-catalog/:name/blobs/:digest", server.ServeBlob)
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("unexpected error while creating listener: %v", err)
	}

	server.httpserver.Listener.Close()
	server.httpserver.Listener = l

	server.httpserver.Start()

	return nil
}

func (server *MockOCIServer) Close() {
	server.httpserver.Close()
}

type ProxyRecorder struct {
	*httptest.ResponseRecorder
	http.CloseNotifier
}

func NewProxyRecorder() *ProxyRecorder {
	return &ProxyRecorder{
		ResponseRecorder: httptest.NewRecorder(),
	}
}

func (rec *ProxyRecorder) CloseNotify() <-chan bool {
	return make(<-chan bool)
}
