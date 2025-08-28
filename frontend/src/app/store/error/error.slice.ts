import { createSlice } from "@reduxjs/toolkit";
import { notifyError } from "../../notifications/notifier";

const initialState = '';

export const errorSlice = createSlice({
    name: "error",
    initialState,
    reducers: {
        putError: (state, action) => {
            if (!action.payload.includes('token') && state !== action.payload) {
                notifyError(action.payload);
            }
            return state;
        }
    }
});

export const { putError } = errorSlice.actions;
export const errorReducer = errorSlice.reducer;