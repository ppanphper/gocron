package models

import "time"

type ProcessWorker struct {
	Id          int64     `json:"id" xorm:"int pk autoincr"`
	HostId      int       `json:"host_id" xorm:"int host_id"`
	ProcessId   int       `json:"process_id" xorm:"int process_id"`
	Pid         int       `json:"pid" xorm:"int pid"`
	IsValid     int8      `json:"is_valid" xorm:"tinyint notnull default 0"`
	State       Status    `json:"state" xorm:"tinyint notnull default 0"`
	CreatedAt   time.Time `json:"created_at" xorm:"datetime notnull default current_timestamp created"`
	StartAt     time.Time `json:"start_at" xorm:"datetime notnull default current_timestamp"`
	LastCheckAt time.Time `json:"last_check_at" xorm:"datetime notnull default current_timestamp"`
	BaseModel   `json:"-" xorm:"-"`
}

func (pw *ProcessWorker) Create() (err error) {
	_, err = Db.Insert(pw)
	return err
}

func (pw *ProcessWorker) GetByProcess(process Process) ([]ProcessWorker, error) {
	var workers []ProcessWorker
	err := Db.Where("process_id = ? AND is_valid = ?", process.Id, 1).Find(&workers)
	return workers, err
}

func (pw *ProcessWorker) GetLimitByProcess(process Process, limit int) ([]ProcessWorker, error) {
	var workers []ProcessWorker
	err := Db.Where("process_id = ? AND is_valid = ?", process.Id, 1).Limit(limit).Find(&workers)
	return workers, err
}

func (pw *ProcessWorker) Update() error {
	_, err := Db.Where("id = ?", pw.Id).Cols(`host_id,process_id,pid,state,start_at,last_check_at,is_valid`).Update(pw)
	return err
}

func (pw *ProcessWorker) SetState(state Status) error {
	_, err := Db.Where("id = ?", pw.Id).Update(ProcessWorker{State: state})
	return err
}
