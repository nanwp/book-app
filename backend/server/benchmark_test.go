package server

import (
	"byfood-interview/book"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkCreateBook(b *testing.B) {
	suite := setupHTTPTestSuite(&testing.T{})
	defer suite.tearDown(&testing.T{})

	testBook := book.Book{
		Title:         "Benchmark Book",
		Author:        "Benchmark Author",
		PublishedYear: 2023,
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			jsonBody, _ := json.Marshal(testBook)
			req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
			}
		}
	})
}

func BenchmarkGetAllBooks(b *testing.B) {
	suite := setupHTTPTestSuite(&testing.T{})
	defer suite.tearDown(&testing.T{})

	// Pre-populate with some books
	for i := 0; i < 10; i++ {
		testBook := book.Book{
			Title:         fmt.Sprintf("Book %d", i),
			Author:        fmt.Sprintf("Author %d", i),
			PublishedYear: 2000 + (i % 24),
		}

		jsonBody, _ := json.Marshal(testBook)
		req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		suite.server.Router.ServeHTTP(rr, req)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/v1/books", nil)
			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
			}
		}
	})
}

func BenchmarkGetBookByID(b *testing.B) {
	suite := setupHTTPTestSuite(&testing.T{})
	defer suite.tearDown(&testing.T{})

	// Create a test book
	testBook := book.Book{
		Title:         "Benchmark Book",
		Author:        "Benchmark Author",
		PublishedYear: 2023,
	}

	jsonBody, _ := json.Marshal(testBook)
	req, _ := http.NewRequest("POST", "/api/v1/books", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	suite.server.Router.ServeHTTP(rr, req)

	var createdBook book.Book
	json.Unmarshal(rr.Body.Bytes(), &createdBook)
	bookID := fmt.Sprintf("%d", createdBook.ID)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/v1/books/"+bookID, nil)
			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
			}
		}
	})
}

func BenchmarkHealthCheck(b *testing.B) {
	suite := setupHTTPTestSuite(&testing.T{})
	defer suite.tearDown(&testing.T{})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, _ := http.NewRequest("GET", "/api/v1/healthcheck", nil)
			rr := httptest.NewRecorder()
			suite.server.Router.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
			}
		}
	})
}
