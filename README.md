# Ledis driver mock for Golang

**ledismock** is a mock library designed to be dropped in place of [ledis](https://github.com/siddontang/ledisdb).  This library has one
purpose - to simulate a functioning ledis database connection in tests, without needing a real database connection.  It helps to maintain
a good **TDD** workflow.

- this library is new and incomplete.
- the only external dependency for this library is ledis.
- this library has strict default expection order matching, except in `HMset` commands, which can be in any order.
- this library does not require any changes to your source code (if you've written testable code).

## Install  
```go get github.com/replicatedhq/ledismock```

## Documentation and examples
Visit [godoc](http://godoc.org/github.com/replicatedhq/ledismock) for general examples and public API reference.

### Something you may want to test
```go
package main

import "github.com/siddontang/ledisdb/ledis"

var sharedDB *ledis.DB

init() {
      l, _ := ledis.Open(cfg)
      sharedDB, _ = l.Select(0)
}

func addNumber(key, val []byte) {
      sharedDB.SAdd(key, val)
}

func AddOneAndTwo() {
      addNumber("numbers", "one")
      addNumber("numbers", "two")
}
```

### Tests with ledismock
```go
package main

import (
      "testing"

      "github.com/replicatedhq/ledismock"
      "github.com/stretchr/testify/assert"
)

func TestMyFunc(t *testing.T) {
      db, mock, err := ledismock.New()
	    assert.Nil(t, err)
      assert.NotNil(t, db)
	    assert.NotNil(t, mock)

      // Replace the global db connection with a mock
	    sharedDB = db

      mock.ExpectSAdd().
		        WithKey("numbers").
		        WithValue("one").
		        WillReturnResult(ledismock.NewResult(1))

      mock.ExpectSAdd().
		        WithKey("numbers").
		        WithValue("two").
		        WillReturnResult(ledismock.NewResult(1))  

      AddOneAndTwo()

      err = mock.ExpectationsWereMet()
	    assert.Nil(t, err)          
}
```
