package updo

import (
	"log"
	"io"
	"context"
	"time"
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/sync/errgroup"
)

func CleanKey(path string) string {
	return strings.TrimLeftFunc(path, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
}

type Uploader struct {
	store Store
	publicKey string
}

func NewUploader(store Store, pubKey string) Uploader {
	return Uploader{
		store: store,
		publicKey: pubKey,
	}
}

type uploadResult struct {
	Input NamedData
	Err error
	Elapsed time.Duration
}

type NamedData struct {
	Data io.ReadCloser
	Key string
}

func (u Uploader) uploadOne(ctx context.Context, d NamedData, c chan<- uploadResult) {
	tic := time.Now()
	g, ctx := errgroup.WithContext(ctx)

	// Buffer?
	pr, pw := io.Pipe()

	g.Go(func() error {
		defer pw.Close()
		return Encrypt(pw, d.Data, u.publicKey)
	})
	g.Go(func() error {
		err := u.store.Put(ctx, pr, d.Key)
		pw.CloseWithError(err) // Encrypt does not respect ctx!
		return err
	})

	err := g.Wait()

	c <- uploadResult{
		Err: err,
		Input: d,
		Elapsed: time.Since(tic),
	}
}

func (u Uploader) Upload(ctx context.Context, data ...NamedData) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	c := make(chan uploadResult, len(data))
	for _, v := range data {
		log.Printf("updo: %q: uploading", v.Key)

		// Concurrency is not upper bounded!
		go u.uploadOne(ctx, v, c)
	}
	for i := 0; i < len(data); i++ {
		res := <-c
		if err := res.Err; err != nil {
			return fmt.Errorf("upload %q: %w", res.Input.Key, err)
		}
		log.Printf("updo: %q: uploaded in %v", res.Input.Key, res.Elapsed)
	}
	return nil
}
