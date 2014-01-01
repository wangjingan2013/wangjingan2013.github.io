package main

import (
	"bufio"
	"flag"
	"github.com/knieriem/markdown"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
)

const htmlPrelog = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta http-equiv="content-type" content="text/html; charset=UTF-8">
    <!-- Meta, title, CSS, favicons, etc. -->
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>写给女儿的日记</title>
    <link href="../css/bootstrap.css" rel="stylesheet">
    <link href="../css/docs.css" rel="stylesheet">
  </head>

  <body data-twttr-rendered="true" class="bs-docs-home">
    <main class="bs-masthead" id="content" role="main">
      <div class="container">
`
const htmlEpilog = `
      </div>
    </main>
  </body>
</html>
`

var contentDir = flag.String("content", "./", "The directory containing markdown files")

func renderMarkdownFile(mdFile, htmlFile string) {
	out, e := os.Create(htmlFile)
	if e != nil {
		log.Fatalf("Cannot create output file: %s: %q\n", htmlFile, e)
	}
	defer out.Close()

	in, e := os.Open(mdFile)
	if e != nil {
		log.Fatalf("Cannot open input file: %s: %q\n", mdFile, e)
	}
	defer in.Close()

	bo := bufio.NewWriter(out)
	bo.WriteString(htmlPrelog)
	markdown.NewParser(nil).Markdown(in, markdown.ToHTML(bo))
	bo.WriteString(htmlEpilog)
	bo.Flush()
}

func main() {
	flag.Parse()

	files, e := ioutil.ReadDir(*contentDir)
	if e != nil {
		log.Fatalf("Cannot read directory: %s : %q\n", *contentDir, e)
	}

	md := regexp.MustCompile("([^\\.]+)\\.md$")
	for _, f := range files {
		subs := md.FindAllStringSubmatch(f.Name(), -1)
		if len(subs) == 1 && len(subs[0]) == 2 && len(subs[0][1]) > 0 {
			renderMarkdownFile(path.Join(*contentDir, f.Name()),
				path.Join(*contentDir, subs[0][1]+".html"))
		}
	}
}
