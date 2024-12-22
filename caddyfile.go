// Copyright 2015 Matthew Holt and The Caddy Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transformencoder

import (
	"strings"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// UnmarshalCaddyfile sets up the module from Caddyfile tokens. Syntax:
//
//	transform [<template>] [{
//	     placeholder	[<placeholder>]
//	}]
//
// If the value of "template" is omitted, Common Log Format is assumed.
// See the godoc on the LogEncoderConfig type for the syntax of
// subdirectives that are common to most/all encoders.
func (se *TransformEncoder) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		args := d.RemainingArgs()
		switch len(args) {
		case 0:
			se.Template = commonLogFormat
		default:
			se.Template = strings.Join(args, " ")
		}

		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "placeholder":
				d.AllArgs(&se.Placeholder)
				// delete the `placeholder` token and the value, and reset the cursor
				d.Delete()
				d.Delete()
			case "unescape_strings":
				if d.NextArg() {
					return d.ArgErr()
				}
				d.Delete()
				se.UnescapeStrings = true
			default:
				d.RemainingArgs() //consume line without getting values
			}
		}
	}

	d.Reset()
	// consume the directive and the template
	d.RemainingArgs()

	return (&se.LogEncoderConfig).UnmarshalCaddyfile(d)
}
