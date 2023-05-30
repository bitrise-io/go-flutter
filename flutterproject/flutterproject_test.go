package flutterproject

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bitrise-io/go-flutter/flutterproject/internal/testassets"
	"github.com/bitrise-io/go-flutter/mocks"
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
	"ASDFFlutterVersion": "3.7.12",
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
