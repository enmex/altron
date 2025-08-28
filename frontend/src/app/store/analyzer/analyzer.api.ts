import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { GetAnalyzerComponentsResponse, GetAnalyzerPayloadRequest, GetAnalyzerPayloadResponse, GetWorkspaceAnalyzerPayloadRequest } from "./analyzer.types";

export const analyzerApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/session-analyzer',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem("auth");
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/analyzer',
    endpoints: (builder) => ({
      getAnalyzerPayload: builder.query<GetAnalyzerPayloadResponse, GetAnalyzerPayloadRequest>({
        query: (data) => ({
          url: `/services/${data.servicePort}`
        })
      }),
      getWorkspaceAnalyzerPayload: builder.query<GetAnalyzerPayloadResponse, GetWorkspaceAnalyzerPayloadRequest>({
        query: (data) => ({
          url: `/workspaces/${data.workspaceId}`
        })
      }),
      getAllComponents: builder.query<GetAnalyzerComponentsResponse, void>({
        query: () => ({
          url: '/components'
        })
      }),
    }),
});

export const {
  useGetAnalyzerPayloadQuery,
  useGetWorkspaceAnalyzerPayloadQuery,
  useGetAllComponentsQuery,
} = analyzerApi;

export const {
  getAnalyzerPayload,
  getWorkspaceAnalyzerPayload,
  getAllComponents
} = analyzerApi.endpoints;