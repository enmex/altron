import { useEffect, useRef } from "react";
import { useNavigate } from "react-router"

export const useAppNavigation = () => {
    const navigate = useNavigate();
    const isInitialRender = useRef(true);

    useEffect(() => {
        (async () => {
            if (isInitialRender.current) {
                isInitialRender.current = false;
                return;
            }
            // await closeConnection((e: CloseEvent) => {
            //     if (e.code === 1006) {
            //         console.log('Disconnected with error');
            //         dispatch(putError({
            //             data: {
            //                 message: "unable to connect"
            //             }
            //         }))
            //     } else {
            //         console.log('Disconnected');
            //     }
            // });
        })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [navigate]);

    return navigate;
}