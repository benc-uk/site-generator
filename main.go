package main

import (
	"errors"
	"flag"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	_ "embed"
)

//go:embed templates/standard.html
var standardTemplate string

type TemplateData struct {
	Title     string
	IndexList []string
	Body      template.HTML
	Footer    template.HTML
}

func main() {
	var outDir string
	flag.StringVar(&outDir, "o", "./html", "Output HTML and site content here")
	var srcDir string
	flag.StringVar(&srcDir, "s", "./src", "Source directory containing Markdown files")
	var templateFile string
	flag.StringVar(&templateFile, "t", "", "Optional, custom template file")
	flag.Parse()

	outDir, err := filepath.Abs(outDir)
	if err != nil {
		log.Fatal(err)
	}
	srcDir, err = filepath.Abs(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(srcDir); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("💥 Source directory %s does not exist!", srcDir)
	}

	_ = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		// Process markdown files
		if d.Type().IsRegular() && strings.HasSuffix(d.Name(), ".md") {
			err := generateHTML(path, outDir, srcDir)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// Process directories
		if d.Type().IsDir() {
			if strings.Contains(path, ".git/") {
				return nil
			}

			err := createIndex(path, outDir, srcDir)
			if err != nil {
				log.Fatalln(err)
			}
		}

		return nil
	})

	_ = copyFile(srcDir+"/style.css", outDir+"/style.css")
}

//
//
//
func createIndex(path string, outDir string, srcDir string) error {
	var outPath string
	var contentTitle string

	if path == srcDir {
		outPath = outDir
		contentTitle = "Main Index"
	} else {
		// Remove the top part of srcDir from the path to form the outPath
		sanePath := strings.Replace(path, srcDir, "", -1)
		outPath = filepath.Join(outDir, sanePath)
		contentTitle = sanePath
	}

	err := os.MkdirAll(outPath, 0755)
	if err != nil {
		return err
	}

	contents, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	indexList := []string{}
	for _, f := range contents {
		fileBase := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		fileTarget := fileBase + ".html"

		if f.IsDir() {
			indexList = append(indexList, f.Name()+"/")
		}

		if strings.HasSuffix(f.Name(), ".md") {
			indexList = append(indexList, fileTarget)
		}
	}

	tmpl, err := template.New("index").Parse(standardTemplate)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(filepath.Join(outPath, "index.html"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	log.Println("📂 Indexing:", outPath)

	return tmpl.Execute(outFile, TemplateData{
		Title:     contentTitle,
		IndexList: indexList,
	})
}

//
//
//
func generateHTML(path string, outDir string, srcDir string) error {
	var outPath string
	if path == srcDir {
		outPath = outDir
	} else {
		sanePath := strings.Replace(path, srcDir, "", -1)
		outPath = filepath.Join(outDir, sanePath)
	}

	outDir, fileName := filepath.Split(outPath)
	fileBase := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	outFileName := filepath.Join(outDir, fileBase+".html")
	log.Println("📜 Generating HTML:", outFileName)

	parser := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	md, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	htmlBody := markdown.ToHTML(md, parser, renderer)

	tmpl, err := template.New("index").Parse(standardTemplate)
	if err != nil {
		return err
	}

	outFile, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tmpl.Execute(outFile, TemplateData{
		Title: fileBase,
		Body:  template.HTML(htmlBody),
	})
}

//
//
//
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
