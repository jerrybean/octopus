package client

import "golang.org/x/crypto/ssh"

type authType byte

const (
	//AuthTypePassword is the type of ssh by password
	AuthTypePassword authType = iota
	//AuthTypePublicKey is the type of ssh by public key
	AuthTypePublicKey
)

type rollbackType byte

//Client is a single ssh client
type Client struct {
	authTyp     authType
	address     string
	user        string
	rollbackTyp rollbackType
	c           *ssh.Client
	err         error
}

//NewClientByPassword creates client by password
func NewClientByPassword(address, user, password string, rollbackTyp rollbackType) *Client {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	}
	c, err := ssh.Dial("tcp", address, config)
	return &Client{
		authTyp:     AuthTypePassword,
		address:     address,
		user:        user,
		rollbackTyp: rollbackTyp,
		c:           c,
		err:         err,
	}
}
