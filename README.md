Transform Encoder Module for Caddy's Logger
===============================================

This module adds logging encoder named `transform`. The module accepts a `template` with the placeholders are surrounded by
braces `{}` and filled by values extracted from the stucture of the JSON log encoder. The JSON configuration looks like this:
```json
{
	"encoder": "transform",
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
	format transform [<template>] {
		placeholder <string>
		message_key <key>
		level_key   <key>
		time_key    <key>
		name_key    <key>
		caller_key  <key>
		stacktrace_key <key>
		line_ending  <char>
		time_format  <format>
		level_format <format>
	}
}
```

The syntax of `template` is defined by the package [github.com/buger/jsonparser](https://github.com/buger/jsonparser). Objects are traversed using the key name. Arrays can be traversed by using the format `[index]`, as in `[0]`. For example, to get the first element in the `User-Agent` array, the template is `{request>headers>User-Agent>[0]}`.

## Examples

### Apache Common Log Format Example 

The module comes with one special value of `{common_log}` for the Apache Common Log format to simplify configuration

```caddyfile
 format transform "{common_log}"
```

The more spelled out way of doing it is:

```caddyfile
format transform `{request>remote_addr} - {request>user_id} [{ts}] "{request>method} {request>uri} {request>proto}" {status} {size}` {
	time_format "02/Jan/2006:15:04:05 -0700"
}
```

### Apache Combined Log Format Example

The more spelled out way of doing it is:

```caddy
format transform `{request>remote_addr} - {request>user_id} [{ts}] "{request>method} {request>uri} {request>proto}" {status} {size} "{request>headers>Referer>[0]}" "{request>headers>User-Agent>[0]}"` {
        time_format "02/Jan/2006:15:04:05 -0700"
}
```

## Install

First, the [xcaddy](https://github.com/caddyserver/xcaddy) command:

```shell
$ go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

Then build Caddy with this Go module plugged in. For example:

```shell
$ xcaddy build --with github.com/caddyserver/transform-encoder
```

