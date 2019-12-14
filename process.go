package lade

import "fmt"

var _ ProcessService = new(ProcessClient)

type ProcessClient struct {
	client *Client
}

type ProcessCreateOpts struct {
	PlanID  string `json:"plan_id"`
	Command string `json:"command"`
}

type ProcessResizeOpts struct {
	Height uint `json:"height"`
	Width  uint `json:"width"`
}

type ProcessUpdateOpts struct {
	Processes []*Process `json:"processes"`
}

func (p *ProcessUpdateOpts) AddProcess(ptype, planID string, replicas int) {
	process := &Process{Type: ptype, PlanID: planID, Replicas: replicas}
	p.Processes = append(p.Processes, process)
}

type ProcessService interface {
	Attach(appID string, number int, handler ConnHandler) error
	Create(appID string, opts *ProcessCreateOpts) (*Process, error)
	List(appID string) ([]*Process, error)
	Resize(appID string, number int, opts *ProcessResizeOpts) error
	Update(appID string, opts *ProcessUpdateOpts) ([]*Process, error)
}

func (p *ProcessClient) Attach(appID string, number int, handler ConnHandler) error {
	path := fmt.Sprintf("apps/%s/processes/%d/attach", appID, number)
	return p.client.doWebsocket(path, handler)
}

func (r *ProcessClient) Create(appID string, opts *ProcessCreateOpts) (process *Process, err error) {
	process = new(Process)
	err = r.client.doCreate("apps/"+appID+"/processes", opts, process)
	return
}

func (p *ProcessClient) List(appID string) (processes []*Process, err error) {
	err = p.client.doList("apps/"+appID+"/processes", nil, &processes)
	return
}

func (p *ProcessClient) Resize(appID string, number int, opts *ProcessResizeOpts) error {
	path := fmt.Sprintf("apps/%s/processes/%d/resize", appID, number)
	return p.client.doUpdate(path, opts, nil)
}

func (p *ProcessClient) Update(appID string, opts *ProcessUpdateOpts) (processes []*Process, err error) {
	err = p.client.doUpdate("apps/"+appID+"/processes", opts, &processes)
	return
}
