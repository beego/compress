package compress

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

func parseTmpl(t *template.Template, data map[string]string) (string, error) {
	buf := bytes.NewBufferString("")
	err := t.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type argString []string

func (a argString) Get(i int, args ...string) (r string) {
	if i >= 0 && i < len(a) {
		r = a[i]
	} else if len(args) > 0 {
		r = args[0]
	}
	return
}

func logError(err string, args ...interface{}) {
	err = fmt.Sprintf(err, args...)
	fmt.Fprintln(os.Stderr, err)
}

func logInfo(info string, args ...interface{}) {
	info = fmt.Sprintf(info, args...)
	fmt.Fprintln(os.Stdout, info)
}

func errHtml(err string, args ...interface{}) string {
	err = fmt.Sprintf("Beego Compress: "+err, args...)
	fmt.Fprintln(os.Stderr, err)
	return "<!-- " + err + " -->"
}
