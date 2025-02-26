package agent1

import (
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

type createAgent func(origin common.Origin, notifier messaging.NotifyFunc, dispatcher messaging.Dispatcher) messaging.Agent

func emissaryAttend(agent *caseOfficer, assignments *assignment1.Assignments, newService createAgent) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false
	createAssignments(agent, assignments, newService)
	agent.startup()

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(agent.ticker, messaging.TickEvent)
			if !paused {
				updateAssignments(agent, assignments, newService)
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.dispatch(agent.emissary, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				agent.finalize()
				return
			case messaging.DataChangeEvent:
			default:
			}
		default:
		}
	}
}
