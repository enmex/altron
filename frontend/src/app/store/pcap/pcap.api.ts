import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { UploadPcapResponse } from "./pcap.types";

export const pcapApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/pcaps',
        prepareHeaders: (headers, { getState }) => {
            const token = localStorage.getItem('auth');
            if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
            }

            return headers;
        },
    }),
    reducerPath: 'api/pcap',
    endpoints: (builder) => ({
        uploadPcap: builder.mutation<UploadPcapResponse, FormData>({
            query: (data) => ({ 
                url: '',
                method: 'POST',
                body: data,
            }),
        }),
    }),
});

export const {
    useUploadPcapMutation,
} = pcapApi;

export const {
    uploadPcap,
} = pcapApi.endpoints;