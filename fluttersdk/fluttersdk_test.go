package fluttersdk

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSDKVersionFinder_FindLatestReleaseFor(t *testing.T) {
	tests := []struct {
		name         string
		platform     Platform
		architecture Architecture
		channel      Channel
		query        SDKQuery
		wantVersion  string
		wantChannel  string
	}{
		{
			name:         "Defaults to the latest stable release",
			platform:     MacOS,
			architecture: ARM64,
			channel:      "",
			query:        SDKQuery{},
			wantVersion:  "3.13.9",
			wantChannel:  "stable",
		},
		{
			name:         "Filters for channel",
			platform:     MacOS,
			architecture: ARM64,
			channel:      "dev",
			wantVersion:  "2.13.0-0.1.pre",
			wantChannel:  "dev",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(flutterSDKsResponse))
		require.NoError(t, err)
	}))
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := SDKVersionFinder{
				SDKVersionLister: defaultSDKVersionLister{baseURLFormat: ts.URL + "/releases_%s.json"},
			}
			got, err := f.FindLatestReleaseFor(tt.platform, tt.architecture, tt.channel, tt.query)
			require.NoError(t, err)
			require.Equal(t, tt.wantVersion, got.Version)
			require.Equal(t, tt.wantChannel, got.Channel)
		})
	}
}

const flutterSDKsResponse = `{
	"base_url": "https://storage.googleapis.com/flutter_infra_release/releases",
	"current_release": {
		"beta": "476aa717cd342d11e16439b71f4f4c9209c50712",
		"dev": "13a2fb10b838971ce211230f8ffdd094c14af02c",
		"stable": "d211f42860350d914a5ad8102f9ec32764dc6d06"
	},
	"releases": [
		{
			"hash": "d211f42860350d914a5ad8102f9ec32764dc6d06",
			"channel": "stable",
			"version": "3.13.9",
			"dart_sdk_version": "3.1.5",
			"dart_sdk_arch": "x64",
			"release_date": "2023-10-25T21:53:10.833612Z",
			"archive": "stable/macos/flutter_macos_3.13.9-stable.zip",
			"sha256": "c14436a8b968d56616d8c99f646470160840f1047fd11e8124493c1c2706c4bf"
		},
		{
			"hash": "d211f42860350d914a5ad8102f9ec32764dc6d06",
			"channel": "stable",
			"version": "3.13.9",
			"dart_sdk_version": "3.1.5",
			"dart_sdk_arch": "arm64",
			"release_date": "2023-10-25T21:49:54.203764Z",
			"archive": "stable/macos/flutter_macos_arm64_3.13.9-stable.zip",
			"sha256": "374615f834f23cff70eaef3ef1c3ebd3f8246ebf4c7b7f100115c98bb32858bb"
		},
		{
			"hash": "d211f42860350d914a5ad8102f9ec32764dc6d06",
			"channel": "beta",
			"version": "3.13.9",
			"dart_sdk_version": "3.1.5",
			"dart_sdk_arch": "arm64",
			"release_date": "2023-10-25T21:49:54.203764Z",
			"archive": "stable/macos/flutter_macos_arm64_3.13.9-stable.zip",
			"sha256": "374615f834f23cff70eaef3ef1c3ebd3f8246ebf4c7b7f100115c98bb32858bb"
		},
		{
			"hash": "13a2fb10b838971ce211230f8ffdd094c14af02c",
			"channel": "dev",
			"version": "2.13.0-0.1.pre",
			"dart_sdk_version": "2.17.0 (build 2.17.0-266.1.beta)",
			"dart_sdk_arch": "arm64",
			"release_date": "2022-04-13T18:48:12.766415Z",
			"archive": "dev/macos/flutter_macos_arm64_2.13.0-0.1.pre-dev.zip",
			"sha256": "17eb685657cee569f8fd3aa2ccd80ec38577918df8bdb9a2552c97662c166f52"
		}
	]
}`
