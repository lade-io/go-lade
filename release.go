package lade

var _ ReleaseService = new(ReleaseClient)

type ReleaseClient struct {
	client *Client
}

type ReleaseCreateOpts struct {
	Source *File `structs:"source,omitnested"`
}

type ReleaseService interface {
	Create(appID string, opts *ReleaseCreateOpts) (*Release, error)
	Latest(appID string) (*Release, error)
	List(appID string) ([]*Release, error)
}

func (r *ReleaseClient) Create(appID string, opts *ReleaseCreateOpts) (release *Release, err error) {
	release = new(Release)
	err = r.client.doForm("apps/"+appID+"/releases", opts, release)
	return
}

func (r *ReleaseClient) Latest(appID string) (release *Release, err error) {
	release = new(Release)
	err = r.client.doGet("apps/"+appID+"/releases/latest", nil, release)
	return
}

func (r *ReleaseClient) List(appID string) (releases []*Release, err error) {
	err = r.client.doList("apps/"+appID+"/releases", nil, &releases)
	return
}
