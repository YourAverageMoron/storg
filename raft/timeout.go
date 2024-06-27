package raft

import (
    "fmt"
    "math/rand/v2"
    "time"
)

func ElectionTimeoutFunc() time.Duration {
  ms := rand.IntN(150) + 150
  return ms * time.Millisecond
}

type Timeout struct {
  durationFunc func() time.Duration
  resetch chan struct{}
  timeoutch chan struct{}
  quitch chan struct{}
}

func NewTimeout(durationFunc func() time.Duration) *Timeout {
  resetch := make(chan struct{})
  timeoutch := make(chan struct{})
  quitch := make(chan struct{})
  return &Timeout{
    durationFunc,
    resetch,
    timeoutch,
    resetch,
  }
}

func (t *Timeout) Start() {
  go t.loop()
}

func (t *Timeout) Consume() chan struct{} {
  return t.timeoutch
}

func (t *Timeout) Stop() {
  t.quitch <- struct{}{}
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
}