package compress

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"html/template"
	"path"
)

func compressFiles(c *compress, force, skip, verbose bool) {
	os.MkdirAll(TmpPath, 0755)

	for name, group := range c.Groups {

		hasError := false
		hasModified := false
		sources := make([]string, 0, len(group.SourceFiles))

		if verbose {
			logInfo("Group '%s'", name)
			logInfo("--------------------------")
		}

		skips := make(map[string]bool, len(group.SkipFiles))
		for _, file := range group.SkipFiles {
			skips[file] = true
		}

		for _, file := range group.SourceFiles {

			modified := false

			var cacheTime *time.Time
			cacheFile := filepath.Join(TmpPath, c.SrcPath, file)
			if info, err := os.Stat(cacheFile); err == nil {
				// get cached file modtime
				t := info.ModTime()
				cacheTime = &t
			}

			sourceFile := filepath.Join(c.SrcPath, file)
			if info, err := os.Stat(sourceFile); err == nil {
				if cacheTime != nil {
					if info.ModTime().Unix() > cacheTime.Unix() {
						// file modified
						modified = true
					}
				} else {
					modified = true
				}
			} else {
				logError("source file %s load error: %s", sourceFile, err.Error())
				hasError = true
				continue
			}

			if skip || modified {
				buf := bytes.NewBufferString("")
				// load content from file
				if f, err := os.Open(sourceFile); err == nil {
					buf.ReadFrom(f)
					f.Close()
				} else {
					logError("source file %s load error: %s", sourceFile, err.Error())
					hasError = true
					continue
				}

				source := buf.String()
				if verbose {
					logInfo("compress file %s ... ", sourceFile)
				}
				if skips[file] {
					if verbose {
						logInfo("skiped ")
					}
				} else {
					for _, filter := range c.filters {
						// compress content
						source = filter(source)
					}
				}
				sources = append(sources, source)

				var writeErr error
				dir, _ := filepath.Split(cacheFile)
				if writeErr = os.MkdirAll(dir, 0755); writeErr == nil {
					if f, err := os.OpenFile(cacheFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
						if _, err := f.WriteString(source); err == nil {
							hasModified = true
							if verbose {
								logInfo("saved")
							}
						} else {
							writeErr = err
						}
						f.Close()
					} else {
						writeErr = err
					}
				}

				if writeErr != nil {
					logError("write error: %s", writeErr.Error())
					hasError = true
				}

			} else {
				buf := bytes.NewBufferString("")
				// load content from file
				if f, err := os.Open(cacheFile); err == nil {
					buf.ReadFrom(f)
				} else {
					logError("cache file %s load error: %s", cacheFile, err.Error())
					hasError = true
					continue
				}

				if verbose {
					logInfo("use cache file %s", cacheFile)
				}
				sources = append(sources, buf.String())
			}
		}

		if !hasError {
			if hasModified || force {
				distFile := filepath.Join(c.DistPath, group.DistFile)
				var writeErr error
				dir, _ := filepath.Split(distFile)
				if writeErr = os.MkdirAll(dir, 0755); writeErr == nil {
					if f, err := os.OpenFile(distFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err == nil {
						if _, err := f.WriteString(strings.Join(sources, "\n\n")); err == nil {
							if verbose {
								logInfo("compressed file %s ... saved", distFile)
							}
						} else {
							writeErr = err
						}
						f.Close()

					} else {
						writeErr = err
					}
				}

				if writeErr != nil {
					logError("compressed file %s write error: %s", distFile, writeErr.Error())
					hasError = true
				}
			} else {
				if verbose {
					logInfo("not modified")
				}
			}
		}

		if verbose {
			logInfo("")
		}
	}
}

func generateHTML(name string, c *compress, t *template.Template) template.HTML {
	if group, ok := c.Groups[name]; ok {
		if c.IsProdMode {

			if c.caches == nil {
				c.caches = make(map[string]template.HTML, len(c.Groups))
			}

			if scripts, ok := c.caches[name]; ok {
				return scripts
			}

			scripts := ""

			filePath := filepath.Join(c.DistPath, group.DistFile)
			if info, err := os.Stat(filePath); err == nil {
				URL := c.StaticURL + path.Join(c.DistURL, group.DistFile) + "?ver=" + fmt.Sprint(info.ModTime().Unix())

				if res, err := parseTmpl(t, map[string]string{"URL": URL}); err != nil {
					errHtml("tempalte execute error: %s", err)

				} else {
					scripts += res
				}

			} else {
				errHtml("load file `%s` for path `%s` error: %s", group.DistFile, filePath, err.Error())
			}

			if len(scripts) > 0 {
				res := template.HTML(scripts + "\n")
				c.caches[name] = res
				return res
			}
		} else {
			scripts := make([]string, 0, len(group.SourceFiles)+2)

			scripts = append(scripts, fmt.Sprintf("<script>/* Beego Compress group `%s` begin */</script>", name))

			for _, file := range group.SourceFiles {
				filePath := filepath.Join(c.SrcPath, file)

				if info, err := os.Stat(filePath); err == nil {
					URL := c.StaticURL + path.Join(c.SrcPath, file) + "?ver=" + fmt.Sprint(info.ModTime().Unix())

					if res, err := parseTmpl(t, map[string]string{"URL": URL}); err != nil {
						scripts = append(scripts, errHtml("tempalte execute error: %s", err))

					} else {
						scripts = append(scripts, res)
					}

				} else {
					scripts = append(scripts, errHtml("load file `%s` for path `%s` error: %s", file, filePath, err.Error()))
				}
			}

			scripts = append(scripts, fmt.Sprintf("<script>/* end */</script>"))

			return template.HTML(strings.Join(scripts, "\n\t"))
		}
	} else {
		return template.HTML(errHtml("not found compress group `%s`", name))
	}

	return ""
}
