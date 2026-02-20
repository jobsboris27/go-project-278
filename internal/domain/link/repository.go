package link

import "context"

type Repository interface {
	Create(ctx context.Context, link *Link) error
	GetByID(ctx context.Context, id int64) (*Link, error)
	GetAll(ctx context.Context) ([]*Link, error)
	Update(ctx context.Context, link *Link) error
	Delete(ctx context.Context, id int64) error
	ExistsByShortName(ctx context.Context, shortName string) (bool, error)
}
