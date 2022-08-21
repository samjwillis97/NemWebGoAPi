package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

type DemandDataPoint struct {
	Time     time.Time `json:"time"`
	RegionID string    `json:"region_id"`
	Value    float64   `json:"value"`
}

type DataPoint struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

type RooftopDataPoint struct {
	Time     time.Time `json:"time"`
	RegionID string    `json:"region_id"`
	Value    float64   `json:"value"`
}

type GenerationDataPoint struct {
	Unit string      `json:"unit"`
	Data []DataPoint `json:"data"`
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

type GeneratorFilter struct {
	Range     RangeFilter     `col:"range"` // col is unused for range but required for parsing
	DuID      StringFilter    `col:"unit" param:"duid"`
	Aggregate AggregateFilter `col:"aggregate"` // col is unused for aggregate but required for parsing
}

type GeneratorGroupedFilter struct {
	Group StringFilter `col:"group"`
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
		return []RooftopDataPoint{}, fmt.Errorf("models.ReadRooftapData: query error: %v", err)
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
		return []RooftopDataPoint{}, fmt.Errorf("models.ReadRooftapData: query parsing error: %v", result.Err())
	}

	return points, nil
}

func ReadGenerationData(db api.QueryAPI, bucket string, filter GeneratorFilter) ([]GenerationDataPoint, error) {

	data := make([]GenerationDataPoint, 0)

	fluxQuery := fmt.Sprintf("from(bucket: \"%s\")", bucket)
	fluxQuery += buildFluxQuery(filter)
	fluxQuery += "\n\t|> filter(fn: (r) => r._measurement == \"generation\")"

	log.Debugln(fluxQuery)

	result, err := db.Query(context.Background(), fluxQuery)

	if err != nil {
		return []GenerationDataPoint{}, fmt.Errorf("models.ReadGenerationData: query error: %v", err)
	}

	var units []string
	unitMap := make(map[string][]DataPoint)

	for result.Next() {
		value, _ := getFloatReflectOnly(result.Record().Value())
		unitName := fmt.Sprintf("%v", result.Record().ValueByKey("unit"))

		if _, ok := unitMap[unitName]; ok {
			unitMap[unitName] = append(unitMap[unitName], DataPoint{
				Time:  result.Record().Time(),
				Value: value,
			})
		} else {
			units = append(units, unitName)
			unitMap[unitName] = []DataPoint{{
				Time:  result.Record().Time(),
				Value: value,
			}}
		}
	}

	if result.Err() != nil {
		return []GenerationDataPoint{}, fmt.Errorf("models.ReadGenerationData: query parsing error: %v", result.Err())
	}

	for _, v := range units {
		data = append(data, GenerationDataPoint{
			Unit: v,
			Data: unitMap[v],
		})
	}

	return data, nil
}

func FilterMaptoDemandFilter(filterMap map[string][]string) DemandFilter {
	var filter DemandFilter
	log.Debugln(filterMap)

	filter.Range.fromFilterMap(filterMap, "range")
	filter.RegionID.fromFilterMap(filterMap, "region_id")
	filter.Aggregate.fromFilterMap(filterMap, "aggregate")

	return filter
}

func FilterMaptoRooftopFilter(filterMap map[string][]string) RooftopFilter {
	var filter RooftopFilter
	log.Debugln(filterMap)

	filter.Range.fromFilterMap(filterMap, "range")
	filter.RegionID.fromFilterMap(filterMap, "region_id")
	filter.Aggregate.fromFilterMap(filterMap, "aggregate")

	return filter
}

func FilterMapToGenerationFilter(filterMap map[string][]string) GeneratorFilter {
	var filter GeneratorFilter
	log.Debugln(filterMap)

	filter.Range.fromFilterMap(filterMap, "range")
	filter.DuID.fromFilterMap(filterMap, "duid")
	filter.Aggregate.fromFilterMap(filterMap, "aggregate")

	return filter
}

func FilterMapToGenerationGroupedFilter(filterMap map[string][]string) GeneratorGroupedFilter {
	var filter GeneratorGroupedFilter
	log.Debugln(filterMap)

	filter.Group.fromFilterMap(filterMap, "group")

	return filter
}

func (g *GeneratorGroupedFilter) GetAllGroupUnitCombinations(db *sql.DB) (map[string][]Unit, error) {
	var unit *Unit

	groupSet := make(map[string]struct{})
	groupedUnits := make(map[string][]Unit)
	groupedFilters := make(map[string]UnitFilter)

	allUnits, err := unit.ReadAll(db, UnitFilter{})
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error retrieving units: %v", err))
	}

	for _, group := range g.Group.GetEq() {
		if _, ok := groupSet[group]; ok {
			continue
		}
		groupSet[group] = struct{}{}

		switch group {
		case "region":
			regions, err := GetUniqueRegions(db)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error retrieving unique regions: %v", err))
			}
			if len(groupedUnits) > 0 && len(groupedFitlers) > 0 {
				// Need to Create a New Filter for Each Region to append to the already existing filters
				for _, filter := range groupedFilters {
				}

			} else {
				for _, region := range regions {
					groupFilter := UnitFilter{}
					groupFilter.RegionID.eq = []string{region}
					groupedFilters[region] = groupFilter
					units := make([]Unit, 0)
					for _, unit := range *allUnits {
						if unit.RegionID == region {
							units = append(units, unit)
						}
					}
					groupedUnits[region] = units
				}
			}
		case "fuel":
		case "technology":
		default:
			return nil, errors.New("unkown grouping")
		}
	}

	log.Debugln("Grouped Filters")
	log.Debugln(groupedFilters)
	log.Debugln("Grouped Units")
	log.Debugln(groupedUnits)

	return groupedUnits, nil
}
