package remotectl

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type CmdMsg struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

type CmdResult struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}

func SendMsg(c net.Conn, msg []byte) (err error) {
	msgLen := make([]byte, 2)
	binary.BigEndian.PutUint16(msgLen, uint16(len(msg)))

	if _, err := c.Write(msgLen); err != nil {
		return err
	}

	buf := make([]byte, 1)
	if _, err := c.Read(buf); err != nil {
		return err
	}
	if buf[0] != byte(6) {
		log.Fatal("ack not received")
		return
	}

	if _, err := c.Write(msg); err != nil {
		return err
	}

	return nil
}

func PrintResult(r CmdResult) {
	os.Stdout.WriteString(r.StdOut)
	os.Stderr.WriteString(r.StdErr)
}

func RecvMsg(c net.Conn) []byte {
	buf := make([]byte, 2)
	if _, err := c.Read(buf); err != nil {
		log.Fatal(err)
	}
	msgLen := binary.BigEndian.Uint16(buf)

	if _, err := c.Write([]byte{6}); err != nil {
		log.Fatal(err)
	}

	msg := make([]byte, msgLen)
	if _, err := c.Read(msg); err != nil {
		log.Fatal(err)
	}
	return msg
}

func RunCmd(cmdMsg CmdMsg) (x CmdResult, err error) {
	log.Print("running: ", cmdMsg)
	cmd := exec.Command(cmdMsg.Cmd, cmdMsg.Args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	stdoutContent, _ := ioutil.ReadAll(stdoutPipe)
	stderrContent, _ := ioutil.ReadAll(stderrPipe)

	if err = cmd.Wait(); err != nil {
		return
	}

	result := CmdResult{
		StdOut: string(stdoutContent),
		StdErr: string(stderrContent),
	}
	return result, nil
}

func GetSockAddr() string {
	if sockVar := os.Getenv("REMOTE_CTL_SOCKET"); sockVar != "" {
		return sockVar
	}

	currentUser, _ := user.Current()
	uid := currentUser.Uid

	displayVar := os.Getenv("DISPLAY")
	display := strings.Replace(strings.Replace(displayVar, ":", "", 1), "unix", "", 1)

	addr := fmt.Sprintf("/run/user/%s/remotectl/display%s.sock", uid, display)

	return addr
}

func OpenSock(addr string) net.Listener {
	log.Printf("using socket on %s", addr)

	sockDir := filepath.Dir(addr)
	var perms os.FileMode = 0700
	if err := os.MkdirAll(sockDir, perms); err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod(sockDir, perms); err != nil {
		log.Fatal(err)
	}

	if err := os.RemoveAll(addr); err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Chmod(addr, perms); err != nil {
		log.Fatal(err)
	}

	return l
}
