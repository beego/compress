# Beego Compress

Beego Compress provides an automated system for compressing JavaScript and CSS files.

By default, it uses the [Google Closure Compiler](https://code.google.com/p/closure-compiler/wiki/BinaryDownloads) for JS, and the [YUI Compressor](https://github.com/yui/yuicompressor/releases) for CSS.
See [Customization Options](#customization-options) for how to changes these defaults.

## Sample Usage with Beego

[After creating a config file](#config-file), you can simply use this library in your beego application by following the steps below:

Move **compiler.jar** and **yuicompressor.jar** to your beego applications main directory.
This is usually the parent directory of your `static` asset directory.

BTW: This library does not depend on the main Beego framework.
Therefore, you can easily integrate it with other frameworks, or use it in a standalone command line tool (see below).

**Usage in your web application:**

```go
func SetupCompression() {
	// Load JSON config file
	isProductionMode := false
	setting, err := compress.LoadJsonConf("conf/compress.json", isProductionMode, "http://127.0.0.1/")
	if err != nil {
		beego.Error(err)
		return
	}

	// Uncomment the next line to enable usage as shown under "Command line usage":
	// setting.RunCommand()

	if isProductionMode {
		// For production mode: Use this method to automatically compress files.
		setting.RunCompress(true, false, true)
	}

	// add func to FuncMap for template use
	beego.AddFuncMap("compress_js", setting.Js.CompressJs)
	beego.AddFuncMap("compress_css", setting.Css.CompressCss)
}
```

Usage in templates:

```html
...
<head>
	...
	{{compress_css "lib"}}
	{{compress_js "lib"}}
	{{compress_js "app"}}
</head>
...
```

#### Congratulations! Let's look at the generated HTML:

Render result when isProductionMode is `false`:

```html
<!-- Beego Compress group `lib` begin -->
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/bootstrap.css?ver=1382331000" />
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/bootstrap-theme.css?ver=1382322974" />
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/font-awesome.min.css?ver=1378615042" />
<link rel="stylesheet" href="http://127.0.0.1/static_source/css/select2.css?ver=1382197742" />
<!-- end -->
<!-- Beego Compress group `lib` begin -->
<script type="text/javascript" src="http://127.0.0.1/static_source/js/jquery.min.js?ver=1378644427"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/bootstrap.js?ver=1382328826"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/lib.min.js?ver=1382328441"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/jStorage.js?ver=1382271840"></script>
<!-- end -->
<!-- Beego Compress group `app` begin -->
<script type="text/javascript" src="http://127.0.0.1/static_source/js/main.js?ver=1382195678"></script>
<script type="text/javascript" src="http://127.0.0.1/static_source/js/editor.js?ver=1382342779"></script>
<!-- end -->

```

Render result when isProductionMode is `true`:

```html
<link rel="stylesheet" href="http://127.0.0.1:8092/static/css/lib.min.css?ver=1382346563" />
<script type="text/javascript" src="http://127.0.0.1:8092/static/js/lib.min.js?ver=1382346557"></script>
<script type="text/javascript" src="http://127.0.0.1:8092/static/js/app.min.js?ver=1382346560"></script>
```

## Config file

Full configuration file example.

Note: All JSON key are not case sensitive.

**compress.json:**

```
{
	"Js": {
		// SrcPath is the path of (uncompressed) source file(s)
		"SrcPath": "static_source/js",
		// DistPath is the path of compressed file(s)
		"DistPath": "static/js",
		// SrcURL is the url prefix for the uncompressed files
		"SrcURL": "static_source/js",
		// DistURL is the url prefix for the compressed files
		"DistURL": "static/js",
		"Groups": {
			// lib is the name of this compression group
			"lib": {
				// All compressed files will be combined and saved to DistFile
				"DistFile": "lib.min.js",
				// Source files of this group
				"SourceFiles": [
					"jquery.min.js",
					"bootstrap.js",
					"lib.min.js",
					"jStorage.js"
				],
				// Files that should not be compressed
				"SkipFiles": [
					"jquery.min.js",
					"lib.min.js"
				]
			},
			"app": {
				"DistFile": "app.min.js",
				"SourceFiles": [
					"main.js",
					"editor.js"
				]
			}
		}
	},
	"Css": {
		// CSS configuration works analogous to JS configuration
		"SrcPath": "static_source/css",
		"DistPath": "static/css",
		"SrcURL": "static_source/css",
		"DistURL": "static/css",
		"Groups": {
			"lib": {
				"DistFile": "lib.min.css",
				"SourceFiles": [
					"bootstrap.css",
					"bootstrap-theme.css",
					"font-awesome.min.css",
					"select2.css"
				],
				"SkipFiles": [
					"font-awesome.min.css",
					"select2.css"
				]
			}
		}
	}
}
```

## Command line usage

When using the API `setting.RunCommand()`:

```
$ go build app.go
$ ./app compress
compress command usage:

    js     - compress all js files
    css    - compress all css files
    all    - compress all files

    Use "compress <command> -h" to get
    more information on a command.

$ ./app compress js -h
Usage of compress command: js:
  -force=false: force recreation of dist file
  -skip=false: force recompression of all files
  -v=false: verbose logging on/off

$ ./app compress css -h
Usage of compress command: css:
  -force=false: force recreation of dist file
  -skip=false: force recompression of all files
  -v=false: verbose logging on/off
```

```
use -force to recreate the dist file (even if there are no changes to it)
use -skip to force recompression of all files (even if they have not changed)
```

Example application:

```go
package main

import (
	"github.com/beego/compress"
)

func main() {
	// Load JSON config file
	isProductionMode := false
	setting, err := compress.LoadJsonConf("conf/compress.json", isProductionMode, "http://127.0.0.1/")
	if err != nil {
		panic(err)
	}

	// Use this method to start the command when called from the shell.
	// The arguments are read from os.Args.
	setting.RunCommand()
}

```

## Customization options

The whole API can be viewed in [GoWalker](http://gowalker.org/github.com/beego/compress).

* [TmpPath](http://gowalker.org/github.com/beego/compress#_variables) is the default path for caching generated files.

* [JsFilters / CssFilters](http://gowalker.org/github.com/beego/compress#_variables) contains the slice of filter functions that are used to compress CSS and JS files, respectively.
  The filters are applied in the same order as they are contained in the slice.
  Each filter gets the output of the previous filter as its input.

* [JsTagTemplate / CssTagTemplate](http://gowalker.org/github.com/beego/compress#_variables) is the [html/template.Template](http://golang.org/pkg/html/template/#Template) used to output `<script>` and `<link>` tags.

##  Contact and Issue Tracking

All beego projects need your support.

Any suggestions are welcome, please [add a new issue](https://github.com/beego/compress/issues/new) to let me know.

## LICENSE

beego compress is licensed under the Apache Licence, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0.html).
