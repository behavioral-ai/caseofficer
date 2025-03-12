package agent1

import (
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/timeseries1"
)

const (
	updateAssignmentEvent = "event:update-assignments"
)

type createAgent func(origin common.Origin, resolver content.Resolution, dispatcher messaging.Dispatcher) messaging.Agent

func emissaryAttend(agent *agentT, assignments *timeseries1.Assigner, newService createAgent, s messaging.Spanner) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false
	updateAssignments(agent, assignments.All, newService)
	agent.reviseTicker(s)

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(agent.ticker, messaging.TickEvent)
			if !paused {
				updateAssignments(agent, assignments.New, newService)
			}
			agent.reviseTicker(s)
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
