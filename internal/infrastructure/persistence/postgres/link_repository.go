package postgres

import (
	"context"
	"database/sql"

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
	_, err := r.queries.CreateLink(ctx, linkEntity.OriginalURL, linkEntity.ShortName)
	return err
}

func (r *LinkRepository) GetByID(ctx context.Context, id int64) (*link.Link, error) {
	dbLink, err := r.queries.GetLinkByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDomainLink(dbLink), nil
}

func (r *LinkRepository) GetAll(ctx context.Context) ([]*link.Link, error) {
	dbLinks, err := r.queries.GetAllLinks(ctx)
	if err != nil {
		return nil, err
	}
	links := make([]*link.Link, len(dbLinks))
	for i, dbLink := range dbLinks {
		links[i] = toDomainLink(dbLink)
	}
	return links, nil
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

func toDomainLink(dbLink sqlc.Link) *link.Link {
	return &link.Link{
		ID:          dbLink.ID,
		OriginalURL: dbLink.OriginalURL,
		ShortName:   dbLink.ShortName,
		CreatedAt:   dbLink.CreatedAt,
	}
}
