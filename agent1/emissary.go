package agent1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/timeseries1"
)

const (
	updateAssignmentEvent = "event:update-assignments"
)

type createAgent func(origin common.Origin, resolver collective.Resolution, dispatcher messaging.Dispatcher) messaging.Agent

func emissaryAttend(agent *caseOfficer, assignments *timeseries1.Assigner, newService createAgent) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false
	updateAssignments(agent, assignments.All, newService)
	agent.startup()

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(agent.ticker, messaging.TickEvent)
			if !paused {
				updateAssignments(agent, assignments.New, newService)
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.dispatch(agent.emissary, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
				agent.serviceAgents.Broadcast(messaging.Pause)
			case messaging.ResumeEvent:
				paused = false
				agent.serviceAgents.Broadcast(messaging.Resume)
			case messaging.ShutdownEvent:
				agent.finalize()
				return
			default:
			}
		default:
		}
	}
}
