package lade

var _ PlanService = new(PlanClient)

type PlanClient struct {
	client *Client
}

type PlanService interface {
	Default() (*Plan, error)
	List() ([]*Plan, error)
	User() ([]*Plan, error)
}

func (p *PlanClient) Default() (plan *Plan, err error) {
	err = p.client.doGet("plans/default", nil, &plan)
	return
}

func (p *PlanClient) List() (plans []*Plan, err error) {
	err = p.client.doList("plans", nil, &plans)
	return
}

func (p *PlanClient) User() (plans []*Plan, err error) {
	err = p.client.doList("plans/user", nil, &plans)
	return
}
