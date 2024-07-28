# gtl
Extendable and minimalistic ETL toolkit in Go, built on generic io (std lib) pipelines

Index:
- [errors](#errors)
- [core interfaces](#core-interfaces)

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