package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/simonrw/circlecicli/internal/circleci"
)

func main() {

	// parse command line arguments
	runPtr := flag.Int("run", 0, "Run ID")
	ownerPtr := flag.String("owner", "", "GitHub owner")
	repoPtr := flag.String("repo", "", "GitHub repo")
	flag.Parse()

	if *runPtr == 0 {
		log.Fatal("no run id given")
	}

	if *ownerPtr == "" {
		log.Fatal("no owner given")
	}

	if *repoPtr == "" {
		log.Fatal("no repo given")
	}

	token := os.Getenv("CIRCLECI_TOKEN")
	client := circleci.New(token)
	pipeline, err := client.GetPipeline("gh", *ownerPtr, *repoPtr, *runPtr)
	if err != nil {
		log.Fatalf("getting pipeline: %v", err)
	}
	workflows, err := client.GetPipelineWorkflows(pipeline.ID)
	if err != nil {
		log.Fatalf("getting workflows: %v", err)
	}
	if len(workflows) == 0 {
		log.Fatalf("no workflows found")
	}
	wID := workflows[0].ID

	for {
		w, err := client.GetWorkflow(wID)
		if err != nil {
			log.Printf("fetching workflow: %v", err)

		} else {
			switch w.Status {
			case circleci.StatusSuccess, circleci.StatusFailed, circleci.StatusError, circleci.StatusCanceled, circleci.StatusUnauthorized:
				fmt.Printf("job finished, status: %s\n", w.Status)
				// terminal state
				if w.Status == circleci.StatusSuccess {
					os.Exit(0)
				} else {
					os.Exit(1)
				}
			default:
				time.Sleep(5 * time.Second)
			}
		}
	}
}
