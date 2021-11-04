package updo

import (
	"context"
	"io"
)

type Store interface {
	Put(ctx context.Context, r io.Reader, key string) error
}
