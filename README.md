# Golang generator

> The code provides a mechanism to create a generator in Go that can be used to process elements via callbacks, suitable for handling tasks such as processing large files or databases. If you are familiar with PHP Generator, it is possible to achieve similar result.

`Note: The next version of golang plans to release a similar iterator with yield`

## Functions
- SetInitFunc(InitFunc)
- SetDeferFunc(DeferFunc)
- Next() <-chan interface{}
- GetLastError() error

### SetInitFunc(InitFunc)
> This is a function which is called first and can return any number of arguments and an error: Signature: 
```func(...interface{}) ([]interface{}, error)```

### SetDeferFunc(DeferFunc)
> This function is called when the process ended or the init function or callback function returned an error: Signature:
```func(...interface{}) error```


### Next() <-chan interface{}
> This is the iterator, in the for loop it returns a channel with ```interface{}```. Examle:
```
for res := range generator.Next() {
    fmt.Println(res)
}
```

### GetLastError() error
> Check at the end of the process if any error occurred


Examples:

- simple list:
```
package main

import (
	"fmt"

	"github.com/olbrichattila/gogenerator"
)

func main() {
	generator := gogenerator.New(callback)

	for res := range generator.Next() {
		fmt.Println(res)
	}

}

func callback(i int, p ...interface{}) (interface{}, error) {
	if i > 10 {
		return nil, nil
	}
	return fmt.Sprintf("hello %d", i), nil
}
```

- Simple list with init func
```
package main

import (
	"fmt"

	"github.com/olbrichattila/gogenerator"
)

func main() {
	generator := gogenerator.New(callback)

	for res := range generator.Next() {
		fmt.Println(res)
	}

}

func callback(i int, p ...interface{}) (interface{}, error) {
	if i > 10 {
		return nil, nil
	}
	return fmt.Sprintf("hello %d", i), nil
}
```

- Error handling from when it comes from init func
```
package main

import (
	"fmt"

	"github.com/olbrichattila/gogenerator"
)

func main() {
	generator := gogenerator.New(callback)
	generator.SetInitFunc(initFunc)

	for res := range generator.Next() {
		fmt.Println(res)
	}

	if generator.GetLastError() != nil {
		fmt.Println(generator.GetLastError())
	}
}

func initFunc(...interface{}) ([]interface{}, error) {
	return nil, fmt.Errorf("error from init func")
}

func callback(i int, p ...interface{}) (interface{}, error) {
	if i > 10 {
		return nil, nil
	}
	return fmt.Sprintf("hello %d", i), nil
}
```

- Error handling when it comes from callback
```
package main

import (
	"fmt"

	"github.com/olbrichattila/gogenerator"
)

func main() {
	generator := gogenerator.New(callback)
	generator.SetInitFunc(initFunc)

	for res := range generator.Next() {
		fmt.Println(res)
	}

	if generator.GetLastError() != nil {
		fmt.Println(generator.GetLastError())
	}
}

func initFunc(...interface{}) ([]interface{}, error) {
	return nil, nil
}

func callback(i int, p ...interface{}) (interface{}, error) {
	return "param", fmt.Errorf("error from callback")
}
```

- Full example with database read:
```
package main

import (
	"database/sql"
	"fmt"

	"github.com/olbrichattila/gogenerator"
	// This blank import is necessary to have the driver
	_ "github.com/mattn/go-sqlite3"
)

type vehicle struct {
	make  string
	model string
}

func main() {
	sql := "select make, model from vehicles"
	generator := gogenerator.New(callback, sql)

	generator.SetInitFunc(initFunc)
	generator.SetDeferFunc(deferFunc)

	for res := range generator.Next() {
		fmt.Println(res)
	}

	if generator.GetLastError() != nil {
		fmt.Println(generator.GetLastError())
	}
}

func initFunc(params ...interface{}) ([]interface{}, error) {
	if len(params) != 1 {
		return nil, fmt.Errorf("init function requires 1 parameter which is the sql string")
	}

	db, err := sql.Open("sqlite3", "./data/database.sqlite")
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(params[0].(string))
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	return []interface{}{db, stmt, rows}, nil
}

func deferFunc(params ...interface{}) error {
	if len(params) != 3 {
		return fmt.Errorf("two parameter required in defer func")
	}

	stmt, ok := params[1].(*sql.Stmt)
	if !ok {
		return fmt.Errorf("parameter 2 is not *sql.Stmt")
	}

	defer stmt.Close()

	db, ok := params[0].(*sql.DB)
	if !ok {
		return fmt.Errorf("parameter 1 is not *sql.DB")
	}
	defer db.Close()

	return nil
}

func callback(i int, params ...interface{}) (interface{}, error) {
	if len(params) != 3 {
		return nil, fmt.Errorf("two parameter required in callback")
	}

	rows, ok := params[2].(*sql.Rows)
	if !ok {
		return nil, fmt.Errorf("parameter 2 is not *sql.Rows")
	}

	hasNextRow := rows.Next()
	if !hasNextRow {
		return nil, nil
	}

	var vehicle vehicle
	rows.Scan(&vehicle.make, &vehicle.model)

	return vehicle, nil
}

```
