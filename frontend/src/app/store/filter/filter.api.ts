import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { CreateFilterRequest, GetAllFiltersResponse, UpdateFilterRequest } from "./filter.types";

export const filterApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/filters',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem('auth');
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/filter',
    endpoints: (builder) => ({
      createFilter: builder.mutation<void, CreateFilterRequest>({
          query: (data) => ({ 
            url: ``,
            method: 'POST',
            body: data
          }),
      }),
      getAllFilters: builder.query<GetAllFiltersResponse, number>({
          query: (servicePort) => ({
              url: `?servicePort=${servicePort}`,
          })
      }),
      deleteFilter: builder.mutation<void, string>({
          query: (filterID: string) => ({ 
              url: `/${filterID}`,
              method: 'DELETE',
            }),
      }),
      updateFilter: builder.mutation<void, UpdateFilterRequest>({
        query: (data) => ({ 
            url: `/${data.id}${data.serviceID ? `?serviceID=${data.serviceID}` : ''}`,
            method: 'PATCH',
            body: data
          }),
      })
    }),
});

export const {
    useCreateFilterMutation,
    useDeleteFilterMutation,
    useGetAllFiltersQuery,
    useUpdateFilterMutation
} = filterApi;

export const {
    createFilter,
    deleteFilter,
    getAllFilters,
    updateFilter
} = filterApi.endpoints;