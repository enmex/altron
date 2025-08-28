export type Stats = {
    id: string;
    name: string;
    status: "up" | "mumble" | "down";
    stats?: {
        pidsStats?: {
            current: number;
            limit: number;
        },
        cpuUsage: number,
        memoryStats: {
            usage: number;
            maxUsage: number;
        },
        growthDynamics: {
            cpu: number;
            memory: number;
        }
        networks: Map<string, NetworkStats>
    }
}

export type NetworkStats = {
	rxBytes: number;
	rxPackets: number;
	txBytes: number;
	txPackets: number;
}