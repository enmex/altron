import { useEffect } from "react";
import { useAppDispatch } from "../../app/store/hooks"
import { unsetService } from "../../app/store/service/service.slice";
import { unsetWorkspace } from "../../app/store/workspace/workspace.slice";
import { AnimatedLayout } from "../../layouts/AnimatedLayout";

export const Home = () => {
    const dispatch = useAppDispatch();

    useEffect(() => {
        dispatch(unsetService());
        dispatch(unsetWorkspace());
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <AnimatedLayout />
    )
}