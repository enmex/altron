import { ReactNode, useEffect, useState } from "react"
import { getService } from "../app/store/service/service.api";
import { useAppDispatch, useAppSelector } from "../app/store/hooks";
import { useLocation } from "react-router";
import { putError } from "../app/store/error/error.slice";
import { setService } from "../app/store/service/service.slice";
import { Loading } from "../components/atoms/Loading";
import { NotFound } from "../pages/404/404";

export const ServiceFallbackLayout = (props:{
    children: ReactNode
}) => {
    const dispatch = useAppDispatch();
    const service = useAppSelector(state => state.rootReducer.service);
    const [isLoading, setIsLoading] = useState(true);
    const [isError, setIsError] = useState(false);
    const [getServiceTrigger] = getService.useLazyQuery();
    const location = useLocation();

    useEffect(() => {
        const port = Number(location.pathname.split('/')[2]);
        if (service.port === port) {
            setIsLoading(false);
            return;
        }
        getServiceTrigger(port).
            unwrap().then((res) => {
                dispatch(setService(res));
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