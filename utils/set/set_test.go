package set

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		want *set[T]
	}
	tests := []testCase[string]{
		{
			name: "New",
			want: New[string](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New[string](); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Add(t *testing.T) {
	type args[T comparable] struct {
		key T
	}
	type testCase[T comparable] struct {
		name string
		set  *set[T]
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "add value",
			set:  New[string](),
			args: args[string]{
				key: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Add(tt.args.key)
		})
	}
}

func Test_set_Clear(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		set  set[T]
	}
	tests := []testCase[string]{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Clear()
		})
	}
}

func Test_set_Contain(t *testing.T) {
	type args[T comparable] struct {
		key T
	}
	type testCase[T comparable] struct {
		name string
		set  set[T]
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Contain(tt.args.key); got != tt.want {
				t.Errorf("Contain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Delete(t *testing.T) {
	type args[T comparable] struct {
		key T
	}
	type testCase[T comparable] struct {
		name string
		set  set[T]
		args args[T]
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.set.Delete(tt.args.key)
		})
	}
}

func Test_set_Len(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		set  set[T]
		want int
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Values(t *testing.T) {
	type testCase[T comparable] struct {
		name string
		set  set[T]
		want []T
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.set.Values(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Values() = %v, want %v", got, tt.want)
			}
		})
	}
}
