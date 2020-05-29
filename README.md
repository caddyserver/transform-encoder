Formtted Encoder Module for Caddy's Logger
===============================================

This plugin adds logging encoder named `formatted`. The module accepts a `template` with the placeholders are surrounded by
braces `{}` and filled by values extracted from the stucture of the JSON log encoder. The JSON configuration looks like this:
```json
{
	"encoder": "formatted",
	"template": "{...}"
}
```

The nesting is traversed using `>`. For example, to print the `uri` field, the traversal is templated as `{request>uri}`.

```json
{
	"request": {
		"method": "GET",
		"uri": "/",
		"proto": "HTTP/2.0",
		...
}
```

The Caddyfile configuration accepts the template immediately following the encoder name, and can be ommitted to assume Apache Common Log Format.

```caddyfile
log {
	format formatted <template>
}
```

## Install

First, the [xcaddy](https://github.com/caddyserver/xcaddy) command:

```shell
$ go get -u github.com/caddyserver/xcaddy/cmd/xcaddy
```

Then build Caddy with this Go module plugged in. For example:

```shell
$ xcaddy build --with github.com/caddyserver/format-encoder
```

