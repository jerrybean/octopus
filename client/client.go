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

const (
	//RollbackTypeNone means no rollback
	RollbackTypeNone rollbackType = iota
	//RollbackTypeOne means just rollback single command,yet last failed one
	RollbackTypeOne
	//RollbackTypeBackTrace means rollback backtrack util one without rollback or the first one
	RollbackTypeBackTrace
	//RollbackTypeAll which is recommended is rollback all,this requires that each command should have rollback command
	RollbackTypeAll
)

//Client is a single ssh client
type Client struct {
	authTyp     authType
	rollbackTyp rollbackType

	address string
	user    string

	sshc *ssh.Client
	err  error
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
		sshc:        c,
		err:         err,
	}
}

func (c *Client) session() (*ssh.Session, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.sshc.NewSession()
}
