package models

import (
	"context"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
)

type DemandDataPoint struct {
	Time     time.Time `json:"time"`
	RegionID string    `json:"region_id"`
	Value    float64   `json:"value"`
}

type RooftopDataPoint struct {
	Time     time.Time `json:"time"`
	RegionID string    `json:"region_id"`
	Value    float64   `json:"value"`
}

type DemandFilter struct {
	Range     RangeFilter     `col:"range"` // col is unused for range but required for parsing
	RegionID  StringFilter    `col:"regionId"`
	Aggregate AggregateFilter `col:"aggregate"` // col is unused for aggregate but required for parsing
}

type RooftopFilter struct {
	Range     RangeFilter     `col:"range"` // col is unused for range but required for parsing
	RegionID  StringFilter    `col:"regionId"`
	Aggregate AggregateFilter `col:"aggregate"` // col is unused for aggregate but required for parsing
}

func ReadDemandData(db api.QueryAPI, bucket string, filter DemandFilter) ([]DemandDataPoint, error) {
	points := make([]DemandDataPoint, 0)

	fluxQuery := fmt.Sprintf("from(bucket: \"%s\")", bucket)
	fluxQuery += buildFluxQuery(filter)
	fluxQuery += "\n\t|> filter(fn: (r) => r._measurement == \"demand\")"

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

func ReadRooftapData(db api.QueryAPI, bucket string, filter RooftopFilter) ([]RooftopDataPoint, error) {
	points := make([]RooftopDataPoint, 0)

	fluxQuery := fmt.Sprintf("from(bucket: \"%s\")", bucket)
	fluxQuery += buildFluxQuery(filter)
	fluxQuery += "\n\t|> filter(fn: (r) => r._measurement == \"rooftop\")"

	result, err := db.Query(context.Background(), fluxQuery)

	if err != nil {
		return []RooftopDataPoint{}, fmt.Errorf("models.ReadDemandData: query error: %v", err)
	}

	for result.Next() {
		value, _ := getFloatReflectOnly(result.Record().Value())
		dataPoint := RooftopDataPoint{
			Time:     result.Record().Time(),
			RegionID: fmt.Sprintf("%v", result.Record().ValueByKey("regionId")),
			Value:    value,
		}
		points = append(points, dataPoint)
	}

	if result.Err() != nil {
		return []RooftopDataPoint{}, fmt.Errorf("models.ReadDemandData: query parsing error: %v", result.Err())
	}

	return points, nil
}

func ReadGenerationData(db api.QueryAPI, bucket string, filter DemandFilter) ([]DemandDataPoint, error) {
	points := make([]DemandDataPoint, 0)

	fluxQuery := fmt.Sprintf("from(bucket: \"%s\")", bucket)
	fluxQuery += buildFluxQuery(filter)
	fluxQuery += "\n\t|> filter(fn: (r) => r._measurement == \"demand\")"

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

	filter.Range.fromFilterMap(filterMap, "range")
	filter.RegionID.fromFilterMap(filterMap, "region_id")
	filter.Aggregate.fromFilterMap(filterMap, "aggregate")

	return filter
}

func FilterMaptoRooftopFilter(filterMap map[string][]string) RooftopFilter {
	var filter RooftopFilter
	fmt.Println(filterMap)

	filter.Range.fromFilterMap(filterMap, "range")
	filter.RegionID.fromFilterMap(filterMap, "region_id")
	filter.Aggregate.fromFilterMap(filterMap, "aggregate")

	return filter
}
