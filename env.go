package lade

var _ EnvService = new(EnvClient)

type EnvClient struct {
	client *Client
}

type EnvSetOpts struct {
	Envs []*Env `json:"envs"`
}

func (e *EnvSetOpts) AddEnv(name, value string) {
	e.Envs = append(e.Envs, &Env{Name: name, Value: value})
}

type EnvUnsetOpts struct {
	Names []string `json:"names"`
}

type EnvService interface {
	List(appID string) ([]*Env, error)
	Set(appID string, opts *EnvSetOpts) ([]*Env, error)
	Unset(appID string, opts *EnvUnsetOpts) error
}

func (e *EnvClient) List(appID string) (envs []*Env, err error) {
	err = e.client.doList("apps/"+appID+"/envs", nil, &envs)
	return
}

func (e *EnvClient) Set(appID string, opts *EnvSetOpts) (envs []*Env, err error) {
	err = e.client.doUpdate("apps/"+appID+"/envs", opts, &envs)
	return
}

func (e *EnvClient) Unset(appID string, opts *EnvUnsetOpts) error {
	return e.client.doDelete("apps/"+appID+"/envs", opts)
}
