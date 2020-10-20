package lade

import "strconv"

var _ AppService = new(AppClient)

type AppClient struct {
	client *Client
}

type AppCreateOpts struct {
	Name     string `json:"name"`
	RegionID string `json:"region_id"`
}

type AppService interface {
	Create(opts *AppCreateOpts) (*App, error)
	Delete(app *App) error
	Get(id string) (*App, error)
	Head(id string) error
	List() ([]*App, error)
}

func (a *AppClient) Create(opts *AppCreateOpts) (app *App, err error) {
	app = new(App)
	err = a.client.doCreate("apps", opts, app)
	return
}

func (a *AppClient) Delete(app *App) error {
	return a.client.doDelete("apps/"+strconv.Itoa(app.ID), nil)
}

func (a *AppClient) Get(id string) (app *App, err error) {
	app = new(App)
	err = a.client.doByID("apps", id, nil, app)
	return
}

func (a *AppClient) Head(id string) error {
	return a.client.doByID("apps", id, nil, nil)
}

func (a *AppClient) List() (apps []*App, err error) {
	err = a.client.doList("apps", nil, &apps)
	return
}
