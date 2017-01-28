# parallel

> Execute many functions in parallel. Immediately return if one returns an
> error.

## Background

A common Go pattern is to execute several independent tasks in parallel and wait
for them to all complete before proceeding. However, the tasks tend to be
all-or-nothing, and you'll want the whole thing to bail if any of them return an
error.


## Examples

Say you have a few long-ish running functions that return errors:

```go
import (
  "errors"
  "fmt"
  parallel "github.com/noffle/parallel"
)

func main() {
  err := parallel.Parallel(nil, nil, sleep, sleep)
  fmt.Println(err)
}

func sleep() error {
  fmt.Printf("sleeping..\n")
  time.Sleep(100 * time.Millisecond)
  fmt.Printf("..slept!\n")
  return errors.New("awoke")
}
```

```
sleeping..
sleeping..
..slept!
..slept!
awoke
```

You can also provide a `done` channel and a cancel function as the first two
parameters to cancel the parallel operation early:

```go

func main() {
  done := make(chan struct{})
  cancelFunc := func() { close(done) }

  parallel.Parallel(done, cancelFunc,
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

This is designed to be used with with `golang.org/x/net/context`!

```go
func main() {
  ctx, cancel := context.WithCancel(context.Background())

  parallel.Parallel(ctx.Done(), cancel,
    sleep,
    func() error {
      time.Sleep(50 * time.Millisecond)
      cancel()
    },
  )
}
```

## API

```go
import "github.com/noffle/parallel"
```

### func Parallel(done <-chan struct{}, cancel context.CancelFunc, ...func() error) error

Accepts an optional channel and cancel function to signal cancellation of the
parallel tasks, and a variable number of functions matching the signature
`func() error`.

The first function to return a non-nil error calls cancel() (if non-nil) and
propogates its error as the return value.

In the case that the done channel is closed, `Parallel` will return 
`parallel.Canceled`.

## Install

```
go get github.com/noffle/parallel
```

## License

MIT
