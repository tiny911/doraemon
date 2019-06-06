package localtracing

import (
	"testing"
)

func TestNewObjectId(t *testing.T) {
	id := NewObjectId().Hex()
	if id == "" {
		t.Errorf("id generate failed!")
	}
}

func BenchmarkNewObjectId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewObjectId().Hex()
	}
}
