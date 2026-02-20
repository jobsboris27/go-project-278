package link

import (
	"context"
	"fmt"

	"app/internal/domain/link"
)

type Service struct {
	repo    link.Repository
	baseURL string
}

func NewService(repo link.Repository, baseURL string) *Service {
	return &Service{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (s *Service) CreateLink(ctx context.Context, originalURL, shortName string) (*link.Link, error) {
	linkEntity, err := link.NewLink(originalURL, shortName)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByShortName(ctx, linkEntity.ShortName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, link.ErrShortNameExists
	}

	if err := s.repo.Create(ctx, linkEntity); err != nil {
		return nil, err
	}

	return linkEntity, nil
}

func (s *Service) GetLink(ctx context.Context, id int64) (*link.Link, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetAllLinks(ctx context.Context) ([]*link.Link, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) UpdateLink(ctx context.Context, id int64, originalURL, shortName string) (*link.Link, error) {
	if originalURL == "" {
		existing, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		originalURL = existing.OriginalURL
	}

	if shortName == "" {
		existing, err := s.repo.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		shortName = existing.ShortName
	}

	linkEntity, err := link.NewLink(originalURL, shortName)
	if err != nil {
		return nil, err
	}
	linkEntity.ID = id

	if err := s.repo.Update(ctx, linkEntity); err != nil {
		return nil, err
	}

	return linkEntity, nil
}

func (s *Service) DeleteLink(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetShortURL(linkEntity *link.Link) string {
	return fmt.Sprintf("%s/r/%s", s.baseURL, linkEntity.ShortName)
}
