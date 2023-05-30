package sdk

import (
	"io"
	"strings"
	"testing"

	"github.com/bitrise-io/go-flutter/flutterproject/internal/testassets"
	"github.com/stretchr/testify/require"
)

func Test_parsePubspecSDKVersions(t *testing.T) {
	tests := []struct {
		name           string
		pubspecReader  io.Reader
		wantFlutterSDK string
		wantDartSDK    string
		wantErr        string
	}{
		{
			name:           "Real pubspec.yaml",
			pubspecReader:  strings.NewReader(testassets.PubspecYaml),
			wantFlutterSDK: "^3.7.12",
			wantDartSDK:    ">=2.19.6 <3.0.0",
		},
		{
			name:          "Empty pubspec.yaml",
			pubspecReader: strings.NewReader(""),
			wantErr:       "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlutterSDK, gotDartSDK, err := parsePubspecSDKVersions(tt.pubspecReader)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				require.Empty(t, gotFlutterSDK)
				require.Empty(t, gotDartSDK)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantFlutterSDK, gotFlutterSDK)
				require.Equal(t, tt.wantDartSDK, gotDartSDK)
			}
		})
	}
}
