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
	s := newStore()
	defer teardown(t, s)
	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("file_name_%d", i)
		data := []byte("some bytes here")
		if _, err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		_, r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}
		b, _ := io.ReadAll(r)
		assert.Equal(t, data, b)
		assert.True(t, s.Has(key))
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}
		assert.True(t, !s.Has(key))
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunction,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
