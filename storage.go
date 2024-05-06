package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootDirectory = "gastore_data"

func CASPathTransformFunction(root string, key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)
	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		Pathname: root + "/" + strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(root string, key string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) Filepath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type StoreOpts struct {
	// Root is the storage root directory
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransportFunc = func(root string, key string) PathKey {
	return PathKey{
		Pathname: root + "/" + key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransportFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootDirectory
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(s.Root, key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	return os.RemoveAll(pathKey.Filepath()) // TODO: this should clear up hanging folders as well...
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(s.Root, key)
	_, err := os.Stat(pathkey.Filepath())
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Read(key string) (int64, io.Reader, error) {
    return s.readStream(key)
}

func (s *Store) Write(key string, r io.Reader) (int64, error) {
	return s.writeStream(key, r)
}

func (s *Store) readStream(key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(s.Root, key)

	file, err := os.Open(pathKey.Filepath())
	if err != nil {
		return 0, nil, err
	}

	stat, err := os.Stat(pathKey.Filepath())
	if err != nil {
		return 0, nil, err
	}
	return stat.Size(), file, nil
}

func (s *Store) writeStream(key string, r io.Reader) (int64, error) {
	pathKey := s.PathTransformFunc(s.Root, key)

	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return 0, err
	}

	pathAndFilename := pathKey.Filepath()

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return 0, err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return 0, err
	}

	return n, nil
}
