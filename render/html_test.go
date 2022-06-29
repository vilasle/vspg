package render

import (
	"io/ioutil"
	"testing"
)

func TestRenderReport(t *testing.T) {
	report := BackupReport{}
	report.Status = 0
	report.ErrorCount = 0
	report.WarningCount = 0
	report.SuccessCount = 1
	
	item := Item{}
	item.Name = "ut"
	item.BackupPath = "somePath"
	item.BackupSize = "123GB"
	item.DatabaseSize = "454GB"
	item.StartTime = "23:30"
	item.FinishTime = "23:50"
	item.OID = 23456
	item.Status = StatusSuccess
	item.Details = ""

	item1 := Item{}
	item1.Name = "ut"
	item1.BackupPath = "somePath"
	item1.BackupSize = "123GB"
	item1.DatabaseSize = "454GB"
	item1.StartTime = "23:30"
	item1.FinishTime = "23:50"
	item1.OID = 123456
	item1.Status = StatusError
	item1.Details = ""

	task := Task{}
	task.Name = "Daily"
	task.Items = append(task.Items, item)
	task.Items = append(task.Items, item1)
	report.Tasks = append(report.Tasks, task)

	if content, err := RenderReport(report); err != nil {
		t.Fatal(err)
	} else {
		ioutil.WriteFile("out.html", content, 0777)
	}

}
