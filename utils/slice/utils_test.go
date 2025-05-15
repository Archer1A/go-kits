package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceToMap(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		slice := []int{}
		result := ToMap(slice, func(item int) int { return item })
		assert.Empty(t, result)
	})

	t.Run("slice of ints", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		result := ToMap(slice, func(item int) int { return item })
		assert.Equal(t, map[int]int{
			1: 1,
			2: 2,
			3: 3,
			4: 4,
			5: 5,
		}, result)
	})

	t.Run("slice of strings", func(t *testing.T) {
		slice := []string{"a", "bb", "ccc"}
		result := ToMap(slice, func(item string) int { return len(item) })
		assert.Equal(t, map[int]string{
			1: "a",
			2: "bb",
			3: "ccc",
		}, result)
	})

	t.Run("struct with custom key", func(t *testing.T) {
		type person struct {
			ID   int
			Name string
		}
		slice := []person{
			{ID: 1, Name: "Alice"},
			{ID: 2, Name: "Bob"},
		}
		result := ToMap(slice, func(p person) int { return p.ID })
		assert.Equal(t, map[int]person{
			1: {ID: 1, Name: "Alice"},
			2: {ID: 2, Name: "Bob"},
		}, result)
	})

	t.Run("duplicate keys - last one wins", func(t *testing.T) {
		slice := []int{1, 2, 3, 2, 1}
		result := ToMap(slice, func(item int) int { return item % 2 }) // keys will be 0 and 1
		assert.Equal(t, map[int]int{
			0: 2, // last even number
			1: 1, // last odd number
		}, result)
	})

	t.Run("custom comparable type as key", func(t *testing.T) {
		type myKey string
		slice := []string{"a", "b", "c"}
		result := ToMap(slice, func(item string) myKey { return myKey(item) })
		assert.Equal(t, map[myKey]string{
			"a": "a",
			"b": "b",
			"c": "c",
		}, result)
	})
}
