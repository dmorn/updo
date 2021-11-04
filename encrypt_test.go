package updo_test

import (
	"testing"
	"bytes"
	"strings"
	"os"

	"github.com/jecoz/updo"
)

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(data)
}

var (
	pubKey = readFile("testdata/id_ed25519.pub")
	privKey = readFile("testdata/id_ed25519")
)

func TestEncrypt(t *testing.T) {
	input := "updo"

	inout := &bytes.Buffer{}
	if err := updo.Encrypt(inout, strings.NewReader(input), pubKey); err != nil {
		t.Fatal(err)
	}
	out := &bytes.Buffer{}
	if err := updo.Decrypt(out, inout, privKey); err != nil {
		t.Fatal(err)
	}

	if out.String() != input {
		t.Fatalf("want: %q, have: %q", input, out.String())
	}
}
