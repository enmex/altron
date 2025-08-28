export type GetLatestLogsRequest = {
    containerID: string;
}

export type GetLatestLogsResponse = {
    logs: string[];
}