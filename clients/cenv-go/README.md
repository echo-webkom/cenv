# cenv-go

Go package for `cenv` runtime validation.

## Install

```sh
go get github.com/echo-webkom/cenv
```

## Use

```go
package main

import "github.com/echo-webkom/cenv"

func main() {
    config := cenv.Config{
        EnvPath: ".env",
        SchemaPath: "cenv.schema.toml",
    }

    if err := cenv.Check(config); err != nil {
        log.Fatal(err)
    }
}
```
