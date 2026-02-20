package sqlc

import (
	"context"
	"database/sql"
	"time"
)

type Link struct {
	ID          int64
	OriginalURL string
	ShortName   string
	CreatedAt   time.Time
}

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

func (q *Queries) DB() *sql.DB {
	return q.db
}

func (q *Queries) GetLinkByShortName(ctx context.Context, shortName string) (Link, error) {
	var link Link
	err := q.db.QueryRowContext(ctx,
		"SELECT id, original_url, short_name, created_at FROM links WHERE short_name = $1",
		shortName).Scan(&link.ID, &link.OriginalURL, &link.ShortName, &link.CreatedAt)
	return link, err
}

func (q *Queries) CreateLink(ctx context.Context, originalURL, shortName string) (Link, error) {
	var link Link
	err := q.db.QueryRowContext(ctx,
		"INSERT INTO links (original_url, short_name) VALUES ($1, $2) RETURNING id, original_url, short_name, created_at",
		originalURL, shortName).Scan(&link.ID, &link.OriginalURL, &link.ShortName, &link.CreatedAt)
	return link, err
}

func (q *Queries) GetLinkByID(ctx context.Context, id int64) (Link, error) {
	var link Link
	err := q.db.QueryRowContext(ctx,
		"SELECT id, original_url, short_name, created_at FROM links WHERE id = $1",
		id).Scan(&link.ID, &link.OriginalURL, &link.ShortName, &link.CreatedAt)
	return link, err
}

func (q *Queries) GetAllLinks(ctx context.Context, offset, limit int) ([]Link, error) {
	rows, err := q.db.QueryContext(ctx,
		"SELECT id, original_url, short_name, created_at FROM links ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var links []Link
	for rows.Next() {
		var link Link
		if err := rows.Scan(&link.ID, &link.OriginalURL, &link.ShortName, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, rows.Err()
}

func (q *Queries) UpdateLink(ctx context.Context, originalURL, shortName string, id int64) (Link, error) {
	var link Link
	err := q.db.QueryRowContext(ctx,
		"UPDATE links SET original_url = $1, short_name = $2 WHERE id = $3 RETURNING id, original_url, short_name, created_at",
		originalURL, shortName, id).Scan(&link.ID, &link.OriginalURL, &link.ShortName, &link.CreatedAt)
	return link, err
}

func (q *Queries) DeleteLink(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM links WHERE id = $1", id)
	return err
}

func (q *Queries) ExistsByShortName(ctx context.Context, shortName string) (bool, error) {
	var exists bool
	err := q.db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM links WHERE short_name = $1)",
		shortName).Scan(&exists)
	return exists, err
}

type LinkVisit struct {
	ID        int64
	LinkID    int64
	IP        string
	UserAgent string
	Referer   string
	Status    int
	CreatedAt time.Time
}

func (q *Queries) CreateLinkVisit(ctx context.Context, linkID int64, ip, userAgent, referer string, status int) (LinkVisit, error) {
	var visit LinkVisit
	err := q.db.QueryRowContext(ctx,
		"INSERT INTO link_visits (link_id, ip, user_agent, referer, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, link_id, ip, user_agent, referer, status, created_at",
		linkID, ip, userAgent, referer, status).Scan(&visit.ID, &visit.LinkID, &visit.IP, &visit.UserAgent, &visit.Referer, &visit.Status, &visit.CreatedAt)
	return visit, err
}

func (q *Queries) GetLinkVisits(ctx context.Context, limit, offset int) ([]LinkVisit, error) {
	rows, err := q.db.QueryContext(ctx,
		"SELECT id, link_id, ip, user_agent, referer, status, created_at FROM link_visits ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var visits []LinkVisit
	for rows.Next() {
		var visit LinkVisit
		if err := rows.Scan(&visit.ID, &visit.LinkID, &visit.IP, &visit.UserAgent, &visit.Referer, &visit.Status, &visit.CreatedAt); err != nil {
			return nil, err
		}
		visits = append(visits, visit)
	}
	return visits, rows.Err()
}

func (q *Queries) CountLinkVisits(ctx context.Context) (int, error) {
	var total int
	err := q.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM link_visits").Scan(&total)
	return total, err
}

func (q *Queries) DeleteLinkVisit(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, "DELETE FROM link_visits WHERE id = $1", id)
	return err
}
