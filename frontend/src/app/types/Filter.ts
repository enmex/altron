export type Filter = {
    id: string;
    serviceId?: string;
    name: string;
    regex?: string;
    ttl?: number;
    totalPackets?: number; 
    inRequest: boolean;
    inResponse: boolean;
    color: string;
    isBlocking: boolean;
}

export type SessionFilter = {
    id: string;
    serviceId?: string;
    name: string;
    regex?: string;
    ttl?: number;
    totalPackets?: number; 
    inRequest: boolean;
    inResponse: boolean;
    color: string;
    isBlocking: boolean;
    matchesCount: number;
    matchedPackets: number[];
}