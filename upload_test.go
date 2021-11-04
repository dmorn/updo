package updo_test

import (
	"context"
	"io"
	"testing"
	"path/filepath"
	"time"
	"math/rand"
	"os"
	"fmt"
	"strings"

	"github.com/jecoz/updo"
)

// NOTE: we might mock the storage in memory.
type Disk struct {
	Root string
}

func (d Disk) Put(ctx context.Context, r io.Reader, key string) error {
	file, err := os.Create(filepath.Join(d.Root, key))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func (d Disk) IsPresent(key string) bool {
	stat, err := os.Stat(filepath.Join(d.Root, key))
	if err != nil {
		return false
	}
	return stat.Size() > 0
}

func TestUpload(t *testing.T) {
	t.Parallel()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	workPath := filepath.Join(os.TempDir(), fmt.Sprintf("updo-%04d", r.Intn(1000)))

	// Test in clean workspace
	os.RemoveAll(workPath)
	os.MkdirAll(workPath, os.ModePerm)

	d := Disk{
		Root: workPath,
	}
	upl := updo.NewUploader(d, pubKey)
	data := []updo.NamedData{
		updo.NamedData{
			Data: io.NopCloser(strings.NewReader("hello")),
			Key: "abc",
		},
		updo.NamedData{
			Data: io.NopCloser(strings.NewReader("world")),
			Key: "cba",
		},
	}

	t.Run("upload", func(t *testing.T) {
		if err := upl.Upload(context.Background(), data...); err != nil {
			t.Fatal(err)
		}

		// NOTE: we're not checking wether the disk contains extra
		// (unwanted) stuff.

		// NOTE 2: each time the content is encrypted it changes, so we
		// cannot check wether disk's contents are actually the ones we
		// expect if we do not decrypt it.

		for _, v := range data {
			if !d.IsPresent(v.Key) {
				t.Fatalf("key %q was not uploaded", v.Key)
			}
		}
	})

	os.RemoveAll(workPath)
}

func TestCleanKey(t *testing.T) {
	t.Parallel()
	tt := []struct{
		Input string
		Want string
	}{
		{"hello.go", "hello.go"},
		{"", ""},
		{"././././file.ex", "file.ex"},
	}

	for _, v := range tt {
		have := updo.CleanKey(v.Input)
		if have != v.Want {
			t.Fatalf("have: %q, want: %q", have, v.Want)
		}
	}
}
