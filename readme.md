# parallel

> Execute many functions in parallel. Immediately return if one returns an
> error.

## Background

A common Go pattern is to execute several independent tasks in parallel and wait
for them to all complete before proceeding. However, the tasks tend to be
all-or-nothing, and you'll want the whole thing to bail if any of them return an
error.


## Example

Say you have a few long-ish running functions:

```go
import (
  "fmt"

  parallel "github.com/noffle/parallel"
)

func main() {
	parallel.Parallel(nil, sleep, sleep)
  fmt.Println("done")
}

func sleep() error {
	fmt.Printf("sleeping..\n")
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("..slept!\n")
	return nil
}
```

```
sleeping..
sleeping..
..slept!
..slept!
done!
```

You can also provide a `chan struct{}` as the first parameter to cancel the
parallel operation early:

```go

func main() {
  var bail = make(chan struct{})

	parallel.Parallel(bail,
    sleep,
    func() error {
      fmt.Printf("going to bail..\n")
      time.Sleep(50 * time.Millisecond)
      close(bail)
      fmt.Printf("..bailed!\n")
    },
  )
}
```

outputs

```
sleeping..
going to bail..
..bailed!
```

## API

```go
import "github.com/noffle/parallel"
```

### func Parallel(chan struct{}, ...func() error) error

Accepts an optional channel to signal cancellation of the parallel tasks, and a
variable number of functions matching the signature `func() error`.

The first function to return a non-nil error closes the cancel channel (if
non-nil) and propogates its error as the return value.

In the case that the cancel channel is closed, `Parallel` will return
`parallel.Canceled`.

## Install

```
go get github.com/noffle/parallel
```

## License

MIT
