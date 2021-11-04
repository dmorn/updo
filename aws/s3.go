package aws

import (
	"context"
	"io"
	"fmt"
)

type Bucket struct {
}

func (b Bucket) Put(ctx context.Context, r io.Reader, key string) error {
	return fmt.Errorf("not implemented!")
}

func NewBucket() (Bucket, error) {
	return Bucket{}, nil
}
