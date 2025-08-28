import { Characteristic } from "./Analyzer";
import { SessionFilter } from "./Filter";
import { Workspace } from "./Workspace";

export type Packet = {
    sentAt: number;
    payload: string;
    isRequest: boolean;
}

export type Session = {
    id: string;
    serviceID: string;
    sentAt: number;
    ttl: number;
    clientHost: string;
    serverPort: number;
    protocol: string;
    packets: Packet[];
    packetsCount: number;
    matchedFilters: SessionFilter[];
    clientUserAgent?: string;
    analyzerMatches?: {
        [componentName: string]: Characteristic;
    };
    isSafe?: boolean;
}

export type Service = {
    id: string;
    name: string;
    link: string;
    port: number;
    containerID?: string;
    plugins: string[];
    workspaces: Workspace[];
}