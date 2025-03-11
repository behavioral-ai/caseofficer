package agent1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/messaging/messagingtest"
	"github.com/behavioral-ai/core/test"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/timeseries1"
	"time"
)

const (
	testDuration = time.Second * 5
)

func operativeNew(origin common.Origin, resolver collective.Resolution, dispatcher messaging.Dispatcher) messaging.Agent {
	return test.NewAgent("resiliency:agent/operative#" + origin.String())
}

func ExampleEmissary() {
	ch := make(chan struct{})
	s := messagingtest.NewTestSpanner(time.Second*2, testDuration)
	dispatcher := messaging.NewFilteredTraceDispatcher([]string{messaging.StartupEvent, messaging.ShutdownEvent}, "")
	agent := newAgent(common.Origin{Region: common.WestRegion}, collective.NewEphemeralResolver(), dispatcher)

	go func() {
		go emissaryAttend(agent, timeseries1.Assignments, operativeNew, s)
		time.Sleep(testDuration * 3)
		agent.Shutdown()
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
