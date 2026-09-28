package sdk

import (
	"io"

	"gopkg.in/yaml.v3"
)

const pubspecLockRelPath = "pubspec.lock"

func NewPubspecLockVersionReader(fileOpener FileOpener) SDKVersionReader {
	return SDKVersionReader{
		fileOpener: fileOpener,
		relPath:    pubspecLockRelPath,
		parse:      parsePubspecLockSDKVersions,
	}
}

func parsePubspecLockSDKVersions(pubspecLockReader io.Reader) (string, string, error) {
	type pubspecLock struct {
		SDKs struct {
			Dart    string `yaml:"dart"`
			Flutter string `yaml:"flutter"`
		} `yaml:"sdks"`
	}

	var config pubspecLock
	d := yaml.NewDecoder(pubspecLockReader)
	if err := d.Decode(&config); err != nil {
		return "", "", err
	}

	return config.SDKs.Flutter, config.SDKs.Dart, nil
}
