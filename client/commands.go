package client

import (
	"bytes"
	"sync"
	"time"
)

//Env defines command environment variables need to set
type Env map[string]string

//CommandStatus defines command status
type CommandStatus byte

const (
	//CommandStatusInit is the default command status
	CommandStatusInit CommandStatus = iota
	//CommandStatusSuccess means this command runs successfully
	CommandStatusSuccess
	//CommandStatusFailed means this command runs failed,and with no rollback
	CommandStatusFailed
	//CommandStatusRollbackSuccess means the command runs failed,the rollback command runs successfully
	CommandStatusRollbackSuccess
	//CommandStatusRollbackFailed means both the command and rollback command run failed
	CommandStatusRollbackFailed
)

var statusMap = map[CommandStatus]string{
	CommandStatusInit:            "init",
	CommandStatusSuccess:         "success",
	CommandStatusFailed:          "failed",
	CommandStatusRollbackSuccess: "rollback success",
	CommandStatusRollbackFailed:  "rollback failed",
}

//Command defines command for run
type Command struct {
	name string
	//cmd is shell command
	cmd    string
	mu     sync.Mutex
	cmdEnv Env
	//rollbackCmd is the rollback command of cmd,empty string if no rollback command
	detectCmd          string
	detectCmdEnv       Env
	detectExpectResult string
	rollbackCmd        string
	rollbackCmdEnv     Env
	timeout            time.Duration
	status             CommandStatus
	output             bytes.Buffer
	errOutput          bytes.Buffer
}

//NewCommand will get a new command by input info
func NewCommand(name, cmd, detectCmd, rollbackCmd string, detectExpectResult string, cmdEnv, detectEnv, rollbackEnv map[string]string, timeout time.Duration) *Command {
	return &Command{
		name:               name,
		cmd:                cmd,
		cmdEnv:             Env(cmdEnv),
		detectCmd:          detectCmd,
		detectCmdEnv:       Env(detectEnv),
		detectExpectResult: detectExpectResult,
		rollbackCmd:        rollbackCmd,
		rollbackCmdEnv:     Env(rollbackEnv),
		timeout:            timeout,
		status:             CommandStatusInit,
	}
}

//Name will get Command's name
func (c *Command) Name() string { return c.name }

//Status will get Command's status
func (c *Command) Status() CommandStatus { return c.status }

//StatusName will get Command's status in string
func (c *Command) StatusName() string { return statusMap[c.status] }

//AddCmdEnv is used for add cmd environment variables
func (c *Command) AddCmdEnv(env map[string]string) {
	c.mu.Lock()
	for k, v := range env {
		c.cmdEnv[k] = v
	}
	c.mu.Unlock()
}

//AddDetectEnv is used for add detect cmd environment variables
func (c *Command) AddDetectEnv(env map[string]string) {
	c.mu.Lock()
	for k, v := range env {
		c.detectCmdEnv[k] = v
	}
	c.mu.Unlock()
}
