package agent1

import (
	"github.com/behavioral-ai/core/aspect"
)

func addFeedback(agent *caseOfficer, newAgent createAgent) {
	a := newAgent(agent.origin, agent, agent.global)
	err := agent.serviceAgents.Register(a)
	if err != nil {
		agent.Notify(aspect.NewStatusError(aspect.StatusInvalidArgument, err))
	} else {
		a.Run()
	}
}
