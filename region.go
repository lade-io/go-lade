package lade

var _ RegionService = new(RegionClient)

type RegionClient struct {
	client *Client
}

type RegionService interface {
	List() ([]*Region, error)
}

func (r *RegionClient) List() (regions []*Region, err error) {
	err = r.client.doList("regions", nil, &regions)
	return
}
