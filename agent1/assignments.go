package agent1

import (
	"fmt"
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
)

func updateAssignments(agent *caseOfficer, resolver collective.Resolution, query assignment1.SelectAssignments, newAgent createAgent) {
	entry, status := query(agent.origin)
	if !status.OK() {
		if !status.NotFound() {
			agent.notify(status)
			return
		}
	}
	resolver.AddActivity(agent, updateAssignmentEvent, agent.emissary.Name(), fmt.Sprintf("added %v assignments", len(entry)))
	addAssignments(agent, entry, newAgent)
}

func addAssignments(agent *caseOfficer, entry []assignment1.Entry, newAgent createAgent) {
	for _, e := range entry {
		a := newAgent(agent, e.Origin, agent.dispatcher)
		err := agent.serviceAgents.Register(a)
		if err != nil {
			agent.notify(messaging.NewStatusError(messaging.StatusInvalidArgument, err, agent.Uri()))
		} else {
			a.Run()
		}
	}
}
