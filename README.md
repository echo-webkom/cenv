<div align="center">

<img src=".github/logo.png" width="30%">

**Keeping your sanity in check by validating your `.env` files!**

</div>

<br>

## How it works

`cenv` uses comments in `.env` files to generate a schema used for checking the env files integrity. When working on a larger project, env files can change a lot and sometimes break your app if you have forgotten to add/edit/update certain fields. With `cenv` you mimize this risk by having a source of truth that makes sure your env is set up correctly.

## Use

Add one or more of the following `@tags` above a field:

- `required`: Marks the field as required. The field has to be present, and have a non-empty value.
- `length [number]`: Requires a specified length for the fields value.

```py
NOT_SO_IMPORTANT=123

# @required
API_KEY=foo-bar-baz

# @length 8
OTHER_KEY=abcdefgh

# Stacking multiple tags
# @required
# @length 4
PIN_CODE=1234
```

Create a schema file from your env:

```sh
# Creates a cenv.schema.json file
env update
```

Check you .env after fetching the latest changes

```sh
# Compares your env with the existing cenv.schema.json
env check
```
