package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/model"
	"github.com/unict-cclab/kubeslo-telemetry/metrics"
)

func getAppsRequestsPerSecond(c *gin.Context) {
	appGroupName := c.Query("app-group")
	appName := c.Query("app")
	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if appName != "" {
		var results model.Vector
		var err error

		results, _, err = metrics.GetAppRequestsPerSecond(appGroupName, appName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			rpsValues := map[string]float64{}
			for _, result := range results {
				if string(result.Metric["source_app"]) == appName {
					rpsValues[string(result.Metric["destination_app"])] = float64(result.Value)
				} else if string(result.Metric["destination_app"]) == appName {
					rpsValues[string(result.Metric["source_app"])] = float64(result.Value)
				}
			}
			c.IndentedJSON(http.StatusOK, rpsValues)
		}
	} else {
		results, _, err := metrics.GetAppsRequestsPerSecond(appGroupName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			rpsValues := map[string]map[string]float64{}
			for _, result := range results {
				_, ok := rpsValues[string(result.Metric["source_app"])]
				if !ok {
					rpsValues[string(result.Metric["source_app"])] = map[string]float64{}
				}
				rpsValues[string(result.Metric["source_app"])][string(result.Metric["destination_app"])] = float64(result.Value)

				_, ok = rpsValues[string(result.Metric["destination_app"])]
				if !ok {
					rpsValues[string(result.Metric["destination_app"])] = map[string]float64{}
				}
				rpsValues[string(result.Metric["destination_app"])][string(result.Metric["source_app"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, rpsValues)
		}
	}
}

func getNodesLatencies(c *gin.Context) {
	nodeName := c.Query("node")

	rangeWidth := c.Query("range-width")

	if rangeWidth == "" {
		rangeWidth = "5m"
	}

	if nodeName != "" {
		latencyValues := map[string]float64{}
		results, _, err := metrics.GetNodeLatencies(nodeName, rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			for _, result := range results {
				latencyValues[string(result.Metric["destination_node"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, latencyValues)
		}
	} else {
		latencyValues := map[string]map[string]float64{}
		results, _, err := metrics.GetNodesLatencies(rangeWidth)

		//fmt.Println(warnings)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			for _, result := range results {
				_, ok := latencyValues[string(result.Metric["origin_node"])]
				if !ok {
					latencyValues[string(result.Metric["origin_node"])] = map[string]float64{}
				}
				latencyValues[string(result.Metric["origin_node"])][string(result.Metric["destination_node"])] = float64(result.Value)
			}
			c.IndentedJSON(http.StatusOK, latencyValues)
		}
	}
}

func main() {
	router := gin.Default()
	router.GET("/metrics/apps/rps", getAppsRequestsPerSecond)
	router.GET("/metrics/nodes/latencies", getNodesLatencies)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Printf("Exiting because of error: %s", err.Error())
		return
	}
}
