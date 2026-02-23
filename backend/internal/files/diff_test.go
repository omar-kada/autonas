package files

import (
	"testing"
)

func TestDiffText(t *testing.T) {
	tests := []struct {
		name     string
		oldStr   string
		newStr   string
		expected string
	}{
		{
			name:     "Simple change",
			oldStr:   "Hello, World!",
			newStr:   "Hello, Go!",
			expected: "- World\n+ Go\n",
		},
		{
			name:     "Multiple line changes",
			oldStr:   "Line 1\nLine 2\nLine 3",
			newStr:   "Line 1\nModified Line 2\nLine 3\nNew Line",
			expected: "+ Modified \n+ \nNew Line\n",
		},
		{
			name:     "Empty strings",
			oldStr:   "",
			newStr:   "New content",
			expected: "+ New content\n",
		},
		{
			name:     "No changes",
			oldStr:   "Same content",
			newStr:   "Same content",
			expected: "",
		},
		{
			name:     "Whitespace changes",
			oldStr:   "Hello  World",
			newStr:   "Hello World",
			expected: "-  \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DiffText(tt.oldStr, tt.newStr)
			if result != tt.expected {
				t.Errorf("DiffText() = '%v', want '%v'", result, tt.expected)
			}
		})
	}
}
