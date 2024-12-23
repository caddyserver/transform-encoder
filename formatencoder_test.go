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

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestEncodeEntry(t *testing.T) {
	tests := []struct {
		name              string
		se                TransformEncoder
		entry             zapcore.Entry
		fields            []zapcore.Field
		expectedLogString string
	}{
		{
			name: "encode entry: no unescape field",
			se: TransformEncoder{
				Encoder:  new(logging.JSONEncoder),
				Template: "{msg} {username}",
			},
			entry: zapcore.Entry{
				Message: "lob\nlaw",
			},
			fields: []zapcore.Field{
				zap.String("username", "john\ndoe"),
			},
			expectedLogString: "lob\\nlaw john\\ndoe\n",
		},
		{
			name: "encode entry: unescape field",
			se: TransformEncoder{
				Encoder:         new(logging.JSONEncoder),
				Template:        "{msg} {username}",
				UnescapeStrings: true,
			},
			entry: zapcore.Entry{
				Message: "lob\nlaw",
			},
			fields: []zapcore.Field{
				zap.String("username", "john\ndoe"),
			},
			expectedLogString: "lob\nlaw john\ndoe\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.se.Provision(caddy.Context{})
			if err != nil {
				t.Fatalf("TransformEncoder.Provision() error = %v", err)
			}

			buf, err := tt.se.EncodeEntry(tt.entry, tt.fields)

			if err != nil {
				t.Fatalf("TransformEncoder.EncodeEntry() error = %v", err)
			}

			if tt.expectedLogString != buf.String() {
				t.Fatalf("Unexpected encoding error: expected = %+v, received: %+v", tt.expectedLogString, buf)
			}

		})
	}
}
