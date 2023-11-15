package net

import (
	"encoding/json"
	"testing"
)

func TestPortRange_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		expect  *PortRange
	}{
		{"SingleIntPortTest", []byte("1080"), false, &PortRange{
			From: 1080,
			To:   1080,
		}},
		{"SingleStringPortTest", []byte("\"1090\""), false, &PortRange{
			From: 1090,
			To:   1090,
		}},
		{"RangePortTest", []byte("\"1080-1090\""), false, &PortRange{
			From: 1080,
			To:   1090,
		}},
		{"RangePortTestFail", []byte("1080-1090"), true, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &PortRange{}
			//err := v.UnmarshalJSON(tt.data)
			err := json.Unmarshal(tt.data, v)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				expectBs, _ := json.Marshal(tt.expect)
				actualBs, _ := json.Marshal(v)
				t.Log(string(actualBs))
				if v.From != tt.expect.From || v.To != tt.expect.To {
					t.Errorf("UnmarshalJSON() expect = %v, actual = %v", string(expectBs), string(actualBs))
				}
			}
		})
	}
}
