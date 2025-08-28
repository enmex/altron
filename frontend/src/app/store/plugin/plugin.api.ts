import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { GetAllPluginsResponse } from "./plugin.types";

export const pluginApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/plugins',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem('auth');
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/plugin',
    endpoints: (builder) => ({
      getAllPlugins: builder.query<GetAllPluginsResponse, void>({
        query: () => ({ 
            url: ''
        }),
      }),
    }),
});

export const {
    useGetAllPluginsQuery,
} = pluginApi;

export const {
    getAllPlugins,
} = pluginApi.endpoints;