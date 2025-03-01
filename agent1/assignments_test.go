package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/test"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/timeseries1"
)

var (
	centralData = []timeseries1.Assignment{
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneA, Host: "host3.com"}},
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneB, Host: "host4.com"}},
	}
)

func testNewAgent(origin common.Origin, resolver collective.Resolution, dispatcher messaging.Dispatcher) messaging.Agent {
	return test.NewAgent(agentUri(origin))
}

func ExampleAddAssignments() {
	origin := common.Origin{Region: common.CentralRegion}
	agent := newAgent(origin, collective.NewEphemeralResolver(), nil)
	addAssignments(agent, centralData, testNewAgent)
	fmt.Printf("test: addAssignments() -> [assignments:%v]\n", agent.serviceAgents.Count())

	agent.finalize()

	//Output:
	//test: addAssignments() -> [assignments:2]

}

func ExampleUpdateAssignments() {
	origin := common.Origin{Region: common.WestRegion}

	agent := newAgent(origin, collective.NewEphemeralResolver(), nil)
	updateAssignments(agent, timeseries1.Assignments.All, testNewAgent)
	fmt.Printf("test: updateAssignments() -> [assignments:%v]\n", agent.serviceAgents.Count())

	agent.finalize()

	//Output:
	//test: addAssignments() -> [assignments:2]

}
