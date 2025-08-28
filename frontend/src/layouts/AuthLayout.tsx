import { ReactNode, useEffect } from "react"
import { INDEX_PATH } from "../config/constants";
import { useAppNavigation } from "../hooks/navigate";

export const AuthLayout = (props:{
    children: ReactNode
}) => {
    const navigate = useAppNavigation();
    
    useEffect(() => {
        const token = localStorage.getItem('auth');
        if (token && token.length > 0 && token !== 'null') {
            navigate(INDEX_PATH);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <>
        {props.children}
        </>
    )
}