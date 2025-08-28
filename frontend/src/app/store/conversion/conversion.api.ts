import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { ConvertToExploitResponse, ConvertSessionToExploitRequest, ConvertPacketToExploitRequest, ExtractFilesRequest, ExtractFilesResponse } from "./conversion.types";

export const conversionApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/conversions',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem('auth');
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/conversion',
    endpoints: (builder) => ({
      convertSessionToExploit: builder.mutation<ConvertToExploitResponse, ConvertSessionToExploitRequest>({
        query: (data) => ({ url: '/sessions', method: 'POST', body: data }),
      }),
      convertPacketToExploit: builder.mutation<ConvertToExploitResponse, ConvertPacketToExploitRequest>({
        query: (data) => ({ url: '/packets', method: 'POST', body: data }),
      }),
      extractFilesFromPacket: builder.mutation<ExtractFilesResponse, ExtractFilesRequest>({
        query: (data) => ({
          url: '/files',
          method: 'POST',
          body: data
        })
      }),
    }),
});

export const {
  useConvertSessionToExploitMutation,
  useConvertPacketToExploitMutation,
  useExtractFilesFromPacketMutation
} = conversionApi;

export const {
  convertSessionToExploit,
  convertPacketToExploit,
  extractFilesFromPacket
} = conversionApi.endpoints;