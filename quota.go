package lade

var _ QuotaService = new(QuotaClient)

type QuotaClient struct {
	client *Client
}

type QuotaService interface {
	Max() (*Quota, error)
}

func (q *QuotaClient) Max() (quota *Quota, err error) {
	quota = new(Quota)
	err = q.client.doGet("quotas/max", nil, quota)
	return
}
