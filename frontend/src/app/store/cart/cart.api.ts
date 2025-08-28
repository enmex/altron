import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { AddSessionsToCartRequest, GetCartRequest, GetCartResponse, MergeSessionsRequest, MergeSessionsResponse, RemoveSessionFromCartRequest } from "./cart.types";

export const cartApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/carts',
        prepareHeaders: (headers, { getState }) => {
            const token = localStorage.getItem('auth');
            if (token !== '') {
                headers.set('Authorization', `Bearer ${token}`);
            }
        
            return headers;
        },
    }),
    reducerPath: 'api/cart',
    endpoints: (builder) => ({
        addWorkspaceSessionsToCart: builder.mutation<void, AddSessionsToCartRequest>({
            query: (data) => ({
                url: `/${data.workspaceId}`,
                body: data,
                method: 'POST',
            })
        }),
        removeWorkspaceSessionFromCart: builder.mutation<void, RemoveSessionFromCartRequest>({
            query: (data) => ({
                url: `/${data.workspaceId}/${data.sessionId}`,
                method: 'DELETE'
            })
        }),
        getWorkspaceCart: builder.query<GetCartResponse, GetCartRequest>({
            query: (data) => ({
                url: `/${data.workspaceID}?pagination=${data.pagination}`
            })
        }),
        clearWorkspaceCart: builder.mutation<void, string>({
            query: (workspaceID) => ({
                url: `/${workspaceID}`,
                method: 'DELETE'
            })
        }),
        mergeSessions: builder.mutation<MergeSessionsResponse, MergeSessionsRequest>({
            query: (data) => ({
                url: `/${data.workspaceID}/sessions/merge`,
                body: data,
                method: 'POST'
            })
        }),
    }),
});

export const {
    useGetWorkspaceCartQuery,
    useAddWorkspaceSessionsToCartMutation,
    useRemoveWorkspaceSessionFromCartMutation,
    useClearWorkspaceCartMutation,
    useMergeSessionsMutation,
} = cartApi;

export const {
    addWorkspaceSessionsToCart,
    removeWorkspaceSessionFromCart,
    getWorkspaceCart,
    clearWorkspaceCart,
    mergeSessions,
} = cartApi.endpoints;