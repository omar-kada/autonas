package storage

import (
	"math"
	"testing"
)

func TestNewIDCursor(t *testing.T) {
	tests := []struct {
		name   string
		limit  int
		offset uint64
		want   Cursor[uint64]
	}{
		{
			name:   "ZeroOffset",
			limit:  10,
			offset: 0,
			want: Cursor[uint64]{
				Limit:  10,
				Offset: math.MaxInt64,
			},
		},
		{
			name:   "NonZeroOffset",
			limit:  10,
			offset: 5,
			want: Cursor[uint64]{
				Limit:  10,
				Offset: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIDCursor(tt.limit, tt.offset)
			if got.Limit != tt.want.Limit || got.Offset != tt.want.Offset {
				t.Errorf("NewIDCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}
