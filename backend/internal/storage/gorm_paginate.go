package storage

import (
	"gorm.io/gorm"
)

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
