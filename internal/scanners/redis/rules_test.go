// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package redis

import (
	"reflect"
	"testing"

	"github.com/Azure/azqr/internal/ref"
	"github.com/Azure/azqr/internal/scanners"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redis/armredis"
)

func TestRedisScanner_Rules(t *testing.T) {
	type fields struct {
		rule        string
		target      interface{}
		scanContext *scanners.ScanContext
	}
	type want struct {
		broken bool
		result string
	}
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "RedisScanner DiagnosticSettings",
			fields: fields{
				rule: "redis-001",
				target: &armredis.ResourceInfo{
					ID: ref.Of("test"),
				},
				scanContext: &scanners.ScanContext{
					DiagnosticsSettings: map[string]bool{
						"test": true,
					},
				},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "RedisScanner Availability Zones",
			fields: fields{
				rule: "redis-002",
				target: &armredis.ResourceInfo{
					Zones: []*string{ref.Of("1"), ref.Of("2"), ref.Of("3")},
				},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "RedisScanner SLA",
			fields: fields{
				rule:        "redis-003",
				target:      &armredis.ResourceInfo{},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "99.9%",
			},
		},
		{
			name: "RedisScanner Private Endpoint",
			fields: fields{
				rule: "redis-004",
				target: &armredis.ResourceInfo{
					Properties: &armredis.Properties{
						PrivateEndpointConnections: []*armredis.PrivateEndpointConnection{
							{
								ID: ref.Of("test"),
							},
						},
					},
				},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "RedisScanner SKU",
			fields: fields{
				rule: "redis-005",
				target: &armredis.ResourceInfo{
					Properties: &armredis.Properties{
						SKU: &armredis.SKU{
							Name: getSKUNamePremium(),
						},
					},
				},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "Premium",
			},
		},
		{
			name: "RedisScanner CAF",
			fields: fields{
				rule: "redis-006",
				target: &armredis.ResourceInfo{
					Name: ref.Of("redis-test"),
				},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "RedisScanner disable non-SSL port",
			fields: fields{
				rule: "redis-008",
				target: &armredis.ResourceInfo{
					Properties: &armredis.Properties{
						EnableNonSSLPort: ref.Of(false),
					},
				},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
		{
			name: "RedisScanner minimum TLS version",
			fields: fields{
				rule: "redis-009",
				target: &armredis.ResourceInfo{
					Properties: &armredis.Properties{
						MinimumTLSVersion: getTLSVersion(),
					},
				},
				scanContext: &scanners.ScanContext{},
			},
			want: want{
				broken: false,
				result: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RedisScanner{}
			rules := s.GetRules()
			b, w := rules[tt.fields.rule].Eval(tt.fields.target, tt.fields.scanContext)
			got := want{
				broken: b,
				result: w,
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedisScanner Rule.Eval() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getSKUNamePremium() *armredis.SKUName {
	s := armredis.SKUNamePremium
	return &s
}

func getTLSVersion() *armredis.TLSVersion {
	s := armredis.TLSVersionOne2
	return &s
}
