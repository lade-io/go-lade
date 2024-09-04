package lade

var _ PlanService = new(PlanClient)

type PlanClient struct {
	client *Client
}

type PlanOpts struct {
	ID   string `qstring:"id,omitempty"`
	Type string `qstring:"type,omitempty"`
}

type PlanService interface {
	Default(ptype string) (*Plan, error)
	List(ptype string) ([]*Plan, error)
	User(id, ptype string) ([]*Plan, error)
}

func (p *PlanClient) Default(ptype string) (plan *Plan, err error) {
	opts := &PlanOpts{Type: ptype}
	err = p.client.doGet("plans/default", opts, &plan)
	return
}

func (p *PlanClient) List(ptype string) (plans []*Plan, err error) {
	opts := &PlanOpts{Type: ptype}
	err = p.client.doList("plans", opts, &plans)
	return
}

func (p *PlanClient) User(id, ptype string) (plans []*Plan, err error) {
	opts := &PlanOpts{ID: id, Type: ptype}
	err = p.client.doList("plans/user", opts, &plans)
	return
}
