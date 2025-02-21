package assignment1

import (
	"github.com/behavioral-ai/domain/common"
	"time"
)

var (
	//safeEntry = guidance.NewSafe()
	westData = []HostEntry{
		{Origin: common.Origin{Region: common.WestRegion, Zone: common.WestZoneA, SubZone: "", Host: "host1.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
		{Origin: common.Origin{Region: common.WestRegion, Zone: common.WestZoneB, SubZone: "", Host: "host2.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
	}

	centralData = []HostEntry{
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneA, SubZone: "", Host: "host3.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneB, SubZone: "", Host: "host4.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
	}

	eastData = []HostEntry{
		{Origin: common.Origin{Region: common.EastRegion, Zone: common.EastZoneA, SubZone: "", Host: "host5.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
	}
)

// HostEntry - host
type HostEntry struct {
	EntryId   int           `json:"entry-id"`
	CreatedTS time.Time     `json:"created-ts"`
	Origin    common.Origin `json:"origin"`
}
