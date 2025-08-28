import { Icons } from "../../config/icons";

export type Item = {
    text?: string;
    icon?: {
        value: keyof Icons,
        color: string;
    };
    onItem: () => void;
    onContextMenu?: () => void;
    children?: Item[];
    color?: string;
}