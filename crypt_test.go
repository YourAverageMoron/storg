package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCopyEncryptDecrypt(t *testing.T) {
	payload := "some info here"
	src := bytes.NewReader([]byte(payload))
	encryptedDst := new(bytes.Buffer)
	key := newEncryptionKey()
	_, err := copyEncrypt(key, src, encryptedDst)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(encryptedDst.String())
	decryptedDst := new(bytes.Buffer)
	nw, err := copyDecrypt(key, encryptedDst, decryptedDst)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payload) {
		t.Error("incorrect byted")
	}

	assert.Equal(t, payload, decryptedDst.String())
	fmt.Println(decryptedDst.String())

}
