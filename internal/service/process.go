package service

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/rpc/grpcpool"
	rpc "github.com/ouqiang/gocron/internal/modules/rpc/proto"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Process struct{}

var ProcessService Process

// Initialize 从数据库中取出所有开启的进程,定时检测进程是否在启动中
func (p Process) Initialize() {
	ticker := time.NewTicker(30 * time.Second)
	go func(t *time.Ticker) {
		for {
			<-t.C
			//检测所有开始状态的进程是否在运行中
			var processes []models.Process
			_ = models.Db.Where("status = ? AND enable = ?", models.Running, 1).Find(&processes)
			for _, process := range processes {
				go p.CheckProcessIsStarted(process)
			}

			//todo 检测所有停止状态的进程是否有在运行中的
		}
	}(ticker)
}

func (p Process) CheckProcessIsStarted(process models.Process) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic#service/process.go:CheckProcessIsStarted#", err)
		}
	}()
	var workers []models.ProcessWorker
	_ = models.Db.Where("process_id = ? AND is_valid = ?", process.Id, 1).Find(&workers)
	ph := models.ProcessHost{}
	hosts := ph.GetByProcess(process)
	if len(hosts) == 0 {
		projectHost := models.ProjectHost{}
		hosts, _ = projectHost.GetHostsByProjectId(process.ProjectId)
	}

	if len(hosts) == 0 {
		return
	}
	index := 0
	for _, worker := range workers {
		host := hosts[index]
		// 确认worker的执行主机
		if worker.HostId == 0 || worker.State == models.Stopped {
			worker.HostId = host.Id
			_ = worker.Update()
			index++
			if index == len(hosts) {
				index = 0
			}
		}

		// 检测worker是否在运行中
		if worker.Pid == 0 {
			worker.LastCheckAt = time.Now()
			workerStart(&worker)
			_ = worker.Update()
		} else {
			p.CheckWorkerIsRunning(worker)
		}
	}
}

func (p Process) CheckWorkerIsRunning(worker models.ProcessWorker) {
	state, err := getWorkerState(worker)
	worker.LastCheckAt = time.Now()
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			//服务不可用
			logger.Info("服务不可用")
			worker.State = models.Unknown
		}
		_ = worker.Update()
		return
	}

	if state != utils.Running {
		workerStart(&worker)
	}
	_ = worker.Update()
}

func getWorkerState(worker models.ProcessWorker) (string, error) {
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		return "", err
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, _ := grpcpool.Pool.GetClient(addr)
	resp, err := client.WorkerStateCheck(context.Background(), &rpc.StateRequest{
		Pid: int64(worker.Pid),
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			grpcpool.Pool.Release(addr) // 链接不可用,释放链接
		}
		return "", err
	}
	return resp.State, nil
}

func workerStart(worker *models.ProcessWorker) {
	logger.Debug("workerStart running")
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		logger.Debug("get worker fail", err)
		return
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, err := grpcpool.Pool.GetClient(addr)

	process := models.Process{}
	_ = process.Get(worker.ProcessId)

	req := rpc.StartRequest{
		Command: process.Command,
		LogFile: process.LogFile,
	}
	resp, err := client.StartWorker(context.Background(), &req)
	worker.State = models.Running
	worker.Pid = int(resp.Pid)
}

func (p Process) StopWorker(worker models.ProcessWorker) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("panic#service/process.go:StopWorker#", err)
		}
	}()
	host := models.Host{}
	err := host.Find(worker.HostId)
	if err != nil {
		return
	}
	addr := fmt.Sprintf("%s:%d", host.Name, host.Port)
	client, err := grpcpool.Pool.GetClient(addr)
	req := rpc.StopRequest{
		Pid: int64(worker.Pid),
	}
	resp, _ := client.StopWorker(context.Background(), &req)
	if resp.Code == "success" {
		_ = worker.SetState(models.Stopped)
	}
}

func (p Process) StopProcess(process models.Process) error {
	pw := models.ProcessWorker{}
	workers, _ := pw.GetByProcess(process)
	for _, worker := range workers {
		p.StopWorker(worker)
	}
	return nil
}
