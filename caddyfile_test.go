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
	"reflect"
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/logging"
	"go.uber.org/zap/zapcore"
)

func TestUnmarshalCaddyfile(t *testing.T) {
	type fields struct {
		LogEncoderConfig logging.LogEncoderConfig
		Encoder          zapcore.Encoder
		Template         string
		Placeholder      string
		UnescapeStrings  bool
	}
	type args struct {
		d *caddyfile.Dispenser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// old name
		{
			name: "formatted: no args",
			fields: fields{
				Template: commonLogFormat,
			},
			args: args{
				d: caddyfile.NewTestDispenser(`formatted`),
			},
			wantErr: false,
		},
		{
			name: "formatted: single argument",
			fields: fields{
				Template: "{obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`formatted "{obj1>obj2>[0]}"`),
			},
			wantErr: false,
		},
		{
			name: "formatted: multiple argument",
			fields: fields{
				Template: "{obj1>obj2>[0]} - {obj3>[2]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`formatted {obj1>obj2>[0]} - {obj3>[2]}`),
			},
			wantErr: false,
		},
		{
			name: "formatted: not template but given placeholder",
			fields: fields{
				Template:    commonLogFormat,
				Placeholder: "|",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`formatted {
					placeholder |
				}`),
			},
			wantErr: false,
		},
		{
			name: "formatted: given template and given placeholder",
			fields: fields{
				Template:    "{obj1>obj2>[0]}",
				Placeholder: "|",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`formatted "{obj1>obj2>[0]}" {
					placeholder |
				}`),
			},
			wantErr: false,
		},
		// new name
		{
			name: "transform: no args",
			fields: fields{
				Template: commonLogFormat,
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform`),
			},
			wantErr: false,
		},
		{
			name: "transform: single argument",
			fields: fields{
				Template: "{obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}"`),
			},
			wantErr: false,
		},
		{
			name: "transform: multiple argument",
			fields: fields{
				Template: "{obj1>obj2>[0]} - {obj3>[2]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform {obj1>obj2>[0]} - {obj3>[2]}`),
			},
			wantErr: false,
		},
		{
			name: "transform: not template but given placeholder",
			fields: fields{
				Template:    commonLogFormat,
				Placeholder: "|",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform {
					placeholder |
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: given template and given placeholder",
			fields: fields{
				Template:    "{obj1>obj2>[0]}",
				Placeholder: "|",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					placeholder |
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: multiple argument with alternative value.",
			fields: fields{
				Template: "{obj1>obj2>[0]:-obj3[0]} - {obj3>[2]:-obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform {obj1>obj2>[0]:-obj3[0]} - {obj3>[2]:-obj1>obj2>[0]}`),
			},
			wantErr: false,
		},
		{
			name: "transform: property `placeholder` set before other properties",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal: true,
				},
				Template:    "{obj1>obj2>[0]}",
				Placeholder: "|",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					placeholder |
					time_local
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: property `placeholder` set after other properties",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal: true,
				},
				Template:    "{obj1>obj2>[0]}",
				Placeholder: "|",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					time_local
					placeholder |
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: delegate unmarshaling in absence of `placeholder`",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal: true,
				},
				Template: "{obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					time_local
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: unmarshal multiple fields of upstream",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal:  true,
					TimeFormat: "iso8601",
				},
				Template: "{obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					time_local
					time_format iso8601
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: `placeholder` propert sitting between other properties",
			fields: fields{
				Placeholder: "-",
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal:  true,
					TimeFormat: "iso8601",
				},
				Template: "{obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					time_local
					placeholder -
					time_format iso8601
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: delegate unmarshaling in absence of template and `placeholder",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal: true,
				},
				Template: commonLogFormat,
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform {
					time_local
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: delegate unmarshaling in presence of template and absence of `placeholder",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal: true,
				},
				Template: "{obj1>obj2>[0]}",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					time_local
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: unquoted template",
			fields: fields{
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal:  true,
					TimeFormat: "iso8601",
				},
				Template: "[{ts}] {request>method} {request>host} {request>uri} {request>headers>User-Agent>[0]} - {request>proto} {status} {size} -",
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform [{ts}] {request>method} {request>host} {request>uri} {request>headers>User-Agent>[0]} - {request>proto} {status} {size} - {
					time_local
					time_format iso8601
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: not template but given unescape_strings",
			fields: fields{
				Template:        commonLogFormat,
				UnescapeStrings: true,
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform {
					unescape_strings
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: given template and given unescape_strings",
			fields: fields{
				Template:        "{obj1>obj2>[0]}",
				UnescapeStrings: true,
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					unescape_strings
				}`),
			},
			wantErr: false,
		},
		{
			name: "transform: `placeholder` `unescape_strings` and  property sitting between other properties",
			fields: fields{
				Template:        "{obj1>obj2>[0]}",
				Placeholder:     "-",
				UnescapeStrings: true,
				LogEncoderConfig: logging.LogEncoderConfig{
					TimeLocal:  true,
					TimeFormat: "iso8601",
				},
			},
			args: args{
				d: caddyfile.NewTestDispenser(`transform "{obj1>obj2>[0]}" {
					time_local
					placeholder -
					unescape_strings
					time_format iso8601
				}`),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := &TransformEncoder{
				Encoder: new(logging.JSONEncoder),
			}
			if err := se.UnmarshalCaddyfile(tt.args.d); (err != nil) != tt.wantErr {
				t.Fatalf("TransformEncoder.UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if se.Template != tt.fields.Template || se.Placeholder != tt.fields.Placeholder || se.UnescapeStrings != tt.fields.UnescapeStrings || !reflect.DeepEqual(se.LogEncoderConfig, tt.fields.LogEncoderConfig) {
				t.Fatalf("Unexpected marshalling error: expected = %+v, received: %+v", tt.fields, *se)
			}
		})
	}
}
