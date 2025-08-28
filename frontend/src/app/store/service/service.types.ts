import { Container } from "../../types/Container";
import { Service } from "../../types/Service";

export type ServiceState = {
    id: string;
    name: string;
    port: number;
    link: string;
    containerID?: string;
    plugins: string[];
    workspaces: {
        id: string;
        name: string;
        servicePort: number;
        status: string;
    }[]
}

export type CreateServiceRequest = {
    name: string;
    link: string;
    port: number;
    plugins: string[];
};

export type CreateServiceResponse = Service;

export type UpdateServiceRequest = {
    id: string;
    name: string | null;
    link: string | null;
    port: number | null;
    containerID: string | null;
    plugins: string[];
};

export type GetServiceResponse = Service;

export type ScanServicesResponse = {
    services: {
        container?: {
            id: string;
            name: string;
            image: string;
        };
        port: number;
        name: string;
        isPublic: boolean;
    }[];
};

export type ScanServiceResponse = Service;

export type GetServiceLogsResponse = {
    logs: string[];
}

export type GetServicePluginsResponse = {
    plugins: string[];
}

export type GetContainersResponse = {
    containers: Container[];
}

export type ScanServicesRequest = {
    scope?: "containers";
}