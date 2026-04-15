package converter

import (
	"path/filepath"
	"testing"
)

const testDir = "../test/"

func TestNodeBounds(t *testing.T) {
	tests := []struct {
		file string
		w, h int
	}{
		{"1_go.svg", 100, 37},
		{"2_php.svg", 100, 50},
		{"3_aws.svg", 99, 99},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			node, err := NewNode(filepath.Join(testDir, tt.file))
			if err != nil {
				t.Fatal(err)
			}
			if node.Width != tt.w || node.Height != tt.h {
				t.Errorf("got %dx%d, want %dx%d", node.Width, node.Height, tt.w, tt.h)
			}
		})
	}
}
