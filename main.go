package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/awesome-gocui/gocui"
	"os"
	"os/exec"
	"time"
)

const (
	prog        = "s3mon"
	version     = "1.2"
	author      = "stuchl4n3k"
	description = "s3mon - simple service status monitor v" + version + " by " + author
)

func main() {
	parser := argparse.NewParser(prog, description)
	printVersion := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "Print version and exit"})
	beVerbose := parser.Flag("V", "verbose", &argparse.Options{Required: false, Help: "Be more verbose"})
	configFilePath := parser.String("c", "config", &argparse.Options{Required: false, Help: "Config file"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Fprint(os.Stderr, parser.Usage(err))
		os.Exit(1)
	}

	if *printVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *configFilePath == "" {
		LogErr("Please specify config file using --config flag. Example: %s --config example.yml", prog)
		os.Exit(1)
	}

	if info, err := os.Stat(*configFilePath); err != nil || info.Size() == 0 {
		LogErr("Given config file '%s' does not exist or is empty", *configFilePath)
		os.Exit(1)
	}

	if *beVerbose {
		LogLevel = LogLevelDebug
	} else {
		LogLevel = LogLevelInfo
	}

	config := NewConfig(*configFilePath)

	// GUI:
	gui := initGui(config)
	defer closeGui(gui)
	go updateItems(config.Items, gui)
	runMainLoop(gui)
}

func updateItems(items []*Item, gui *gocui.Gui) {
	for {
		for _, item := range items {
			updateItem(item, gui)
			time.Sleep(1 * time.Second)
		}
	}
}

func updateItem(item *Item, gui *gocui.Gui) {
	result := monitorItem(item)
	updateItemView(item, result, gui)
}

func monitorItem(item *Item) string {
	output, err := runScript(item.Script, item.Label)
	if err != nil {
		return ResultErr
	}

	return string(output)
}

func runScript(script, label string) ([]byte, error) {
	LogDebug("Running script for %s: %s", label, script)
	cmd := exec.Command("sh", "-c", script)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("label=%v", label))
	return cmd.Output()
}
