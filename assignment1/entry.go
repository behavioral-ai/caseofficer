package assignment1

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/domain/common"
)

var (
	westData = []Entry{
		{Origin: common.Origin{Region: common.WestRegion, Zone: common.WestZoneA, Host: "host1.com"}},
		{Origin: common.Origin{Region: common.WestRegion, Zone: common.WestZoneB, Host: "host2.com"}},
	}

	centralData = []Entry{
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneA, Host: "host3.com"}},
		{Origin: common.Origin{Region: common.CentralRegion, Zone: common.CentralZoneB, Host: "host4.com"}},
	}

	/*
		eastData = []HostEntry{
			{Origin: common.Origin{Region: common.EastRegion, Zone: common.EastZoneA, SubZone: "", Host: "host5.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
		}

	*/
)

func get(origin common.Origin) ([]Entry, *messaging.Status) {
	if origin.Region == common.WestRegion {
		return westData, messaging.StatusOK()
	}
	if origin.Region == common.CentralRegion {
		return centralData, messaging.StatusOK()
	}
	return []Entry{}, messaging.StatusNotFound()
}
