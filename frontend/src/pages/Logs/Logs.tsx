import { useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Log } from "../../app/types/Log";
import { MAX_LOGS_IN_CACHE, WEBSOCKET_SERVER_HOST } from "../../config";
import { useAppNavigation } from "../../hooks/navigate";
import { useLocation } from "react-router-dom";
import { Button } from "../../components/atoms/Button";
import { Icon } from "../../components/atoms/Icon";
import { Loading } from "../../components/atoms/Loading";
import { LogsPanel } from "../../components/organisms/LogsPanel";
import { Select } from "../../components/molecules/Select";
import { useTranslation } from "react-i18next";
import { Overlay } from "../../components/atoms/Overlay";
import { LogsContainerUpdatePanel } from "../../components/organisms/LogsContainerUpdatePanel";
import { Text } from "../../components/atoms/Text";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { getLatestLogs } from "../../app/store/logs/logs.api";

export const Logs = () => {
    const { t } = useTranslation();
    const dispatch = useAppDispatch();
    const location = useLocation();
    const service = useAppSelector(state => state.rootReducer.service);
    const [getLatestLogsTrigger] = getLatestLogs.useLazyQuery();
    const [isWaiting, setIsWaiting] = useState(false);
    const [connected, setConnected] = useState(false);
    const [updateLogsContainerPanelActive, setUpdateLogsContainerPanelActive] = useState(false);
    const [paginationIndex, setPaginationIndex] = useState(0);

    const [logs, setLogs] = useState<Log[]>([]);
    const socketRef = useRef<{
        socket: WebSocket | null
    }>({
        socket: null
    });

    const navigate = useAppNavigation();

    const onMessage = (e: MessageEvent) => {
        const message = e.data as string;
        if (message === "connect") {
            setIsWaiting(false);
            setConnected(true);
            return;
        }
        const log = JSON.parse(message) as Log;
        setLogs(prev => {
            if (prev.length >= MAX_LOGS_IN_CACHE) {
                return [...prev, log].slice(1, prev.length + 1);
            }
            return [...prev, log];
        });
    }

    const connectWebSocket = (disconnectedByError?: boolean): WebSocket => {
        setIsWaiting(true);
        const socket = new WebSocket(`${WEBSOCKET_SERVER_HOST}/connect/logs/${service.containerID}`);

        socket.onopen = () => {
            console.log('Connected');
            if (disconnectedByError) {
                notifyInfo('Reconnected successfully');
            }
        };

        socket.onclose = (e: CloseEvent) => {
            if (e.code === 1006) {
                console.log('Disconnected with error: ', e.reason);
                dispatch(putError('unable to connect. Reconnecting...'));
                socketRef.current.socket = connectWebSocket(true);
            } else {
                console.log('Disconnected');
            }
        }

        socket.onmessage = onMessage;

        return socket;
    }

    const handleConnectionSwitched = () => {
        if (!service.containerID) {
            return;
        }
        socketRef.current.socket?.send(JSON.stringify({
            action: 'pause'
        }));
        setConnected(prev => !prev);
    }

    useEffect(() => {
        if (!service.containerID) {
            return;
        }
        if (updateLogsContainerPanelActive) {
            return;
        }
        getLatestLogsTrigger({
            containerID: service.containerID
        })
        .unwrap()
        .then((res) => {
            setLogs(prev => [...prev, ...res.logs.map(logJson => JSON.parse(logJson) as Log)]);
        }).catch((err) => {
            dispatch(putError(err));
        })
        const socket = connectWebSocket();
        socketRef.current.socket = socket;

        return () => {
            (async () => {
                setLogs(prev => []);

                await new Promise<void>((resolve, reject) => {
                    socket.onclose = (e: CloseEvent) => {
                        if (e.code === 1006) {
                            console.log(e.reason);
                            console.log('Disconnected with error');
                            dispatch(putError('unable to connect'));
                        } else {
                            console.log('Disconnected');
                        }
                        resolve();
                    }
                    socket.close(1000);
                });
            })()
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [location.pathname, updateLogsContainerPanelActive]);

    return (
        <div className="flex flex-col h-full w-full mx-4">
            <div className="flex mt-4 flex-row w-full">
                <div className="flex flex-row w-full justify-start">
                    <Button
                        onClick={handleConnectionSwitched}
                    >
                        {
                            isWaiting 
                            ? <Loading hidden={false} size={30}/>
                            : connected
                            ? <Icon 
                                tip="continue"
                                type="negative"
                                name="debugPause"
                                size={30}
                            />
                            : <Icon 
                                tip="pause"
                                type="positive"
                                name="debugStart"
                                size={30}
                            />
                        }
                    </Button>
                    {
                        service.port > 0 && (
                            <>
                            <Button
                                onClick={() => navigate(`/services/${service.port}`)}
                            >
                                <Icon type="contrast" name="server" size={30}/>
                            </Button>
                                <Select 
                                className="relative"
                                icon="options"
                                items={[
                                    {
                                        text: t('update') + "",
                                        onItem: () => setUpdateLogsContainerPanelActive(true)
                                    }
                                ]}
                            />
                            </>
                        )
                    }
                    
                    {
                        service.containerID && (
                            <Text
                                className="text-xl font-bold w-5/6 flex justify-end"
                            >{`container ID ${service.containerID}`}</Text>
                        )
                    }
                </div>
            </div>
            <div className="flex h-[87%] w-full justify-center">
                <LogsPanel 
                    logs={
                        logs.length > 100 
                            ? logs.slice(logs.length - (paginationIndex + 1) * 100, logs.length) 
                            : logs
                    }
                    onExpand={() => setPaginationIndex(paginationIndex + 1)}
                />
            </div>
            {
                updateLogsContainerPanelActive && (
                    <Overlay>
                        <LogsContainerUpdatePanel 
                            onClose={() => setUpdateLogsContainerPanelActive(false)}
                        />
                    </Overlay>
                )
            }
        </div>
    );
}