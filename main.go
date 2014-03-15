package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dynport/gocloud/aws/cloudformation"
	"github.com/fraenkel/candiedyaml"
	"github.com/mgutz/ansi"

	"github.com/vito/cloudformer/aws"
	"github.com/vito/cloudformer/aws/deployer"
)

func main() {
	var spec DeploymentSpec

	var source io.Reader

	var err error

	if len(os.Args) > 1 {
		file, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}

		source = file
	} else {
		source = os.Stdin
	}

	err = candiedyaml.NewDecoder(source).Decode(&spec)
	if err != nil {
		panic(err)
	}

	former := aws.New(spec.Description)

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	builder := Builder{
		Region: region,
		Spec:   spec,
	}

	err = builder.Build(former)
	if err != nil {
		panic(err)
	}

	payload, err := json.Marshal(former.Template)
	if err != nil {
		panic(err)
	}

	deployer := deployer.New(cloudformation.NewFromEnv())

	events, err := deployer.Deploy(spec.Name, payload)
	if err != nil {
		panic(err)
	}

	ok := renderEvents(events)
	if !ok {
		fmt.Println()
		fmt.Println(ansi.Color("formation failed and was rolled back", "yellow"))
		os.Exit(1)
	}
}
