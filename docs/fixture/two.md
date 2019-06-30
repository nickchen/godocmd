# two

import "github.com/nickchen/godocmd/fixture/two"

- [Overview](#Overview)
- [Index](#Index)
- [Examples](#Examples)

## Index

* [Constants](#Constants)

* [Functions](#Functions)
    * [func BenchmarkFibCacheTen(b *testing.B)](#func-benchmarkfibcacheten)
    * [func ExampleFTwo_FibCache()](#func-exampleftwo_fibcache)
    * [func FibFunction(n int)(int)](#func-fibfunction)
    * [func TestFibCache(t *testing.T)](#func-testfibcache)

* [Types](#Types)
    * [type Cache](#type-cache)
        * [func NewCache()(*Cache)](#func-newcache)
    * [type FOneSub](#type-fonesub)
    * [type FTwo](#type-ftwo)
        * [func New()(*FTwo)](#func-new)


### Package files

 [two.go](../../fixture/two/two.go)  [two_test.go](../../fixture/two/two_test.go) 


## Constants

```
// A package level constants
const (
	// One is 1
	One	int	= iota
	// Two is 2
	Two
	Three
	Four
)
```




## Functions

### func BenchmarkFibCacheTen

```
func BenchmarkFibCacheTen(b *testing.B)
```

### func ExampleFTwo_FibCache

```
func ExampleFTwo_FibCache()
```

### func FibFunction
FibFunction uses no objects

```
func FibFunction(n int)(int)
```

### func TestFibCache

```
func TestFibCache(t *testing.T)
```




## Types

### type Cache
```
// Cache fixture for type
type Cache map[int]int
```
 



#### func NewCache()(*Cache)
NewCache doc string

```
func NewCache()(*Cache)
```
 
 

### type FOneSub
```
// FOneSub use FOne as parent
type FOneSub struct {
	FTwo
}
```
 

 

### type FTwo
```
// FTwo struct doc here
type FTwo struct {
	// Name of the instance object
	Name	string
	// Value of the instance object
	Value	string
	// contains filtered or unexported fields
}
```
 



#### func New()(*FTwo)
New get a new FTwo object

```
func New()(*FTwo)
```
 
 
 
 