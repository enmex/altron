import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { CreateSessionRequest, GetSessionResponse } from "./session.types";

export const sessionApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/sessions',
        prepareHeaders: (headers, { getState }) => {
          const token = localStorage.getItem('auth');
          if (token !== '') {
            headers.set('Authorization', `Bearer ${token}`);
          }
    
          return headers;
        },
    }),
    reducerPath: 'api/session',
    endpoints: (builder) => ({
      createSession: builder.mutation<void, CreateSessionRequest>({
        query: (data) => ({ url: '', method: 'POST', body: data }),
      }),
      getSession: builder.query<GetSessionResponse, string>({
        query: (sessionId: string) => ({
            url: `/${sessionId}`,
        })
      }),
    }),
});

export const {
    useCreateSessionMutation,
    useGetSessionQuery,
} = sessionApi;

export const {
  createSession,
  getSession,
} = sessionApi.endpoints;