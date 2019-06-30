<h1 id="two">two</h1>

<p>import “github.com/nickchen/godocmd/fixture/two”</p>

<ul>
  <li><a href="#Overview">Overview</a></li>
  <li><a href="#Index">Index</a></li>
  <li><a href="#Examples">Examples</a></li>
</ul>

<h2 id="index">Index</h2>

<ul>
  <li>
    <p><a href="#Constants">Constants</a></p>
  </li>
  <li><a href="#Functions">Functions</a>
    <ul>
      <li><a href="#func-BenchmarkFibCacheTen">func BenchmarkFibCacheTen(b *testing.B)</a></li>
      <li><a href="#func-ExampleFTwo_FibCache">func ExampleFTwo_FibCache()</a></li>
      <li><a href="#func-TestFibCache">func TestFibCache(t *testing.T)</a></li>
    </ul>
  </li>
  <li><a href="#Types">Types</a></li>
</ul>

<h3 id="package-files">Package files</h3>

<p><a href="../../fixture/two/two.go">two.go</a>  <a href="../../fixture/two/two_test.go">two_test.go</a></p>

<h2 id="functions">Functions</h2>

<h3 id="func-benchmarkfibcachetena-hreffunc-benchmarkfibcachetena">func BenchmarkFibCacheTen<a href="func-BenchmarkFibCacheTen"></a></h3>

<p><code>
func BenchmarkFibCacheTen(b *testing.B)
</code></p>

<h3 id="func-exampleftwofibcachea-hreffunc-exampleftwofibcachea">func ExampleFTwo_FibCache<a href="func-ExampleFTwo_FibCache"></a></h3>

<p><code>
func ExampleFTwo_FibCache()
</code></p>

<h3 id="func-testfibcachea-hreffunc-testfibcachea">func TestFibCache<a href="func-TestFibCache"></a></h3>

<p><code>
func TestFibCache(t *testing.T)
</code></p>

<h2 id="types">Types</h2>

<h3 id="type-cache">type Cache</h3>
<p><code>
// Cache fixture for type
type Cache map[int]int
</code></p>

<p><code>
const (
	// InitCache is the initial cache for fib calculation
	InitCache Cache = map[int]int{0: 0, 1: 1}
)
</code></p>

<h4 id="func-newcachecache">func NewCache()(*Cache)</h4>
<p><code>func NewCache()(*Cache)</code></p>

<p>NewCache doc string</p>

<h3 id="type-fonesub">type FOneSub</h3>
<p><code>
// FOneSub use FOne as parent
type FOneSub struct {
	FTwo
}
</code></p>

<h3 id="type-ftwo">type FTwo</h3>
<p><code>
// FTwo struct doc here
type FTwo struct {
	// Name of the instance object
	Name	string
	// Value of the instance object
	Value	string
	// contains filtered or unexported fields
}
</code></p>

<h4 id="func-newftwo">func New()(*FTwo)</h4>
<p><code>func New()(*FTwo)</code></p>

