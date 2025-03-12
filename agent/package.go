package agent

import (
	"github.com/behavioral-ai/caseofficer/agent1"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

func New(origin common.Origin, resolver content.Resolution, dispatcher messaging.Dispatcher) messaging.Agent {
	return agent1.New(origin, resolver, dispatcher)
}
