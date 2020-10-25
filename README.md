Formtted Encoder Module for Caddy's Logger
===============================================

This plugin adds logging encoder named `formatted`. The module accepts a `template` with the placeholders are surrounded by
braces `{}` and filled by values extracted from the stucture of the JSON log encoder. The JSON configuration looks like this:
```json
{
	"encoder": "formatted",
	"template": 
```

Then build Caddy with this Go module plugged in. For example:

```shell
$ xcaddy build --with github.com/caddyserver/format-encoder
```

