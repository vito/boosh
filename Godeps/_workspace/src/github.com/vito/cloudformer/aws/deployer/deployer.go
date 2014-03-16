package deployer

import (
	"sort"
	"strings"
	"time"

	"github.com/dynport/gocloud/aws/cloudformation"
)

type AWSDeployer struct {
	client *cloudformation.Client
}

func New(client *cloudformation.Client) *AWSDeployer {
	return &AWSDeployer{client}
}

func (deployer *AWSDeployer) Deploy(name string, template []byte) (<-chan *cloudformation.StackEvent, error) {
	stacks, err := deployer.client.DescribeStacks(&cloudformation.DescribeStacksParameters{
		StackName: name,
	})

	update := false
	if err == cloudformation.ErrorNotFound {
		update = false
	} else if err != nil {
		return nil, err
	} else if stacks.DescribeStacksResult.Stacks[0].StackStatus == "ROLLBACK_COMPLETE" {
		err := deployer.client.DeleteStack(name)
		if err != nil {
			return nil, err
		}
	} else {
		update = true
	}

	events := make(chan *cloudformation.StackEvent)

	go deployer.watchEvents(events, name)

	if update {
		_, err = deployer.client.UpdateStack(cloudformation.UpdateStackParameters{
			BaseParameters: cloudformation.BaseParameters{
				StackName:    name,
				TemplateBody: string(template),
			},
		})
		if err != nil {
			if err.Error() == "No updates are to be performed." {
				close(events)
				return events, nil
			}

			return nil, err
		}
	} else if !update {
		_, err = deployer.client.CreateStack(cloudformation.CreateStackParameters{
			BaseParameters: cloudformation.BaseParameters{
				StackName:    name,
				TemplateBody: string(template),
			},
		})
		if err != nil {
			return nil, err
		}
	}

	return events, nil
}

func (deployer *AWSDeployer) watchEvents(events chan<- *cloudformation.StackEvent, name string) {
	start := time.Now()

	nextToken := ""

	seen := make(map[string]bool)

	for {
		params := &cloudformation.DescribeStackEventsParameters{
			StackName: name,
			NextToken: nextToken,
		}

		res, err := deployer.client.DescribeStackEvents(params)
		if err != nil {
			continue
		}

		completed := false

		sorted := Events(res.DescribeStackEventsResult.StackEvents)
		sort.Sort(sorted)

		for _, ev := range sorted {
			if seen[ev.EventId] {
				continue
			}

			seen[ev.EventId] = true

			if ev.Timestamp.Before(start) {
				continue
			}

			events <- ev

			if ev.ResourceType == "AWS::CloudFormation::Stack" {
				completed = !strings.HasSuffix(ev.ResourceStatus, "_IN_PROGRESS")
			}
		}

		if completed {
			close(events)
			return
		}

		nextToken = res.DescribeStackEventsResult.NextToken

		time.Sleep(1 * time.Second)
	}
}

type Events []*cloudformation.StackEvent

func (list Events) Len() int {
	return len(list)
}

func (list Events) Swap(a, b int) {
	list[a], list[b] = list[b], list[a]
}

func (list Events) Less(a, b int) bool {
	return list[a].Timestamp.Before(list[b].Timestamp)
}
