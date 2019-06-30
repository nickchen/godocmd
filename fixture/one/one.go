
// Package one provide some novel functionality, for use in package testing
package one

type FOne struct {
	// Name of the instance object
	Name string
	// Value of the instance object
	Value string
}

// FOneSub use FOne as parent
type FOneSub struct {
	FOne
}

// _calculateSum should be 
func _calculateSum(a, b int) int {
	return a + b
}

// CalculateSum return sum of two integers
func (one *FOne) CalculateSum(a, b int) int {
	return _calculateSum(a, b)
}

// CalculateSum return sum of two integers
func (one *FOne) FibNoCache(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return n
	}
	return one.FibNoCache(n - 1) + one.FibNoCache(n - 2)
}