package sdk

import (
	"io"
	"path/filepath"
)

// sdkVersionParser extracts the raw Flutter and Dart SDK version constraints from a project file.
type sdkVersionParser func(io.Reader) (string, string, error)

// SDKVersionReader reads the SDK version constraints declared by a single project file.
type SDKVersionReader struct {
	fileOpener FileOpener
	relPath    string
	parse      sdkVersionParser
}

func (r SDKVersionReader) ReadSDKVersions(projectRootDir string) (*VersionConstraint, *VersionConstraint, error) {
	f, err := r.fileOpener.OpenReaderIfExists(filepath.Join(projectRootDir, r.relPath))
	if err != nil {
		return nil, nil, err
	}

	if f == nil {
		return nil, nil, nil
	}

	flutterVersionStr, dartVersionStr, err := r.parse(f)
	if err != nil {
		return nil, nil, err
	}

	var flutterVersion *VersionConstraint
	if flutterVersionStr != "" {
		flutterVersion, err = NewVersionConstraint(flutterVersionStr)
		if err != nil {
			return nil, nil, err
		}
	}

	var dartVersion *VersionConstraint
	if dartVersionStr != "" {
		dartVersion, err = NewVersionConstraint(dartVersionStr)
		if err != nil {
			return nil, nil, err
		}
	}

	return flutterVersion, dartVersion, nil
}
