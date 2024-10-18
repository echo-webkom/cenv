<div align="center">

<img src=".github/logo.png" width="30%">

**Keeping your sanity in check by validating your `.env` files!**

</div>

<br>

## How it works

`cenv` uses comments in `.env` files to generate a schema used for checking the env files integrity. When working on a larger project, env files can change a lot and sometimes break your app if you have forgotten to add/edit/update certain fields. With `cenv` you mimize this risk by having a source of truth that makes sure your env is set up correctly.

## Use

Add `@required` above fields that have to exist and have a non-empty value:

```py
NOT_SO_IMPORTANT=123

# @required
API_KEY=foo-bar-baz

# @length 8
OTHER_KEY=abcdefgh
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
