package flutterproject

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bitrise-io/go-flutter/flutterproject/internal/testassets"
	"github.com/bitrise-io/go-flutter/fluttersdk"
	"github.com/bitrise-io/go-flutter/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProject_FlutterAndDartSDKVersions(t *testing.T) {
	fileOpener := new(mocks.FileManager)
	fileOpener.On("OpenReaderIfExists", ".fvm/fvm_config.json").Return(strings.NewReader(testassets.FVMConfigJSON), nil)
	fileOpener.On("OpenReaderIfExists", ".tool-versions").Return(strings.NewReader(testassets.ToolVersions), nil)
	fileOpener.On("OpenReaderIfExists", "pubspec.lock").Return(strings.NewReader(testassets.PubspecLock), nil)
	fileOpener.On("OpenReaderIfExists", "pubspec.yaml").Return(strings.NewReader(testassets.PubspecYaml), nil)

	proj := Project{
		fileManager: fileOpener,
	}

	sdkVersions, err := proj.FlutterAndDartSDKVersions()
	require.NoError(t, err)

	b, err := json.MarshalIndent(sdkVersions, "", "\t")
	require.NoError(t, err)

	require.Equal(t, string(b), `{
	"FVMFlutterVersion": "3.7.12",
	"FVMFlutterChannel": "",
	"ASDFFlutterVersion": "3.7.12",
	"ASDFFlutterChannel": "",
	"PubspecFlutterVersion": {
		"Version": null,
		"Constraint": "^3.7.12"
	},
	"PubspecDartVersion": {
		"Version": null,
		"Constraint": "\u003e=2.19.6 \u003c3.0.0"
	},
	"PubspecLockFlutterVersion": null,
	"PubspecLockDartVersion": {
		"Version": null,
		"Constraint": "\u003e=2.19.6 \u003c3.0.0"
	}
}`)
}

func TestProject_FlutterSDKVersionToUse(t *testing.T) {
	tests := []struct {
		name                       string
		availableSDKsByChannel     map[string][]fluttersdk.Release
		projectSDKFromToolVersions string
		wantVersion                string
		wantChannel                string
	}{
		{
			name: "Project required version is available",
			availableSDKsByChannel: map[string][]fluttersdk.Release{"stable": {{
				Channel:        "stable",
				Version:        "3.13.8",
				DartSdkVersion: "3.1.4",
			}}},
			projectSDKFromToolVersions: "3.13.8",
			wantVersion:                "3.13.8",
			wantChannel:                "stable",
		},
		{
			name: "Project required version is not available",
			availableSDKsByChannel: map[string][]fluttersdk.Release{"stable": {{
				Channel:        "stable",
				Version:        "3.13.8",
				DartSdkVersion: "3.1.4",
			}}},
			projectSDKFromToolVersions: "3.13.9",
			wantVersion:                "",
			wantChannel:                "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			availableSDKLister := new(mocks.SDKVersionLister)
			availableSDKLister.On("ListReleasesByChannel", mock.Anything, mock.Anything).Return(tt.availableSDKsByChannel, nil)

			sdkVersionFinder := fluttersdk.SDKVersionFinder{SDKVersionLister: availableSDKLister}

			fileOpener := new(mocks.FileManager)
			fileOpener.On("OpenReaderIfExists", ".tool-versions").Return(strings.NewReader("flutter "+tt.projectSDKFromToolVersions), nil)
			fileOpener.On("OpenReaderIfExists", ".fvm/fvm_config.json").Return(nil, nil)
			fileOpener.On("OpenReaderIfExists", "pubspec.lock").Return(nil, nil)
			fileOpener.On("OpenReaderIfExists", "pubspec.yaml").Return(nil, nil)

			p := &Project{
				fileManager:      fileOpener,
				sdkVersionFinder: sdkVersionFinder,
			}
			gotVersion, gotChannel, err := p.FlutterSDKVersionToUse()
			require.NoError(t, err)
			require.Equal(t, tt.wantVersion, gotVersion)
			require.Equal(t, tt.wantChannel, gotChannel)
		})
	}
}
