package agent1

import (
	"fmt"
	"github.com/behavioral-ai/collective/content"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

func ExampleNewAgent() {
	a := New(common.Origin{Region: common.CentralRegion}, content.NewEphemeralResolver(), messaging.NewTraceDispatcher())
	fmt.Printf("test: NewAgent() -> [%v] [%v]\n", a.Uri(), a.Name())

	//Output:
	//test: NewAgent() -> [resiliency:agent/caseofficer1#us-central1] [resiliency:agent/caseofficer]

}
