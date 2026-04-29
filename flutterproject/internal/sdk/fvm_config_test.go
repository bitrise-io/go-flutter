package sdk

import (
	"io"
	"strings"
	"testing"

	"github.com/bitrise-io/go-flutter/flutterproject/internal/testassets"
	"github.com/stretchr/testify/require"
)

func Test_parseFVMFlutterVersion(t *testing.T) {
	tests := []struct {
		name            string
		fvmConfigReader io.Reader
		wantFlutterSDK  string
		wantChannel     string
		wantErr         string
	}{
		{
			name:            "Real fvm_config.json",
			fvmConfigReader: strings.NewReader(testassets.FVMConfigJSON),
			wantFlutterSDK:  "3.7.12",
		},
		{
			name: "Real fvm_config.json with channel",
			fvmConfigReader: strings.NewReader(`{
  "flutterSdkVersion": "2.2.2@beta",
  "flavors": {}
}`),
			wantFlutterSDK: "2.2.2",
			wantChannel:    "beta",
		},
		{
			name:            "Empty fvm_config.json",
			fvmConfigReader: strings.NewReader(""),
			wantErr:         "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlutterSDK, gotChannel, err := parseFVMFlutterVersion(tt.fvmConfigReader)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotFlutterSDK)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantFlutterSDK, gotFlutterSDK)
				require.Equal(t, tt.wantChannel, gotChannel)
			}
		})
	}
}

func Test_parseFVMRCFlutterVersion(t *testing.T) {
	tests := []struct {
		name           string
		fvmrcReader    io.Reader
		wantFlutterSDK string
		wantChannel    string
		wantErr        string
	}{
		{
			name:           "Real .fvmrc",
			fvmrcReader:    strings.NewReader(testassets.FVMRC),
			wantFlutterSDK: "3.41.7",
		},
		{
			name:           "Real .fvmrc with channel",
			fvmrcReader:    strings.NewReader(testassets.FVMRCWithChannel),
			wantFlutterSDK: "3.41.7",
			wantChannel:    "stable",
		},
		{
			name:        "Empty .fvmrc",
			fvmrcReader: strings.NewReader(""),
			wantErr:     "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlutterSDK, gotChannel, err := parseFVMRCFlutterVersion(tt.fvmrcReader)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotFlutterSDK)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantFlutterSDK, gotFlutterSDK)
				require.Equal(t, tt.wantChannel, gotChannel)
			}
		})
	}
}
