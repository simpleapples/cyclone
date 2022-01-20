## Cyclone

[![goreportcard for simpleapples/cyclone][1]][2]
[![License: MIT][3]][4]

Cyclone is a go package to spawn, reuse and manage a number of goroutines.

## Install

Go Get:

``` sh
go get github.com/simpleapples/cyclone
```

Dep:

``` sh
dep ensure -add github.com/simpleapples/cyclone
```

## Examples

```go
size := runtime.NumCPU()
pool := NewWithClosure(int64(size), func(payload interface{}) interface{} {
    intV := payload.(int)
    fmt.Println(intV)
    return intV
})
defer pool.Close()

for i := 0; i < size; i++ {
    _, err := pool.Run(i)
}
```
Parallel jobs:

```go
size := 5
total := 20

wg := sync.WaitGroup{}

pool := NewWithCallback(int64(size), func(payload interface{}) interface{} {
    intV := payload.(int)
    time.Sleep(1 * time.Second)
    return intV
}, func(result interface{}) {
    intV := result.(int)
    fmt.Println(intV)
    wg.Done()
})

for i := 0; i < total; i++ {
    wg.Add(1)
    pool.Run(i)
}
wg.Wait()
```

Change pool size:

```go
pool.SetSize(100)
pool.SetSize(10000)
```

[1]: https://goreportcard.com/badge/github.com/simpleapples/cyclone
[2]: https://goreportcard.com/report/simpleapples/cyclone
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: https://opensource.org/licenses/MIT
[5]: https://img.shields.io/badge/license-NPL%20(The%20996%20Prohibited%20License)-blue.svg
[6]: https://github.com/996icu/996.ICU/blob/master/LICENSE
[7]: https://img.shields.io/badge/link-996.icu-red.svg
[8]: https://996.icu
