<div align="center">

<img src=".github/logo.png" width="30%">

**Keeping your sanity in check by validating your `.env` files!**

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

Once installed, you can self-update with `cenv upgrade` or the `cenv-install` binary (separate download).

### Go package

You can also use cenv as a single-util package. See [the package source](cenv.go).

```sh
go get github.com/echo-webkom/cenv
```

## Use

Add one or more of the following `@tags` above a field:

- `public`: Marks the field as public. The value will be included in the schema. This is for required static values.
- `required`: Marks the field as required. The field has to be present, and have a non-empty value.
- `length [number]`: Requires a specified length for the fields value.
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
```

Create a schema file from your env:

```sh
# Creates a cenv.schema.json file
cenv update
```

Check you .env after fetching the latest changes

```sh
# Compares your env with the existing cenv.schema.json
cenv check
```

You can fix and outdated .env with the fix command. Note that this will overwrite the existing .env, but use the values that were there before, like secret API keys etc. This may not work correctly if the .env is formatted incorrectly.

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
