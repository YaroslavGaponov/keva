package client

import (
	"bytes"
	"compress/lzw"
	"io/ioutil"
)

type LzwCompressor struct {
	order    lzw.Order
	litWidth int
}

func NewLzwCompressor() LzwCompressor {
	return LzwCompressor{
		order:    lzw.MSB,
		litWidth: 8,
	}
}

func (c LzwCompressor) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := lzw.NewWriter(&b, c.order, c.litWidth)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

func (c LzwCompressor) Decompress(data []byte) ([]byte, error) {
	r := lzw.NewReader(bytes.NewReader(data), c.order, c.litWidth)
	defer r.Close()
	return ioutil.ReadAll(r)
}
