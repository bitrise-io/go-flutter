package sdk

import (
	"io"
	"strings"
	"testing"

	"github.com/bitrise-io/go-flutter/flutterproject/internal/testassets"
	"github.com/stretchr/testify/require"
)

func Test_parsePubspecLockSDKVersions(t *testing.T) {
	tests := []struct {
		name              string
		pubspecLockReader io.Reader
		wantFlutterSDK    string
		wantDartSDK       string
		wantErr           string
	}{
		{
			name:              "Real pubspec.lock",
			pubspecLockReader: strings.NewReader(testassets.PubspecLock),
			wantFlutterSDK:    "",
			wantDartSDK:       ">=2.19.6 <3.0.0",
		},
		{
			name: "pubspec.lock with Flutter SDK",
			pubspecLockReader: strings.NewReader(`sdks:
  dart: ">=2.19.6 <3.0.0"
  flutter: ">=3.7.12"`),
			wantFlutterSDK: ">=3.7.12",
			wantDartSDK:    ">=2.19.6 <3.0.0",
		},
		{
			name:              "Empty pubspec.lock",
			pubspecLockReader: strings.NewReader(""),
			wantErr:           "EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlutterSDK, gotDartSDK, err := parsePubspecLockSDKVersions(tt.pubspecLockReader)
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
