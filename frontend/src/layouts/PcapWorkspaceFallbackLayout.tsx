import { ReactNode, useEffect, useState } from "react";
import { useAppDispatch, useAppSelector } from "../app/store/hooks"
import { getPcapWorkspace } from "../app/store/pcapworkspace/pcap-workspace.api";
import { useLocation } from "react-router";
import { setPcapWorkspace } from "../app/store/pcap/pcap.slice";
import { putError } from "../app/store/error/error.slice";
import { Loading } from "../components/atoms/Loading";
import { NotFound } from "../pages/404/404";

export const PcapWorkspaceFallbackLayout = (props:{
    children: ReactNode
}) => {
    const dispatch = useAppDispatch();
    const location = useLocation();
    const pcapWorkspace = useAppSelector(state => state.rootReducer.pcap);
    const [isLoading, setIsLoading] = useState(true);
    const [isError, setIsError] = useState(false);
    const [getPcapWorkspaceTrigger] = getPcapWorkspace.useLazyQuery();

    useEffect(() => {
        const pcapWorkspaceId = location.pathname.split("/")[2];
        if (pcapWorkspace.id === pcapWorkspaceId) {
            setIsLoading(false);
            return;
        }
        getPcapWorkspaceTrigger(pcapWorkspaceId).
            unwrap().then((res) => {
                dispatch(setPcapWorkspace(res));
            }).catch((err) => {
                setIsError(true);
                dispatch(putError(err.data.message));
            }).finally(() => {
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