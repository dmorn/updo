package updo

import (
	"io"
	"fmt"

	"filippo.io/age"
	"filippo.io/age/armor"
	"filippo.io/age/agessh"
)

func Encrypt(dst io.Writer, src io.Reader, pubKey string) error {
	recipient, err := agessh.ParseRecipient(pubKey)
	if err != nil {
		return err
	}

	armorWriter := armor.NewWriter(dst)
	w, err := age.Encrypt(armorWriter, recipient)
	if err != nil {
		return err
	}

	if _, err := io.Copy(w, src); err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("close encrypted file: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return fmt.Errorf("close armor file: %w", err)
	}
	return nil
}

func Decrypt(dst io.Writer, src io.Reader, privKey string) error {
	identity, err := agessh.ParseIdentity([]byte(privKey))
	if err != nil {
		return err
	}
	armorReader := armor.NewReader(src)

	r, err := age.Decrypt(armorReader, identity)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, r); err != nil {
		return fmt.Errorf("descrypt: %w", err)
	}
	return nil
}
