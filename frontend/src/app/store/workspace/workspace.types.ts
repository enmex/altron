import { Session } from "../../types/Service";
import { Workspace } from "../../types/Workspace";

export type WorkspaceState = {
    id: string;
    name: string;
    servicePort: number;
    status: string;
}

export type CreateWorkspaceRequest = {
    name: string;
    servicePort: number;
    timeout: string;
    startTime: number | null;
}

export type CreateWorkspaceResponse = Workspace;

export type GetWorkspaceResponse = Workspace;

export type UpdateWorkspaceRequest = {
    id: string;
    name: string;
}

export type GetAllWorkspacesResponse = {
    workspaces: Workspace[];
};

export type GetSessionsResponse = {
    sessions: Session[];
}

export type GetPaginatedSessionsRequest = {
    workspaceId: string;
    paginationIndex: number;
}

export type ClearWorkspaceSessionsRequest = {
    workspaceId: string;
}

export type SearchWorkspaceSessionsRequest = {
    workspaceId: string;
    filterId?: string;
    pagination: number;
    searchValue?: string;
    selectedCharacteristics?: {
        [k: string]: {
            value: string;
            selected: boolean;
            blocked: boolean;
        }[];
    };
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