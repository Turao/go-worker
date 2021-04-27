package main

import (
	"log"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/turao/go-worker/client"
)

type Client interface {
	Dispatch(name string, args ...string) (interface{}, error)
	Stop(jobID string) (interface{}, error)
	Query(jobID string) (interface{}, error)
}

func NewCLI() *cobra.Command {
	url, _ := url.Parse("http://localhost:8080")
	c := client.New(url)

	cmd := &cobra.Command{
		Use: "cli",
	}

	cmd.AddCommand(makeDispatchCommand(c))
	cmd.AddCommand(makeStopCommand(c))
	cmd.AddCommand(makeQueryCommand(c))

	return cmd
}

func makeDispatchCommand(client Client) *cobra.Command {
	return &cobra.Command{
		Use:   "dispatch [command name] [...args]",
		Short: "dispatch a new job",
		Long:  "dispatch a new job with a unix process to be executed",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res, err := client.Dispatch(args[0], args[1:]...)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(res)
		},
	}
}

func makeStopCommand(client Client) *cobra.Command {
	return &cobra.Command{
		Use:   "stop [job id]",
		Short: "stop a running job",
		Long:  "stop a running job. a job MUST be running for this command to work",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res, err := client.Stop(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(res)
		},
	}
}

func makeQueryCommand(client Client) *cobra.Command {
	return &cobra.Command{
		Use:   "query [job id]",
		Short: "query job information",
		Long:  "query job information. a -1 exit code means that the job has not finished yet, or has been interrupted",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res, err := client.Query(args[0])
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(res)
		},
	}
}

func main() {
	err := NewCLI().Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
