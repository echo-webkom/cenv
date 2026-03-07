<div align="center">

<img src=".github/logo.png" width="30%">

**Keeping your sanity in check by validating your `.env` files!**

**<a href="https://github.com/echo-webkom/cenv/releases/latest">Latest release ➜</a>**

</div>

<br>

## How it works

`cenv` uses a TOML schema file (`cenv.schema.toml`) to validate and manage your `.env` files. When working on a larger project, env files can change a lot and sometimes break your app if you have forgotten to add/edit/update certain fields. With `cenv` you minimize this risk by having a source of truth that makes sure your env is set up correctly.

## Install

### From source

Make sure you have [Rust](https://www.rust-lang.org/tools/install) installed, then run:

```sh
cargo install --path .
```

## Use

### Schema file

Create a `cenv.schema.toml` file that defines your environment variables:

```toml
[[entries]]
key = "DATABASE_URL"
required = true
kind = "Url"

[[entries]]
key = "API_PORT"
hint = "Port number as an integer"
required = true
default = "3000"
kind = { Integer = { min = 1, max = 65535 } }

[[entries]]
key = "LOG_LEVEL"
required = false
default = "info"
legal_values = ["debug", "info", "warn", "error"]
```

### Entry fields

Each entry supports the following fields:

- `key`: The environment variable name (required)
- `hint`: Human-readable description of the field
- `required`: If `true`, the field must be present and non-empty
- `default`: Default value used when generating the `.env` file
- `legal_values`: List of allowed values for the field
- `required_length`: Exact required length for the value
- `regex_match`: Regex pattern the value must match
- `kind`: Type validation - one of:
    - `String`
    - `Integer` (with optional `min`/`max`)
    - `Float` (with optional `min`/`max`)
    - `Bool`
    - `Url`
    - `Email`
    - `IpAddress`
    - `Path`

### Commands

Check your `.env` against the schema:

```sh
# Validates .env against cenv.schema.toml
cenv check
```

Generate or update your `.env` file from the schema:

```sh
# Creates .env or fills in missing values from schema defaults
cenv fix
```

The `fix` command preserves any existing values in your `.env` file (like API keys) while adding missing fields with their default values from the schema.

Both commands support custom paths:

```sh
cenv check --schema my-schema.toml --env .env.local
cenv fix --schema my-schema.toml --env .env.local
```

