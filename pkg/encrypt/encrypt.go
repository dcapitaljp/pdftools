package encrypt

import (
	"bytes"
	"io"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

func Encrypt(infile string) (err error) {
	// func Optimize(rs io.ReadSeeker, w io.Writer, conf *pdfcpu.Configuration) error
	api.Optimize()
}
