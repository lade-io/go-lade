package lade

var _ ContainerService = new(ContainerClient)

type ContainerClient struct {
	client *Client
}

type ContainerService interface {
	List(appID string) ([]*Container, error)
}

func (d *ContainerClient) List(appID string) (containers []*Container, err error) {
	err = d.client.doList("apps/"+appID+"/containers", nil, &containers)
	return
}
