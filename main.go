package main

import (
	"os"

	"os/exec"
	"path/filepath"

	"fmt"

	"github.com/ngaut/log"
	"github.com/spf13/cobra"
	"github.com/sryanyuan/pithyblog/blog"
)

var (
	rootCommand = &cobra.Command{
		Use:   "pithyblog",
		Short: "pithyblog provides some commands to run the blog site",
	}
	setupCommand = &cobra.Command{
		Use:   "setup",
		Short: "setup initialize the site",
		Run:   siteSetupFunc,
	}
	runCommand = &cobra.Command{
		Use:   "run",
		Short: "run the site",
		Run:   siteRunFunc,
	}
)

// Run options
var (
	runConfigPath   string
	setupConfigPath string
)

func init() {
	runCommand.PersistentFlags().StringVar(&runConfigPath, "config", "", "site config file")
	setupCommand.PersistentFlags().StringVar(&setupConfigPath, "config", "", "site config file")

	rootCommand.AddCommand(setupCommand)
	rootCommand.AddCommand(runCommand)
}

func getModulePath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dir := filepath.Dir(path)
	return dir
}

func main() {
	if err := rootCommand.Execute(); nil != err {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}

func siteSetupFunc(cmd *cobra.Command, args []string) {
	var err error

	// Load config
	if setupConfigPath == "" {
		cmd.Println("config mustn't be empty")
		cmd.Usage()
		return
	}
	config, err := blog.ReadTOMLConfig(setupConfigPath)
	if nil != err {
		cmd.Println(err)
		return
	}

	site := blog.NewSite(config)
	if err = site.Setup(true); nil != err {
		cmd.Println(err)
		return
	}

	// Create default user
	var password string
	if password, err = site.NewAdmin(); nil != err {
		cmd.Println(err)
		return
	}
	cmd.Println("Setup ok, account:admin password:", password)
}

func siteRunFunc(cmd *cobra.Command, args []string) {
	var err error

	// Load config
	if runConfigPath == "" {
		cmd.Println("config mustn't be empty")
		cmd.Usage()
		return
	}
	config, err := blog.ReadTOMLConfig(runConfigPath)
	if nil != err {
		cmd.Println(err)
		return
	}

	if config.Debug {
		log.SetLevelByString("debug")
	} else {
		log.SetLevelByString("info")
	}

	// Clean up
	defer func() {
		e := recover()
		if nil != e {
			log.Error("Main routine quit with error:", e)
		} else {
			log.Info("Main routine quit normally")
		}
	}()

	site := blog.NewSite(config)
	if err = site.Start(); nil != err {
		log.Error("Blog start error:", err)
	}
}
