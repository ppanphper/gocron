//go:build !windows
// +build !windows

package utils

import (
	"errors"
	"fmt"
	"github.com/ouqiang/gocron/internal/modules/logger"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

type Result struct {
	output string
	err    error
}

// ExecShell 执行shell命令，可设置执行超时时间
func ExecShell(ctx context.Context, command string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	resultChan := make(chan Result)
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- Result{string(output), err}
	}()
	select {
	case <-ctx.Done():
		if cmd.Process.Pid > 0 {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return "", errors.New("timeout killed")
	case result := <-resultChan:
		return result.output, result.err
	}
}

// StartWorker 启动一个worker进程
func StartWorker(_ context.Context, req *rpc.StartRequest) (int, error) {
	cmd := exec.Command("/bin/bash", "-c", req.Command)

	//		ParentProcess: 1,
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	if req.LogFile != "" {
		logFile := req.LogFile
		d, err := filepath.Abs(path.Dir(logFile))
		if err != nil {
			return 0, err
		}
		_, err = os.Stat(d)
		if err != nil || os.IsNotExist(err) {
			_ = os.MkdirAll(d, 0666)
		}
		stdout, err := os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(os.Getpid(), ": 打开日志文件错误:", err)
			return 0, err
		}
		cmd.Stderr = stdout
		cmd.Stdout = stdout
	}

	err := cmd.Start()
	if err != nil {
		logger.Error("进程启动失败:%s", err)
		return 0, err
	}

	pid := cmd.Process.Pid
	return pid, err
}

// StopWorker 通过 pid 停止指定进程
func StopWorker(pid int) error {
	//_ = exec.Command("kill", "-9", strconv.Itoa(int(pid))).Run()
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = process.Kill()
	if err != nil {
		return err
	}
	// 避免产生僵尸进程 https://blog.csdn.net/qq_27068845/article/details/78816995
	_, err = process.Wait()
	return err
}

func isRunning(pid int) bool {
	pathString := fmt.Sprintf("/proc/%d/status", pid)
	_, err := os.Stat(pathString)
	if err != nil {
		return !os.IsNotExist(err)
	}

	content, err := ioutil.ReadFile(pathString)
	if err != nil {
		return false
	}
	regex, _ := regexp.Compile("State:.*")
	str := regex.Find(content)
	if strings.Contains(strings.ToLower(string(str)), "zombie") {
		// 僵尸进程
		return false
	}

	return true
}

func WorkerStateCheck(pid int) (string, error) {
	if pid == 0 {
		return Error, errors.New("pid 不能为0")
	}
	cmd := exec.Command("which", "ps")
	bytes, err := cmd.CombinedOutput()
	if err == nil && strings.Contains(string(bytes), "/bin/ps") {
		// 使用PS命令校验
		cmd := exec.Command("ps", "-p", fmt.Sprintf("%d", pid))
		output, err := cmd.CombinedOutput()
		if err != nil {
			return Stop, err
		}
		rows := strings.Split(string(output), "\n")
		if len(rows) < 2 { //异常
			return Stop, errors.New(fmt.Sprintf("output exception,process %d not found", pid))
		}
		if strings.Contains(rows[1], "<defunct>") {
			//僵尸进程
			return Stop, errors.New(fmt.Sprintf("process %d is zombie", pid))
		}
		if !strings.Contains(rows[1], fmt.Sprintf("%d", pid)) {
			return Stop, errors.New(fmt.Sprintf("process %d not found", pid))
		}
		return Running, nil
	} else {
		// 通过文件校验
		if isRunning(pid) {
			return Running, nil
		} else {
			return Stop, nil
		}
	}
}
