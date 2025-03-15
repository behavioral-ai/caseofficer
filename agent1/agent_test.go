package agent1

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
	"github.com/behavioral-ai/operative/test"
	"time"
)

func ExampleNewAgent() {
	a := New(common.Origin{Region: common.CentralRegion}, content.NewEphemeralResolver(), messaging.NewTraceDispatcher())
	fmt.Printf("test: NewAgent() -> [%v] [%v]\n", a.Uri(), a.Name())

	//Output:
	//test: NewAgent() -> [resiliency:agent/behavioral-ai/caseofficer1#us-central1] [resiliency:agent/caseofficer]

}

func ExampleAgent_Run() {
	ch := make(chan struct{})
	//s := messagingtest.NewTestSpanner(time.Second*2, testDuration)
	dispatcher := messaging.NewFilteredTraceDispatcher([]string{messaging.ResumeEvent, messaging.PauseEvent}, "")
	r, _ := test.NewResiliencyResolver()
	agent := newAgent(common.Origin{Region: common.WestRegion}, r, dispatcher)

	go func() {
		agent.Run()
		time.Sleep(testDuration * 6)
		agent.Message(messaging.Pause)
		time.Sleep(testDuration * 6)
		agent.Message(messaging.Resume)
		time.Sleep(testDuration * 6)
		agent.Shutdown()
		time.Sleep(testDuration * 4)
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//fail

}
