const DEV_MODE = false

export const SERVER_HOST = DEV_MODE ? `http://192.168.0.132:64001/api` : `http://${window.location.host}/api`;
export const WEBSOCKET_SERVER_HOST = DEV_MODE ? `ws://192.168.0.132:64005` : `ws://${window.location.host}`;

export const MAX_SESSIONS_IN_CACHE = 500;
export const MAX_LOGS_IN_CACHE = 1000;