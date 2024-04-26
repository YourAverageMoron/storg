package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "file"
	pathkey := CASPathTransformFunction(key)
	expectedOriginal := "971c419dd609331343dee105fffd0f4608dc0bf2"
	expectedPathName := "971c4/19dd6/09331/343de/e105f/ffd0f/4608d/c0bf2"
	expectedPathKey := PathKey{
		Original: expectedOriginal,
		Pathname: expectedPathName,
	}
	assert.Equal(t, pathkey, expectedPathKey)
	fmt.Println(pathkey)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunction,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some bytes here"))

	if err := s.writeStream("test_dir_name", data); err != nil {
		t.Error(err)
	}
}
