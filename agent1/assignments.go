package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/timeseries1"
)

func updateAssignments(agent *caseOfficer, query timeseries1.SelectAssignments, newAgent createAgent) {
	entry, status := query(agent.origin)
	if !status.OK() {
		if !status.NotFound() {
			agent.resolver.Notify(status)
			return
		}
	}
	agent.resolver.AddActivity(agent, updateAssignmentEvent, agent.emissary.Name(), fmt.Sprintf("added %v assignments", len(entry)))
	addAssignments(agent, entry, newAgent)
}

func addAssignments(agent *caseOfficer, entry []timeseries1.Assignment, newAgent createAgent) {
	for _, e := range entry {
		a := newAgent(e.Origin, agent.resolver, agent.dispatcher)
		err := agent.serviceAgents.Register(a)
		if err != nil {
			agent.resolver.Notify(messaging.NewStatusError(messaging.StatusInvalidArgument, err, agent.Uri()))
		} else {
			a.Run()
		}
	}
}
