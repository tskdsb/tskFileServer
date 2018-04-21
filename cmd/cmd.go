package cmd

import (
  "os/exec"
  "golang.org/x/crypto/ssh"
  "strings"
  "syscall"
)

type CmdObject struct {
  Cmd  string
  Args []string
}

type SshCmdObject struct {
  CmdObject
  User     string
  Password string
  IP       string
  Port     string
}

type ReturnInfo struct {
  Code    int
  Message string
}

func RunLocal(c CmdObject) (r ReturnInfo) {
  cmd := exec.Command(c.Cmd, c.Args...)
  combinedOutput, err := cmd.CombinedOutput()
  if cmd.ProcessState == nil {
    r.Code = -1
    r.Message = err.Error()
  } else {
    // r.Code = tool.Field(cmd.ProcessState, "status").(int)
    if exitCode, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
      r.Code = exitCode.ExitStatus()
    } else {
      r.Code = -2
    }
    r.Message = string(combinedOutput)
  }

  return
}

func RunSsh(s SshCmdObject) (r ReturnInfo) {
  config := &ssh.ClientConfig{
    User: s.User,
    Auth: []ssh.AuthMethod{
      ssh.Password(s.Password),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    Timeout:         0,
  }

  client, err := ssh.Dial("tcp", s.IP+":"+s.Port, config)
  if err != nil {
    return
  }
  defer client.Close()

  // Each ClientConn can support multiple interactive sessions,
  // represented by a Session.
  session, err := client.NewSession()
  if err != nil {
    return
  }
  defer session.Close()

  // Once a Session is created, you can execute a single command on
  // the remote side using the Run method.
  argsString := strings.Join(s.Args, " ")
  cmdArgs := s.Cmd + " " + argsString
  result, err := session.CombinedOutput(cmdArgs)
  if err != nil {
    if waitMsg, ok := err.(*ssh.ExitError); ok {
      r.Code = waitMsg.ExitStatus()
    } else {
      r.Code = -1
    }
  }
  r.Message = string(result)

  return
}
