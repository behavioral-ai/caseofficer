package agent1

import (
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
)

const (
	updateAssignmentEvent = "event:update-assignments"
)

type createAgent func(handler messaging.Agent, origin common.Origin, dispatcher messaging.Dispatcher) messaging.Agent

func emissaryAttend(agent *caseOfficer, resolver collective.Resolution, assignments *assignment1.Assignments, newService createAgent) {
	agent.dispatch(agent.emissary, messaging.StartupEvent)
	paused := false
	updateAssignments(agent, resolver, assignments.All, newService)
	agent.startup()

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(agent.ticker, messaging.TickEvent)
			if !paused {
				updateAssignments(agent, resolver, assignments.New, newService)
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
