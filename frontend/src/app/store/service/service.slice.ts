import { createSlice } from "@reduxjs/toolkit";
import { ServiceState } from "./service.types";

const initialState: ServiceState = {
    id: "",
    name: "",
    link: "",
    port: 0,
    plugins: [],
    workspaces: []
}

export const serviceSlice = createSlice({
    name: "service",
    initialState,
    reducers: {
        setService: (state, action) => {
            state = {
                ...action.payload
            }
            return state;
        },
        unsetService: (state) => {
            state = {
                ...initialState
            };
            return state;
        },
        setContainerID: (state, action) => {
            state.containerID = action.payload;
            return state;
        },
        putWorkspace: (state, action) => {
            state = {
                ...state,
                workspaces: [...state.workspaces, action.payload]
            }
            return state;
        }  
    }
});

export const { setService, unsetService, setContainerID, putWorkspace } = serviceSlice.actions;
export const serviceReducer = serviceSlice.reducer;