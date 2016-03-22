package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/vito/cloudformer/aws"

	"github.com/vito/boosh/builder"
)

func generate(source io.Reader) {
	var spec builder.DeploymentSpec

	var err error

	err = candiedyaml.NewDecoder(source).Decode(&spec)
	if err != nil {
		panic(err)
	}

	former := aws.New(spec.Description)

	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		region = "us-east-1"
	}

	builder := builder.New(spec, region)

	err = builder.Build(former)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(os.Stdout).Encode(former.Template)
	if err != nil {
		fatal(err)
	}
}
