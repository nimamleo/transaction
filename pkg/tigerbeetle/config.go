package tigerbeetle

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ClusterID uint64
	Port      string
}

func LoadConfig() Config {
	clusterIDStr := getEnv("TIGERBEETLE_CLUSTER_ID")
	port := getEnv("TIGERBEETLE_PORT")

	clusterID, err := strconv.ParseUint(clusterIDStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("invalid TIGERBEETLE_CLUSTER_ID value: %s", clusterIDStr))
	}

	return Config{
		ClusterID: clusterID,
		Port:      port,
	}
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s is required but not set", key))
	}
	return value
}
