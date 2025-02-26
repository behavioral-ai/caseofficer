package agent

import (
	"github.com/behavioral-ai/caseofficer/agent1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

func New(handler messaging.OpsAgent, origin common.Origin, dispatcher messaging.Dispatcher) messaging.Agent {
	return agent1.New(handler, origin, dispatcher)
}
