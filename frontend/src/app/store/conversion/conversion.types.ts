import { Packet, Session } from "../../types/Service";

export type ConvertSessionToExploitRequest = {
    session: Session
    exportType: string
};

export type ConvertPacketToExploitRequest = {
    packet: Packet
    exportType: string
    servicePort: number
}

export type ConvertToExploitResponse = {
    exploit: string;
}

export type ExtractFilesRequest = {
    sessionID: string;
    packetNumber: number;
    packet: Packet;
}

export type ExtractFilesResponse = {
    data: string;
}