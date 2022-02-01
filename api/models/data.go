package models

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

// {
// 	[
// 	time:
// 	region_id:
// 	value:
// 	]
// }

type DemandDataPoint struct {
	Time     time.Time `json:"time"`
	RegionID string    `json:"region_id"`
	Value    float64   `json:"value"`
}

type DemandFilter struct {
	RegionID StringFilter `col:"regionId"`
}

// Filters include
// [] of regions
// start string
// stop string
// aggregate (every + fn)

func ReadDemandData(db api.QueryAPI, bucket string, filter DemandFilter) ([]DemandDataPoint, error) {
	fmt.Println(buildFluxQuery(filter))
	points := make([]DemandDataPoint, 0)

	fluxQuery := fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: -7d)
			|> filter(fn: (r) => r.regionId == "NSW1")
			|> filter(fn: (r) => r._measurement == "demand")
	`, bucket)

	result, err := db.Query(context.Background(), fluxQuery)

	if err != nil {
		return []DemandDataPoint{}, fmt.Errorf("models.ReadDemandData: query error: %v", err)
	}

	for result.Next() {
		value, _ := getFloatReflectOnly(result.Record().Value())
		dataPoint := DemandDataPoint{
			Time:     result.Record().Time(),
			RegionID: fmt.Sprintf("%v", result.Record().ValueByKey("regionId")),
			Value:    value,
		}
		points = append(points, dataPoint)
	}

	if result.Err() != nil {
		return []DemandDataPoint{}, fmt.Errorf("models.ReadDemandData: query parsing error: %v", result.Err())
	}

	return points, nil
}

func FilterMaptoDemandFilter(filterMap map[string][]string) DemandFilter {
	var filter DemandFilter
	fmt.Println(filterMap)

	filter.RegionID.fromFilterMap(filterMap, "region_id")

	return filter
}
