package assignment1

import (
	"fmt"
	"github.com/behavioral-ai/domain/common"
)

func ExampleGet() {
	o := common.Origin{Region: common.EastRegion, Zone: common.WestZoneA}
	e, status := get(o)
	fmt.Printf("test: Get(\"%v\") -> [status:%v] [%v]\n", o, status, e)

	o = common.Origin{Region: common.WestRegion, Zone: common.WestZoneA}
	e, status = get(o)
	fmt.Printf("test: Get(\"%v\") -> [status:%v] [%v]\n", o, status, e)

	o = common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneA}
	e, status = get(o)
	fmt.Printf("test: Get(\"%v\") -> [status:%v] [%v]\n", o, status, e)

	//Output:
	//test: Get("us-east1.w-zone-a.") -> [status:Not Found] [[]]
	//test: Get("us-west1.w-zone-a.") -> [status:OK] [[{us-west1.w-zone-a.host1.com} {us-west1.w-zone-b.host2.com}]]
	//test: Get("us-central1.c-zone-a.") -> [status:OK] [[{us-central1.c-zone-a.host3.com} {us-central1.c-zone-b.host4.com}]]

}
