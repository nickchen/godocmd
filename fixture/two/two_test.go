package two

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFibCache(t *testing.T) {

	two := New()
	for i, expected := range []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55} {
		assert.Equal(t, expected, two.FibCache(i), fmt.Sprintf("fib %d should be %d", i, expected))
	}
}

func ExampleFTwo_FibCache() {
	fmt.Println(New().FibCache(10))
	// Output: 55
}

func BenchmarkFibCacheTen(b *testing.B) {
	two := New()
	for i := 0; i < b.N; i++ {
		two.FibCache(10)
    }
}
