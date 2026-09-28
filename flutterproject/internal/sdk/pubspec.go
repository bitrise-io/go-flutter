package sdk

import (
	"io"

	"gopkg.in/yaml.v3"
)

const PubspecRelPath = "pubspec.yaml"

func NewPubspecVersionReader(fileOpener FileOpener) SDKVersionReader {
	return SDKVersionReader{
		fileOpener: fileOpener,
		relPath:    PubspecRelPath,
		parse:      parsePubspecSDKVersions,
	}
}

func parsePubspecSDKVersions(pubspecReader io.Reader) (string, string, error) {
	type pubspec struct {
		Environment struct {
			Dart    string `yaml:"sdk"`
			Flutter string `yaml:"flutter"`
		} `yaml:"environment"`
	}

	var config pubspec
	d := yaml.NewDecoder(pubspecReader)
	if err := d.Decode(&config); err != nil {
		return "", "", err
	}

	return config.Environment.Flutter, config.Environment.Dart, nil
}
