package hashicorpreleases

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type ReleaseOptions struct {
	// Limit is the number of results returned. Maximum 20.
	Limit int
	// After is a timestamp used as a pagination marker,
	// indicating that only releases that occurred prior to
	// it should be retrieved. When fetching subsequent pages,
	// this parameter should be set to the creation
	// timestamp of the oldest release listed on the current page.
	// This needs to be a RFC3339 timestamp in string form.
	After string
	// LicenseClass can be either "enterprise" or "oss", used for returning
	// either enterprise versions or open source versions of HashiCorp
	// products.
	LicenseClass string
}

// ReleasesResponse is a list of Release
type ReleasesResponse []Release

//ReleaseMetadataResponse is a Release
type ReleaseMetadataResponse Release

// Release represents a single release and its metadata
type Release struct {
	// Builds is a list of builds
	Builds []Build `json:"builds"`
	// A docker image name and tag for this release in the format name:tag
	DockerNameTag string `json:"docker_name_tag"`
	// True if and only if this product release is a prerelease.
	IsPrerelease bool `json:"is_prerelease"`
	// The license class indicates if this is an enterprise product or an open source product.
	LicenseClass string `json:"license_class"`
	// The product name
	Name string `json:"name"`
	// Status of the product release
	Status Status `json:"status"`
	// Timestamp at which this product release was created (RFC3339 string)
	TimestampCreated string `json:"timestamp_created"`
	// Timestamp at which this product release was most recently updated.
	// This does not consider release status changes, such as when a release
	// transitions from supported to unsupported status --
	// that is tracked within the release status. (RFC3339 string)
	TimestampUpdated string `json:"timestamp_updated"`
	// URL for a blog post announcing this release;
	// Note that patch releases typically are not announced on the blog,
	// so this may refer to a major or minor parent release.
	BlogpostURL string `json:"url_blogpost"`
	// URL for the changelog covering this release
	ChangelogURL string `json:"url_changelog"`
	// URL for this product's docker image(s) on DockerHub
	DockerhubURL string `json:"url_docker_registry_dockerhub"`
	// URL for this product's docker image(s) on Amazon ECR-Public
	AmazonECRURL string `json:"url_docker_registry_ecr"`
	// URL for the software license applicable to this release
	LicenseURL string `json:"url_license"`
	// The project's website URL
	WebsiteURL string `json:"url_project_website"`
	// URL for this release's change notes
	ReleaseNotesURL string `json:"url_release_notes"`
	// URL for this release's file containing checksums of all the included build artifacts
	ShaSumsURL string `json:"url_shasums"`
	// An array of URLs, each pointing to a signature file. Each signature file
	// is a detached signature of the checksums file (see field url_shasums).
	// Signature files may or may not embed the signing key ID in the filename.
	ShaSumsSignaturesURL []string `json:"url_shasums_signatures"`
	// URL for the product's source code repository. This field is empty for enterprise products.
	SourceRepositoryURL string `json:"url_sorce_repository"`
	// The version of this release
	Version string `json:"version"`
}

// Build represents the architecture, OS, support status, and URL of a released binary
type Build struct {
	// The target architecture for this build
	Architecture string `json:"arch"`
	// The targeted operating system for this build
	OperatingSystem string `json:"os"`
	// True if this build is not supported by HashiCorp.
	// Some os/arch combinations are built by HashiCorp for
	// customer convenience but not officially supported.
	Unsupported bool `json:"unsupported"`
	// The URL where this build can be downloaded.
	URL string `json:"url"`
}

type Status struct {
	// Provides information about the most recent change; required when state="withdrawn"
	Message string
	// The state name of the release
	State string
	// The timestamp when the release status was last updated
	TimestampUpdated time.Time
}

// GetReleases retrieves the release metadata for multiple releases.
// This endpoint uses pagination for products with many releases.
// Results are ordered by release creation time from newest to oldest.
func (c *Client) GetReleases(product string, options *ReleaseOptions) (ReleasesResponse, error) {

	// Create the URL with ReleaseOptions as query parameters
	u := fmt.Sprintf("%s/releases/%s", c.URL, product)
	fullURL, err := handleReleaseOptions(u, options)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	setJSONHeader(req)

	// Issue the request against the API
	res := ReleasesResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// GetReleaseMetadata returns all metadata for a single product release
func (c *Client) GetReleaseMetadata(product string, version string) (*ReleaseMetadataResponse, error) {

	// Create the request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/releases/%s/%s", c.URL, product, version), nil)
	if err != nil {
		return nil, err
	}
	setJSONHeader(req)

	// Issue the request against the API
	res := ReleaseMetadataResponse{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func handleReleaseOptions(u string, options *ReleaseOptions) (string, error) {
	limit := 10
	after := time.Now().UTC().Format(time.RFC3339)
	if options != nil {
		if options.Limit != 0 {
			limit = options.Limit
		}
		if options.After != "" {
			after = options.After
		}
	}

	urlA, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	values := urlA.Query()
	values.Add("limit", strconv.Itoa(limit))
	values.Add("after", after)
	if options.LicenseClass != "" {
		values.Add("license_class", options.LicenseClass)
	}
	urlA.RawQuery = values.Encode()
	return urlA.String(), nil
}
