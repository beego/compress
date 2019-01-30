package compress

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

var (
	TmpPath           = "tmp"
	JsTagTemplate, _  = template.New("").Parse(`<script type="text/javascript" src="{{.URL}}"></script>`)
	CssTagTemplate, _ = template.New("").Parse(`<link rel="stylesheet" href="{{.URL}}" />`)
)

type Compresser interface {
	setType(t string)
	Compress(name string) template.HTML
	SetProdMode(isProd bool)
	SetStaticURL(url string)
	loadFilters()
	compressFiles(force, skip, verbose bool)
}

type Settings struct {
	Js  Compresser
	Css Compresser
}

func (s *Settings) RunCompress(force, skip, verbose bool) {
	s.Js.compressFiles(force, skip, verbose)
	s.Css.compressFiles(force, skip, verbose)
}

func LoadJsonConf(filePath string, prodMode bool, staticURL string) (setting *Settings, err error) {
	type Conf struct {
		Js  *compress
		Css *compress
	}

	var data []byte
	if file, err := os.Open(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Beego Compress: Conf Load %s", err.Error())
	} else {
		data, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("Beego Compress: Conf Read %s", err.Error())
		}
	}

	conf := Conf{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return nil, fmt.Errorf("Beego Compress: Conf Parse %s", err.Error())
	}

	setting = new(Settings)
	if conf.Js != nil {
		setting.Js = conf.Js
	} else {
		setting.Js = new(compress)
	}
	setting.Js.setType(Js)

	if conf.Css != nil {
		setting.Css = conf.Css
	} else {
		setting.Css = new(compress)
	}
	setting.Css.setType(Css)

	if staticURL == "" {
		staticURL = "/"
	}

	setting.Js.SetProdMode(prodMode)
	setting.Css.SetProdMode(prodMode)

	setting.Js.SetStaticURL(staticURL)
	setting.Css.SetStaticURL(staticURL)

	setting.Js.loadFilters()
	setting.Css.loadFilters()

	return setting, nil
}
