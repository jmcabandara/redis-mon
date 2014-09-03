package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

// vars set by build
var Version, Gitcommit, Gitbranch, Date, Buildtag string

var options struct {
	RedisPort      int
	MetricEndpoint string
}

// query and post metrics
func post(c *cli.Context) {
	sr := NewRedisStats(c.String("port"))
	stats, err := sr.Read()
	if err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
	if err = Post(c.String("consumer-url"), stats); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}

// print version infomation
func printVersion(c *cli.Context) {
	fmt.Printf("Version: %s\n", Version)
	if c.Bool("verbose") {
		fmt.Printf("Gitbranch: %s\n", Gitbranch)
		fmt.Printf("Date: %s\n", Date)
		fmt.Printf("Buildtag: %s\n", Buildtag)
	}
		
}

func main() {

	app := cli.NewApp()
	app.Name = "redis-mon"
	app.Version = Version
	app.Usage = "a simple client that posts redis-server status to Zenoss's metric service"
	app.Commands = []cli.Command{
		{
			Name:  "post",
			Usage: "post metrics to the metric service",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Usage: "redis-server port",
					Value: "localhost:6379",
				},
				cli.StringFlag{
					Name:  "consumer-url",
					Usage: "metric consumer url",
					Value: "http://localhost:22350/api/metrics/store",
				},
			},
			Action: post,
		},
		{
			Name:        "version",
			Description: "print version",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "verbose",
					Usage: "print extended version info",
				},
			},
			Action: printVersion,
		},
	}
	app.Run(os.Args)
}
