package agent1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/domain/timeseries1"
	"github.com/behavioral-ai/operative/agent"
	"time"
)

const (
	testDuration = time.Second * 5
)

func ExampleEmissary() {
	ch := make(chan struct{})
	dispatcher := messaging.NewFilteredTraceDispatcher([]string{messaging.StartupEvent, messaging.ShutdownEvent}, "")
	officer := newAgent(common.Origin{Region: common.WestRegion}, collective.NewEphemeralResolver(), dispatcher)

	go func() {
		go emissaryAttend(officer, timeseries1.Assignments, agent.New)
		time.Sleep(testDuration * 3)
		officer.Shutdown()
		time.Sleep(testDuration * 2)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
