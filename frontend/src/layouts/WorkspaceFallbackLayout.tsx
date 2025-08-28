import { ReactNode, useEffect, useState } from "react"
import { useLocation } from "react-router"
import { useAppDispatch, useAppSelector } from "../app/store/hooks";
import { getWorkspace } from "../app/store/workspace/workspace.api";
import { setWorkspace } from "../app/store/workspace/workspace.slice";
import { putError } from "../app/store/error/error.slice";
import { Loading } from "../components/atoms/Loading";
import { NotFound } from "../pages/404/404";

export const WorkspaceFallbackLayout = (props:{
    children: ReactNode
}) => {
    const dispatch = useAppDispatch();
    const location = useLocation();
    const workspace = useAppSelector(state => state.rootReducer.workspace);
    const [isLoading, setIsLoading] = useState(true);
    const [isError, setIsError] = useState(false);
    const [getWorkspaceTrigger] = getWorkspace.useLazyQuery();

    useEffect(() => {
        const workspaceId = location.pathname.split("/")[2];
        if (workspace.id === workspaceId) {
            setIsLoading(false);
            return;
        }
        getWorkspaceTrigger(workspaceId)
            .unwrap()
            .then((res) => {
                dispatch(setWorkspace(res));
            })
            .catch((err) => {
                setIsError(true);
                dispatch(putError(err));
            })
            .finally(() => {
                setIsLoading(false);
            });
    }, [location.pathname]);

    return (
        <>
        {
            isLoading 
                ? <Loading hidden={false} size={40}/> 
                : isError
                ? <NotFound />
                : props.children
        }
        </>
    )
}