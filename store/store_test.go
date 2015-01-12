package store

import (
	"github.com/hawx/tw-linkfeed/stream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStore(t *testing.T) {
	s := New(5)

	s.Add(stream.Tweet{Id: 1}) // 1 - - - -
	s.Add(stream.Tweet{Id: 2}) // 1 2 - - -
	s.Add(stream.Tweet{Id: 3}) // 1 2 3 - -
	s.Add(stream.Tweet{Id: 4}) // 1 2 3 4 -
	s.Add(stream.Tweet{Id: 5}) // 1 2 3 4 5
	s.Add(stream.Tweet{Id: 6}) // 6 2 3 4 5
	s.Add(stream.Tweet{Id: 7}) // 6 7 3 4 5

	arr := s.Latest()
	assert.Equal(t, 5, len(arr))
	assert.Equal(t, 3, arr[0].Id)
	assert.Equal(t, 4, arr[1].Id)
	assert.Equal(t, 5, arr[2].Id)
	assert.Equal(t, 6, arr[3].Id)
	assert.Equal(t, 7, arr[4].Id)
}

func TestNonFullStore(t *testing.T) {
	s := New(5)

	s.Add(stream.Tweet{Id: 1}) // 1 - - - -
	s.Add(stream.Tweet{Id: 2}) // 1 2 - - -

	arr := s.Latest()
	assert.Equal(t, 2, len(arr))
	assert.Equal(t, 1, arr[0].Id)
	assert.Equal(t, 2, arr[1].Id)
}

func TestEmptyStore(t *testing.T) {
	s := New(5)

	arr := s.Latest()
	assert.Equal(t, 0, len(arr))
}
