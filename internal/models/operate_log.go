package models

import (
	"fmt"
	"time"
)

type OperateLog struct {
	Id         int       `json:"id" xorm:"pk autoincr notnull "`
	Username   string    `json:"username" xorm:"varchar(32) notnull"`
	Ip         string    `json:"ip" xorm:"varchar(15) not null"`
	Uri        string    `json:"uri" xorm:"varchar(128) not null"`
	HttpMethod string    `json:"http_method" xorm:"varchar(8) not null"`
	Data       string    `json:"data" xorm:"text not null"`
	Created    time.Time `json:"created" xorm:"datetime notnull created"`
	BaseModel  `json:"-" xorm:"-"`
}

type ActiveUser struct {
	Username string `json:"username"`
	Count    uint   `json:"count"`
}

var _ = make(map[string]string)

func (log *OperateLog) Create() (insertId int, err error) {
	_, err = Db.Insert(log)
	if err == nil {
		insertId = log.Id
	}

	return
}

// GetActiveUsers 获取30天内最活跃用户12名用户
func (log *OperateLog) GetActiveUsers() ([]ActiveUser, error) {
	date := time.Unix(time.Now().Unix()-3600*24*30, 0).Format("2006-01-02")
	var users = make([]ActiveUser, 0)
	err := Db.SQL(fmt.Sprintf("SELECT username, count(0) AS count FROM operate_log WHERE created > '%s' AND username != '' GROUP BY username ORDER BY count DESC LIMIT 12", date)).Find(&users)
	return users, err
}
