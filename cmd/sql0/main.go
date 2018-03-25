package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"encoding/json"
	"fmt"
	"github.com/k-kawa/sql0/app"
	cli "gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type config struct {
	ProjectID string
	Query     string
}

func loadConfig(filepath string) *config {
	var err error
	var f *os.File

	if f, err = os.Open(filepath); err != nil {
		o := newOutput(false, 0)
		o.Errors = append(o.Errors, fmt.Sprintf("Failed to open file: %s", err.Error()))
		printOutput(o)
		os.Exit(1)
	}
	defer f.Close()

	var b []byte
	if b, err = ioutil.ReadAll(f); err != nil {
		o := newOutput(false, 0)
		o.Errors = append(o.Errors, fmt.Sprintf("Failed to read file: %s", err.Error()))
		printOutput(o)
		os.Exit(1)
	}

	var c config
	if err = json.Unmarshal(b, &c); err != nil {
		o := newOutput(false, 0)
		o.Errors = append(o.Errors, fmt.Sprintf("Failed to parse config file: %s", err.Error()))
		printOutput(o)
		os.Exit(1)
	}

	return &c
}

type output struct {
	Success  bool
	RowCount int
	Time     time.Time
	Errors   []string
}

func newOutput(success bool, rowCount int) *output {
	return &output{
		Success:  success,
		RowCount: rowCount,
		Time:     time.Now(),
	}
}

func printOutput(o *output) {
	var err error
	var oByte []byte

	if oByte, err = json.MarshalIndent(o, "", "  "); err != nil {
		log.Panicf("%+v", err)
	}
	fmt.Println(string(oByte))
}

func NewCliApp() *cli.App {
	cliApp := cli.NewApp()
	cliApp.Name = "SQL0"
	cliApp.Usage = "The Simplest tool to check if the query returns 0"
	cliApp.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", EnvVar: "CONFIG_FILE"},
	}
	cliApp.Author = "Kohei Kawasaki"

	cliApp.Action = func(clic *cli.Context) error {
		ctx := context.Background()
		var err error

		c := loadConfig(clic.String("config"))

		var bqClient *bigquery.Client
		if bqClient, err = bigquery.NewClient(ctx, c.ProjectID); err != nil {
			o := newOutput(false, 0)
			o.Errors = append(o.Errors, err.Error())
			printOutput(o)
			os.Exit(1)
		}
		a := app.NewApp(bqClient)

		var result *app.Result
		if result, err = a.Run(ctx, c.Query); err != nil {
			o := newOutput(false, 0)
			o.Errors = append(o.Errors, err.Error())
			printOutput(o)
			os.Exit(1)
		}

		if result.IsOK() {
			return nil
		}

		o := newOutput(false, result.RowCount)

		if result.JobStatus.Err() != nil {
			o.Errors = append(o.Errors, result.JobStatus.Err().Error())
			for _, e := range result.JobStatus.Errors {
				o.Errors = append(o.Errors, e.Error())
			}
		}
		printOutput(o)
		os.Exit(1)
		return nil
	}
	return cliApp
}

func main() {
	cliApp := NewCliApp()
	_ = cliApp.Run(os.Args)
}
