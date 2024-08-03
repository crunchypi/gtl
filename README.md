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
There are two core interfaces and they are shown below. They are simply generic variants of `io.Reader`and `io.Writer` combined with a `context.Context`. 
```go
type Reader[T any] interface {
	Read(context.Context) (T, error)
}
```

```go
type Writer[T any] interface {
	Write(context.Context, T) error
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

type WriteCloser[T any] interface {
	io.Closer
	Writer[T]
}

type ReadWriter[T, U any] interface {
	Reader[T]
	Writer[U]
}

type ReadWriteCloser[T, U any] interface {
	io.Closer
	Reader[T]
	Writer[U]
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

```go
type WriterImpl[T any] struct {
	Impl func(context.Context, T) error
}

// Calls impl.Impl.
func (impl WriterImpl[T]) Write(ctx context.Context, v T) (err error)
```

```go
type WriteCloserImpl[T any] struct {
	ImplC func() error
	ImplW func(context.Context, T) error
}

// Calls impl.ImplC
func (impl WriteCloserImpl[T]) Close() error 

// Calls impl.ImplW
func (impl WriteCloserImpl[T]) Write(ctx context.Context, v T) (err error)
```

```go
type ReadWriterImpl[T, U any] struct {
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

// Calls impl.ImplR
func (impl ReadWriterImpl[T, U]) Read(ctx context.Context) (r T, err error)

// Calls impl.ImplW
func (impl ReadWriterImpl[T, U]) Write(ctx context.Context, v U) (err error)
```

```go
type ReadWriteCloserImpl[T, U any] struct {
	ImplC func() error
	ImplR func(context.Context) (T, error)
	ImplW func(context.Context, U) error
}

// Calls impl.ImplC
func (impl ReadWriteCloserImpl[T, U]) Close() (err error)

// Calls impl.ImplR
func (impl ReadWriteCloserImpl[T, U]) Read(ctx context.Context) (r T, err error)

// Calls impl.ImplW
func (impl ReadWriteCloserImpl[T, U]) Write(ctx context.Context, v U) (err error)
```

</details>



## Core constructors
Core constructors for the most part facilitates interoperability between core interfaces and the `io` package. I.e conversion of `io.Reader` (bytes) into `core.Reader[T]`(generic vals) and back, and `io.Writer` (bytes) into `core.Writer[T]` (vals) and back. 
- `func NewReaderFrom[T any](vs ...T) Reader[T]`
- `func NewReaderFromBytes[T any](r io.Reader) func(f func(io.Reader) Decoder) Reader[T]`
- `func NewReaderFromValues[T any](r Reader[T]) func(f func(io.Writer) Encoder) io.Reader`
* `func NewWriterFromValues[T any](w io.Writer) func(f func(io.Writer) Encoder) Writer[T]`
* `func NewWriterFromBytes[T any](w Writer[T]) func(f func(io.Reader) Decoder) io.Writer`

Also, there are additional constructors for manipulating streams.
- `func NewReaderWithBatching[T any](r Reader[T], size int) Reader[[]T]`
- `func NewReaderWithUnbatching[T any](r Reader[[]T]) Reader[T]`
- `func NewWriterWithBatching[T any](w Writer[[]T], size int) Writer[T]`
- `func NewWriterWithUnbatching[T any](w Writer[T]) Writer[[]T]`