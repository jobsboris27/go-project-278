package link

import (
	"errors"
	"math/rand"
	"net/url"
	"time"
)

var (
	ErrInvalidURL      = errors.New("invalid URL")
	ErrEmptyURL        = errors.New("URL cannot be empty")
	ErrShortNameExists = errors.New("short name already exists")
)

type Link struct {
	ID          int64
	OriginalURL string
	ShortName   string
	CreatedAt   time.Time
}

func NewLink(originalURL string, shortName string) (*Link, error) {
	if originalURL == "" {
		return nil, ErrEmptyURL
	}

	if _, err := url.ParseRequestURI(originalURL); err != nil {
		return nil, ErrInvalidURL
	}

	if shortName == "" {
		shortName = GenerateShortName()
	}

	return &Link{
		OriginalURL: originalURL,
		ShortName:   shortName,
		CreatedAt:   time.Now(),
	}, nil
}

func (l *Link) Validate() error {
	if l.OriginalURL == "" {
		return ErrEmptyURL
	}

	if _, err := url.ParseRequestURI(l.OriginalURL); err != nil {
		return ErrInvalidURL
	}

	if l.ShortName == "" {
		return errors.New("short name cannot be empty")
	}

	return nil
}

func GenerateShortName() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
