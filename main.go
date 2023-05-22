package main

import (
	"fmt"

	"github.com/jszwedko/go-circleci"
)

func main() {
	client := &circleci.Client{Token: ""}

	limit := 10
	offset := 0
	builds, err := client.ListRecentBuildsForProject("localstack", "localstack", "master", "", limit, offset)
	if err != nil{
		panic(err)
	}

	for _, build := range builds {
		fmt.Printf("%d: %s\n", build.BuildNum, build.Status)
	}
}
