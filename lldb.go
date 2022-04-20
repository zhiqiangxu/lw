package lw

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

type LLDB struct {
	config Config
	cmd    *exec.Cmd
	stdinW io.WriteCloser
}

type Config struct {
	Stdout io.Writer
	Stderr io.Writer
}

func AttachPID(pid int, config Config) (l *LLDB, err error) {
	cmd := exec.Command("lldb", "-p", strconv.Itoa(pid))
	stdinW, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr

	l = &LLDB{config: config, cmd: cmd, stdinW: stdinW}
	return
}

func AttachName(name string, config Config) (l *LLDB, err error) {
	cmd := exec.Command("lldb", "-n", name)
	stdinW, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	cmd.Stdout = config.Stdout
	cmd.Stderr = config.Stderr
	l = &LLDB{config: config, cmd: cmd, stdinW: stdinW}
	return
}

func (l *LLDB) Start() (err error) {
	return l.cmd.Start()
}

func (l *LLDB) Stop() (err error) {
	return l.cmd.Wait()
}

// with Expr, you can do various things including memory update:
// 		expr -- *(char*)0x2ae4 = 0x12
func (l *LLDB) Expr(expr string) {
	io.WriteString(l.stdinW, fmt.Sprintf("expr -- %s\n", expr))
}

// find a byte pattern in memory
func (l *LLDB) MemoryFind(search, eval string) {
	switch {
	case search != "":
		io.WriteString(l.stdinW, fmt.Sprintf("memory find -s %s\n", search))
	case eval != "":
		io.WriteString(l.stdinW, fmt.Sprintf("memory find -e %s\n", eval))
	}
}
