# one

import "github.com/nickchen/godocmd/fixture/one"

- [Overview](#Overview)
- [Index](#Index)
- [Examples](#Examples)

## Index


* [Types](#Types)
    * [type FOne](#type-fone)
    * [type FOneSub](#type-fonesub)


### Package files

 [one.go](/Users/nickchen/Documents/GitHub/godocmd/fixture/one/one.go)  [one_test.go](/Users/nickchen/Documents/GitHub/godocmd/fixture/one/one_test.go) 





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
 

 



#### func CalculateSum(a, b int)(int)
CalculateSum return sum of two integers

```
func CalculateSum(a, b int)(int)
```
 

#### func FibNoCache(n int)(int)
CalculateSum return sum of two integers

```
func FibNoCache(n int)(int)
```


##### Example (FibNoCache)
```
func ExampleFOne_FibNoCache() {
	fmt.Println((&FOne{}).FibNoCache(10))

}
```
 
  


### type FOneSub
```
// FOneSub use FOne as parent
type FOneSub struct {
	FOne
}
```
 

 

 

 
 
<p align="center" ><small>automatically generated</small></p>
