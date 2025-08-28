import { createSlice } from "@reduxjs/toolkit";

const initialState: number = 0;

export const metricsSlice = createSlice({
    name: "metrics",
    initialState,
    reducers: {
        setMetrics: (state, action) => {
            state = action.payload;
            return state;
        },
        unsetMetrics: (state) => {
            state = 0;
            return state;
        }
    }
});

export const { setMetrics, unsetMetrics } = metricsSlice.actions;
export const metricsReducer = metricsSlice.reducer;