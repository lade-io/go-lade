package lade

var _ ServiceService = new(ServiceClient)

type ServiceClient struct {
	client *Client
}

type ServiceService interface {
	Get(id string) (*Service, error)
	List() ([]*Service, error)
}

func (a *ServiceClient) Get(id string) (service *Service, err error) {
	service = new(Service)
	err = a.client.doByID("services", id, nil, service)
	return
}

func (s *ServiceClient) List() (services []*Service, err error) {
	err = s.client.doList("services", nil, &services)
	return
}
