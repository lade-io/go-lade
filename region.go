package lade

var _ RegionService = new(RegionClient)

type RegionClient struct {
	client *Client
}

type RegionService interface {
	List() ([]*Region, error)
}

func (a *RegionClient) List() (regions []*Region, err error) {
	err = a.client.doList("regions", nil, &regions)
	return
}
