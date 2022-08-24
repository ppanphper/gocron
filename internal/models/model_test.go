package models

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/modules/app"
	"github.com/ouqiang/gocron/internal/modules/logger"
	"github.com/ouqiang/gocron/internal/modules/setting"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"testing"
)

func init() {
	fmt.Println("setup")
	app.InitEnv("1.5")
	config, err := setting.Read(app.AppConfig)
	if err != nil {
		logger.Fatal("读取应用配置失败", err)
	}
	app.Setting = config

	// 初始化DB
	Db = CreateDb()
}

func TestCreateTable(t *testing.T) {
	err := Db.CreateTables(Project{})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAlertTable(t *testing.T) {
	results, err := Db.Query("alter table `task` add project_id int default 0 not null;")
	t.Log(results, err)
}

func TestProjectSetHosts(t *testing.T) {
	projects := make([]Project, 0)
	_ = Db.Find(&projects)
	p := Project{}
	projects = p.setHostsForProjects(projects)
	t.Log(projects)
}

func TestProcessWorker_GetByProcess(t *testing.T) {
	pw := ProcessWorker{}
	workers, err := pw.GetByProcess(Process{Id: 1})
	t.Log(err)
	t.Log(workers)
}

func TestProcessHost_GetByProcess(t *testing.T) {
	ph := ProcessHost{}
	hosts := ph.GetByProcess(Process{Id: 1})
	//t.Log(err)
	t.Log(hosts)
}

func TestProcessHost_DeleteForProcess(t *testing.T) {
	ph := ProcessHost{}
	ph.DeleteForProcess(Process{Id: 1})
}

func TestGetTasks(t *testing.T) {
	task := Task{}
	tasks, _ := task.List(CommonMap{})
	for _, task := range tasks {
		t.Log(task.Status)
	}
	t.Log(utils.JsonResp.Success("", tasks))
}

func TestTaskHost_GetHostsByTaskId(t *testing.T) {
	th := TaskHost{}
	hosts, err := th.GetHostsByTaskId(1)
	t.Log(hosts, err)
}

func TestGetActionUsers(t *testing.T) {
	_, err := Db.Query(fmt.Sprintf("alter table `%suser` add source varchar(16) default 'system' not null after email;", TablePrefix))
	t.Log(err)

}

func TestProjectHost_GetHostsByProjectId(t *testing.T) {
	ph := ProjectHost{}
	t.Log(ph.GetHostsByProjectId(1))
}

func TestTaskLog_Clear(t *testing.T) {
	l := TaskLog{}

	t.Log(l.Clear(CommonMap{"taskId": "1"}))
}

type ChartData struct {
	Name string
	Type string
	Data []int
}

func TestDashboard(t *testing.T) {
	times := utils.GetMondayTimes()

	for _, time := range times {
		year, week := time.ISOWeek()
		t.Log(fmt.Sprintf("%d-%d", year, week))
	}
	p := Process{}
	charts := p.GetChartDataForDashboard(times[0])

	t.Log(charts)

	var maps = make(map[int]map[string]int)
	for _, chart := range charts {
		_, ok := maps[chart.ProjectId]
		if !ok {
			maps[chart.ProjectId] = make(map[string]int)
		}
		maps[chart.ProjectId][chart.Week] = chart.Count
	}

	for projectId, m := range maps {
		var data = ChartData{Name: fmt.Sprintf("%d 新增进程", projectId), Type: "line", Data: make([]int, len(times))}
		for index, d := range times {
			year, week := d.ISOWeek()
			s := fmt.Sprintf("%d-%d", year, week)
			c, ok := m[s]
			if ok {
				data.Data[index] = c
			} else {
				data.Data[index] = 0
			}
		}
		t.Log(data)
	}

	/**
	{
	            name: 'ims新增进程',
	            type: 'line',
	            data: [2200, 182, 191, 234, 290, 330, 310]
	          }
	*/

}
