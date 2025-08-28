import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { ClearWorkspaceSessionsRequest, CreateWorkspaceRequest, CreateWorkspaceResponse, GetPaginatedSessionsRequest, GetSessionsResponse, GetWorkspaceResponse, SearchWorkspaceSessionsRequest, UpdateWorkspaceRequest } from "./workspace.types";

export const workspaceApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/workspaces',
        prepareHeaders: (headers, { getState }) => {
            const token = localStorage.getItem('auth');
            if (token !== '') {
                headers.set('Authorization', `Bearer ${token}`);
            }
        
            return headers;
        },
    }),
    reducerPath: 'api/workspace',
    endpoints: (builder) => ({
        createWorkspace: builder.mutation<CreateWorkspaceResponse, CreateWorkspaceRequest>({
            query: (data) => ({ url: '', method: 'POST', body: data }),
        }),
        getWorkspace: builder.query<GetWorkspaceResponse, string>({
            query: (workspaceId) => ({
                url: `/${workspaceId}`
            }),
        }),
        updateWorkspace: builder.mutation<void, UpdateWorkspaceRequest>({
            query: (data) => ({ url: `/${data.id}`, method: 'PATCH', body: data }),
        }),
        deleteWorkspace: builder.mutation<void, string>({
            query: (workspaceId) => ({ url: `/${workspaceId}`, method: 'DELETE'}),
        }),
        getPaginatedSessions: builder.query<GetSessionsResponse, GetPaginatedSessionsRequest>({
            query: (data) => ({ url: `/${data.workspaceId}/sessions?pagination=${data.paginationIndex}` })
        }),
        clearWorkspaceSessions: builder.mutation<void, ClearWorkspaceSessionsRequest>({
            query: (data) => ({ url: `/${data.workspaceId}/sessions`, method: 'DELETE', body: data }), 
        }),
        searchWorkspaceSessions: builder.mutation<GetSessionsResponse, SearchWorkspaceSessionsRequest>({
            query: (data) => ({ 
                url: `/${data.workspaceId}/sessions/search?pagination=${data.pagination}`,
                body: data,
                method: 'POST',
            }) 
        }),
    }),
});

export const {
    useCreateWorkspaceMutation,
    useGetWorkspaceQuery,
    useUpdateWorkspaceMutation,
    useDeleteWorkspaceMutation,
    useGetPaginatedSessionsQuery,
    useClearWorkspaceSessionsMutation,
    useSearchWorkspaceSessionsMutation,
} = workspaceApi;

export const {
    createWorkspace,
    getWorkspace,
    updateWorkspace,
    deleteWorkspace,
    getPaginatedSessions,
    clearWorkspaceSessions,
    searchWorkspaceSessions,
} = workspaceApi.endpoints;