package raft

import "ryan-jones.io/storg/transport"

type StoredState[T any] struct {
	transport.Encoder
	filePath string
}

// TODO: SHOULD d BE AN IO.WRITER?
func (s *StoredState[T]) Update(d T) error {
	//TODO: MARSHALL DATA AND STORE
	return nil
}

// TODO: ANY SHOULD BE T
func (s *StoredState[T]) Get() (any, error) {
	// TODO: SHOULD THIS BE AN IO.READER?
	// TODO: read file and unmarshall data
	return nil, nil
}

type Log struct {
	filePath          string
	logLength         int64
	transport.Encoder //TODO SHOULD THE ENCODER BE IN TRANSPORT?
}

func (l *Log) Append(e []LogEntry) error {
	// TODO IMPLEMENT this
	return nil
}

func (l *Log) Read() ([]LogEntry, error) {
	//TODO implement this
	return nil, nil
}

