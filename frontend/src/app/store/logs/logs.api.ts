import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { GetLatestLogsRequest, GetLatestLogsResponse } from "./logs.types";

export const logsApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/logs',
        prepareHeaders: (headers, { getState }) => {
            const token = localStorage.getItem('auth');
            if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
            }

            return headers;
        },
    }),
    reducerPath: 'api/logs',
    endpoints: (builder) => ({
        getLatestLogs: builder.query<GetLatestLogsResponse, GetLatestLogsRequest>({
            query: (data) => ({
                url: `/${data.containerID}`
            })
        }),
    }),
});

export const {
    useGetLatestLogsQuery,
} = logsApi;

export const {
    getLatestLogs,
} = logsApi.endpoints;