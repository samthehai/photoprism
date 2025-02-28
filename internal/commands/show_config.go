package commands

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/photoprism/photoprism/internal/config"
	"github.com/photoprism/photoprism/internal/get"
	"github.com/photoprism/photoprism/pkg/report"
)

// ShowConfigCommand configures the command name, flags, and action.
var ShowConfigCommand = cli.Command{
	Name:   "config",
	Usage:  "Displays global config options and their current values",
	Flags:  report.CliFlags,
	Action: showConfigAction,
}

// ConfigReports specifies which configuration reports to display.
var ConfigReports = []Report{
	{Title: "Config Options", NoWrap: true, Report: func(conf *config.Config) ([][]string, []string) {
		return conf.Report()
	}},
}

// showConfigAction shows global config option names and values.
func showConfigAction(ctx *cli.Context) error {
	conf := config.NewConfig(ctx)
	conf.SetLogLevel(logrus.FatalLevel)
	get.SetConfig(conf)

	if err := conf.Init(); err != nil {
		log.Debug(err)
	}

	fmt.Println("")

	for _, rep := range ConfigReports {
		// Get values.
		rows, cols := rep.Report(conf)

		// Render report.
		opt := report.Options{Format: report.CliFormat(ctx), NoWrap: rep.NoWrap}
		result, _ := report.Render(rows, cols, opt)

		// Show report.
		if opt.Format == report.Default {
			fmt.Printf("### %s ###\n\n", rep.Title)
		}

		fmt.Println(result)
	}

	return nil
}
