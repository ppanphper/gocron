package dashboard

import (
	"fmt"
	"github.com/ouqiang/gocron/internal/models"
	"github.com/ouqiang/gocron/internal/modules/utils"
	"gopkg.in/macaron.v1"
	"time"
)

type NewData struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data []int  `json:"data"`
}

type Dashboard struct {
	TotalGroup       map[string]int64          `json:"totalGroup"`
	ActiveUsers      []models.ActiveUser       `json:"activeUsers"`
	ProjectTasks     []models.ProjectTaskChart `json:"projectTasks"`
	ProjectNewX      []string                  `json:"projectNewX"`
	ProjectNewCharts []NewData                 `json:"projectNewCharts"`
}

func Index(ctx *macaron.Context) string {
	data := Dashboard{TotalGroup: map[string]int64{}}

	t := models.Task{}
	taskCount, _ := t.Total(models.CommonMap{})
	data.TotalGroup["taskCount"] = taskCount

	process := models.Process{}
	processCount, _ := process.Total(models.CommonMap{})
	data.TotalGroup["processCount"] = processCount
	u := models.User{}
	userCount, _ := u.Total()
	data.TotalGroup["userCount"] = userCount

	pro := models.Project{}
	projectCount, _ := pro.Total(models.CommonMap{})
	data.TotalGroup["projectCount"] = projectCount

	log := models.OperateLog{}
	users, _ := log.GetActiveUsers()
	data.ActiveUsers = users

	project := models.Project{}
	data.ProjectTasks = project.GetProjectTasksChart()

	//折线图数据组装
	times := utils.GetMondayTimes()

	data.ProjectNewX = make([]string, len(times))
	for i, d := range times {
		year, w := d.ISOWeek()
		data.ProjectNewX[i] = fmt.Sprintf("%d年 第%d周", year, w)
	}

	// 处理新增的进程数据
	parseNewDataForChart(&data, times, process.GetChartDataForDashboard(times[0]), "新增进程")
	// 处理新增的任务数据
	parseNewDataForChart(&data, times, t.GetChartDataForDashboard(times[0]), "新增定时任务")

	return utils.JsonResp.Success(utils.SuccessContent, data)
}

func parseNewDataForChart(data *Dashboard, times []time.Time, charts []models.ChartNew, title string) {
	var maps = make(map[int]map[string]int)
	var projects = make(map[int]string)
	for _, chart := range charts {
		projects[chart.ProjectId] = chart.ProjectName
		_, ok := maps[chart.ProjectId]
		if !ok {
			maps[chart.ProjectId] = make(map[string]int)
		}
		maps[chart.ProjectId][chart.Week] = chart.Count
	}
	
	projects[0] = "项目待定"

	for projectId, m := range maps {
		line := NewData{Name: fmt.Sprintf("%s %s", projects[projectId], title), Type: "line", Data: make([]int, len(times))}

		for index, d := range times {
			year, week := d.ISOWeek()
			s := fmt.Sprintf("%d-%d", year, week)
			c, ok := m[s]
			if ok {
				line.Data[index] = c
			} else {
				line.Data[index] = 0
			}
		}
		data.ProjectNewCharts = append(data.ProjectNewCharts, line)
	}
}
