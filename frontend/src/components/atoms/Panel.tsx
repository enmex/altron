import { ReactNode } from "react";
import { useAppSelector } from "../../app/store/hooks";

export const Panel = (props: {
    children: ReactNode
    withBorder?: boolean
    className?: string
    color?: string
    borderColor?: string
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);

    return (
        <div 
            className={props.className ?? "animate-fade flex flex-col items-center p-4 m-2 border-2 rounded-md"}
            style={{
                backgroundColor: props.color ?? theme.secondary,
                borderColor: props.withBorder ? props.borderColor ?? theme.accents.neutral : "transparent"
            }}
        >
            {props.children}
        </div>
    );
}