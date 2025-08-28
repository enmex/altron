import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react"
import { SERVER_HOST } from "../../../config";
import { SignInRequest, LoginResponse as SignInResponse } from "./auth.types";

export const authApi = createApi({
    baseQuery: fetchBaseQuery({
        baseUrl: SERVER_HOST + '/auth',
    }),
    reducerPath: 'api/auth',
    endpoints: (builder) => ({
      signIn: builder.mutation<SignInResponse, SignInRequest>({
        query: (data) => ({ url: '/signIn', method: 'POST', body: data }),
      }),
    }),
});

export const {
  useSignInMutation,
} = authApi;

export const {
  signIn: login,
} = authApi.endpoints;