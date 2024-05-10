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

func CASPathTransformFunction(root string, id string, key string) PathKey {
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
		Pathname: root + "/" + id + "/" + strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(root string, id string, key string) PathKey

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

var DefaultPathTransportFunc = func(root string, id string, key string) PathKey {
	return PathKey{
		Pathname: root + "/" + id + "/" + key,
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

func (s *Store) Delete(id string, key string) error {
	pathKey := s.PathTransformFunc(s.Root, id, key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	return os.RemoveAll(pathKey.Filepath()) // TODO: this should clear up hanging folders as well...
}

func (s *Store) Has(id string, key string) bool {
	pathkey := s.PathTransformFunc(s.Root, id, key)
	_, err := os.Stat(pathkey.Filepath())
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Read(id string, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}
func (s *Store) WriteDecrypt(id string, encKey []byte, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	n, err := copyDecrypt(encKey, r, f)
	return int64(n), err

}

func (s *Store) writeStream(id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, r)
}

func (s *Store) openFileForWriting(id string, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(s.Root, id, key)
	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return nil, err
	}
	pathAndFilename := pathKey.Filepath()
	return os.Create(pathAndFilename)
}

func (s *Store) readStream(id string, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(s.Root, id, key)

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
