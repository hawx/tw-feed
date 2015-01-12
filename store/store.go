package store

import (
	"github.com/hawx/tw-stream"
)

type Store interface {
	Add(stream.Tweet)
	Latest() []stream.Tweet
}

type store struct {
	arr  []stream.Tweet
	here int
	fill bool
}

func New(size int) Store {
	return &store{make([]stream.Tweet, size), -1, false}
}

func (s *store) next() int {
	if s.here == -1 {
		s.here = 0
	} else if s.here == len(s.arr) - 1 {
		s.fill = true
		s.here = 0
	} else {
		s.here++
	}

	return s.here
}

func (s *store) Add(tweet stream.Tweet) {
	s.arr[s.next()] = tweet
}

func (s *store) Latest() []stream.Tweet {
	if !s.fill {
		return s.arr[0:s.here+1]
	}

	return append(s.arr[s.here+1:len(s.arr)], s.arr[0:s.here+1]...)
}
