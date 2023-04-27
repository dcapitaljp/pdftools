package crypto

import (
	"bytes"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func EncryptoInplace(infile string, conf *model.Configuration) error {
	var infd, tmpfd *os.File
	infd, err := os.OpenFile(infile, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = func() error {
		defer func() { infd.Close() }()
		if err := api.Optimize(infd, buf, conf); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return err
	}
	result := buf.Bytes()
	tmpFile := infile + ".tmp"
	tmpfd, err = os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	var writeErr error
	defer func() {
		tmpfd.Close()
		if writeErr != nil {
			os.Remove(tmpFile)
			return
		}
		os.Rename(tmpFile, infile)
	}()
	if _, writeErr = tmpfd.Write(result); writeErr != nil {
		return err
	}
	return nil
}
