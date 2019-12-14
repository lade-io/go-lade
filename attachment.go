package lade

var _ AttachmentService = new(AttachmentClient)

type AttachmentClient struct {
	client *Client
}

type AttachmentCreateOpts struct {
	Name string `json:"name"`
}

type AttachmentService interface {
	Create(appID, addonID string, opts *AttachmentCreateOpts) ([]*Attachment, error)
	Delete(appID, addonID string) error
	List(appID, addonID string) ([]*Attachment, error)
}

func (a *AttachmentClient) Create(appID, addonID string, opts *AttachmentCreateOpts) (
	attachments []*Attachment, err error) {
	err = a.client.doCreate("apps/"+appID+"/attachments/"+addonID, opts, attachments)
	return
}

func (a *AttachmentClient) Delete(appID, addonID string) error {
	return a.client.doDelete("apps/"+appID+"/attachments/"+addonID, nil)
}

func (a *AttachmentClient) List(appID, addonID string) (attachments []*Attachment, err error) {
	err = a.client.doList("apps/"+appID+"/attachments/"+addonID, nil, &attachments)
	return
}
