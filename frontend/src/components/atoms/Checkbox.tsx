import { ReactNode, useState } from "react";
import { useAppSelector } from "../../app/store/hooks";
import { darkenColor } from "../../utils/utils";

export const Checkbox = (props: {
    children: ReactNode;
    onChange: () => void;
    checked: boolean;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [hovered, setHovered] = useState(false);
    
    return (
        <button
            type="button"
            className="cursor-pointer border-2 duration-200 m-2 p-2"
            style={{
                backgroundColor: props.checked ? theme.accents.neutral : "transparent",
                borderColor: hovered ? darkenColor(theme.accents.neutral, 0.7) : theme.accents.neutral
            }}
            onMouseEnter={() => setHovered(true)}
            onMouseLeave={() => setHovered(false)}
            onClick={props.onChange}
        >
            {props.children}
        </button>
    );
}