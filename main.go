package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unict-cclab/kubeslo-telemetry/metrics"
)

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
	router.GET("/metrics/nodes/latencies", getNodesLatencies)

	err := router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Printf("Exiting because of error: %s", err.Error())
		return
	}
}
