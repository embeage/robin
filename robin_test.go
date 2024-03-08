package robin_test

import (
	"reflect"
	"testing"

	"github.com/embeage/robin"
)

func TestRobin(t *testing.T) {
	tests := []struct {
		name       string
		maxLen     int
		options    []robin.BoundedOption[int]
		operations []func(*robin.Robin[int]) interface{}
		want       []interface{}
	}{
		{
			name: "basic round-robin",
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { r.Add(4); v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
			},
			want: []interface{}{1, 2, 4, 3, 1, 2},
		},
		{
			name: "removing from robin and next on empty robin",
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); r.Remove(1); v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { r.Remove(2); v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { r.Remove(1, 3); _, ok := r.Next(); return ok },
			},
			want: []interface{}{2, 3, 3, false},
		},
		{
			name: "duplicates should be ignored",
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { return r.Len() },
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); return r.Len() },
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); return r.Len() },
				func(r *robin.Robin[int]) interface{} { r.Remove(1); return r.Len() },
				func(r *robin.Robin[int]) interface{} { r.Remove(2); return r.Len() },
				func(r *robin.Robin[int]) interface{} { r.Remove(3); return r.Len() },
			},
			want: []interface{}{0, 3, 3, 2, 1, 0},
		},
		{
			name: "basic contains",
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); return r.Contains(3) },
				func(r *robin.Robin[int]) interface{} { return r.Contains(0) },
				func(r *robin.Robin[int]) interface{} { r.Remove(3); return r.Contains(3) },
			},
			want: []interface{}{true, false, false},
		},
		{
			name:    "adding to full bounded robin without buffer should be no-op",
			maxLen:  2,
			options: []robin.BoundedOption[int]{},
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); return r.Len() },
				func(r *robin.Robin[int]) interface{} { r.Add(4); return r.Len() },
				func(r *robin.Robin[int]) interface{} { return r.Contains(3) },
				func(r *robin.Robin[int]) interface{} { return r.Contains(4) },
				func(r *robin.Robin[int]) interface{} { r.Remove(1); return r.Len() },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
			},
			want: []interface{}{2, 2, false, false, 1, 2, 2},
		},
		{
			name:    "adding and removing from full bounded robin with buffer should push and pop buffer",
			maxLen:  2,
			options: []robin.BoundedOption[int]{robin.WithBuffer[int](robin.NewLIFOBuffer[int](2))},
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); return r.Len() },
				func(r *robin.Robin[int]) interface{} { return r.Contains(3) },
				func(r *robin.Robin[int]) interface{} { return r.BufferLen() },
				func(r *robin.Robin[int]) interface{} { return r.BufferContains(3) },
				func(r *robin.Robin[int]) interface{} { r.Remove(2); return r.Len() },
				func(r *robin.Robin[int]) interface{} { return r.BufferContains(3) },
				func(r *robin.Robin[int]) interface{} { return r.BufferLen() },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
				func(r *robin.Robin[int]) interface{} { v, _ := r.Next(); return v },
			},
			want: []interface{}{2, false, 1, true, 2, false, 0, 1, 3},
		},
		{
			name:    "basic reset",
			maxLen:  2,
			options: []robin.BoundedOption[int]{robin.WithBuffer[int](robin.NewLIFOBuffer[int](2))},
			operations: []func(*robin.Robin[int]) interface{}{
				func(r *robin.Robin[int]) interface{} { r.Add(1, 2, 3); r.Reset(); return r.Len() },
				func(r *robin.Robin[int]) interface{} { return r.BufferLen() },
				func(r *robin.Robin[int]) interface{} { _, ok := r.Next(); return ok },
			},
			want: []interface{}{0, 0, false},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var r *robin.Robin[int]
			if tc.maxLen > 0 {
				r = robin.NewBounded[int](tc.maxLen, tc.options...)
			} else {
				r = robin.NewUnbounded[int]()
			}
			var got []interface{}
			for _, op := range tc.operations {
				got = append(got, op(r))
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Test %q failed: got %v, want %v", tc.name, got, tc.want)
			}
		})
	}
}
