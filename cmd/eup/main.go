package main

import (
	"flag"
	"os"
	"context"
	"fmt"

	"github.com/jecoz/updo"
	"github.com/jecoz/updo/aws"
)

var (
	publicKey = flag.String("r", "", "Encrypt to the specified recipient. Required")
)

func Main(ctx context.Context, paths []string) error {
	bucket, err := aws.NewBucket()
	if err != nil {
		return err
	}

	files := make([]updo.NamedData, len(paths))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		files[i] = updo.NamedData{
			Data: file,
			Key: updo.CleanKey(file.Name()),
		}
	}

	upl := updo.NewUploader(bucket, *publicKey)
	return upl.Upload(ctx, files...)
}

func main() {
	flag.Parse()
	if err := Main(context.Background(), flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}
