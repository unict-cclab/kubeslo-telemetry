package metrics

import (
	"context"
	"fmt"
	"os"
	"time"

	prometheusapi "github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var prometheusAddress = os.Getenv("PROMETHEUS_ADDRESS")

func newPrometheusClient(serverAddress string) (prometheus.API, error) {
	client, err := prometheusapi.NewClient(prometheusapi.Config{
		Address: serverAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating metrics client: %v", err)
	}
	return prometheus.NewAPI(client), nil
}

func GetAppRequestsPerSecond(appGroupName, appName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		rate(istio_requests_total{app_group="`+appGroupName+`", app="`+appName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetAppsRequestsPerSecond(appGroupName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		rate(istio_requests_total{reporter="source", app_group="`+appGroupName+`", source_app!="unknown", destination_app!="unknown"}[`+rangeWidth+`])
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodeLatencies(nodeName, rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum{origin_node="`+nodeName+`"}[`+rangeWidth+`]) / rate(node_latency_count{origin_node="`+nodeName+`"}[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}

func GetNodesLatencies(rangeWidth string) (model.Vector, prometheus.Warnings, error) {
	prometheusClient, err := newPrometheusClient(prometheusAddress)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := prometheusClient.Query(ctx, `
		(rate(node_latency_sum[`+rangeWidth+`]) / rate(node_latency_count[`+rangeWidth+`])) * 1000
	`, time.Now())

	if err != nil {
		return nil, nil, fmt.Errorf("error during query execution: %v", err)
	}

	vector, ok := result.(model.Vector)

	if !ok {
		return nil, nil, fmt.Errorf("query result is not a vector: %v", err)
	}

	return vector, warnings, err
}
