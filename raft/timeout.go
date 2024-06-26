package raft

type Timeout struct {
  afterFunc func() Duration
  resetch chan struct{}
  timeoutch chan struct{}
}

func (t *Timeout) Start() chan struct{} {
  go t.loop()
  return t.timeoutch
}

func (t *Timeout) loop() {
  for {
    select {
    case reset := <-t.resetch:
        fmt.Println("received reset command, resetting timeout")
    case <-time.After(t.afterFunc()):
        fmt.Println("timeout exceeded, sending timeout message")
        t.timeoutch <- struct{}{}
  }
}