package agent1

import (
	"fmt"
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/test"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
)

var (
	centralData = []assignment1.Entry{
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneA, Host: "host3.com"}},
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneB, Host: "host4.com"}},
	}
)

func testNewAgent(handler messaging.Agent, origin common.Origin, dispatcher messaging.Dispatcher) messaging.Agent {
	return test.NewAgent(agentUri(origin))
}

func ExampleAddAssignments() {
	origin := common.Origin{Region: common.CentralRegion}
	agent := newAgent(nil, origin, messaging.Notify, nil)
	addAssignments(agent, centralData, testNewAgent)
	fmt.Printf("test: addAssignments() -> [assignments:%v]\n", agent.serviceAgents.Count())

	agent.finalize()

	//Output:
	//test: addAssignments() -> [assignments:2]

}

func ExampleUpdateAssignments() {
	origin := common.Origin{Region: common.WestRegion}

	agent := newAgent(nil, origin, messaging.Notify, nil)
	updateAssignments(agent, collective.NewEphemeralResolver(), assignment1.Entries.All, testNewAgent)
	fmt.Printf("test: updateAssignments() -> [assignments:%v]\n", agent.serviceAgents.Count())

	agent.finalize()

	//Output:
	//test: addAssignments() -> [assignments:2]

}
