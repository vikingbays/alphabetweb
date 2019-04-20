// Copyright 2019 The VikingBays(in Nanjing , China) . All rights reserved.
// Released under the Apache license : http://www.apache.org/licenses/LICENSE-2.0 .
//
// authors:   VikingBays
// email  :   vikingbays@gmail.com

package service

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"sort"
	"strings"
)

// A Writer generates multipart messages.
type RpcUploadWriter struct {
	boundary string
	lastpart *rpcuploadpart
}

// NewWriter returns a new multipart Writer with a random boundary,
// writing to w.
func NewRpcUploadWriter() *RpcUploadWriter {
	return &RpcUploadWriter{
		boundary: randomBoundary(),
	}
}

// Boundary returns the Writer's boundary.
func (w *RpcUploadWriter) Boundary() string {
	return w.boundary
}

// SetBoundary overrides the Writer's default randomly-generated
// boundary separator with an explicit value.
//
// SetBoundary must be called before any parts are created, may only
// contain certain ASCII characters, and must be non-empty and
// at most 69 bytes long.
func (w *RpcUploadWriter) SetBoundary(boundary string) error {
	if w.lastpart != nil {
		return errors.New("mime: SetBoundary called after write")
	}
	// rfc2046#section-5.1.1
	if len(boundary) < 1 || len(boundary) > 69 {
		return errors.New("mime: invalid boundary length")
	}
	for _, b := range boundary {
		if 'A' <= b && b <= 'Z' || 'a' <= b && b <= 'z' || '0' <= b && b <= '9' {
			continue
		}
		switch b {
		case '\'', '(', ')', '+', '_', ',', '-', '.', '/', ':', '=', '?':
			continue
		}
		return errors.New("mime: invalid boundary character")
	}
	w.boundary = boundary
	return nil
}

// FormDataContentType returns the Content-Type for an HTTP
// multipart/form-data with this Writer's Boundary.
func (w *RpcUploadWriter) FormDataContentType() string {
	return "multipart/form-data; boundary=" + w.boundary
}

func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

// CreatePart creates a new multipart section with the provided
// header. The body of the part should be written to the returned
// Writer. After calling CreatePart, any previous part may no longer
// be written to.
func (w *RpcUploadWriter) CreatePart(header textproto.MIMEHeader) ([]byte, error) {
	if w.lastpart != nil {
		if err := w.lastpart.close(); err != nil {
			return nil, err
		}
	}
	var b bytes.Buffer
	if w.lastpart != nil {
		fmt.Fprintf(&b, "\r\n--%s\r\n", w.boundary)
	} else {
		fmt.Fprintf(&b, "--%s\r\n", w.boundary)
	}

	keys := make([]string, 0, len(header))
	for k := range header {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		for _, v := range header[k] {
			fmt.Fprintf(&b, "%s: %s\r\n", k, v)
		}
	}
	fmt.Fprintf(&b, "\r\n")
	//_, err := io.Copy(w.w, &b)
	d1 := b.Bytes()

	p := &rpcuploadpart{
		mw: w,
	}
	w.lastpart = p
	return d1, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// CreateFormFile is a convenience wrapper around CreatePart. It creates
// a new form-data header with the provided field name and file name.
func (w *RpcUploadWriter) CreateFormFileStart(fieldname, filename string) ([]byte, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", "application/octet-stream")
	return w.CreatePart(h)
}

// CreateFormField calls CreatePart with a header using the
// given field name.
func (w *RpcUploadWriter) CreateFormField(fieldname string) ([]byte, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"`, escapeQuotes(fieldname)))
	return w.CreatePart(h)
}

// WriteField calls CreateFormField and then writes the given value.
func (w *RpcUploadWriter) WriteField(fieldname, value string) ([]byte, []byte, error) {
	d1, err := w.CreateFormField(fieldname)
	if err != nil {
		return nil, nil, err
	}
	return d1, []byte(value), nil
}

// Close finishes the multipart message and writes the trailing
// boundary end line to the output.
func (w *RpcUploadWriter) End() ([]byte, error) {
	if w.lastpart != nil {
		if err := w.lastpart.close(); err != nil {
			return nil, err
		}
		w.lastpart = nil
	}
	s1 := fmt.Sprintf("\r\n--%s--\r\n", w.boundary)
	//_, err := fmt.Fprintf(w.w, "\r\n--%s--\r\n", w.boundary)
	return []byte(s1), nil
}

type rpcuploadpart struct {
	mw     *RpcUploadWriter
	closed bool
	we     error // last error that occurred writing
}

func (p *rpcuploadpart) close() error {
	p.closed = true
	return p.we
}

func (p *rpcuploadpart) Write(d []byte) (d1 []byte, err error) {
	if p.closed {
		return nil, errors.New("multipart: can't write to finished part")
	}
	return d1, nil
}
