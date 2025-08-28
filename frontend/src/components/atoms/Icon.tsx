import { useState } from "react";
import { useAppSelector } from "../../app/store/hooks";
import { IconType } from "react-icons";
import { Icons as IconTypes, icons } from "../../config/icons";
import { Tooltip } from "./Tooltip";

export const Icon = (props: {
    type?: "positive" | "negative" | "neutral" | "contrast";
    color?: string;
    name: keyof IconTypes;
    size: number;
    tip?: string;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [hovered, setHovered] = useState(false);

    const IconComponent: IconType = icons[props.name];

    let hoverColor;
    switch (props.type) {
        case "contrast" : {
            hoverColor = theme.accents.contrast;
            break;
        }
        case "neutral" : {
            hoverColor = theme.accents.neutral;
            break;
        }
        case "positive" : {
            hoverColor = theme.accents.positive;
            break;
        }
        case "negative" : {
            hoverColor = theme.accents.negative;
            break;
        }
        default: {
            hoverColor = props.color ?? theme.text;
            break;
        }
    }

    return (
        <Tooltip
            tip={props.tip}
        >
            <IconComponent 
                color={hovered ? hoverColor : props.color ?? theme.text}
                size={props.size}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
            />
        </Tooltip>
    );
}