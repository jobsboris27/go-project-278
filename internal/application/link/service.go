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

func (s *Service) GetLinkByShortName(ctx context.Context, shortName string) (*link.Link, error) {
	return s.repo.GetByShortName(ctx, shortName)
}

func (s *Service) GetAllLinks(ctx context.Context, offset, limit int) ([]*link.Link, int, error) {
	return s.repo.GetAll(ctx, offset, limit)
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

func (s *Service) RecordVisit(ctx context.Context, linkID int64, ip, userAgent, referer string, status int) error {
	visit := link.NewLinkVisit(linkID, ip, userAgent, referer, status)
	return s.repo.CreateVisit(ctx, visit)
}

func (s *Service) GetVisits(ctx context.Context, offset, limit int) ([]*link.LinkVisit, int, error) {
	return s.repo.GetVisits(ctx, offset, limit)
}

func (s *Service) DeleteVisit(ctx context.Context, id int64) error {
	return s.repo.DeleteVisit(ctx, id)
}
