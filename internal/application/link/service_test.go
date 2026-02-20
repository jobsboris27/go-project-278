package link

import (
	"context"
	"errors"
	"testing"

	"app/internal/domain/link"
)

type mockRepository struct {
	links           map[int64]*link.Link
	shortNameExists map[string]bool
	nextID          int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		links:           make(map[int64]*link.Link),
		shortNameExists: make(map[string]bool),
		nextID:          1,
	}
}

func (m *mockRepository) Create(ctx context.Context, link *link.Link) error {
	link.ID = m.nextID
	m.nextID++
	m.links[link.ID] = link
	m.shortNameExists[link.ShortName] = true
	return nil
}

func (m *mockRepository) GetByID(ctx context.Context, id int64) (*link.Link, error) {
	link, ok := m.links[id]
	if !ok {
		return nil, errors.New("link not found")
	}
	return link, nil
}

func (m *mockRepository) GetAll(ctx context.Context) ([]*link.Link, error) {
	links := make([]*link.Link, 0, len(m.links))
	for _, link := range m.links {
		links = append(links, link)
	}
	return links, nil
}

func (m *mockRepository) Update(ctx context.Context, link *link.Link) error {
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

func TestServiceCreateLink(t *testing.T) {
	t.Run("create link with custom short name", func(t *testing.T) {
		repo := newMockRepository()
		service := NewService(repo, "https://short.io")

		linkEntity, err := service.CreateLink(context.Background(), "https://example.com", "exmpl")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if linkEntity.OriginalURL != "https://example.com" {
			t.Errorf("expected original URL %q, got %q", "https://example.com", linkEntity.OriginalURL)
		}
		if linkEntity.ShortName != "exmpl" {
			t.Errorf("expected short name %q, got %q", "exmpl", linkEntity.ShortName)
		}
		if linkEntity.ID != 1 {
			t.Errorf("expected ID 1, got %d", linkEntity.ID)
		}
	})

	t.Run("create link with generated short name", func(t *testing.T) {
		repo := newMockRepository()
		service := NewService(repo, "https://short.io")

		linkEntity, err := service.CreateLink(context.Background(), "https://example.com", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if linkEntity.ShortName == "" {
			t.Error("expected generated short name")
		}
	})

	t.Run("create link with existing short name", func(t *testing.T) {
		repo := newMockRepository()
		repo.shortNameExists["existing"] = true
		service := NewService(repo, "https://short.io")

		_, err := service.CreateLink(context.Background(), "https://example.com", "existing")
		if err != link.ErrShortNameExists {
			t.Errorf("expected error %v, got %v", link.ErrShortNameExists, err)
		}
	})

	t.Run("create link with invalid URL", func(t *testing.T) {
		repo := newMockRepository()
		service := NewService(repo, "https://short.io")

		_, err := service.CreateLink(context.Background(), "not-a-url", "exmpl")
		if err != link.ErrInvalidURL {
			t.Errorf("expected error %v, got %v", link.ErrInvalidURL, err)
		}
	})
}

func TestServiceGetLink(t *testing.T) {
	t.Run("get existing link", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &link.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		service := NewService(repo, "https://short.io")

		linkEntity, err := service.GetLink(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if linkEntity.ID != 1 {
			t.Errorf("expected ID 1, got %d", linkEntity.ID)
		}
	})

	t.Run("get non-existing link", func(t *testing.T) {
		repo := newMockRepository()
		service := NewService(repo, "https://short.io")

		_, err := service.GetLink(context.Background(), 999)
		if err == nil {
			t.Error("expected error for non-existing link")
		}
	})
}

func TestServiceGetAllLinks(t *testing.T) {
	t.Run("get all links", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &link.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		repo.links[2] = &link.Link{ID: 2, OriginalURL: "https://example.org", ShortName: "exmp2"}
		service := NewService(repo, "https://short.io")

		links, err := service.GetAllLinks(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(links) != 2 {
			t.Errorf("expected 2 links, got %d", len(links))
		}
	})
}

func TestServiceUpdateLink(t *testing.T) {
	t.Run("update existing link", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &link.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		service := NewService(repo, "https://short.io")

		linkEntity, err := service.UpdateLink(context.Background(), 1, "https://new-example.com", "newex")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if linkEntity.OriginalURL != "https://new-example.com" {
			t.Errorf("expected original URL %q, got %q", "https://new-example.com", linkEntity.OriginalURL)
		}
	})

	t.Run("update non-existing link", func(t *testing.T) {
		repo := newMockRepository()
		service := NewService(repo, "https://short.io")

		_, err := service.UpdateLink(context.Background(), 999, "https://example.com", "exmpl")
		if err == nil {
			t.Error("expected error for non-existing link")
		}
	})
}

func TestServiceDeleteLink(t *testing.T) {
	t.Run("delete existing link", func(t *testing.T) {
		repo := newMockRepository()
		repo.links[1] = &link.Link{ID: 1, OriginalURL: "https://example.com", ShortName: "exmpl"}
		service := NewService(repo, "https://short.io")

		err := service.DeleteLink(context.Background(), 1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := repo.links[1]; ok {
			t.Error("expected link to be deleted")
		}
	})

	t.Run("delete non-existing link", func(t *testing.T) {
		repo := newMockRepository()
		service := NewService(repo, "https://short.io")

		err := service.DeleteLink(context.Background(), 999)
		if err == nil {
			t.Error("expected error for non-existing link")
		}
	})
}

func TestServiceGetShortURL(t *testing.T) {
	repo := newMockRepository()
	service := NewService(repo, "https://short.io")

	linkEntity := &link.Link{ShortName: "exmpl"}
	url := service.GetShortURL(linkEntity)

	expected := "https://short.io/r/exmpl"
	if url != expected {
		t.Errorf("expected %q, got %q", expected, url)
	}
}
