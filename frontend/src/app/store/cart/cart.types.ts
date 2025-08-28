import { Session } from "../../types/Service";

export type GetSessionsResponse = {
    sessions: Session[];
}

export type AddSessionsToCartRequest = {
    workspaceId: string;
    sessions: string[];
}

export type RemoveSessionFromCartRequest = {
    workspaceId: string;
    sessionId: string;
}

export type GetCartResponse = GetSessionsResponse;

export type GetCartRequest = {
    workspaceID: string;
    pagination: number;
}

export type MergeSessionsRequest = {
    workspaceID: string;
    sessions: string[];
}

export type MergeSessionsResponse = {
    sessionID: string;
}