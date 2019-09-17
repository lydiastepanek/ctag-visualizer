package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type TagData struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Pattern string `json:"pattern"`
	Kind    string `json:"kind"`
	Scope   string `json:"scope"`
	Package string `json:"package"`
	LineNum int    `json:"lineNum"`
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "in", "", "input log file")

	flag.Parse()

	tagsFile, err := ioutil.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	raw := string(tagsFile)
	lines := strings.Split(raw, "\n")
	tags := make([]TagData, 0)
	for _, line := range lines {
		var tag TagData
		json.Unmarshal([]byte(line), &tag)
		if tag.Path == "" {
			continue
		}
		fmt.Println("checking ", tag.Path)
		// read the first two lines to extract package
		contents, err := ioutil.ReadFile(tag.Path)
		if err != nil {
			panic("err reading:" + tag.Path + ":" + err.Error())
		}
		rawCont := string(contents)
		contLines := strings.Split(rawCont, "\n")
	INNER:
		for i, contLine := range contLines {
			if i < 2 {
				if strings.Contains(contLine, "package") {
					tokens := strings.Fields(contLine)
					tag.Package = tokens[1]
				}
			}
			if tag.Pattern != "" {
				pattern := tag.Pattern[2 : len(tag.Pattern)-2]
				pattern = fmt.Sprintf("^%s$", regexp.QuoteMeta(pattern))
				matched, err := regexp.MatchString(pattern, contLine)
				if err != nil {
					panic("failed to check regex: " + pattern + ". err: " + err.Error())
				}
				if matched {
					tag.LineNum = i + 1
					break INNER
				}
			}
		}
		tags = append(tags, tag)
	}
	newJson, err := json.Marshal(tags)
	if err != nil {
		panic("failed to unmarshal. err: " + err.Error())
	}
	fmt.Println(string(newJson))

	fo, err := os.Create(inputFile)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := fo.Write(newJson); err != nil {
		panic(err)
	}

}
