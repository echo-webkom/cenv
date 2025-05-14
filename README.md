<div align="center">

<img src=".github/logo.png" width="30%">

**Keeping your sanity in check by validating your `.env` files!**

**<a href="https://github.com/echo-webkom/cenv/releases/latest">Latest release ➜</a>**

</div>

<br>

## How it works

`cenv` uses comments in `.env` files to generate a schema used for checking the env files integrity. When working on a larger project, env files can change a lot and sometimes break your app if you have forgotten to add/edit/update certain fields. With `cenv` you mimize this risk by having a source of truth that makes sure your env is set up correctly.

## Install

### CLI app

Copy and run the following command. cenv will be put in `$HOME/.local/bin`. Make sure that the directory is in your `$PATH`.

```sh
curl -fsSL https://raw.githubusercontent.com/echo-webkom/cenv/refs/heads/main/install.sh | bash
```

Once installed, you can self-update with `cenv upgrade`.

### Go package

You can also use cenv as a single-util package. See [the package source](cenv.go).

```sh
go get github.com/echo-webkom/cenv
```

```go
func main() {
    err := cenv.Load()
    if err != nil {
        log.Fatal(err)
    }

    ...
}
```

## Use

Add one or more of the following `@tags` above a field:

- `public`: Marks the field as public. The value will be included in the schema. This is for required static values.
- `required`: Marks the field as required. The field has to be present, and have a non-empty value.
- `length [number]`: Requires a specified length for the fields value.
- `default [value]`: Set a default value. Running `cenv fix` will automatically fill this in if the field is empty.
- `enum [value1] | [value2] | ...`: Require that the field value is one of the given enum values, separated by `|`.
- `format [format]`: Requires a specified format for the value. Uses [gokenizer patterns](https://github.com/jesperkha/gokenizer).

```py
NOT_SO_IMPORTANT=123

# @required
API_KEY=foo-bar-baz

# @length 8
OTHER_KEY=abcdefgh

# Stacking multiple tags
# @required
# @length 4
# @format {number}
PIN_CODE=1234

# @enum user | guest | admin
# @default user
ROLE=user
```

Create a schema file from your .env:

```sh
# Creates a cenv.schema.json file
cenv update
```

Check your .env after fetching the latest changes

```sh
# Compares your env with the existing cenv.schema.json
cenv check
```

The `fix` command creates a .env file based on the schema, or fills in missing fields in an existing one. Any values already in your .env, like API keys will be kept, while any missing values will be added if a default or public one is provided.

```sh
cenv fix
```

## Building

To build the project, you need to have Go installed. Run the following command to build the project (make sure the `bin` directory exists):

```sh
go build -o bin/cenv app/main.go
```

If you want to overwrite the `Version` variable in `main.go` you have add the following flags:

```sh
go build -o bin/cenv -ldflags "-X 'github.com/echo-webkom/cenv/cmd.Version=<your-version>'" app/main.go
```

