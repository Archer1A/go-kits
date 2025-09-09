package set

import (
	"fmt"
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
func myFunc(m map[string]string) {

	m["2"] = "2"
	fmt.Printf("myfunc %v %p\n", m, m)

}
func setNil(m map[string]string) {
	m = make(map[string]string)
}

func TestMap(t *testing.T) {
	newMap := map[string]string{"hello": "world"}
	fmt.Printf("Before %v %p\n", newMap, newMap)
	myFunc(newMap)
	fmt.Printf("After %v %p\n", newMap, newMap)
	setNil(newMap)
	fmt.Printf("setNil %v %p\n", newMap, newMap)

}

func appendSlice(arr []string) {

	arr = append(arr, "world")
	fmt.Printf("func %v %p \n", arr, arr)
}

func updateIndex(arr []string) {
	arr[0] = "world"
}
func TestSlice(t *testing.T) {
	slice := []string{"hello"}
	fmt.Printf("Start %v %p \n", slice, slice)
	appendSlice(slice)
	fmt.Printf("append %v %p \n", slice, slice)
	updateIndex(slice)
	fmt.Printf("update %v %p \n", slice, slice)

}
