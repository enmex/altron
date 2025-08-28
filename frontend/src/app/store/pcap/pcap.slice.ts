import { createSlice } from "@reduxjs/toolkit";

const initialState: {
    id: string;
    fileName: string;
    status: string;
} = {
    id: "",
    fileName: "",
    status: "",
};

export const pcapSlice = createSlice({
    name: "pcap",
    initialState,
    reducers: {
        setPcapWorkspace: (state, action) => {
            state = {
                ...action.payload
            };
            return state;
        },
        unsetPcapWorkspace: (state) => {
            state = {
                id: "",
                fileName: "",
                status: "",
            }
            return state;
        }
    }
});

export const {
    setPcapWorkspace,
    unsetPcapWorkspace
} = pcapSlice.actions;

export const pcapReducer = pcapSlice.reducer;