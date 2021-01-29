package main

import (
	"flag"
	"github.com/go-xmlfmt/xmlfmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var duplicateLineBreakRep = regexp.MustCompile(`(\s*?\r?\n)+`)
var firstLineLineBreakRep = regexp.MustCompile(`(^\r?\n)+`)
var endLineLineBreakRep = regexp.MustCompile(`(\r?\n$)+`)

func main() {
	flag.Parse()
	paths := flag.Args()

	for _, rootPath := range paths {
		fileInfos, err := ioutil.ReadDir(rootPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, fileInfo := range fileInfos {
			format(rootPath, fileInfo)
		}
	}
}

func format(parentDirPath string, fileInfo os.FileInfo) {
	if fileInfo.IsDir() {
		dirPath := filepath.Join(parentDirPath, fileInfo.Name())
		fileInfos, err := ioutil.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, fileInfo := range fileInfos {
			format(dirPath, fileInfo)
		}
	} else if strings.HasSuffix(fileInfo.Name(), ".xml") {
		targetFilePath := filepath.Join(parentDirPath, fileInfo.Name())
		content, err := ioutil.ReadFile(targetFilePath)
		if err != nil {
			log.Fatal(err)
		}

		formattedXML := xmlfmt.FormatXML(string(content), "", "    ")
		formattedXML = duplicateLineBreakRep.ReplaceAllString(formattedXML, "\r\n")
		formattedXML = firstLineLineBreakRep.ReplaceAllString(formattedXML, "")
		formattedXML = endLineLineBreakRep.ReplaceAllString(formattedXML, "")
		if err = ioutil.WriteFile(targetFilePath, []byte(formattedXML), fileInfo.Mode().Perm()); err != nil {
			log.Fatal(err)
		}
	}
}
