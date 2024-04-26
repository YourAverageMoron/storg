package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "file"
	pathkey := CASPathTransformFunction("root_dir", key)
	expectedOriginal := "971c419dd609331343dee105fffd0f4608dc0bf2"
	expectedPathName := "root_dir/971c4/19dd6/09331/343de/e105f/ffd0f/4608d/c0bf2"
	expectedPathKey := PathKey{
		Filename: expectedOriginal,
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
	key := "some_key"
	data := []byte("some bytes here")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}
	b, _ := io.ReadAll(r)
	assert.Equal(t, data, b)
	assert.True(t, s.Has(key))
	s.Delete(key)
	assert.True(t, !s.Has(key))
}
