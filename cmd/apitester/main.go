package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"os"
)

func main() {
	// 创建一个新的CLI应用程序
	app := &cli.Command{
		Name:        "API Tester",
		Usage:       "Welcome to API Testing Tool",
		Description: "This tool is developed in Go language and provides functionalities such as parsing OpenApi documents, automated interface testing, and generating test reports. You can use this tool to quickly and conveniently conduct API testing and generate detailed test reports for further analysis.",
		Commands: []*cli.Command{
			{
				Name:    "openapi",
				Aliases: []string{"swagger"},
				Usage:   "Parse OpenApi document and generate apis configuration",
				Action: func(ctx context.Context, c *cli.Command) error {
					configFile := c.String("c")
					outputFile := c.String("o")

					fmt.Printf("Run mode: %s\n", c.Name)
					fmt.Printf("Configuration file path: %s\n", configFile)
					fmt.Printf("Output path: %s\n", outputFile)

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
					outputFile := c.String("o")
					threads := c.Int("t")

					fmt.Printf("Run mode: %s\n", c.Name)
					fmt.Printf("Configuration file path: %s\n", configFile)
					fmt.Printf("Output path: %s\n", outputFile)
					fmt.Printf("Threads count: %d\n", threads)

					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "c",
						Usage: "Configuration file path",
					},
					&cli.StringFlag{
						Name:  "o",
						Usage: "Output path",
					},
					&cli.IntFlag{
						Name:  "t",
						Usage: "Threads count",
						Value: 1,
					},
				},
			},
		},
	}

	// 运行CLI应用程序
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
