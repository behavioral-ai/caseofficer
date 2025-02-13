package agent1

func addFeedback(agent *caseOfficer, newAgent createAgent) {
	a := newAgent(agent.origin, agent, agent.global)
	err := agent.serviceAgents.Register(a)
	if err != nil {
		agent.Notify(err)
	} else {
		a.Run()
	}
}
