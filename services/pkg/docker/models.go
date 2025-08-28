package docker

import "github.com/docker/docker/api/types"

type DockerStats struct {
	PidsStats   types.PidsStats               `json:"pids_stats"`
	CPUStats    types.CPUStats                `json:"cpu_stats"`
	MemoryStats types.MemoryStats             `json:"memory_stats"`
	Networks    map[string]types.NetworkStats `json:"networks"`
}

type MemoryStats struct {
	Usage    uint64 `json:"usage"`
	MaxUsage uint64 `json:"maxUsage"`
}

type NetworkStats struct {
	RxBytes   uint64 `json:"rxBytes"`
	RxPackets uint64 `json:"rxPackets"`
	TxBytes   uint64 `json:"txBytes"`
	TxPackets uint64 `json:"txPackets"`
}

// -1 means down
// 0 means no changes
// 1 means up
type GrowthDynamics struct {
	CPU    int8 `json:"cpu"`
	Memory int8 `json:"memory"`
}

type Stats struct {
	PidsStats      types.PidsStats         `json:"pidsStats"`
	CPUUsage       float64                 `json:"cpuUsage"`
	MemoryStats    MemoryStats             `json:"memoryStats"`
	Networks       map[string]NetworkStats `json:"networks"`
	GrowthDynamics GrowthDynamics          `json:"growthDynamics"`
}

type ContainerHealth struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Stats  Stats  `json:"stats"`
}
