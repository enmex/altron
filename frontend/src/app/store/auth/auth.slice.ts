import { createSlice } from "@reduxjs/toolkit";
import { AuthState } from "./auth.types";

const initialState: AuthState = {
    token: "",
}

export const authSlice = createSlice({
    name: "auth",
    initialState,
    reducers: {
        setAuth: (state, action) => {
            state.token = action.payload.token;
            localStorage.setItem("auth", action.payload.token);
            return state;
        },
        clearAuth: (state) => {
            state.token = "";
            localStorage.removeItem("auth");
            return state;
        }
    }
});

export const { setAuth, clearAuth } = authSlice.actions;
export const authReducer = authSlice.reducer;