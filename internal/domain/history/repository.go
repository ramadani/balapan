package history

import "context"

type Repository interface {
	Create(ctx context.Context, data *History) error
}
