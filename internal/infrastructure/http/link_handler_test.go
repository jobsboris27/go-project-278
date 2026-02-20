package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/internal/application/link"
	domainLink "app/internal/domain/link"

	"github.com/gin-gonic/gin"
)

type mockRepository struct {
	links           map[int64]*domainLink.Link
	shortNameExists map[string]bool
	nextID          int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		links:           make(map[int64]*domainLink.Link),
		shortNameExists: make(map[string]bool),
		nextID:          1,
	}
}

func (m *mockRepository) Create(ctx context.Context, link *domainLink.Link) error {
	link.ID = m.nextID
	m.nextID++
	m.links[link.ID] = link
	m.shortNameExists[link.ShortName] = true
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*domainLink.Link, error) {
	link, ok := m.links[id]
	if !ok {
		return nil, errors.New("link not found")
	}
	return link, nil
}

func (m *mockRepository) GetAll(ctx context.Context) ([]*domainLink.Link, error) {
	links := make([]*domainLink.Link, 0, len(m.links))
	for _, link := range m.links {
		links = append(links, link)
	}
	return links, nil
}

func (m *mockRepository) Update(ctx context.Context, link *domainLink.Link) error {
	if _, ok := m.links[link.ID]; !ok {
		return errors.New("link not found")
	}
	m.links[link.ID] = link
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id int64) error {
	if _, ok := m.links[id]; !ok {
		return errors.New("link not found")
	}
	delete(m.links, id)
	return nil
}

func (m *mockRepository) ExistsByShortName(ctx context.Context, shortName string) (bool, error) {
	return m.shortNameExists[shortName], nil
}

func setupRouter(service *link.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(service)
	handler.RegisterRoutes(router)
	return router
}

func TestGetAllLinks(t *testing.T) {
	repo := newMockRepository()
	repo.links[1] = &domainLink.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
	service := link.NewService(repo, "https://short.io")
	router := setupRouter(service)

	req := httptest.NewRequest(http.MethodGet, "/api/links", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []LinkResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("expected 1 link, got %d", len(response))
	}
}

func TestCreateLink(t *testing.T) {
	t.Run("create link with custom short name", func(t *testing.T) {
		repo := newMockRepository()
		service := link.NewService(repo, "https://short.io")
		router := setupRouter(service)

		body := CreateLinkRequest{
			OriginalURL: "https://example.com",
			ShortName:   "exmpl",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response LinkResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.OriginalURL != "https://example.com" {
			t.Errorf("expected original URL %q, got %q", "https://example.com", response.OriginalURL)
		}
	})

	t.Run("create link with invalid JSON", func(t *testing.T) {
		repo := newMockRepository()
		service := link.NewService(repo, "https://short.io")
		router := setupRouter(service)

		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewReader([]byte("invalid")))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

func TestGetLinkByID(t *testing.T) {
	t.Run("get existing link", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &domainLink.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		service := link.NewService(repo, "https://short.io")
		router := setupRouter(service)

		req := httptest.NewRequest(http.MethodGet, "/api/links/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response LinkResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.ID != 1 {
			t.Errorf("expected ID 1, got %d", response.ID)
		}
	})

	t.Run("get non-existing link", func(t *testing.T) {
		repo := newMockRepository()
		service := link.NewService(repo, "https://short.io")
		router := setupRouter(service)

		req := httptest.NewRequest(http.MethodGet, "/api/links/999", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

func TestUpdateLink(t *testing.T) {
	t.Run("update existing link", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &domainLink.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		service := link.NewService(repo, "https://short.io")
		router := setupRouter(service)

		body := CreateLinkRequest{
			OriginalURL: "https://new-example.com",
			ShortName:   "newex",
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/links/1", bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response LinkResponse
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if response.OriginalURL != "https://new-example.com" {
			t.Errorf("expected original URL %q, got %q", "https://new-example.com", response.OriginalURL)
		}
	})
}

func TestDeleteLink(t *testing.T) {
	t.Run("delete existing link", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &domainLink.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		service := link.NewService(repo, "https://short.io")
		router := setupRouter(service)

		req := httptest.NewRequest(http.MethodDelete, "/api/links/1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}
