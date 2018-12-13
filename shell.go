package main

import (
	"strings"

	"github.com/chzyer/readline"
)

type Shell interface {
	Read() <-chan []byte
	Write([]byte)
	Close()
}

type shell struct {
	input chan []byte

	readline *readline.Instance
}

func NewShell() (Shell, error) {
	instance, err := readline.New("> ")
	if err != nil {
		return nil, err
	}

	s := &shell{
		readline: instance,
		input:    make(chan []byte),
	}

	go s.read()

	return s, nil
}

func (s *shell) Read() <-chan []byte {
	return s.input
}

func (s *shell) Write(b []byte) {
	if len(b) > 0 {
		println(string(b))
	}
}

func (s *shell) Close() {
	s.readline.Close()
}

func (s *shell) read() {
	for {
		line, err := s.readline.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSpace(line)
		s.input <- []byte(line)
	}
}
