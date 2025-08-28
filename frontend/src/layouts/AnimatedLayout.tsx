import { url } from "inspector";
import { ReactNode } from "react"
import { useAppSelector } from "../app/store/hooks";

export const AnimatedLayout = (props: {
    children?: ReactNode
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);

    return (
        <div
            className="bg-cover w-full h-full"
            style={{
                backgroundImage: `url(${theme.webpPath})`
            }}
        >
            {props.children}
        </div>
    )
}