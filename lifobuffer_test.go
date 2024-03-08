package robin_test

import (
	"reflect"
	"testing"

	"github.com/embeage/robin"
)

func TestLIFOBuffer(t *testing.T) {
	tests := []struct {
		name       string
		capacity   int
		operations []func(*robin.LIFOBuffer[int]) interface{}
		want       []interface{}
	}{
		{
			name:     "basic push and pop",
			capacity: 2,
			operations: []func(*robin.LIFOBuffer[int]) interface{}{
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(1); b.Push(2); v, _ := b.Pop(); return v },
				func(b *robin.LIFOBuffer[int]) interface{} { v, _ := b.Pop(); return v },
				func(b *robin.LIFOBuffer[int]) interface{} { _, ok := b.Pop(); return ok },
			},
			want: []interface{}{2, 1, false},
		},
		{
			name:     "basic len",
			capacity: 2,
			operations: []func(*robin.LIFOBuffer[int]) interface{}{
				func(b *robin.LIFOBuffer[int]) interface{} { return b.Len() },
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(1); return b.Len() },
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(2); return b.Len() },
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(3); return b.Len() },
			},
			want: []interface{}{0, 1, 2, 2},
		},
		{
			name:     "basic contains",
			capacity: 2,
			operations: []func(*robin.LIFOBuffer[int]) interface{}{
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(1); return b.Contains(1) },
				func(b *robin.LIFOBuffer[int]) interface{} { return b.Contains(2) },
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(1); b.Pop(); return b.Contains(1) },
				func(b *robin.LIFOBuffer[int]) interface{} { b.Pop(); return b.Contains(1) },
			},
			want: []interface{}{true, false, true, false},
		},
		{
			name:     "pushing to full buffer should overwrite oldest value",
			capacity: 2,
			operations: []func(*robin.LIFOBuffer[int]) interface{}{
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(1); b.Push(2); b.Push(3); return b.Contains(1) },
				func(b *robin.LIFOBuffer[int]) interface{} { v, _ := b.Pop(); return v },
			},
			want: []interface{}{false, 3},
		},
		{
			name:     "basic reset",
			capacity: 2,
			operations: []func(*robin.LIFOBuffer[int]) interface{}{
				func(b *robin.LIFOBuffer[int]) interface{} { b.Push(1); b.Reset(); return b.Len() },
				func(b *robin.LIFOBuffer[int]) interface{} { return b.Contains(1) },
			},
			want: []interface{}{0, false},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			b := robin.NewLIFOBuffer[int](tc.capacity)
			var got []interface{}
			for _, op := range tc.operations {
				got = append(got, op(b))
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Test %q failed: got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
