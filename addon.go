package lade

import "strconv"

var _ AddonService = new(AddonClient)

type AddonClient struct {
	client *Client
}

type AddonCreateOpts struct {
	Name     string `json:"name"`
	Service  string `json:"service"`
	PlanID   string `json:"plan_id"`
	RegionID string `json:"region_id"`
	Release  string `json:"release"`
	Public   bool   `json:"public"`
}

type AddonUpdateOpts struct {
	PlanID  string `json:"plan_id"`
	Release string `json:"release"`
	Public  bool   `json:"public"`
}

type AddonService interface {
	Create(opts *AddonCreateOpts) (*Addon, error)
	Delete(addon *Addon) error
	Get(id string) (*Addon, error)
	List() ([]*Addon, error)
	Update(id string, opts *AddonUpdateOpts) (*Addon, error)
}

func (a *AddonClient) Create(opts *AddonCreateOpts) (addon *Addon, err error) {
	addon = new(Addon)
	err = a.client.doCreate("addons", opts, addon)
	return
}

func (a *AddonClient) Delete(addon *Addon) error {
	return a.client.doDelete("addons/"+strconv.Itoa(addon.ID), nil)
}

func (a *AddonClient) Get(id string) (addon *Addon, err error) {
	addon = new(Addon)
	err = a.client.doByID("addons", id, nil, addon)
	return
}

func (a *AddonClient) List() (addons []*Addon, err error) {
	err = a.client.doList("addons", nil, &addons)
	return
}

func (a *AddonClient) Update(id string, opts *AddonUpdateOpts) (addon *Addon, err error) {
	addon = new(Addon)
	err = a.client.doUpdate("addons/"+id, opts, addon)
	return
}
