import { createListenerMiddleware } from "@reduxjs/toolkit";
import { putError } from "./error/error.slice";
import { clearAuth } from "./auth/auth.slice";

const errorListenerMiddleware = createListenerMiddleware();

errorListenerMiddleware.startListening({
    actionCreator: putError,
    effect: (action, listenerApi) => {
        const error = action.payload;
        if (error === "invalid security token") {
            listenerApi.dispatch(clearAuth());
        }
    }
});

export default errorListenerMiddleware;