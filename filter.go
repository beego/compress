package compress

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
)

type Filter func(source string) string

var (
	Filters = map[string]Filter{
		"JsFilter":      JsFilter,
		"CssFilter":     CssFilter,
		"ClosureFilter": ClosureFilter,
		"YuiFilter":     YuiFilter,
	}

	DefaultJsFilters  = []Filter{JsFilter}
	DefaultCssFilters = []Filter{CssFilter}
	closureBin        = "java -jar compiler.jar"
	closureArgs       = map[string]string{
		"compilation_level": "SIMPLE_OPTIMIZATIONS",
		"warning_level":     "QUIET",
	}

	yuiBin  = "java -jar yuicompressor.jar"
	yuiArgs = map[string]string{
		"type": Css,
	}
	minifier *minify.M
)

func JsFilter(source string) string {
	return runMinifier("text/javascript", source)
}

func CssFilter(source string) string {
	return runMinifier("text/css", source)
}

func runMinifier(mediaType, source string) string {
	if minifier == nil {
		minifier = minify.New()
		minifier.AddFunc("text/css", css.Minify)
		minifier.AddFunc("text/javascript", js.Minify)
	}
	s, err := minifier.String(mediaType, source)
	if err != nil {
		logError(err.Error())
		return source
	}
	return s
}

func ClosureFilter(source string) string {
	args := strings.Fields(closureBin)
	for arg, value := range closureArgs {
		args = append(args, "--"+arg)
		args = append(args, value)
	}
	return runFilter(args[0], args[1:], source)
}

func YuiFilter(source string) string {
	args := strings.Fields(yuiBin)
	for arg, value := range yuiArgs {
		args = append(args, "--"+arg)
		args = append(args, value)
	}
	return runFilter(args[0], args[1:], source)
}

func runFilter(bin string, args []string, source string) string {
	buf := bytes.NewBufferString(source)
	out := bytes.NewBufferString("")

	cmd := exec.Command(bin, args...)
	cmd.Stdin = buf
	cmd.Stderr = os.Stderr
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		logError(err.Error())
		return source
	} else {
		return out.String()
	}
}
