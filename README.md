# flogger
Function call logger. Wrapper around zerolog

## Motivation
Automate a process of printing function's input parameters and function execution time.

## Usage
```
import "github.com/axkit/flogger"

type CustomerRepo struct {
    flog flogger.Logger
    ...    
}

type Filter struct {
    Name    string 
    Balance *int64
    ...
}

func NewCustomerRepo(l *zerolog.Logger) *CustomerRepo {
    return &CustomerRepo{flog : l.With("repo", "customer").Logger()}
}

 
func (repo *CustomerRepo)Customers(filter *Filter) []Customer {
    
    // filter = struct {Name : "ego", Balance: 300}
    fc := repo.flog.Enter(filter)
    defer fc.Exit()
    ...
    ..
    fc.Debug().Int("stage", 1).Msg("done")
    .
}

Output:
{"level":"debug","repo":"customer","func":"Customers","params":"[{Name:ego Balance:300}]","message":"enter"}
{"level":"debug","repo":"customer","func":"Customers","stage":1,"message":"done"}
{"level":"debug","repo":"customer","func":"Customers","dur":"12ms","message":"exit"}

func (repo *CustomerRepo)NewCustomers() []Customer {
    
    defer repo.flog.Enter().Exit()
  
    ...
    ..
    .    
}

Output:
{"level":"debug","repo":"customer","func":"NewCustomers","message":"enter"}
{"level":"debug","repo":"customer","func":"NewCustomers","dur":"12ms","message":"exit"}

func (repo *CustomerRepo)OldCustomers(age int) []Customer {
    
    defer repo.flog.Enter("ageLimit", age).Exit()
  
    ...
    ..
    .    
}
Output:
{"level":"debug","repo":"customer","func":"OldCustomers","params":"[ageLimit 42]","message":"enter"}
{"level":"debug","repo":"customer","func":"OldCustomers","dur":"12ms","message":"exit"}


func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    
    // Does not write "enter" message to the log.
    // Writes function's input params to the log on exit. 
    defer repo.flog.EnterSilent("id", id).Flush().Exit()
   
    ...
    ..
    .    
}

Output:
{"level":"debug","repo":"customer","func":"CustomerByID","params":"[id 2]","dur":"3ms","message":"enter/exit"}

func (repo *CustomerRepo)CustomerByID(id int) *Customer {
    
    // Does not write enter and exit messages to the log if everything fine. 
    // Write function's input params if error failed.
    fc := repo.flog.EnterSilent("id", id)
    defer fc.Exit()

    if c, err := orm.SelectOne(&Customer{}, id); err != nil {
        fc.Flush()
        fc.Error().Str("errmsg", err.Error()).Msg("ORM SelectOne() failed")
    }
   
    ...
    ..
    .    
}

Output:
if no error: 
    no line in the log
if error:
    {"level":"error","repo":"customer","func":"CustomerByID","errmsg" : "table not found", "message":"ORM SelectOne() failed"}
    {"level":"debug","repo":"customer","func":"CustomerByID","params":"[id 2]","dur":"3ms","message":"enter/exit"}

```



