package link

import (
	"testing"
)

func TestNewLink(t *testing.T) {
	t.Run("valid link with custom short name", func(t *testing.T) {
		link, err := NewLink("https://example.com/long-url", "exmpl")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if link.OriginalURL != "https://example.com/long-url" {
			t.Errorf("expected original URL %q, got %q", "https://example.com/long-url", link.OriginalURL)
		}
		if link.ShortName != "exmpl" {
			t.Errorf("expected short name %q, got %q", "exmpl", link.ShortName)
		}
		if link.ID != 0 {
			t.Errorf("expected ID 0, got %d", link.ID)
		}
	})

	t.Run("valid link with generated short name", func(t *testing.T) {
		link, err := NewLink("https://example.com/long-url", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if link.ShortName == "" {
			t.Error("expected generated short name")
		}
		if len(link.ShortName) != 6 {
			t.Errorf("expected short name length 6, got %d", len(link.ShortName))
		}
	})

	t.Run("empty URL", func(t *testing.T) {
		_, err := NewLink("", "exmpl")
		if err != ErrEmptyURL {
			t.Errorf("expected error %v, got %v", ErrEmptyURL, err)
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := NewLink("not-a-url", "exmpl")
		if err != ErrInvalidURL {
			t.Errorf("expected error %v, got %v", ErrInvalidURL, err)
		}
	})
}

func TestLinkValidate(t *testing.T) {
	t.Run("valid link", func(t *testing.T) {
		link := &Link{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortName:   "exmpl",
		}
		if err := link.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty original URL", func(t *testing.T) {
		link := &Link{
			OriginalURL: "",
			ShortName:   "exmpl",
		}
		if err := link.Validate(); err != ErrEmptyURL {
			t.Errorf("expected error %v, got %v", ErrEmptyURL, err)
		}
	})

	t.Run("invalid original URL", func(t *testing.T) {
		link := &Link{
			OriginalURL: "not-a-url",
			ShortName:   "exmpl",
		}
		if err := link.Validate(); err != ErrInvalidURL {
			t.Errorf("expected error %v, got %v", ErrInvalidURL, err)
		}
	})

	t.Run("empty short name", func(t *testing.T) {
		link := &Link{
			OriginalURL: "https://example.com",
			ShortName:   "",
		}
		if err := link.Validate(); err == nil {
			t.Error("expected error for empty short name")
		}
	})
}

func TestGenerateShortName(t *testing.T) {
	t.Run("generates non-empty string", func(t *testing.T) {
		name := GenerateShortName()
		if name == "" {
			t.Error("expected non-empty short name")
		}
	})

	t.Run("generates unique names", func(t *testing.T) {
		names := make(map[string]bool)
		for i := 0; i < 100; i++ {
			name := GenerateShortName()
			if names[name] {
				t.Error("expected unique short names")
			}
			names[name] = true
		}
	})

	t.Run("length is 6", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			name := GenerateShortName()
			if len(name) != 6 {
				t.Errorf("expected length 6, got %d", len(name))
			}
		}
	})
}
