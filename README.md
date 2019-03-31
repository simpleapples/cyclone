## Cyclone

Cyclone is a go package to spawn, reuse and manage a number of goroutines by a pool.

[![goreportcard for simpleapples/cyclone][1]][2]
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)


## Features

- Dynamically changing pool size.
- Automatically reusing a fixed number of goroutines.
- Succinct APIs make the pool is easy to use.

## Install

Go Get:

``` sh
go get github.com/simpleapples/cyclone
```

Dep:

``` sh
dep ensure -add github.com/simpleapples/cyclone
```

### Examples

```go
size := runtime.NumCPU()
pool := NewWithClosure(int64(size), func(payload interface{}) interface{} {
    intV := payload.(int)
    return intV
})
defer pool.Close()

for i := 0; i < size; i++ {
    result, err := pool.Run(i)
}
```

[1]: https://goreportcard.com/badge/github.com/simpleapples/cyclone
[2]: https://goreportcard.com/report/simpleapples/cyclone
