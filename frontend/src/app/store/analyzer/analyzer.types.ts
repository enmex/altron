import { Characteristic } from "../../types/Analyzer";

export type GetAnalyzerPayloadRequest = {
    servicePort?: number;
    workspaceId?: string;
}

export type GetWorkspaceAnalyzerPayloadRequest = {
    workspaceId: string;
}

export type GetAnalyzerPayloadResponse = {
    hasChecker: boolean;
    analyzer: {
        [componentName: string]: Characteristic[]
    };
};

export type GetAnalyzerComponentsResponse = {
    components: string[];
}