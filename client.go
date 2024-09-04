package lade

import (
	"net/http"

	"golang.org/x/oauth2"
)

type Client struct {
	httpClient *http.Client
	apiURL     string
	userAgent  string
	Addon      AddonService
	App        AppService
	Attachment AttachmentService
	Container  ContainerService
	Disk       DiskService
	Domain     DomainService
	Env        EnvService
	Log        LogService
	Plan       PlanService
	Process    ProcessService
	Quota      QuotaService
	Region     RegionService
	Release    ReleaseService
	Service    ServiceService
	User       UserService
}

var (
	DefaultScopes = []string{"app", "user", "offline"}
	Endpoint      = oauth2.Endpoint{
		AuthURL:   defaultAPIURL + "/login/oauth2/authorize",
		TokenURL:  defaultAPIURL + "/login/oauth2/token",
		AuthStyle: oauth2.AuthStyleInHeader,
	}
)

const (
	DefaultClientID  = "lade-client"
	defaultAPIURL    = "https://api.lade.io"
	apiVersion       = "/v1/"
	libraryUserAgent = "go-lade/" + libraryVersion
	libraryVersion   = "1.0"
)

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{
		httpClient: httpClient,
		apiURL:     defaultAPIURL + apiVersion,
		userAgent:  libraryUserAgent,
	}
	c.Addon = &AddonClient{client: c}
	c.App = &AppClient{client: c}
	c.Attachment = &AttachmentClient{client: c}
	c.Container = &ContainerClient{client: c}
	c.Disk = &DiskClient{client: c}
	c.Domain = &DomainClient{client: c}
	c.Env = &EnvClient{client: c}
	c.Log = &LogClient{client: c}
	c.Plan = &PlanClient{client: c}
	c.Process = &ProcessClient{client: c}
	c.Quota = &QuotaClient{client: c}
	c.Region = &RegionClient{client: c}
	c.Release = &ReleaseClient{client: c}
	c.Service = &ServiceClient{client: c}
	c.User = &UserClient{client: c}
	return c
}

func (c *Client) SetAPIURL(apiURL string) {
	c.apiURL = apiURL + apiVersion
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent + " " + libraryUserAgent
}
