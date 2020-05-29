package lade

var _ UserService = new(UserClient)

type UserClient struct {
	client *Client
}

type UserService interface {
	Me() (*User, error)
}

func (u *UserClient) Me() (user *User, err error) {
	user = new(User)
	err = u.client.doGet("users/me", nil, user)
	return
}
