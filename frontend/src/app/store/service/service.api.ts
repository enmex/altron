import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { CreateServiceRequest, CreateServiceResponse, GetServiceLogsResponse, GetServiceResponse, ScanServiceResponse, ScanServicesRequest, ScanServicesResponse, UpdateServiceRequest } from "./service.types";

export const serviceApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/services',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem('auth');
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/service',
    endpoints: (builder) => ({
      createService: builder.mutation<CreateServiceResponse, CreateServiceRequest>({
        query: (data) => ({ url: '', method: 'POST', body: data }),
      }),
      updateService: builder.mutation<void, UpdateServiceRequest>({
        query: (data) => ({ url: `/${data.id}`, method: 'PATCH', body: data }),
      }),
      getService: builder.query<GetServiceResponse, number>({
        query: (servicePort: number) => ({
            url: `/${servicePort}`,
        })
      }),
      deleteService: builder.mutation<void, string>({
        query: (serviceId) => ({ url: `/${serviceId}`, method: 'DELETE' }),
      }),
      scanServices: builder.query<ScanServicesResponse, ScanServicesRequest>({
        query: (data) => ({
          url: `/scan${data.scope ? `?scope=${data.scope}` : ''}`,
        })
      }),
      scanService: builder.query<ScanServiceResponse, number>({
        query: (port) => ({
          url: `/scan/${port}`
        })
      }),
      getServiceLogs: builder.query<GetServiceLogsResponse, number>({
        query: (servicePort: number) => ({
            url: `/${servicePort}/logs`,
        })
      }),
    }),
});

export const {
  useCreateServiceMutation,
  useUpdateServiceMutation,
  useDeleteServiceMutation,
  useScanServicesQuery,
  useScanServiceQuery,
  useGetServiceLogsQuery,
  useGetServiceQuery
} = serviceApi;

export const {
  createService,
  updateService,
  getService,
  deleteService,
  scanServices,
  scanService,
  getServiceLogs,
} = serviceApi.endpoints;