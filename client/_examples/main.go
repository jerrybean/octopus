package main

import (
	"fmt"
	"time"

	"github.com/jerrybean/octopus/client"
)

func main() {
	c := client.NewClientByPassword("127.0.0.1:22", "root", "yf@123", client.RollbackTypeNone)
	if c.Err() != nil {
		panic(c.Err())
	}
	cmds := []*client.Command{client.NewCommand("ls root file", "ls /root", "", "", "", nil, nil, nil, time.Second)}
	newCmds, err := c.RunCmds(cmds)
	fmt.Println(newCmds)
	fmt.Println(err)
	fmt.Println(newCmds[0].StatusName())
	fmt.Println(newCmds[0].Output())
	fmt.Println(newCmds[0].ErrOutput())
}
