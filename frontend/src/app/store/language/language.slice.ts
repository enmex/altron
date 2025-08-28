import { createSlice } from "@reduxjs/toolkit";

const initialState = 'en';

export const languageSlice = createSlice({
    name: "language",
    initialState,
    reducers: {
        setLanguage: (state, action) => {
            state = action.payload;
            return state;
        }
    }
});

export const { setLanguage } = languageSlice.actions;
export const languageReducer = languageSlice.reducer;