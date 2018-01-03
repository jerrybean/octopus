package client

import (
	"errors"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type (
	authType byte

	rollbackType byte

	//Client is a single ssh client
	Client struct {
		authTyp     authType
		rollbackTyp rollbackType

		address string
		user    string

		sshc *ssh.Client
		err  error
	}
)

const (
	//AuthTypePassword is the type of ssh by password
	AuthTypePassword authType = iota
	//AuthTypePublicKey is the type of ssh by public key
	AuthTypePublicKey
)

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

var (
	//ErrEmptyCommand defines empty command error
	ErrEmptyCommand = errors.New("empty command")
	//ErrRollbackTypeAllWithNoRollbackCmd defines rollback type all with no rollback command
	ErrRollbackTypeAllWithNoRollbackCmd = errors.New("rollback type all with no rollback command")

	//ErrRunWithEmptyCommand defines run with empty command
	ErrRunWithEmptyCommand = errors.New("run with empty command")
)

//NewClientByPassword creates client by password
func NewClientByPassword(address, user, password string, rollbackTyp rollbackType) *Client {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
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

//Err will return client's error
func (c *Client) Err() error { return c.err }

func checkCommandsByRollbackType(cmds []*Command, rt rollbackType) error {
	if len(cmds) == 0 {
		return ErrEmptyCommand
	}
	switch rt {
	case RollbackTypeNone, RollbackTypeOne, RollbackTypeBackTrace:
		return nil
	case RollbackTypeAll:
		for _, cmd := range cmds {
			if cmd.rollbackCmd == "" {
				return ErrRollbackTypeAllWithNoRollbackCmd
			}
		}
		return nil
	}
	return nil
}

func sessionRunCmdWithEnv(c *Client, cmd string, env Env, output, errOutput io.Writer) error {
	if cmd == "" {
		return ErrRunWithEmptyCommand
	}
	s, err := c.session()
	if err != nil {
		return err
	}
	for k, v := range env {
		if err := s.Setenv(k, v); err != nil {
			return err
		}
	}
	s.Stdout = output
	s.Stderr = errOutput
	return s.Run(cmd)
}

//RunCmds will run command step by step
//if error occurs,return the error
//will rollback if one command occurs according to the client rollback type
func (c *Client) RunCmds(cmds []*Command) ([]*Command, error) {
	if err := checkCommandsByRollbackType(cmds, c.rollbackTyp); err != nil {
		return nil, err
	}
	for i := range cmds {
		if err := sessionRunCmdWithEnv(c, cmds[i].cmd, cmds[i].cmdEnv, &cmds[i].output, &cmds[i].errOutput); err != nil {
			cmds[i].status = CommandStatusFailed
		} else {
			cmds[i].status = CommandStatusSuccess
		}
		var checkResult bool
		if cmds[i].needCheck() {
			//TODO check
		}
		if cmds[i].needRollback() {
			//TODO rollback
		}
	}
	return cmds, nil
}
