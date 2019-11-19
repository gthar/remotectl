package main

import (
	"encoding/json"
	"log"
	"net"
	"remotectl/internal"
)

func cmdServer(c net.Conn) {
	msg := remotectl.RecvMsg(c)
	var cmdMsg remotectl.CmdMsg
	json.Unmarshal(msg, &cmdMsg)
	result, err := remotectl.RunCmd(cmdMsg)
	if err != nil {
		log.Print(err)
	}
	resultMsg, _ := json.Marshal(result)

	if err := remotectl.SendMsg(c, resultMsg); err != nil {
		remotectl.PrintResult(result)
	}
}

func main() {
	addr := remotectl.GetSockAddr()
	l := remotectl.OpenSock(addr)
	defer l.Close()
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("Accept error: ", err)
		}
		go cmdServer(fd)

	}
}
