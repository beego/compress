package compress

import (
	"html/template"
)

const (
	Js  = "js"
	Css = "css"
)

type Group struct {
	DistFile    string
	SourceFiles []string
	SkipFiles   []string
}

type compress struct {
	Type       string
	StaticURL  string
	SrcPath    string
	DistPath   string
	SrcURL     string
	DistURL    string
	Groups     map[string]Group
	IsProdMode bool
	FilterList []string
	filters    []Filter
	caches     map[string]template.HTML
}

func (c *compress) setType(t string) {
	if t != Js && t != Css {
		logError("Beego Compress: Invalid compress type " + t)
		return
	}
	c.Type = t
}

func (c *compress) SetProdMode(isProd bool) {
	c.IsProdMode = isProd
}

func (c *compress) SetStaticURL(url string) {
	c.StaticURL = url
}

func (c *compress) loadFilters() {
	for _, v := range c.FilterList {
		if f, exists := Filters[v]; exists {
			c.filters = append(c.filters, f)
		}
	}
	if len(c.filters) == 0 {
		switch c.Type {
		case Js:
			c.filters = DefaultJsFilters
		case Css:
			c.filters = DefaultCssFilters
		}
	}
}

func (c *compress) compressFiles(force, skip, verbose bool) {
	compressFiles(c, force, skip, verbose)
}

func (c *compress) Compress(name string) template.HTML {
	var templ *template.Template
	switch c.Type {
	case Js:
		templ = JsTagTemplate
	case Css:
		templ = CssTagTemplate
	}
	return generateHTML(name, c, templ)
}
