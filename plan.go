package lade

var _ PlanService = new(PlanClient)

type PlanClient struct {
	client *Client
}

type PlanService interface {
	List() ([]*Plan, error)
}

func (a *PlanClient) List() (plans []*Plan, err error) {
	err = a.client.doList("plans", nil, &plans)
	return
}
