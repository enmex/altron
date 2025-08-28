export const INDEX_PATH = "/home";
export const SIGN_IN_PATH = "/signIn";
export const SERVICE_PATH = "/services/*";
export const WORKSPACE_PATH = "/workspaces/*";
export const SESSION_PATH = "/sessions/*";
export const LOGS_PATH = "/logs/*";
export const ALTRON_LOGS_PATH = "/logs/altron/*";
export const HEALTH_PATH = "/health";
export const CART_PATH = "/carts/*";
export const PCAP_PATH = "/pcaps/*";

export const conversionTypes: {
    text: string,
    type: string,
    inRequest: boolean,
    inResponse: boolean,
}[] = [
    {
        text: "as ascii bytes",
        type: "ascii",
        inRequest: true,
        inResponse: false,
    },
    {
        text: "as python bytes",
        type: "python_bytes",
        inRequest: true,
        inResponse: false,
    },
    {
        text: "as pwntools",
        type: "pwntools",
        inRequest: true,
        inResponse: false,
    },
    {
        text: "as requests",
        type: "requests",
        inRequest: true,
        inResponse: false,
    },
    {
        text: "as curl",
        type: "curl",
        inRequest: true,
        inResponse: false,
    },
    // {
    //     text: "as plugin",
    //     type: "plugin",
    //     inRequest: true,
    //     inResponse: false,
    // },
    {
        text: "export files",
        type: "files",
        inRequest: true,
        inResponse: true,
    }
];

export const CONTAINERS = [
    "agent",
    "core",
    "session",
    "connection",
    "plugin",
    "converter",
    "frontend",
    "database",
	"redis",
	"rabbitmq",
    "ftp"
];