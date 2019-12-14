package lade

var _ UserService = new(UserClient)

type UserClient struct {
	client *Client
}

type UserService interface {
	Get() (*User, error)
}

func (u *UserClient) Get() (user *User, err error) {
	user = new(User)
	err = u.client.doGet("users/me", nil, user)
	return
}
