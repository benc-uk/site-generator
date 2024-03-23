package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	_ "embed"
)

var version string = "0.0.0"
var buildInfo string = "Local build"

//go:embed templates/default.html
var htmlTemplate string

type TemplateData struct {
	Title     string
	IndexList []IndexEntry
	Body      template.HTML
	Footer    template.HTML
	IsTop     bool
	IsIndex   bool
}

type IndexEntry struct {
	ShortName string
	FullName  string
	IsFile    bool
}

func main() {
	fmt.Printf("🧵 Simple Site Generator v%s (%s)\n\n", version, buildInfo)

	var outDir, srcDir, templateFile string

	flag.StringVar(&outDir, "o", "./html", "Output directory for generated HTML files")
	flag.StringVar(&srcDir, "s", "./src", "Source directory containing markdown files")
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

	if templateFile != "" {
		log.Println("🧩 Using custom template:", templateFile)

		if _, err := os.Stat(templateFile); err != nil {
			log.Fatal(err)
		}

		newTemplate, err := os.ReadFile(templateFile)
		if err != nil {
			log.Fatal(err)
		}

		htmlTemplate = string(newTemplate)
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
}

// Generate index.html for a directory
func createIndex(path string, outDir string, srcDir string) error {
	var outPath, contentTitle string

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

	indexList := []IndexEntry{}

	for _, f := range contents {
		fileBase := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
		fileTarget := fileBase + ".html"

		if f.IsDir() {
			indexList = append(indexList, IndexEntry{
				ShortName: fileBase,
				FullName:  fileBase + "/",
				IsFile:    false,
			})
		}

		if strings.HasSuffix(f.Name(), ".md") {
			indexList = append(indexList, IndexEntry{
				ShortName: fileBase,
				FullName:  fileTarget,
				IsFile:    true,
			})
		}
	}

	sort.Slice(indexList, func(i, j int) bool {
		return !indexList[i].IsFile
	})

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	//nolint
	outFile, err := os.OpenFile(filepath.Join(outPath, "index.html"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	log.Println("📂 Indexing:", outPath)

	return tmpl.Execute(outFile, TemplateData{
		Title:     contentTitle,
		IndexList: indexList,
		IsIndex:   true,
		IsTop:     path == srcDir,
	})
}

// Convert markdown to HTML and write to a file
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

	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return err
	}

	//nolint
	outFile, err := os.OpenFile(outFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tmpl.Execute(outFile, TemplateData{
		Title: fileBase,
		// nolint
		Body: template.HTML(htmlBody),
	})
}
