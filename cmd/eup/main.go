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
	publicKey = flag.String("r", "", "Encrypt to the specified recipient")
	publicKeyFile = flag.String("R", "", "Encrypt to the specified recipient file")
)

func readPublicKey() (string, error) {
	if *publicKey != "" {
		return *publicKey, nil
	}
	if *publicKeyFile == "" {
		return "", fmt.Errorf("at least -r or -R must be specified")
	}
	data, err := os.ReadFile(*publicKeyFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Main(ctx context.Context, paths []string) error {
	key, err := readPublicKey()
	if err != nil {
		return err
	}

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

	upl := updo.NewUploader(bucket, key)
	return upl.Upload(ctx, files...)
}

func main() {
	flag.Parse()
	if err := Main(context.Background(), flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}
