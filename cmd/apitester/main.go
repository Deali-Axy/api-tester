package main

import (
	"api-tester/pkg/api/tester"
	"api-tester/pkg/utilities/exporter"
	tlog "api-tester/pkg/utilities/logger"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/nleeper/goment"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	var logger *zap.SugaredLogger

	app := &cli.Command{
		Name:        "API Tester",
		Usage:       "Welcome to API Testing Tool",
		Description: "This tool is developed in Go language and provides functionalities such as parsing OpenApi documents, automated interface testing, and generating test reports. You can use this tool to quickly and conveniently conduct API testing and generate detailed test reports for further analysis.",
		Before: func(ctx context.Context, command *cli.Command) error {
			logToFile := command.Bool("l")
			logPath := command.String("p")
			var err error
			logger, err = tlog.BuildLogger(&tlog.Options{
				WriteToFile: logToFile,
				Folder:      logPath,
			})
			if err != nil {
				return err
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "l",
				Aliases: []string{"log-to-file"},
				Usage:   "write log to file",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "p",
				Aliases: []string{"log-path"},
				Usage:   "The path to store log files",
				Value:   "logs",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "parse",
				Aliases: []string{"openapi", "swagger"},
				Usage:   "Parse OpenApi document and generate apis configuration",
				Action: func(ctx context.Context, c *cli.Command) error {
					configFile := c.String("c")
					outputFile := c.String("o")

					logger.Debugf("Run mode: %s", c.Name)
					logger.Debugf("Configuration file path: %s", configFile)
					logger.Debugf("Output path: %s", outputFile)

					parser := tester.Parser{
						Logger: logger,
					}
					if _, err := parser.LoadFromFile(configFile); err != nil {
						return err
					}

					if err := parser.SaveConfig(outputFile); err != nil {
						return err
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "c",
						Aliases:  []string{"openapi-doc", "swagger-doc"},
						Usage:    "Configuration file path",
						Required: true,
					},
					&cli.StringFlag{
						Name:    "o",
						Aliases: []string{"output"},
						Usage:   "Output path",
					},
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Run Apis test",
				Action: func(ctx context.Context, c *cli.Command) error {
					configFile := c.String("c")
					outputPath := c.String("o")
					baseUrl := c.String("base-url")
					threads := c.Int("t")
					proxy := c.String("proxy")
					timeout := c.Duration("timeout")
					authToken := c.String("auth-token")

					logger.Debugf("Run mode: %s", c.Name)
					logger.Debugf("Configuration file path: %s", configFile)
					logger.Debugf("Output path: %s", outputPath)
					logger.Debugf("Threads count: %d", threads)

					t := tester.Tester{
						Logger:      logger,
						RestyClient: resty.New(),
						BaseURL:     baseUrl,
					}

					t.RestyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

					if len(proxy) > 0 {
						t.RestyClient.SetProxy(proxy)
					}
					t.RestyClient.SetTimeout(timeout)

					if len(authToken) > 0 {
						t.AuthToken = authToken
					}

					apis, err := tester.ReadConfig(configFile)
					if err != nil {
						logger.Errorln(err)
						return err
					}

					reports, err := t.TestApis(apis, int(threads))
					if err != nil {
						logger.Errorln(err)
						return err
					}

					g, _ := goment.New()
					timeStr := g.Format("YYYY-MM-DD_HH-mm-ss-x")
					e := exporter.Exporter{}
					if err := e.ToExcel(reports,
						filepath.Join(outputPath, fmt.Sprintf("reports_%s.xlsx", timeStr)),
					); err != nil {
						logger.Errorln(err)
						return err
					}

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "c",
						Usage:    "Configuration file path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "o",
						Usage:    "Output folder path",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "base-url",
						Required: true,
					},
					&cli.DurationFlag{
						Name:  "timeout",
						Usage: "Set request timeout",
					},
					&cli.StringFlag{
						Name:  "proxy",
						Usage: "Set proxy url",
					},
					&cli.IntFlag{
						Name:  "t",
						Usage: "Threads count",
						Value: 1,
					},
					&cli.StringFlag{
						Name:  "auth-token",
						Usage: "Set auth token",
					},
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err)
	}
}
