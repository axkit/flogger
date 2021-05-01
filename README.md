# flogger [![GoDoc](https://pkg.go.dev/badge/github.com/axkit/flogger?status.svg)](https://pkg.go.dev/github.com/axkit/flogger) [![Build Status](https://travis-ci.org/axkit/flogger.svg?branch=main)](https://travis-ci.org/axkit/flogger) [![Coverage Status](https://coveralls.io/repos/github/axkit/flogger/badge.svg)](https://coveralls.io/github/flogger/flogger) [![Go Report Card](https://goreportcard.com/badge/github.com/axkit/flogger)](https://goreportcard.com/report/github.com/axkit/flogger)

# flogger
The flogger package provides a simple function call logging by wrapping [zerolog](https://github.com/rs/zerolog) package.

## Motivation
Simplify a routine of logging function's input parameters and execution duration. 

## Installation
```
go get -u github.com/axkit/flogger
```

## Getting Started
For simple logging use global logger flogger.G 

```
package main

import (
    "github.com/axkit/flogger"
)

func main() {
    
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    log.Info().Msg("application started")

    task("John", 100, 10)
}

func task(name string, a, b int) {
    fc := flogger.G.Enter(name, a, b)
    defer fc.Exit()

    // function logic
    // write intermediate lines to the log
    fc.Debug().Int("subtask", 1).Msg("completed")
}

// Output: 
{"time":1516134303,"level":"info","message":"application started"}
{"time":1516134304,"level":"debug","func":"task","params":"[John 100 10]","message":"enter"}
{"time":1516134310,"level":"debug","func":"task","subtask":1,"message":"completed"}
{"time":1516134315,"level":"debug","func":"task","dur":500,"message":"exit"}
```

## Usage 
Let's say we have following structure.
```
import "github.com/axkit/flogger"

type CustomerRepo struct {
    flog flogger.Logger
    ...    
}

func NewCustomerRepo(l *zerolog.Logger) *CustomerRepo {
    return &CustomerRepo{flog: flogger.New(l, "repo", "customer")}
}
```

### Examples
Writes 1 line: the function invocation fact to the log.
```
func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    repo.flog.Enter()
}
```

Writes 1 line: the function invocation fact with input parameters.
```
func (repo *CustomerRepo)NewCustomer(id int, name string, age int, ssn string) *Customer {
    repo.flog.Enter(id, name, age)
}
```

Writes 2 lines: the function invocation fact with input parameters and function execution duration on exit.
```
func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    defer repo.flog.Enter(id).Exit()
}
```

Writes 3 lines: the function invocation fact with input parameters. Writes the function execution duration on exit. Get possibility to write additional lines to the log.
```
func (repo *CustomerRepo)CustomerByID(id int) *Customer {

    // fc stands for function call
    fc := repo.flog.Enter(id)
    defer fc.Exit()

    fc.Debug().Int("subtask", 1).Msg("completed)
}
```

Writes 1 line: the function execution duration on exit. 
```
func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    defer repo.flog.Enter().Exit()
}
```

Writes 1 line: the function execution duration and input params on exit. 
```
func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    defer repo.flog.EnterSilent().Exit(id)
}
```

Writes 0 or 1 line: the function execution duration and input params on exit conditionally.
Writes nothing if OnExit() was not called.
```
func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    fc := repo.flog.EnterSilent()
    defer fc.Exit()

    if err := sql.Query(...); err != nil {
        fc.OnExit(id)
    }
    return nil
}
```
### Function Call Logging Customization
Exit method calls floggger.ExitHandler (declared as public variable) function to write log item on exit. 
As instance, we need extend exit call and save performance metrics.
```
flog := flogger.New(l, "repo", "customer")
flog.SetSecondExitHandler(func(fc *flogger.FuncCall){
    // call Prometheus
    // 
})

...
..

// write 1 line to the log on exit together with pushing 
// information to Prometeus. 
defer flog.EnterSilent().Exit(id, name, age)
```

## Performance
```
cpu: Intel(R) Core(TM) i7-6700HQ CPU @ 2.60GHz
BenchmarkLogger_Enter1-8                  714012      1660 ns/op  360 B/op    6 allocs/op
BenchmarkLogger_Enter3-8                  557878      2063 ns/op  400 B/op    6 allocs/op
BenchmarkLogger_EnterSilentExit1-8        614272      1995 ns/op  360 B/op    6 allocs/op
BenchmarkLogger_EnterSilentExit3-8        412462      2650 ns/op  400 B/op    6 allocs/op
BenchmarkLogger_EnterExit3-8              477002      2426 ns/op  400 B/op    6 allocs/op
```


## License
MIT






