package main

import (
	"fmt"
	"github.com/ankyra/escape-client/cmd"
	"github.com/spf13/cobra/doc"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const fmTemplate = `---
date: %s
title: "%s"
slug: %s
type: "docs"
---
`

func main() {
	filePrepender := func(filename string) string {
		now := time.Now().Format(time.RFC3339)
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "../" + strings.ToLower(base) + "/"
	}
	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "./docs/cmd", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
}