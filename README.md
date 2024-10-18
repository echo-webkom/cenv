# cenv

Keeping your sanity in check by validating your `.env` files!

<br>

## How it works

`cenv` uses comments in `.env` files to generate a schema used for checking the env files integrity. When working on a larger project, env files can change a lot and sometimes break your app if you have forgotten to add certain values to your config. With `cenv` you mimize this risk by having a source of truth that makes sure you have your env set up correctly.

## Use

Add `@required` above fields that have to exist and have a non-empty value:

```py
NOT_SO_IMPORTANT=123

# @required
API_KEY=foo-bar-baz
```

Create a schema file from your env:

`cenv update`

Check you .env after fetching the latest changes

`cenv check`
