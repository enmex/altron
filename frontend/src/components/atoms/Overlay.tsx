import { ReactNode } from "react"

export const Overlay = (props: {
    children: ReactNode
}) => {
    return (
        <div className="fixed inset-0 flex items-center justify-center z-50 bg-opacity-50 backdrop-filter backdrop-blur-sm">
            {props.children}
        </div>
    );
}