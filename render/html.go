package render

import (
	"bytes"
	"fmt"
	"html/template"
)

const (
	StatusSuccess = 0
	StatusWarning = 1
	StatusError   = 2
)

func RenderReport(report BackupReport, templatesPath string) (content []byte, err error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s\\%s", templatesPath, "report.html"))

	if err != nil {
		return content, err
	}

	w := bytes.NewBufferString("")
	if err = tmpl.Execute(w, report); err == nil {
		return w.Bytes(), err
	} else {
		return content, err
	}
}
