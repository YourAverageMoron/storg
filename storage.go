package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

const defaultRootDirectory = "gastore"

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

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(s.Root, key)
	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.Filename)
	}()
	return os.RemoveAll(pathKey.Filepath())
	// TODO this should clear up hanging folders as well...
}

func (s *Store) Has(key string) bool {
	pathkey := s.PathTransformFunc(s.Root, key)
	fmt.Println(pathkey.Filepath())
	f, err := os.Stat(pathkey.Filepath())
	if err == fs.ErrNotExist || f == nil {
		return false
	}

	return true
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(s.Root, key)
	return os.Open(pathKey.Filepath())
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(s.Root, key)

	if err := os.MkdirAll(pathKey.Pathname, os.ModePerm); err != nil {
		return err
	}

	pathAndFilename := pathKey.Filepath()

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)

	return nil
}
