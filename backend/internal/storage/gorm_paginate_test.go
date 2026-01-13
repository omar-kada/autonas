package storage

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPaginate(t *testing.T) {
	tests := []struct {
		name          string
		cursor        Cursor[uint64]
		expectedCount int
	}{
		{
			name: "ValidCursor",
			cursor: Cursor[uint64]{
				Limit:  10,
				Offset: 100,
			},
			expectedCount: 10,
		},
		{
			name: "LargeLimit",
			cursor: Cursor[uint64]{
				Limit:  200,
				Offset: 1000,
			},
			expectedCount: 50,
		},
		{
			name: "ZeroLimit",
			cursor: Cursor[uint64]{
				Limit:  0,
				Offset: 100,
			},
			expectedCount: 10,
		},
		{
			name: "OffsetSmall",
			cursor: Cursor[uint64]{
				Limit:  10,
				Offset: 5,
			},
			expectedCount: 4,
		},
		{
			name: "NoResults",
			cursor: Cursor[uint64]{
				Limit:  10,
				Offset: 1,
			},
			expectedCount: 0,
		},
	}
	// Create a mock DB instance
	db, _ := gorm.Open(sqlite.Open(":memory:"))

	// Use a simple User model and migrate
	type User struct {
		ID   uint64 `gorm:"primaryKey"`
		Name string
	}
	db.AutoMigrate(&User{})

	// Seed rows so pagination has enough data
	for i := 1; i <= 60; i++ {
		db.Create(&User{Name: fmt.Sprintf("u%d", i)})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run the paginated query and verify semantics (IDs < offset and count == expected limit)
			var users []User
			db.Model(&User{}).Scopes(Paginate(tt.cursor)).Order("id desc").Find(&users)

			// Assert expected count and that every returned ID is < offset
			assert.Equal(t, tt.expectedCount, len(users))
			for _, u := range users {
				assert.Less(t, u.ID, tt.cursor.Offset)
			}
		})
	}
}
