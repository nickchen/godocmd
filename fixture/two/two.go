
// Package two provide some novel functionality, for use in package testing
package two

// A package level constants
const (
	// One is 1
	One int = iota
	// Two is 2
	Two
	Three
	Four
)

// Cache fixture for type
type Cache map[int]int

// InitCache is the initial cache for fib calculation 
var InitCache Cache = map[int]int{0:0, 1:1}

// FTwo struct doc here
type FTwo struct {
	// Name of the instance object
	Name string
	// Value of the instance object
	Value string
	
	cache *Cache
}

// FOneSub use FOne as parent
type FOneSub struct {
	FTwo
}

// NewCache doc string
func NewCache() *Cache {
	r := make(Cache)
	for k, v := range InitCache {
		r[k] = v
	}
	return &r
}

// _calculateSum should be unexported
func _calculateSum(a, b int) int {
	return a + b
}

// CalculateSum return sum of two integers
func (one *FOneSub) CalculateSum(a, b int) int {
	return _calculateSum(a, b)
}

// New get a new FTwo object
func New() *FTwo {
	return &FTwo{cache: NewCache()}
}
// FibNoCache return fibonacci value at sequence index n
func (two *FTwo) FibCache(n int) int {
	if n < 0 {
		return 0
	}
	if r, ok := (*two.cache)[n]; !ok {
		r = two.FibCache(n - 1) + two.FibCache(n - 2)
		(*two.cache)[n] = r
		return r
	} else {
		return r
	}
}

// FibFunction uses no objects
func FibFunction(n int) int {
	return New().FibCache(n)
}