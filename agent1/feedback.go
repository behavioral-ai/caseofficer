package agent1

import "github.com/behavioral-ai/core/messaging"

func addFeedback(agent *caseOfficer, newAgent createAgent) {
	a := newAgent(agent.origin, agent.notifier, agent.dispatcher)
	err := agent.serviceAgents.Register(a)
	if err != nil {
		agent.Notify(messaging.NewStatusError(0, err))
	} else {
		a.Run()
	}
}
