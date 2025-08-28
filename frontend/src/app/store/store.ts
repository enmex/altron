import { combineReducers, configureStore } from "@reduxjs/toolkit";
import { authApi } from "./auth/auth.api";
import { authReducer } from "./auth/auth.slice";
import { serviceReducer } from "./service/service.slice";
import { serviceApi } from "./service/service.api";
import { workspaceApi } from "./workspace/workspace.api";
import { workspaceReducer } from "./workspace/workspace.slice";
import { dashboardApi } from "./dashboard/dashboard.api";
import { pluginApi } from "./plugin/plugin.api";
import { languageReducer } from "./language/language.slice";
import { filterApi } from "./filter/filter.api";
import { themeReducer } from "./theme/theme.slice";
import { conversionApi } from "./conversion/conversion.api";
import { analyzerApi } from "./analyzer/analyzer.api";
import { sessionApi } from "./session/session.api";
import { metricsReducer } from "./metrics/metrics.slice";
import { pcapReducer } from "./pcap/pcap.slice";
import { pcapApi } from "./pcap/pcap.api";
import errorListenerMiddleware from "./middleware";
import { errorReducer } from "./error/error.slice";
import { pcapWorkspaceApi } from './pcapworkspace/pcap-workspace.api';
import { cartApi } from './cart/cart.api';
import { logsApi } from './logs/logs.api';

const rootReducer = combineReducers({
    auth: authReducer,
    service: serviceReducer,
    workspace: workspaceReducer,
    language: languageReducer,
    theme: themeReducer,
    metrics: metricsReducer,
    pcap: pcapReducer,
    error: errorReducer
});

export const store = configureStore({
    reducer: {
        rootReducer,
        [authApi.reducerPath]: authApi.reducer,
        [serviceApi.reducerPath]: serviceApi.reducer,
        [workspaceApi.reducerPath]: workspaceApi.reducer,
        [dashboardApi.reducerPath]: dashboardApi.reducer,
        [pluginApi.reducerPath]: pluginApi.reducer,
        [filterApi.reducerPath]: filterApi.reducer,
        [conversionApi.reducerPath]: conversionApi.reducer,
        [analyzerApi.reducerPath]: analyzerApi.reducer,
        [sessionApi.reducerPath]: sessionApi.reducer,
        [pcapApi.reducerPath]: pcapApi.reducer,
        [pcapWorkspaceApi.reducerPath]: pcapWorkspaceApi.reducer,
        [cartApi.reducerPath]: cartApi.reducer,
        [logsApi.reducerPath]: logsApi.reducer,
    },

    middleware: (middleware) => 
        middleware({
            serializableCheck: false,
        }).prepend(errorListenerMiddleware.middleware)
            .concat(authApi.middleware)
            .concat(serviceApi.middleware)
            .concat(workspaceApi.middleware)
            .concat(dashboardApi.middleware)
            .concat(pluginApi.middleware)
            .concat(filterApi.middleware)
            .concat(conversionApi.middleware)
            .concat(analyzerApi.middleware)
            .concat(sessionApi.middleware)
            .concat(pcapApi.middleware)
            .concat(pcapWorkspaceApi.middleware)
            .concat(cartApi.middleware)
            .concat(logsApi.middleware)
});

export type RootState = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;