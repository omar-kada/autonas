package storage

import "math"

// Cursor represents a pagination cursor with limit and offset values.
type Cursor[T any] struct {
	Limit  int
	Offset T
}

// NewIDCursor creates a new cursor for uint64 offsets with the given limit and offset.
// If the offset is 0, it will be set to math.MaxUint64.
func NewIDCursor(limit int, offset uint64) Cursor[uint64] {
	if offset == 0 {
		offset = math.MaxInt64
	}
	return Cursor[uint64]{
		Limit:  limit,
		Offset: offset,
	}
}
