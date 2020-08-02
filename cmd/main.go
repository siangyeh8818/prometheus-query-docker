package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	utils "github.com/siangyeh8818/prometheus-query-docker/internal"
	prometheus "github.com/siangyeh8818/prometheus-query-docker/internal/client"
	nexus "github.com/siangyeh8818/prometheus-query-docker/internal/nexus"
	"github.com/ymotongpoo/datemaki"
)

type options struct {
	format string
	server string
	query  string
	start  string
	end    string
	step   string
}

func main() {
	options := parseFlags()
	err := validateOptions(options)
	if err != nil {
		onError(err)
	}

	start, err := datemaki.Parse(options.start)
	if err != nil {
		onError(err)
	}

	end, err := datemaki.Parse(options.end)
	if err != nil {
		onError(err)
	}

	step, err := time.ParseDuration(options.step)
	if err != nil {
		onError(err)
	}

	client, err := prometheus.NewClient(options.server)
	if err != nil {
		onError(err)
	}

	resp, err := client.QueryRange(options.query, start, end, step)
	if err != nil {
		onError(err)
	}

	err = utils.PrintResp(resp, options.format, "report.csv")
	if err != nil {
		onError(err)
	}
	if options.format == "csv" {
		nTime := time.Now()
		local1, _ := time.LoadLocation("Asia/Taipei") //等同于"CST"

		logDay := nTime.In(local1).Format("20060102")
		NexusServer := os.Getenv("NEXUS_SERVER")
		fmt.Printf("NEXUS_SERVER : %s", NexusServer)
		NexusUser := os.Getenv("NEXUS_USER")
		NexusPassword := os.Getenv("NEXUS_PASSWORD")
		NexusRepository := os.Getenv("NEXUS_REPOSITORY")
		fmt.Printf("NEXUS_REPOSITORY : %s", NexusRepository)
		PostURL := "curl -X POST  '" + NexusServer + "/service/rest/v1/components?repository=" + NexusRepository + "' --user  " + NexusUser + ":" + NexusPassword + " -F 'raw.directory=" + logDay + "' -F 'raw.asset1=@report.csv;type=text/csv' -F 'raw.asset1.filename=report.csv' -H 'accept: application/json' -H 'Content-Type: multipart/form-data'"
		fmt.Println(PostURL)

		result, _ := nexus.ExecShell(PostURL)
		fmt.Println(result)
		//nexus.POSTForm_NesusAPI("report.csv")
		//nexus.PostNesusAPI(PostURL, NexusUser, NexusPassword, "")
	}

}

func onError(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func parseFlags() options {
	format := flag.String("format", "json", "Format (available formats are json, tsv and csv)")
	server := flag.String("server", os.Getenv("PROMETHEUS_SERVER"), "Prometheus server URL like 'https://prometheus.example.com' (can be set by PROMETHEUS_SERVER environment variable)")
	query := flag.String("query", "", "Query")
	start := flag.String("start", "1 hour ago", "Start time")
	end := flag.String("end", "now", "End time")
	step := flag.String("step", "15s", "Step")

	flag.Parse()

	return options{
		format: *format,
		server: *server,
		query:  *query,
		start:  *start,
		end:    *end,
		step:   *step,
	}
}

func validateOptions(options options) error {
	if options.server == "" {
		return errors.New("-server is mandatory")
	}
	if options.query == "" {
		return errors.New("-query is mandatory")
	}

	return nil
}
