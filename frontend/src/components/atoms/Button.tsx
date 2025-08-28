import { ReactNode, useState } from "react";
import { darkenColor } from "../../utils/utils";

export const Button = (props: {
    type?: "button" | "submit" | null,
    children: ReactNode,
    className?: string,
    borderColor?: string,
    backgroundColor?: string,
    onMouseEnter?: () => void,
    onMouseLeave?: () => void,
    onClick?: () => void,
}) => {
    const [hovered, setHovered] = useState(false);

    const onMouseEnter = () => {
        setHovered(true);
        if (props.onMouseEnter) {
            props.onMouseEnter();
        }
    }

    const onMouseLeave = () => {
        setHovered(false);
        if (props.onMouseLeave) {
            props.onMouseLeave();
        }
    }

    return (
        <button 
            type={props.type ? props.type : "button"} 
            onClick={props.onClick} 
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
            className={props.className ?? "flex items-center duration-200 mx-2"}
            style={{
                backgroundColor: props.backgroundColor ? hovered ? darkenColor(props.backgroundColor, 0.7) : props.backgroundColor : "transparent",
                borderColor: props.borderColor ? hovered ? darkenColor(props.borderColor, 0.7) : props.borderColor : "transparent",
            }}
        >{ props.children }</button>
    );
}