import { Filter } from "../../types/Filter";

export type CreateFilterRequest = {
    name: string;
    regex?: string;
    ttl?: number;
    inRequest: boolean;
    inResponse: boolean;
    isBlocking: boolean;
    serviceID?: string;
}

export type FiltersResponse = Filter[];

export type GetAllFiltersResponse = {
    filters: FiltersResponse;
}

export type UpdateFilterRequest ={
    id: string;
    name?: string;
    regex?: string;
    ttl?: number;
    inRequest?: boolean;
    inResponse?: boolean;
    serviceID?: string;
}