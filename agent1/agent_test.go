package agent1

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

func ExampleNewAgent() {
	a := New(nil, common.Origin{Region: "us-central"}, messaging.Notify, messaging.NewTraceDispatcher())
	fmt.Printf("test: NewAgent() -> [%v] [%v]\n", a.Uri(), a.Name())

	//Output:
	//test: NewAgent() -> [resiliency:agent/caseofficer/agent1#] [resiliency:agent/caseofficer/agent]

}
