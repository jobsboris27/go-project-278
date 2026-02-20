package postgres

import (
	"context"
	"database/sql"
	"errors"

	"app/db/sqlc"
	"app/internal/domain/link"
)

type LinkRepository struct {
	queries *sqlc.Queries
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{
		queries: sqlc.New(db),
	}
}

func (r *LinkRepository) Create(ctx context.Context, linkEntity *link.Link) error {
	dbLink, err := r.queries.CreateLink(ctx, linkEntity.OriginalURL, linkEntity.ShortName)
	if err != nil {
		return err
	}
	linkEntity.ID = dbLink.ID
	linkEntity.CreatedAt = dbLink.CreatedAt
	return nil
}

func (r *LinkRepository) GetByID(ctx context.Context, id int64) (*link.Link, error) {
	dbLink, err := r.queries.GetLinkByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("link not found")
		}
		return nil, err
	}
	return toDomainLink(dbLink), nil
}

func (r *LinkRepository) GetByShortName(ctx context.Context, shortName string) (*link.Link, error) {
	dbLink, err := r.queries.GetLinkByShortName(ctx, shortName)
	if err != nil {
		return nil, err
	}
	return toDomainLink(dbLink), nil
}

func (r *LinkRepository) GetAll(ctx context.Context, offset, limit int) ([]*link.Link, int, error) {
	total, err := r.count(ctx)
	if err != nil {
		return nil, 0, err
	}

	dbLinks, err := r.queries.GetAllLinks(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	links := make([]*link.Link, len(dbLinks))
	for i, dbLink := range dbLinks {
		links[i] = toDomainLink(dbLink)
	}
	return links, total, nil
}

func (r *LinkRepository) count(ctx context.Context) (int, error) {
	var total int
	err := r.queries.DB().QueryRowContext(ctx, "SELECT COUNT(*) FROM links").Scan(&total)
	return total, err
}

func (r *LinkRepository) Update(ctx context.Context, linkEntity *link.Link) error {
	_, err := r.queries.UpdateLink(ctx, linkEntity.OriginalURL, linkEntity.ShortName, linkEntity.ID)
	return err
}

func (r *LinkRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteLink(ctx, id)
}

func (r *LinkRepository) ExistsByShortName(ctx context.Context, shortName string) (bool, error) {
	return r.queries.ExistsByShortName(ctx, shortName)
}

func (r *LinkRepository) CreateVisit(ctx context.Context, visit *link.LinkVisit) error {
	_, err := r.queries.CreateLinkVisit(ctx, visit.LinkID, visit.IP, visit.UserAgent, visit.Referer, visit.Status)
	return err
}

func (r *LinkRepository) GetVisits(ctx context.Context, offset, limit int) ([]*link.LinkVisit, int, error) {
	total, err := r.queries.CountLinkVisits(ctx)
	if err != nil {
		return nil, 0, err
	}

	dbVisits, err := r.queries.GetLinkVisits(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	visits := make([]*link.LinkVisit, len(dbVisits))
	for i, dbVisit := range dbVisits {
		visits[i] = toDomainVisit(dbVisit)
	}
	return visits, total, nil
}

func (r *LinkRepository) DeleteVisit(ctx context.Context, id int64) error {
	return r.queries.DeleteLinkVisit(ctx, id)
}

func toDomainVisit(dbVisit sqlc.LinkVisit) *link.LinkVisit {
	return &link.LinkVisit{
		ID:        dbVisit.ID,
		LinkID:    dbVisit.LinkID,
		IP:        dbVisit.IP,
		UserAgent: dbVisit.UserAgent,
		Referer:   dbVisit.Referer,
		Status:    dbVisit.Status,
		CreatedAt: dbVisit.CreatedAt,
	}
}

func toDomainLink(dbLink sqlc.Link) *link.Link {
	return &link.Link{
		ID:          dbLink.ID,
		OriginalURL: dbLink.OriginalURL,
		ShortName:   dbLink.ShortName,
		CreatedAt:   dbLink.CreatedAt,
	}
}
