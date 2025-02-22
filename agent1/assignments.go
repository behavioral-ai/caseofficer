package agent1

import (
	"errors"
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

func createAssignments(agent *caseOfficer, assignments *assignment1.Assignments, newAgent createAgent) {
	if newAgent == nil {
		agent.notify(messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("error: create assignments newAgent is nil")))
		return
	}
	entry, status := assignments.All(agent.origin)
	if status == nil {
		addAssignments(agent, entry, newAgent)
	}
	//if !status.NotFound() {
	agent.notify(status)
	//}
}

func updateAssignments(agent *caseOfficer, assignments *assignment1.Assignments, newAgent createAgent) {
	if newAgent == nil {
		agent.notify(messaging.NewStatusError(http.StatusBadRequest, errors.New("error: update assignments newAgent is nil")))
		return
	}
	entry, status := assignments.New(agent.origin)
	if !status.OK() {
		addAssignments(agent, entry, newAgent)
	}
	//	if !status.NotFound() {
	agent.notify(status)
	//	}
}

func addAssignments(agent *caseOfficer, entry []assignment1.Entry, newAgent createAgent) {
	for _, e := range entry {
		a := newAgent(e.Origin, agent.notifier, agent.dispatcher)
		err := agent.serviceAgents.Register(a)
		if err != nil {
			agent.notify(messaging.NewStatusError(messaging.StatusInvalidArgument, err))
		} else {
			a.Run()
		}
	}
}
