package one

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFibNoCache(t *testing.T) {

	one := &FOne{}
	for i, expected := range []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55} {
		assert.Equal(t, expected, one.FibNoCache(i), fmt.Sprintf("fib %d should be %d", i, expected))
	}
}

func ExampleFOne_FibNoCache() {
	fmt.Println((&FOne{}).FibNoCache(10))
	// Output: 55
}

func BenchmarkFibNoCacheTen(b *testing.B) {
	one := &FOne{}
	for i := 0; i < b.N; i++ {
		one.FibNoCache(10)
    }
}
