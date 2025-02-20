package agent1

import (
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

type createAgent func(origin common.Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) messaging.Agent

/*
func newFeedbackAgent(origin core.Origin, handler messaging.OpsAgent, global messaging.Dispatcher) messaging.Agent {
	return feedback.NewAgent(origin, handler, global)
}
*/
//type newServiceAgent func(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher) messaging.Agent
//type newFeedbackAgent func(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher) messaging.Agent

func emissaryAttend(agent *caseOfficer, assignments *assignment1.Assignments, newService createAgent, newFeedback createAgent) {
	paused := false
	createAssignments(agent, assignments, newService)
	//addFeedback(agent, newFeedback)
	agent.startup()
	//agent.dispatch(messaging.StartupEvent)

	for {
		select {
		case <-agent.ticker.C():
			if !paused {
				updateAssignments(agent, assignments, newService)
				agent.dispatch(agent.ticker, messaging.TickEvent)
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				agent.finalize()
				agent.dispatch(agent.emissary, msg.Event())
				return
			case messaging.DataChangeEvent:

			default:

			}
			agent.dispatch(agent.emissary, msg.Event())
		default:
		}
	}
}
