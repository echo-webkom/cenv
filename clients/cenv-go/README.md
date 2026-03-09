# cenv-go

Go package for `cenv` runtime validation.

## Install

```sh
go get github.com/echo-webkom/cenv/clients/cenv-go
```

## Use

```go
package main

import cenv "github.com/echo-webkom/cenv/clients/cenv-go"

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
