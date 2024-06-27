package raft

type Timeout struct {
  durationFunc func() Duration
  resetch chan struct{}
  timeoutch chan struct{}
  quitch chan struct{}
}

func (t *Timeout) Start() chan struct{} {
  go t.loop()
  return t.timeoutch
}

func (t *Timeout) Stop() {
  quitch <- strict{}{}
}

func (t *Timeout) Reset() {
  t.resetch <- struct{}{}
}

func (t *Timeout) loop() {
  for {
    select {
    case quit := <- t.quitch:
      return
    case reset := <-t.resetch:
        fmt.Println("received reset command, resetting timeout")
    case <-time.After(t.durationFunc()):
        fmt.Println("timeout exceeded, sending timeout message")
        t.timeoutch <- struct{}{}
  }
}