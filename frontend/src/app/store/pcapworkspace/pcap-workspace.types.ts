import { Session } from "../../types/Service";

export type GetSessionsResponse = {
    sessions: Session[];
}

export type GetPaginatedSessionsRequest = {
    workspaceId: string;
    paginationIndex: number;
}

export type GetPcapWorkspaceResponse = {
    id: string;
    filename: string;
    status: string;
}