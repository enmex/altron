import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { GetPaginatedSessionsRequest, GetPcapWorkspaceResponse, GetSessionsResponse } from "./pcap-workspace.types";

export const pcapWorkspaceApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/pcap-workspaces',
        prepareHeaders: (headers, { getState }) => {
            const token = localStorage.getItem('auth');
            if (token !== '') {
                headers.set('Authorization', `Bearer ${token}`);
            }
        
            return headers;
        },
    }),
    reducerPath: 'api/pcapworkspace',
    endpoints: (builder) => ({
        getPaginatedPcapSessions: builder.query<GetSessionsResponse, GetPaginatedSessionsRequest>({
            query: (data) => ({
                url: `/${data.workspaceId}/sessions?pagination=${data.paginationIndex}`
            })
        }),
        getPcapWorkspace: builder.query<GetPcapWorkspaceResponse, string>({
            query: (pcapWorkspaceId) => ({
                url: `/${pcapWorkspaceId}`
            })
        }),
        deletePcapWorkspace: builder.mutation<void, string>({
            query: (pcapWorkspaceId) => ({ url: `/${pcapWorkspaceId}`, method: 'DELETE'}),
        }),
    }),
});

export const {
    useGetPaginatedPcapSessionsQuery,
    useGetPcapWorkspaceQuery,
    useDeletePcapWorkspaceMutation
} = pcapWorkspaceApi;

export const {
    getPaginatedPcapSessions,
    getPcapWorkspace,
    deletePcapWorkspace
} = pcapWorkspaceApi.endpoints;