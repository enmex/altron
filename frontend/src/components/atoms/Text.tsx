import { ReactNode } from "react";
import { useAppSelector } from "../../app/store/hooks";

export const Text = (props: {
    className?: string;
    children: ReactNode;
    color?: string;
    backgroundColor?: string;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const className = props.className + " duration-200";

    return (
        <span
            className={className}
            style={{
                color: props.color ?? theme.text,
                backgroundColor: props.backgroundColor
            }}
        >
            {props.children}
        </span>
    );
}