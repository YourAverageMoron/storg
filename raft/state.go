package raft

//TODO: MAKE THIS A GENERIC FOR DATA
type StoredState[T] struct {
  transport.Encoder
  filePath string
}

func (s *StoredState) Update(d T) error{
  //TODO: MARSHALL DATA AND STORE
}

func (s *StoredState) Get() (T, error) {
  // TODO: read file and unmarshall data
}