package server

import (
	"byfood-interview/book"
	"byfood-interview/helper"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type HTTPTestSuite struct {
	server    *Server
	container testcontainers.Container
	db        *sqlx.DB
}

func setupHTTPTestSuite(t *testing.T) *HTTPTestSuite {
	ctx := context.Background()

	// Setup test container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "book_test",
			"POSTGRES_USER":     "test_user",
			"POSTGRES_PASSWORD": "test_password",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithPollInterval(5 * time.Second).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Set environment variables for database connection
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_NAME", "book_test")

	// Wait a bit more for database to be fully ready
	time.Sleep(5 * time.Second) // Adjust as necessary for your environment
	// Create server instance
	server := NewServer("../migration/file")

	return &HTTPTestSuite{
		server:    server,
		container: container,
		db:        server.DB,
	}
}

func (suite *HTTPTestSuite) tearDown(t *testing.T) {
	ctx := context.Background()
	if suite.db != nil {
		suite.db.Close()
	}
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func TestHealthCheck(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	req, err := http.NewRequest("GET", "/api/v1/healthcheck", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestCreateBook(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	tests := []struct {
		name           string
		book           book.Book
		expectedStatus int
		expectError    bool
	}{
		{
			name: "Valid book creation",
			book: book.Book{
				Title:         "Test Book",
				Author:        "Test Author",
				PublishedYear: 2023,
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "Missing title",
			book: book.Book{
				Author:        "Test Author",
				PublishedYear: 2023,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Missing author",
			book: book.Book{
				Title:         "Test Book",
				PublishedYear: 2023,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name: "Invalid published year",
			book: book.Book{
				Title:         "Test Book",
				Author:        "Test Author",
				PublishedYear: 0,
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.book)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if !tt.expectError {
				var response helper.Response
				err = json.Unmarshal(rr.Body.Bytes(), &response)

				books := response.Data.(map[string]interface{})

				require.NoError(t, err)
				assert.Equal(t, tt.book.Title, books["title"])
				assert.Equal(t, tt.book.Author, books["author"])
				assert.Equal(t, tt.book.PublishedYear, int(books["published_year"].(float64)))
				assert.NotZero(t, books["id"])
			}
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	// Create test books first
	testBooks := []book.Book{
		{Title: "Book 1", Author: "Author 1", PublishedYear: 2021},
		{Title: "Book 2", Author: "Author 2", PublishedYear: 2022},
		{Title: "Book 3", Author: "Author 3", PublishedYear: 2023},
	}

	var createdBooksCount int
	for _, b := range testBooks {
		jsonBody, err := json.Marshal(b)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusOK, rr.Code)

		createdBooksCount++
	}

	// Test get all books
	req, err := http.NewRequest("GET", "/api/v1/books", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	response := helper.Response{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	books := response.Data.([]interface{})
	assert.Len(t, books, len(testBooks))
}

func TestGetBookByID(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	// Create a test book first
	testBook := book.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		PublishedYear: 2023,
	}

	jsonBody, err := json.Marshal(testBook)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var response helper.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	createdData := response.Data.(map[string]interface{})

	tests := []struct {
		name           string
		bookID         string
		expectedStatus int
	}{
		{
			name:           "Valid book ID",
			bookID:         strconv.FormatFloat(createdData["id"].(float64), 'f', 0, 64),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid book ID",
			bookID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Non-numeric book ID",
			bookID:         "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/api/v1/books/"+tt.bookID, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response helper.Response
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				fetchedBook := response.Data.(map[string]interface{})

				require.NoError(t, err)
				assert.Equal(t, createdData["id"], fetchedBook["id"])
				assert.Equal(t, createdData["title"], fetchedBook["title"])
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	// Create a test book first
	testBook := book.Book{
		Title:         "Original Title",
		Author:        "Original Author",
		PublishedYear: 2023,
	}

	jsonBody, err := json.Marshal(testBook)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var response helper.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	createdData := response.Data.(map[string]interface{})

	tests := []struct {
		name           string
		bookID         string
		updateData     book.Book
		expectedStatus int
	}{
		{
			name:   "Valid update",
			bookID: strconv.FormatFloat(createdData["id"].(float64), 'f', 0, 64),
			updateData: book.Book{
				Title:         "Updated Title",
				Author:        "Updated Author",
				PublishedYear: 2024,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid book ID",
			bookID: "999999",
			updateData: book.Book{
				Title:         "Updated Title",
				Author:        "Updated Author",
				PublishedYear: 2024,
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.updateData)
			require.NoError(t, err)

			req, err := http.NewRequest("PUT", "/api/v1/books/"+tt.bookID, bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response helper.Response
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				fetchedBook := response.Data.(map[string]interface{})

				assert.Equal(t, tt.updateData.Title, fetchedBook["title"])
				assert.Equal(t, tt.updateData.Author, fetchedBook["author"])
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	// Create a test book first
	testBook := book.Book{
		Title:         "Book to Delete",
		Author:        "Author",
		PublishedYear: 2023,
	}

	jsonBody, err := json.Marshal(testBook)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var response helper.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	createdData := response.Data.(map[string]interface{})

	tests := []struct {
		name           string
		bookID         string
		expectedStatus int
	}{
		{
			name:           "Valid deletion",
			bookID:         strconv.FormatFloat(createdData["id"].(float64), 'f', 0, 64),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid book ID",
			bookID:         "999999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/api/v1/books/"+tt.bookID, nil)
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			// If deletion was successful, verify book is gone
			if tt.expectedStatus == http.StatusOK {
				req, err := http.NewRequest("GET", "/api/v1/books/"+tt.bookID, nil)
				require.NoError(t, err)

				rr := httptest.NewRecorder()
				suite.server.Router.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusNotFound, rr.Code)
			}
		})
	}
}

func TestProcessURL(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid URL processing",
			requestBody: map[string]interface{}{
				"url":       "https://example.com/path?param=value",
				"operation": "all",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid URL",
			requestBody: map[string]interface{}{
				"url": "not-a-valid-url",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing URL",
			requestBody:    map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/process-url", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

func TestRateLimitingAndConcurrency(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	// Test concurrent requests
	concurrentRequests := 10
	results := make(chan int, concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func(id int) {
			testBook := book.Book{
				Title:         fmt.Sprintf("Concurrent Book %d", id),
				Author:        "Concurrent Author",
				PublishedYear: 2023,
			}

			jsonBody, err := json.Marshal(testBook)
			if err != nil {
				results <- http.StatusInternalServerError
				return
			}

			req, err := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
			if err != nil {
				results <- http.StatusInternalServerError
				return
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			results <- rr.Code
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrentRequests; i++ {
		statusCode := <-results
		if statusCode == http.StatusOK {
			successCount++
		}
	}

	// All concurrent requests should succeed
	assert.Equal(t, concurrentRequests, successCount)
}

func TestAPIErrorHandling(t *testing.T) {
	suite := setupHTTPTestSuite(t)
	defer suite.tearDown(t)

	tests := []struct {
		name           string
		method         string
		url            string
		body           string
		expectedStatus int
	}{
		{
			name:           "Invalid JSON format",
			method:         "POST",
			url:            "/api/v1/books",
			body:           `{"title": "Test", "author":}`, // Invalid JSON
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Non-existent endpoint",
			method:         "GET",
			url:            "/api/v1/nonexistent",
			body:           "",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
			require.NoError(t, err)

			if tt.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
