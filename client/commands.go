package client

import "time"

//Env defines command environment variables need to set
type Env struct {
	k string
	v string
}

//Command defines command for run
type Command struct {
	//cmd is shell command
	cmd    string
	cmdEnv []Env
	//rollbackCmd is the rollback command of cmd,empty string if no rollback command
	rollbackCmd    string
	rollbackCmdEnv []Env
	timeout        time.Duration
	// maybe need err chan or output chan?
}
