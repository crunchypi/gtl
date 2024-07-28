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

There is also an impl struct which lets you implement a `core.Reader`with a function. 

```go
type ReaderImpl[T any] struct {
	Impl func(context.Context) (T, error)
}

// Calls impl.Impl.
func (impl ReaderImpl[T]) Read(ctx context.Context) (r T, err error)
```