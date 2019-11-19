package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"remotectl/internal"
)

func main() {
	wait := flag.Bool("wait", false, "wait for output")
	flag.Parse()

	args := flag.Args()
	msg, _ := json.Marshal(remotectl.CmdMsg{args[0], args[1:]})

	c, err := net.Dial("unix", remotectl.GetSockAddr())
	if err != nil {
		log.Fatal("dial error: ", err)
	}

	remotectl.SendMsg(c, msg)

	if *wait {
		result := remotectl.RecvMsg(c)
		var cmdResult remotectl.CmdResult
		json.Unmarshal(result, &cmdResult)
		remotectl.PrintResult(cmdResult)
	}

	c.Close()
}
