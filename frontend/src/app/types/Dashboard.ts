import { PcapWorkspace } from "./PcapWorkspace";
import { Service } from "./Service";

export type Dashboard = {
    services: Service[];
    pcapWorkspaces: PcapWorkspace[];
}