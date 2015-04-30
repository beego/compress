package compress

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Group struct {
	DistFile    string
	SourceFiles []string
	SkipFiles   []string
}

type compress struct {
	StaticURL        string
	SrcPath          string
	DistPath         string
	SrcURL           string
	DistURL          string
	Groups           map[string]Group
	IsProductionMode bool
	caches           map[string]template.HTML
}

func (c *compress) SetProMode(isProductionMode bool) {
	c.IsProductionMode = isProductionMode
}

func (c *compress) SetStaticURL(url string) {
	c.StaticURL = url
}

func errHtml(err string, args ...interface{}) string {
	err = fmt.Sprintf("Beego Compress: "+err, args...)
	fmt.Fprintln(os.Stderr, err)
	return "<!-- " + err + " -->"
}

type compressJs struct {
	compress
}

func (c *compressJs) CompressJs(name string) template.HTML {
	return generateHTML(name, &c.compress, JsTagTemplate)
}

type compressCss struct {
	compress
}

func (c *compressCss) CompressCss(name string) template.HTML {
	return generateHTML(name, &c.compress, CssTagTemplate)
}

func generateHTML(name string, c *compress, t *template.Template) template.HTML {
	if group, ok := c.Groups[name]; ok {
		if c.IsProductionMode {

			if c.caches == nil {
				c.caches = make(map[string]template.HTML, len(c.Groups))
			}

			if scripts, ok := c.caches[name]; ok {
				return scripts
			}

			if scripts := generateHTMLForFile(c.DistPath, c.DistURL, group.DistFile, c, t); len(scripts) > 0 {
				// NOTE: Because the generated HTMl might be inside an XML comment like <!--[if IE 7]>,
				//       it cannot contain another XML comment (<!--...-->).
				//       JavaScript comments provide a good workaround for this.
				preamble := "<script>/* Powered by Beego Compress */</script>\n\t"
				res := template.HTML(preamble + scripts + "\n")
				c.caches[name] = res
				return res
			}
		} else {
			scripts := make([]string, 0, len(group.SourceFiles)+2)
			scripts = append(scripts, fmt.Sprintf("<script>/* Beego Compress group `%s` begin */</script>", name))

			for _, file := range group.SourceFiles {
				scripts = append(scripts, generateHTMLForFile(c.SrcPath, c.SrcURL, file, c, t))
			}

			scripts = append(scripts, "<script>/* end */</script>")
			return template.HTML(strings.Join(scripts, "\n\t"))
		}
	} else {
		return template.HTML(errHtml("compression group `%s` was not found", name))
	}

	return ""
}

func generateHTMLForFile(basePath string, baseUrl string, file string, c *compress, t *template.Template) string {
	filePath := filepath.Join(basePath, file)

	if info, err := os.Stat(filePath); err == nil {
		URL := c.StaticURL + path.Join(baseUrl, file) + "?ver=" + fmt.Sprint(info.ModTime().Unix())

		if res, err := parseTmpl(t, map[string]string{"URL": URL}); err != nil {
			return errHtml("template execution error: %s", err)
		} else {
			return res
		}

	} else {
		return errHtml("loading file `%s` from path `%s` failed: %s", file, filePath, err.Error())
	}
}
