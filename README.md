# gtl
Extendable and minimalistic ETL toolkit in Go, built on generic io (std lib) pipelines

Index:
- [errors](#errors)
- [core interfaces](#core-interfaces)
- [core constructors](#core-constructors)
- [etl components](#etl-components)

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

<br>

Signatures are links to the Go playground (examples).
- [`type ReaderImpl[T any] struct`](https://go.dev/play/p/B_OXoh8V6Y-)
- [`type ReadCloserImpl[T any] struct`](https://go.dev/play/p/5GSJ1TZf2n5)
- [`type WriterImpl[T any] struct`](https://go.dev/play/p/ER8VOQ6VwRO)
- [`type WriteCloserImpl[T any] struct`](https://go.dev/play/p/rKTDQxIJgKf)
- [`type ReadWriterImpl[T, U any] struct`](https://go.dev/play/p/Ky2IE72bifw)
- [`type ReadWriteCloserImpl[T, U any] struct`](https://go.dev/play/p/DJ3AXmOpUJc)

</details>



## Core constructors
Core constructors for the most part facilitates interoperability between core interfaces and the `io` package. I.e conversion of `io.Reader` (bytes) into `core.Reader[T]`(generic vals) and back, and `io.Writer` (bytes) into `core.Writer[T]` (vals) and back. 
- [`func NewReaderFrom[T any](vs ...T) Reader[T]]`](
	https://go.dev/play/p/MAoiD4GNKVF
)
- [`func NewReaderFromBytes[T any](r io.Reader) func(f func(io.Reader) Decoder) Reader[T]`](
	https://go.dev/play/p/ud3zj4YT5QI
)
- [`func NewReaderFromValues[T any](r Reader[T]) func(f func(io.Writer) Encoder) io.Reader`](
	https://go.dev/play/p/FqjcyoRdASp
)
* [`func NewWriterFromValues[T any](w io.Writer) func(f func(io.Writer) Encoder) Writer[T]`](
	https://go.dev/play/p/2jYSMjo5Epr
)
* [`func NewWriterFromBytes[T any](w Writer[T]) func(f func(io.Reader) Decoder) io.Writer`](
	https://go.dev/play/p/P5Cp4piAWES
)
- [`func NewReadWriterFrom[T any](vs ...T) ReadWriter[T, T]`](
	https://go.dev/play/p/aS8fln6RiH2
)

Also, there are additional constructors for manipulating streams.
- [`func NewReaderWithBatching[T any](r Reader[T], size int) Reader[[]T]`](
	https://go.dev/play/p/5WhRAXCTBx9
)
- [`func NewReaderWithUnbatching[T any](r Reader[[]T]) Reader[T]`](
	https://go.dev/play/p/Bvzn7fRzqzF
)
- [`func NewWriterWithBatching[T any](w Writer[[]T], size int) Writer[T]`](
	https://go.dev/play/p/rlwM47TKAdr
)
- [`func NewWriterWithUnbatching[T any](w Writer[T]) Writer[[]T]`](
	https://go.dev/play/p/93GgwXIly5_V
)



## ETL Components
This section covers ETL components, they wrap [core interfaces](#core-interfaces) in order to provide some useful functionality. All links here go to the Go playground.

Logging
- [log.NewStreamedReader](https://go.dev/play/p/SY19PSZrr0a)
- [log.NewBatchedReader](https://go.dev/play/p/jYS_Zs3v7zw)
- [log.NewStreamedWriter](https://go.dev/play/p/NPztmctsrbQ)
- [log.NewBatchedWriter](https://go.dev/play/p/acwrPXfGrre)

Stats
- [stats.NewStreamedTeeReader](https://go.dev/play/p/xQOOBB9vG0A)
- [stats.NewBatchedTeeReader](https://go.dev/play/p/8T-eN52RPoE)