import { createSlice } from "@reduxjs/toolkit";
import { WorkspaceState } from "./workspace.types";

const initialState: WorkspaceState = {
    id: "",
    name: "",
    servicePort: 0,
    status: "",
}

export const workspaceSlice = createSlice({
    name: "workspace",
    initialState,
    reducers: {
        setWorkspace: (state, action) => {
            state.id = action.payload.id;
            state.name = action.payload.name;
            return state;
        },
        unsetWorkspace: (state) => {
            state = {
                ...initialState,
            };
            return state;
        }
    }
});

export const { setWorkspace, unsetWorkspace } = workspaceSlice.actions;
export const workspaceReducer = workspaceSlice.reducer;