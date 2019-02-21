package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mithrandie/csvq/lib/query"
)

const githubApiLatestReleaseURL = "https://api.github.com/repos/mithrandie/csvq/releases/latest"

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

var CurrentVersion *Version

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) IsLaterThan(v2 *Version) bool {
	if v == nil || v2 == nil {
		return false
	}

	if v.Major > v2.Major {
		return true
	}
	if v.Minor > v2.Minor {
		return true
	}
	if v.Patch > v2.Patch {
		return true
	}
	return false
}

func (v *Version) String() string {
	if v == nil {
		return ""
	}
	return strings.Join([]string{strconv.Itoa(v.Major), strconv.Itoa(v.Minor), strconv.Itoa(v.Patch)}, ".")
}

func ParseVersion(s string) (*Version, error) {
	s = PickVersionNumber(s)
	a := strings.Split(s, ".")

	if len(a) != 3 {
		return nil, errors.New("cannot parse to version")
	}

	major, err := strconv.Atoi(a[0])
	if err != nil {
		return nil, errors.New("cannot parse to version")
	}

	minor, err := strconv.Atoi(a[1])
	if err != nil {
		return nil, errors.New("cannot parse to version")
	}

	patch, err := strconv.Atoi(a[2])
	if err != nil {
		return nil, errors.New("cannot parse to version")
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

type GithubClient interface {
	GetLatestRelease() (*GithubRelease, error)
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
	body, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body, release)
	return release, err
}

func PickVersionNumber(s string) string {
	if 0 < len(s) && s[0] == 'v' {
		s = s[1:]
	}
	return s
}

func CheckUpdate() error {
	msg, err := CheckForUpdates(NewClient(), runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}

	query.Log(msg, false)
	return nil
}

func CheckForUpdates(client GithubClient, goos string, goarch string) (string, error) {
	rel, err := client.GetLatestRelease()
	if err != nil {
		return "", err
	}

	latestVersion, _ := ParseVersion(rel.TagName)
	if latestVersion == nil {
		return "", errors.New(fmt.Sprintf("Invalid release number: %s", rel.TagName))
	}

	publishedAt := ""
	if publishedTime, err := time.Parse(time.RFC3339, rel.PublishedAt); err == nil {
		publishedAt = publishedTime.Format("Jan 02, 2006")
	}

	if CurrentVersion == nil {
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
