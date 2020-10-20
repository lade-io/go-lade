package lade

import "fmt"

var _ DomainService = new(DomainClient)

type DomainClient struct {
	client *Client
}

type DomainCreateOpts struct {
	Hostname string `json:"hostname"`
}

type DomainService interface {
	Create(appID string, opts *DomainCreateOpts) (*Domain, error)
	Delete(domain *Domain) error
	Get(appID, domainID string) (*Domain, error)
	Head(appID, domainID string) error
	List(appID string) ([]*Domain, error)
}

func (d *DomainClient) Create(appID string, opts *DomainCreateOpts) (domain *Domain, err error) {
	domain = new(Domain)
	err = d.client.doCreate("apps/"+appID+"/domains", opts, domain)
	return
}

func (d *DomainClient) Delete(domain *Domain) error {
	path := fmt.Sprintf("apps/%d/domains/%d", domain.AppID, domain.ID)
	return d.client.doDelete(path, nil)
}

func (d *DomainClient) Get(appID, domainID string) (domain *Domain, err error) {
	domain = new(Domain)
	err = d.client.doByID("apps/"+appID+"/domains", domainID, nil, domain)
	return
}

func (d *DomainClient) Head(appID, domainID string) error {
	return d.client.doByID("apps/"+appID+"/domains", domainID, nil, nil)
}

func (d *DomainClient) List(appID string) (domains []*Domain, err error) {
	err = d.client.doList("apps/"+appID+"/domains", nil, &domains)
	return
}
