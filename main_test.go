package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app/internal/application/link"
	linkhttp "app/internal/infrastructure/http"
	domainLink "app/internal/domain/link"

	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != "pong" {
		t.Errorf("expected body %q, got %q", "pong", w.Body.String())
	}
}

func TestAPIEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	repo := &mockRepository{
		links:           make(map[int64]*domainLink.Link),
		shortNameExists: make(map[string]bool),
		nextID:          1,
	}
	service := link.NewService(repo, "https://short.io")
	handler := linkhttp.NewHandler(service)
	handler.RegisterRoutes(router)

	t.Run("GET /api/links returns empty list", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/links", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("POST /api/links creates link", func(t *testing.T) {
		body := `{"original_url": "https://example.com", "short_name": "test"}`
		req := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})
}

type mockRepository struct {
	links           map[int64]*domainLink.Link
	shortNameExists map[string]bool
	nextID          int64
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

func (m *mockRepository) GetByShortName(ctx context.Context, shortName string) (*domainLink.Link, error) {
	for _, link := range m.links {
		if link.ShortName == shortName {
			return link, nil
		}
	}
	return nil, errors.New("link not found")
}

func (m *mockRepository) GetAll(ctx context.Context, offset, limit int) ([]*domainLink.Link, int, error) {
	all := make([]*domainLink.Link, 0, len(m.links))
	for _, l := range m.links {
		all = append(all, l)
	}

	total := len(all)

	end := offset + limit
	if end > total {
		end = total
	}

	if offset >= total {
		return []*domainLink.Link{}, total, nil
	}

	return all[offset:end], total, nil
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

func (m *mockRepository) CreateVisit(ctx context.Context, visit *domainLink.LinkVisit) error {
	return nil
}

func (m *mockRepository) GetVisits(ctx context.Context, offset, limit int) ([]*domainLink.LinkVisit, int, error) {
	return []*domainLink.LinkVisit{}, 0, nil
}
