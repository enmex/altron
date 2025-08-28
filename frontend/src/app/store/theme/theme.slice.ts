import { THEMES } from './../../../config/themes';
import { createSlice } from "@reduxjs/toolkit";
import { ThemeState } from "./theme.types";

const initialState: ThemeState = {
    ...THEMES.find(t => t.name === localStorage.getItem("theme")) ?? THEMES[0],
}

export const themeSlice = createSlice({
    name: "theme",
    initialState,
    reducers: {
        setTheme: (state, action) => {
            state = {
                ...action.payload
            }
            localStorage.setItem("theme", action.payload.name);
            return state;
        },
    }
});

export const { setTheme } = themeSlice.actions;
export const themeReducer = themeSlice.reducer;