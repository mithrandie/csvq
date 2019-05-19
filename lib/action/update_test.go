package action

import (
	"errors"
	"reflect"
	"testing"
)

var versionIsLaterThanTests = []struct {
	Version1 *Version
	Version2 *Version
	Result   bool
}{
	{
		Version1: nil,
		Version2: nil,
		Result:   false,
	},
	{
		Version1: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
		Version2: nil,
		Result:   false,
	},
	{
		Version1: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
		Version2: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
		Result: false,
	},
	{
		Version1: &Version{
			Major: 1,
			Minor: 2,
			Patch: 1,
		},
		Version2: &Version{
			Major: 1,
			Minor: 3,
			Patch: 0,
		},
		Result: false,
	},
	{
		Version1: &Version{
			Major: 2,
			Minor: 2,
			Patch: 3,
		},
		Version2: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
		Result: true,
	},
	{
		Version1: &Version{
			Major: 1,
			Minor: 3,
			Patch: 3,
		},
		Version2: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
		Result: true,
	},
	{
		Version1: &Version{
			Major: 1,
			Minor: 2,
			Patch: 4,
		},
		Version2: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
		Result: true,
	},
	{
		Version1: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 0,
		},
		Version2: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 4,
		},
		Result: true,
	},
	{
		Version1: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 4,
		},
		Version2: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 0,
		},
		Result: false,
	},
	{
		Version1: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 4,
		},
		Version2: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 5,
		},
		Result: false,
	},
}

func TestVersion_IsLaterThan(t *testing.T) {
	for _, v := range versionIsLaterThanTests {
		result := v.Version1.IsLaterThan(v.Version2)
		if result != v.Result {
			t.Errorf("result = %t, want %t for %v, %v", result, v.Result, v.Version1, v.Version2)
		}
	}
}

func TestVersion_String(t *testing.T) {
	v := &Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}
	expect := "1.2.3"
	result := v.String()
	if result != expect {
		t.Errorf("result = %q, want %q", result, expect)
	}
}

var parseVersionTests = []struct {
	Input  string
	Result *Version
	Error  string
}{
	{
		Input: "1.2.3",
		Result: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
	},
	{
		Input: "v1.2.3",
		Result: &Version{
			Major: 1,
			Minor: 2,
			Patch: 3,
		},
	},
	{
		Input: "v1.2.3-pr.1",
		Result: &Version{
			Major:      1,
			Minor:      2,
			Patch:      3,
			PreRelease: 1,
		},
	},
	{
		Input: "",
		Error: "cannot parse to version",
	},
	{
		Input: "v1.2",
		Error: "cannot parse to version",
	},
	{
		Input: "va.2.3",
		Error: "cannot parse to version",
	},
	{
		Input: "v1.a.3",
		Error: "cannot parse to version",
	},
	{
		Input: "v1.2.a",
		Error: "cannot parse to version",
	},
}

func TestParseVersion(t *testing.T) {
	for _, v := range parseVersionTests {
		result, err := ParseVersion(v.Input)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Input, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Input, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Input, v.Error)
			continue
		}
		if !reflect.DeepEqual(result, v.Result) {
			t.Errorf("%s: result = %v, want %v", v.Input, result, v.Result)
		}
	}
}

var checkForUpdatesTests = []struct {
	Name           string
	CurrentVersion string
	IncludePR      bool
	GoOS           string
	GoArch         string
	ClientError    string
	LatestVersion  string
	PublishedAt    string
	Result         string
	Error          string
}{
	{
		Name:           "Github Client Error",
		CurrentVersion: "v1.2.3",
		GoOS:           "darwin",
		GoArch:         "amd64",
		ClientError:    "client error",
		Error:          "client error",
	},
	{
		Name:           "Invalid release number",
		CurrentVersion: "v1.2.3",
		GoOS:           "darwin",
		GoArch:         "amd64",
		LatestVersion:  "invalid",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Error:          "Invalid release number: invalid",
	},
	{
		Name:           "Current Version is Invalid",
		CurrentVersion: "invalid",
		GoOS:           "darwin",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.4",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result: "The current version is an invalid number.\n" +
			"The latest version is 1.2.4, released on Feb 13, 2019.\n" +
			"  Release URL: https://example.com",
	},
	{
		Name:           "Current is the latest version",
		CurrentVersion: "v1.2.3",
		GoOS:           "darwin",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.3",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result:         "The current version 1.2.3 is up to date.",
	},
	{
		Name:           "Latest version is available",
		CurrentVersion: "v1.2.3",
		GoOS:           "darwin",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.4",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result: "Version 1.2.4 is now available.\n" +
			"  Release Date: Feb 13, 2019\n" +
			"  Release URL:  https://example.com\n" +
			"  Download URL: https://example.com/download",
	},
	{
		Name:           "Latest version is available and executable binary does not exist",
		CurrentVersion: "v1.2.3",
		GoOS:           "unknown",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.4",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result: "Version 1.2.4 is now available.\n" +
			"  Release Date: Feb 13, 2019\n" +
			"  Release URL:  https://example.com",
	},
	{
		Name:           "Latest pre-release version is available",
		CurrentVersion: "v1.2.3",
		IncludePR:      true,
		GoOS:           "darwin",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.4-pr.1",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result: "Version 1.2.4-pr.1 is now available.\n" +
			"  Release Date: Feb 13, 2019\n" +
			"  Release URL:  https://example.com\n" +
			"  Download URL: https://example.com/download",
	},
	{
		Name:           "Latest pre-release version is available and executable binary does not exist",
		CurrentVersion: "v1.2.3",
		IncludePR:      true,
		GoOS:           "unknown",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.4-pr.1",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result: "Version 1.2.4-pr.1 is now available.\n" +
			"  Release Date: Feb 13, 2019\n" +
			"  Release URL:  https://example.com",
	},
	{
		Name:           "Latest pre-release version does not exist",
		CurrentVersion: "v1.2.4",
		IncludePR:      true,
		GoOS:           "darwin",
		GoArch:         "amd64",
		LatestVersion:  "v1.2.4-pr.1",
		PublishedAt:    "2019-02-13T00:00:00Z",
		Result:         "The current version 1.2.4 is up to date.",
	},
}

type GithubClientMock struct {
	LatestVersion string
	PublishedAt   string
	Error         string
}

func (c GithubClientMock) GetLatestRelease() (*GithubRelease, error) {
	if 0 < len(c.Error) {
		return nil, errors.New(c.Error)
	}

	htmlurl := "https://example.com"
	assetName := "csvq-v1.2.4-darwin-amd64.tar.gz"
	downloadURL := htmlurl + "/download"
	return &GithubRelease{
		TagName:     c.LatestVersion,
		PublishedAt: c.PublishedAt,
		HTMLURL:     htmlurl,
		Assets: []GithubReleaseAsset{
			{
				Name:               assetName,
				BrowserDownloadURL: downloadURL,
			},
		},
	}, nil
}

func (c GithubClientMock) GetLatestReleaseIncludingPreRelease() (*GithubRelease, error) {
	if 0 < len(c.Error) {
		return nil, errors.New(c.Error)
	}

	htmlurl := "https://example.com"
	assetName := "csvq-v1.2.4-pr.1-darwin-amd64.tar.gz"
	downloadURL := htmlurl + "/download"
	return &GithubRelease{
		TagName:     c.LatestVersion,
		PublishedAt: c.PublishedAt,
		HTMLURL:     htmlurl,
		Assets: []GithubReleaseAsset{
			{
				Name:               assetName,
				BrowserDownloadURL: downloadURL,
			},
		},
	}, nil
}

func TestCheckForUpdates(t *testing.T) {
	for _, v := range checkForUpdatesTests {
		CurrentVersion, _ = ParseVersion(v.CurrentVersion)
		client := &GithubClientMock{
			LatestVersion: v.LatestVersion,
			PublishedAt:   v.PublishedAt,
			Error:         v.ClientError,
		}

		result, err := CheckForUpdates(v.IncludePR, client, v.GoOS, v.GoArch)
		if err != nil {
			if len(v.Error) < 1 {
				t.Errorf("%s: unexpected error %q", v.Name, err)
			} else if err.Error() != v.Error {
				t.Errorf("%s: error %q, want error %q", v.Name, err.Error(), v.Error)
			}
			continue
		}
		if 0 < len(v.Error) {
			t.Errorf("%s: no error, want error %q", v.Name, v.Error)
			continue
		}
		if result != v.Result {
			t.Errorf("%s: result = %v, want %v", v.Name, result, v.Result)
		}
	}

	CurrentVersion = &Version{}
}
