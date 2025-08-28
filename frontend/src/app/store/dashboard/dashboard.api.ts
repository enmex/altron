import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { GetDashboardResponse } from "./dashboard.types";

export const dashboardApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/dashboard',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem('auth');
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/dashboard',
    endpoints: (builder) => ({
      getDashboard: builder.query<GetDashboardResponse, void>({
        query: () => ({ 
          url: ''
        }),
      }),
    }),
});

export const {
    useGetDashboardQuery,
} = dashboardApi;

export const {
    getDashboard,
} = dashboardApi.endpoints;