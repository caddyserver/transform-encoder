Transform Encoder Module for Caddy's Logger
===============================================

This module adds logging encoder named `transform`. The module accepts a `template` with the placeholders are surrounded
by
braces `{}` and filled by values extracted from the stucture of the JSON log encoder. The JSON configuration looks like
this:

```json
{
  "encoder": "transform",
  "template": "{...}"
}
```

The nesting is traversed using `>`. For example, to print the `uri` field, the traversal is templated as `{request>uri}`
.

```json
{
  "request": {
    "method": "GET",
    "uri": "/",
    "proto": "HTTP/2.0",
    ...
  }
```

The Caddyfile configuration accepts the template immediately following the encoder name, and can be ommitted to assume
Apache Common Log Format. The body of the block accepts the custom `placeholder` property in addition to all the properties listed within the [format modules section](https://caddyserver.com/docs/caddyfile/directives/log#format-modules).

```caddyfile
log {
	format transform [<template>] {
		placeholder <string>
		unescape_strings
		# other fields accepted by JSON encoder
	}
}
```

The syntax of `template` is defined by the package [github.com/buger/jsonparser](https://github.com/buger/jsonparser).
Objects are traversed using the key name. Arrays can be traversed by using the format `[index]`, as in `[0]`. For
example, to get the first element in the `User-Agent` array, the template is `{request>headers>User-Agent>[0]}`.

## Examples

### Apache Common Log Format Example

The module comes with one special value of `{common_log}` for the Apache Common Log format to simplify configuration

```caddyfile
 format transform "{common_log}"
```

The more spelled out way of doing it is:

```caddyfile
format transform `{request>client_ip} - {user_id} [{ts}] "{request>method} {request>uri} {request>proto}" {status} {size}` {
	time_format "02/Jan/2006:15:04:05 -0700"
}
```

### Apache Combined Log Format Example

The more spelled out way of doing it is:

```caddy
format transform `{request>client_ip} - {user_id} [{ts}] "{request>method} {request>uri} {request>proto}" {status} {size} "{request>headers>Referer>[0]}" "{request>headers>User-Agent>[0]}"` {
        time_format "02/Jan/2006:15:04:05 -0700"
}
```

# Alternative value

You can use an alternative value by using the following syntax `{val1:val2}`. For example, to show the `X-Forwarded-For`
header as `remote_ip` replacement you can do the following

```caddy
format transform `{request>headers>X-Forwarded-For>[0]:request>remote_ip} - {user_id} [{ts}] "{request>method} {request>uri} {request>proto}" {status} {size} "{request>headers>Referer>[0]}" "{request>headers>User-Agent>[0]}"` {
        time_format "02/Jan/2006:15:04:05 -0700"
}
```

The character `:` act as indicator for alternative value if the preceding key is not set.

For example, `{request>headers>X-Forwarded-For>[0]:request>remote_ip}` means that if `X-Forwarded-For` first array value is not empty use that otherwise fallback on `remote_ip`.

## Install

First, the [xcaddy](https://github.com/caddyserver/xcaddy) command:

```shell
$ go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

Then build Caddy with this Go module plugged in. For example:

```shell
$ xcaddy build --with github.com/caddyserver/transform-encoder
```

Alternatively, Caddy 2.4.4 and up supports adding a module directly:
```shell
$ caddy add-package github.com/caddyserver/transform-encoder
```

Note that in all cases Caddy should be stopped and started to get the module loaded, and future updates should be conducted through either xcaddy (if built with xcaddy) or `caddy upgrade` (if added with `add-package`).
