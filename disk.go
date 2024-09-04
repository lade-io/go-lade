package lade

import "fmt"

var _ DiskService = new(DiskClient)

type DiskClient struct {
	client *Client
}

type DiskCreateOpts struct {
	Name   string `json:"name"`
	PlanID string `json:"plan_id"`
	Path   string `json:"path"`
}

type DiskUpdateOpts struct {
	PlanID string `json:"plan_id"`
}

type DiskService interface {
	Create(appID string, opts *DiskCreateOpts) (*Disk, error)
	Delete(disk *Disk) error
	Get(appID, diskID string) (*Disk, error)
	Head(appID, diskID string) error
	List(appID string) ([]*Disk, error)
	Update(appID, diskID string, opts *DiskUpdateOpts) (*Disk, error)
}

func (d *DiskClient) Create(appID string, opts *DiskCreateOpts) (disk *Disk, err error) {
	disk = new(Disk)
	err = d.client.doCreate("apps/"+appID+"/disks", opts, disk)
	return
}

func (d *DiskClient) Delete(disk *Disk) error {
	path := fmt.Sprintf("apps/%d/disks/%d", disk.AppID, disk.ID)
	return d.client.doDelete(path, nil)
}

func (d *DiskClient) Get(appID, diskID string) (disk *Disk, err error) {
	disk = new(Disk)
	err = d.client.doByID("apps/"+appID+"/disks", diskID, nil, disk)
	return
}

func (d *DiskClient) Head(appID, diskID string) error {
	return d.client.doByID("apps/"+appID+"/disks", diskID, nil, nil)
}

func (d *DiskClient) List(appID string) (disks []*Disk, err error) {
	err = d.client.doList("apps/"+appID+"/disks", nil, &disks)
	return
}

func (d *DiskClient) Update(appID, diskID string, opts *DiskUpdateOpts) (disk *Disk, err error) {
	disk = new(Disk)
	err = d.client.doUpdate("apps/"+appID+"/disks/"+diskID, opts, disk)
	return
}
