package lade

var _ PlanService = new(PlanClient)

type PlanClient struct {
	client *Client
}

type PlanOpts struct {
	ID string `qstring:"id,omitempty"`
}

type PlanService interface {
	Default() (*Plan, error)
	List() ([]*Plan, error)
	User(id string) ([]*Plan, error)
}

func (p *PlanClient) Default() (plan *Plan, err error) {
	err = p.client.doGet("plans/default", nil, &plan)
	return
}

func (p *PlanClient) List() (plans []*Plan, err error) {
	err = p.client.doList("plans", nil, &plans)
	return
}

func (p *PlanClient) User(id string) (plans []*Plan, err error) {
	opts := &PlanOpts{ID: id}
	err = p.client.doList("plans/user", opts, &plans)
	return
}
