# gtl
Extendable and minimalistic ETL toolkit in Go, built on generic io (std lib) pipelines

Index:
- [errors](#errors)
- [core interfaces](#core-interfaces)
- [core constructors](#core-constructors)

## Errors
<details>
    <summary>Expand/collapse section. </summary>

GTL tries to get out of your way and so only two errors are used in the core pkg, both inherited from `io` in the std lib:
```go
io.EOF              // Stop reading/pulling/consuming.
io.ErrClosedPipe    // Stop writing/pushing/producing.
```
</details>

## Core interfaces
There is one core interface and it is shown below. It is simply a generic variant of `io.Reader`combined with a `context.Context`. 
```go
type Reader[T any] interface {
	Read(context.Context) (T, error)
}
```


<br>
<details>
<summary>
As with the io package, there are varying combinations of basic interfaces, e.g io.ReadCloser. These groupings are for the most part mirrored here and can be viewed by clicking on this section.
</summary>

```go
type ReadCloser[T any] interface {
	io.Closer
	Reader[T]
}
```
</details>

<br>
<details>
<summary>
There are also "impl structs" which lets you implement most core interfaces with a function, allowing you to dodge boilerplate-y code. These can be viewed by clicking on this section.
</summary>

```go
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Calls impl.Impl.
func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error)
```

```go
type ReadCloserImpl[T any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
}

// Calls impl.ImplC.
func (impl ReadCloserImpl[T]) Close() (err error)

// Calls impl.ImplR.
func (impl ReadCloserImpl[T]) Read(ctx context.Context) (r T, err error)
```

</details>



## Core constructors
Core constructors for the most part convert `io.Reader` (bytes) into `core.Reader` (generic values) and back, here is a list of signatures:
- `func NewReaderFrom[T any](vs ...T) Reader[T]`
- `func NewReaderFromBytes[T any](r io.Reader) func(f func(io.Reader) Decoder) Reader[T]`
- `func NewReaderFromValues[T any](r Reader[T]) func(f func(io.Writer) Encoder) io.Reader`

Also, there are additional constructors for manipulating streams.
- `func NewReaderWithBatching[T any](r Reader[T], size int) Reader[[]T]`
- `func NewReaderWithUnbatching[T any](r Reader[[]T]) Reader[T]`