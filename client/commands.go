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
	//CommandStatusCheckFailed means check cmd runs failed
	CommandStatusCheckFailed
	//CommandStatusCheckSuccess means check cmd runs successfully
	CommandStatusCheckSuccess
)

var statusMap = map[CommandStatus]string{
	CommandStatusInit:            "init",
	CommandStatusSuccess:         "success",
	CommandStatusFailed:          "failed",
	CommandStatusRollbackSuccess: "rollback success",
	CommandStatusRollbackFailed:  "rollback failed",
	CommandStatusCheckFailed:     "check failed",
}

//Command defines command for run
type Command struct {
	name string
	//cmd is shell command
	cmd    string
	mu     sync.Mutex
	cmdEnv Env
	//rollbackCmd is the rollback command of cmd,empty string if no rollback command
	checkCmd          string
	checkCmdEnv       Env
	checkExpectResult string
	rollbackCmd       string
	rollbackCmdEnv    Env
	timeout           time.Duration
	status            CommandStatus
	output            bytes.Buffer
	errOutput         bytes.Buffer
	checkOutput       bytes.Buffer
	checkErrOutput    bytes.Buffer
	rollbackOutput    bytes.Buffer
	rollbackErrOutput bytes.Buffer
}

//NewCommand will get a new command by input info
func NewCommand(name, cmd, checkCmd, rollbackCmd string, checkExpectResult string, cmdEnv, checkEnv, rollbackEnv map[string]string, timeout time.Duration) *Command {
	return &Command{
		name:              name,
		cmd:               cmd,
		cmdEnv:            Env(cmdEnv),
		checkCmd:          checkCmd,
		checkCmdEnv:       Env(checkEnv),
		checkExpectResult: checkExpectResult,
		rollbackCmd:       rollbackCmd,
		rollbackCmdEnv:    Env(rollbackEnv),
		timeout:           timeout,
		status:            CommandStatusInit,
	}
}

//Name will get Command's name
func (c *Command) Name() string { return c.name }

//Status will get Command's status
func (c *Command) Status() CommandStatus { return c.status }

//StatusName will get Command's status in string
func (c *Command) StatusName() string { return statusMap[c.status] }

//Output will get command run result if cmd runs success
func (c *Command) Output() string { return c.output.String() }

//ErrOutput will get command err info if cmd runs failed
func (c *Command) ErrOutput() string { return c.errOutput.String() }

//CheckOutput will get check command output
func (c *Command) CheckOutput() string { return c.checkOutput.String() }

//CheckErrOutput will get check command err info if check command runs failed
func (c *Command) CheckErrOutput() string { return c.checkErrOutput.String() }

//RollbackOutput will get rollback command output
func (c *Command) RollbackOutput() string { return c.rollbackOutput.String() }

//RollbackErrOutput will get rollback commanderr info if rollback command runs failed
func (c *Command) RollbackErrOutput() string { return c.rollbackErrOutput.String() }

// only command's run command runs successfully and checkCmd is not empty need check
func (c *Command) needCheck() bool {
	return c.Status() == CommandStatusSuccess && c.checkCmd != ""
}

func (c *Command) needRollback(checkResult bool) bool {
	if c.rollbackCmd == "" {
		return false
	}
	return !checkResult
}

//AddCmdEnv is used for add cmd environment variables
func (c *Command) AddCmdEnv(env map[string]string) {
	c.mu.Lock()
	for k, v := range env {
		c.cmdEnv[k] = v
	}
	c.mu.Unlock()
}

//AddCheckEnv is used for add check cmd environment variables
func (c *Command) AddCheckEnv(env map[string]string) {
	c.mu.Lock()
	for k, v := range env {
		c.checkCmdEnv[k] = v
	}
	c.mu.Unlock()
}
