package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/timeseries1"
)

func updateAssignments(agent *agentT, query timeseries1.SelectAssignments, newAgent createAgent) {
	entry, status := query(agent.origin)
	if !status.OK() {
		if !status.NotFound() {
			agent.notify(status)
			return
		}
	}
	agent.addActivity(messaging.ActivityItem{
		Agent:   agent,
		Event:   updateAssignmentEvent,
		Source:  agent.emissary.Name(),
		Content: fmt.Sprintf("added %v assignments from %v", len(entry), agent.origin),
	})
	addAssignments(agent, entry, newAgent)
}

func addAssignments(agent *agentT, entry []timeseries1.Assignment, newAgent createAgent) {
	for _, e := range entry {
		a := newAgent(e.Origin, agent.activity, agent.notifier, agent.dispatcher)
		err := agent.serviceAgents.Register(a)
		if err != nil {
			agent.notify(messaging.NewStatusError(messaging.StatusInvalidArgument, err, agent.Uri()))
		} else {
			a.Run()
		}
	}
}
