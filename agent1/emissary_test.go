package agent1

import (
	"github.com/behavioral-ai/caseofficer/assignment1"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/collective"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/agent"
	"time"
)

const (
	testDuration = time.Second * 5
)

func ExampleEmissary() {
	ch := make(chan struct{})
	officer := newAgent(nil, common.Origin{Region: common.WestRegion}, messaging.Notify, messaging.NewTraceDispatcher())

	go func() {
		go emissaryAttend(officer, collective.NewEphemeralResolver(), assignment1.Entries, agent.New)
		officer.Shutdown()
		time.Sleep(testDuration)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
