package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/query"
)

const githubApiLatestReleaseURL = "https://api.github.com/repos/mithrandie/csvq/releases/latest"
const githubApiLatestPreReleaseURL = "https://api.github.com/repos/mithrandie/csvq/releases?per_page=1"
const preReleaseIdentifier = "pr"

type GithubRelease struct {
	HTMLURL     string               `json:"html_url"`
	TagName     string               `json:"tag_name"`
	PublishedAt string               `json:"published_at"`
	Assets      []GithubReleaseAsset `json:"assets"`
}

type GithubReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

var CurrentVersion = &Version{}

type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease int
}

func (v *Version) IsEmpty() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0 && v.PreRelease == 0
}

func (v *Version) IsLaterThan(v2 *Version) bool {
	if v == nil || v2 == nil {
		return false
	}

	if v.Major != v2.Major {
		return v.Major > v2.Major
	}
	if v.Minor != v2.Minor {
		return v.Minor > v2.Minor
	}
	if v.Patch != v2.Patch {
		return v.Patch > v2.Patch
	}
	if v.PreRelease != v2.PreRelease {
		if v.PreRelease == 0 {
			return true
		}
		if v2.PreRelease == 0 {
			return false
		}
		return v.PreRelease > v2.PreRelease
	}
	return false
}

func (v *Version) String() string {
	if v.PreRelease == 0 {
		return strings.Join([]string{strconv.Itoa(v.Major), strconv.Itoa(v.Minor), strconv.Itoa(v.Patch)}, ".")
	}
	return strings.Join([]string{strconv.Itoa(v.Major), strconv.Itoa(v.Minor), strconv.Itoa(v.Patch)}, ".") +
		"-" +
		strings.Join([]string{preReleaseIdentifier, strconv.Itoa(v.PreRelease)}, ".")

}

func ParseVersion(s string) (*Version, error) {
	v := &Version{}

	s = PickVersionNumber(s)
	words := strings.Split(s, "-")
	rVer := strings.Split(words[0], ".")

	if len(rVer) != 3 {
		return v, errors.New("cannot parse to version")
	}

	major, err := strconv.Atoi(rVer[0])
	if err != nil {
		return v, errors.New("cannot parse to version")
	}

	minor, err := strconv.Atoi(rVer[1])
	if err != nil {
		return v, errors.New("cannot parse to version")
	}

	patch, err := strconv.Atoi(rVer[2])
	if err != nil {
		return v, errors.New("cannot parse to version")
	}

	preRelease := 0
	if 1 < len(words) {
		prVer := strings.Split(words[1], ".")
		if len(prVer) != 2 {
			return v, errors.New("cannot parse to version")
		}
		if prVer[0] != preReleaseIdentifier {
			return v, errors.New("cannot parse to version")
		}
		preRelease, err = strconv.Atoi(prVer[1])
		if err != nil {
			return v, errors.New("cannot parse to version")
		}
	}

	v.Major = major
	v.Minor = minor
	v.Patch = patch
	v.PreRelease = preRelease
	return v, nil
}

type GithubClient interface {
	GetLatestRelease() (*GithubRelease, error)
	GetLatestReleaseIncludingPreRelease() (*GithubRelease, error)
}

type Client struct{}

func NewClient() GithubClient {
	return &Client{}
}

func (c Client) GetLatestRelease() (*GithubRelease, error) {
	res, err := http.Get(githubApiLatestReleaseURL)
	if err != nil {
		return nil, err
	}

	release := &GithubRelease{}
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &release)
	return release, err
}

func (c Client) GetLatestReleaseIncludingPreRelease() (*GithubRelease, error) {
	res, err := http.Get(githubApiLatestPreReleaseURL)
	if err != nil {
		return nil, err
	}

	release := []*GithubRelease{{}}
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &release)
	return release[0], err
}

func PickVersionNumber(s string) string {
	if 0 < len(s) && s[0] == 'v' {
		s = s[1:]
	}
	return s
}

func CheckUpdate(includePreRelaese bool) error {
	msg, err := CheckForUpdates(includePreRelaese, NewClient(), runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	return query.NewSession().WriteToStdoutWithLineBreak(msg)
}

func CheckForUpdates(includePreRelease bool, client GithubClient, goos string, goarch string) (string, error) {
	var rel *GithubRelease
	var err error

	if includePreRelease {
		rel, err = client.GetLatestReleaseIncludingPreRelease()
	} else {
		rel, err = client.GetLatestRelease()
	}
	if err != nil {
		return "", err
	}

	latestVersion, _ := ParseVersion(rel.TagName)
	if latestVersion.IsEmpty() {
		return "", errors.New(fmt.Sprintf("Invalid release number: %s", rel.TagName))
	}

	publishedAt := ""
	if publishedTime, err := time.Parse(time.RFC3339, rel.PublishedAt); err == nil {
		publishedAt = publishedTime.Format("Jan 02, 2006")
	}

	if CurrentVersion.IsEmpty() {
		return fmt.Sprintf("The current version is an invalid number.\nThe latest version is %s, released on %s.\n  Release URL: %s", latestVersion.String(), publishedAt, rel.HTMLURL), nil
	}

	if !latestVersion.IsLaterThan(CurrentVersion) {
		return fmt.Sprintf("The current version %s is up to date.", CurrentVersion.String()), nil
	}

	msg := fmt.Sprintf("Version %s is now available.\n  Release Date: %s\n  Release URL:  %s", latestVersion.String(), publishedAt, rel.HTMLURL)

	archiveName := fmt.Sprintf("csvq-v%s-%s-%s.tar.gz", latestVersion.String(), goos, goarch)
	for _, assets := range rel.Assets {
		if archiveName == assets.Name {
			msg = msg + fmt.Sprintf("\n  Download URL: %s", assets.BrowserDownloadURL)
		}
	}

	return msg, nil
}
