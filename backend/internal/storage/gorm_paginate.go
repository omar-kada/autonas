package storage

import (
	"math"

	"gorm.io/gorm"
)

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


// Paginate returns a GORM scope function that applies pagination based on the given cursor.
// It uses the id field for cursor-based pagination with the condition "id < offset".
func Paginate[T any](c Cursor[T]) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := c.Offset

		pageSize := c.Limit
		switch {
		case pageSize > 50:
			pageSize = 50
		case pageSize <= 0:
			pageSize = 10
		}
		return db.Where("ID < ?", offset).Limit(pageSize)
	}
}
