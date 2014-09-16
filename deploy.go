package main

import (
	"io"
  "fmt"
	"io/ioutil"

	"github.com/dynport/gocloud/aws/cloudformation"
	"github.com/vito/cloudformer/aws/deployer"
  "github.com/mgutz/ansi"
)

func deploy(name string, source io.Reader) {
	template, err := ioutil.ReadAll(source)
	if err != nil {
		fatal(err)
	}

	deployer := deployer.New(cloudformation.NewFromEnv())

	events, err := deployer.Deploy(name, template)
	if err != nil {
		panic(err)
	}

  ok := streamEvents(events)
  if !ok {
    fmt.Println()
    fatal(ansi.Color("formation failed and was rolled back", "yellow"))
  }
}
