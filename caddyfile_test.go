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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			se := &TransformEncoder{
				Encoder: new(logging.JSONEncoder),
			}
			if err := se.UnmarshalCaddyfile(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("TransformEncoder.UnmarshalCaddyfile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if se.Template != tt.fields.Template || se.Placeholder != tt.fields.Placeholder {
				t.Errorf("Unexpected marshalling error: expected = %+v, received: %+v", tt.fields, se)
			}
		})
	}
}
