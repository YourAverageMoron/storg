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
	if _, err := copyDecrypt(key, encryptedDst, decryptedDst); err != nil {
		t.Error(err)
	}

	assert.Equal(t, payload, decryptedDst.String())
	fmt.Println(decryptedDst.String())

}
