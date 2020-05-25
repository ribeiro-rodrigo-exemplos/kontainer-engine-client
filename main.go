package main

import (
	"context"
	"fmt"
	"kontainer-engine-client/types"
	"time"

	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
)

func getDriverOptions() *types.DriverOptions {
	options := &types.DriverOptions{
		BoolOptions: map[string]bool{
			"node-pool-autoscale": true,
		},
		StringOptions: map[string]string{
			"token":          "token-aqui",
			"display-name":   "cluster-test2",
			"name":           "my-cluster",
			"region-slug":    "nyc3",
			"version-slug":   "1.17.5-do.0",
			"node-pool-name": "node-pool-1",
			"node-pool-size": "s-2vcpu-2gb",
		},
		IntOptions: map[string]int64{
			"node-pool-min":   int64(2),
			"node-pool-max":   int64(4),
			"node-pool-count": int64(3),
		},
		StringSliceOptions: map[string]*types.StringSlice{
			"node-pool-labels": {Value: []string{"owner=rodrigo,projeto=dorancher"}},
		},
	}

	return options
}

func getClusterInfo() *types.ClusterInfo {
	return &types.ClusterInfo{}
}

func getCreateRequest() *types.CreateRequest {
	clusterInfo := getClusterInfo()
	options := getDriverOptions()

	return &types.CreateRequest{ClusterInfo: clusterInfo, DriverOptions: options}
}

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())

	defer conn.Close()

	if err != nil {
		log.Fatalf("Error in connection server %v", err)
	}

	client := types.NewDriverClient(conn)

	createRequest := getCreateRequest()

	clusterInfo, err := client.Create(context.Background(), createRequest)

	if err != nil {
		log.Fatalf("Error in create %v", err)
	}

	fmt.Println(clusterInfo)
	fmt.Println()

	clusterInfo, err = client.PostCheck(context.Background(), clusterInfo)

	if err != nil {
		log.Fatalf("Error in postcheck %v", err)
	}

	fmt.Println(clusterInfo)
	fmt.Println()

	time.Sleep(10 * time.Second)

	_, err = client.Remove(context.Background(), clusterInfo)

	if err != nil {
		log.Fatalf("Error in delete %v", err)
	}
}
