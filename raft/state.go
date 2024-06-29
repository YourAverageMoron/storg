package raft


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


type Log struct {
  filePath string
  logLength int64
  transport.Encoder //TODO SHOULD THE ENCODER BE IN TRANSPORT?
}

func (l *Log) Append(e []LogEntry) error{
  // TODO IMPLEMENT this 
}

func (l *Log) Read() ([]LogEntry, error) {
  //TODO implement this
}