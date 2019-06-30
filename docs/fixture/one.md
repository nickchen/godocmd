# one

import "github.com/nickchen/godocmd/fixture/one"

- [Overview](#Overview)
- [Index](#Index)
- [Examples](#Examples)

## Index


* [Functions](#Functions)
    * [func BenchmarkFibNoCacheTen(b *testing.B)](#func-benchmarkfibnocacheten)
    * [func ExampleFOne_FibNoCache()](#func-examplefone_fibnocache)
    * [func TestFibNoCache(t *testing.T)](#func-testfibnocache)

* [Types](#Types)
    * [type FOne](#type-fone)
    * [type FOneSub](#type-fonesub)


### Package files

 [one.go](../../fixture/one/one.go)  [one_test.go](../../fixture/one/one_test.go) 




## Functions

### func BenchmarkFibNoCacheTen

```
func BenchmarkFibNoCacheTen(b *testing.B)
```

### func ExampleFOne_FibNoCache

```
func ExampleFOne_FibNoCache()
```

### func TestFibNoCache

```
func TestFibNoCache(t *testing.T)
```




## Types

### type FOne
```
type FOne struct {
	// Name of the instance object
	Name	string
	// Value of the instance object
	Value	string
}
```
 

 

### type FOneSub
```
// FOneSub use FOne as parent
type FOneSub struct {
	FOne
}
```
 

 
 
 