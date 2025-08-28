import { useEffect, useRef, useState } from "react";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Session } from "../../app/types/Service";
import { Button } from "../../components/atoms/Button";
import { SessionListPanel } from "../../components/organisms/SessionListPanel";
import { SessionPanel } from "../../components/organisms/SessionPanel";
import { unsetService } from "../../app/store/service/service.slice";
import { Overlay } from "../../components/atoms/Overlay";
import { useDeleteServiceMutation } from "../../app/store/service/service.api";
import { ServiceUpdatePanel } from "../../components/organisms/ServiceUpdatePanel";
import { MAX_SESSIONS_IN_CACHE, WEBSOCKET_SERVER_HOST } from "../../config";
import { FilterPanel } from "../../components/molecules/FilterPanel";
import { Filter } from "../../app/types/Filter";
import { useTranslation } from "react-i18next";
import { Form } from "../../components/molecules/Form";
import { useAppNavigation } from "../../hooks/navigate";
import { Select } from "../../components/molecules/Select";
import { Text } from "../../components/atoms/Text";
import { Loading } from "../../components/atoms/Loading";
import { Icon } from "../../components/atoms/Icon";
import { WorkspaceCreationPanel } from "../../components/organisms/WorkspaceCreationPanel";
import { setWorkspace } from "../../app/store/workspace/workspace.slice";
import { useLocation } from "react-router";
import { Clipboard } from "../../components/atoms/Clipboard";
import { Icons } from "../../config/icons";
import { Workspace } from "../../app/types/Workspace";
import { AnalyzerPanel } from "../../components/molecules/AnalyzerPanel";
import { setMetrics } from "../../app/store/metrics/metrics.slice";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { Characteristic } from "../../app/types/Analyzer";
import { useGetAnalyzerPayloadQuery } from "../../app/store/analyzer/analyzer.api";

export const Service = () => {
    const dispatch = useAppDispatch();
    const location = useLocation();
    const theme = useAppSelector(state => state.rootReducer.theme);
    const { t } = useTranslation();
    const service = useAppSelector(state => state.rootReducer.service);
    const [sessions, setSessions] = useState<{
        session: Session,
        visible: boolean
    }[]>([]);
    const { data: analyzerPayloadData, isLoading: isAnalyzerPayloadLoading } = useGetAnalyzerPayloadQuery({
        servicePort: service.port
    });
    const [currentSession, setCurrentSession] = useState<Session | null>(null);
    const currentFilterRef = useRef<Filter>();
    const analyzerPayloadRef = useRef<{
        hasChecker: boolean,
        analyzer: {
            [componentName: string]: Characteristic[]
        },
    }>(analyzerPayloadData ? {
        ...analyzerPayloadData
    } : {
        hasChecker: false,
        analyzer: {}
    });
    const [connected, setConnected] = useState(false);
    const [isWaiting, setIsWaiting] = useState(false);
    const [updateServicePanelActive, setUpdateServicePanelActive] = useState(false);
    const [createWorkspacePanelActive, setCreateWorkspacePanelActive] = useState(false);
    const [paginationIndex, setPaginationIndex] = useState(0);
    const [confirmPanel, setConfirmPanel] = useState({
        active: false,
        message: "",
        onConfirm: () => {},
    });
    const [deleteService] = useDeleteServiceMutation();
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
        if (message.match("rpm")) {
            const metricsResponse = JSON.parse(message) as {
                rpm: number;
            }
            dispatch(setMetrics(metricsResponse.rpm));
            return;
        }
        const filter = currentFilterRef.current;
        const session = JSON.parse(message) as Session;
        const analyzerMatches = session.analyzerMatches;
        if (analyzerMatches) countSessionAnalyzerMatches(analyzerMatches);

        if ((filter && session.matchedFilters.find(f => f.name === filter.name)) || !filter) {
            setSessions(prev => {
                if (prev.length >= MAX_SESSIONS_IN_CACHE) {
                    return [...prev, {
                        session: session,
                        visible: true,
                    }].slice(1, prev.length + 1);
                }
                return [...prev, {
                    session: session,
                    visible: true,
                }];
            });
        }
    }

    const countSessionAnalyzerMatches = (analyzerMatches: {
        [componentName: string]: Characteristic;
    }) => {
        if (!analyzerPayloadRef.current.analyzer) return;

        let analyzer = {
            ...analyzerPayloadRef.current.analyzer
        };
        Object.keys(analyzerMatches).forEach(componentName => {
            const matched = analyzerMatches[componentName];
            
            if (!analyzer[componentName]) {
                analyzer[componentName] = [...analyzer[componentName], matched];
                return;
            }
            
            if (!analyzer[componentName].find(a => a.value === matched.value)) {
                analyzer[componentName] = [...analyzer[componentName], matched];
                return;
            }

            analyzer[componentName] = analyzer[componentName].map(ch => {
                return {
                    ...ch,
                    number: ch.value === matched.value ? ch.number + 1 : ch.number
                }
            })
        }); 
        analyzerPayloadRef.current = {
            ...analyzerPayloadRef.current,
            analyzer: analyzer,
        }
    }
    
    const handleConnectionSwitched = () => {
        socketRef.current.socket?.send(JSON.stringify({
            action: 'pause'
        }));
        setConnected(prev => !prev);
    }

    const onClickLink = () => {
        if (service.link.startsWith('http://')) {
            window.open(service.link, '_blank');
        }
    }

    const onDeleteService = () => {
        deleteService(service.id)
        .unwrap()
        .then(() => {
            notifyInfo(t('delete_service_success', {serviceName: service.name}));
            dispatch(unsetService());
            navigate('/home');
        }).catch((err) => {
            dispatch(putError(err.data.message));
        })
    }

    const onRefreshConnection = () => {
        socketRef.current.socket?.send(JSON.stringify({
            action: 'refresh'
        }));
    }

    const onClickFilter = (filter: Filter) => {
        if (currentFilterRef.current && currentFilterRef.current.id === filter.id) {
            setSessions(prev => prev.map((s) => {
                return {
                    ...s,
                    session: {
                        ...s.session,
                        selected: true,
                    }
                }
            }));
            currentFilterRef.current = undefined;
            return;
        }
        currentFilterRef.current = filter;
        setSessions(prev => prev.map((s) => {
            return {
                ...s,
                selected: !!(s.session.matchedFilters.find(s => s.name === filter.name))
            }
        }));
    }

    const onClickWorkspace = (workspace: Workspace) => {
        dispatch(setWorkspace(workspace));
        navigate("/workspaces/" + workspace.id);
    }

    const onClickCharacteristic = (characteristic: {
        componentName: string;
        value: string;
    }, action: string, reset: boolean) => {
        socketRef.current.socket?.send(JSON.stringify({
            action: action,
            payload: {
                key: characteristic.componentName,
                value: characteristic.value
            }
        }));
        
        // setSessions(prev => prev.map(s => {
        //     const visible = 
        //     (
        //         characteristic.value === 'analyzer-pass' &&
        //         !!(s.session.analyzerMatches) &&
        //         s.session.analyzerMatches[characteristic.componentName].value === characteristic.value
        //     ) 
        //     ||
        //     (
        //         characteristic.value === 'analyzer-block' &&
        //         (
        //             !s.session.analyzerMatches || 
        //             (
        //                 s.session.analyzerMatches && 
        //                 s.session.analyzerMatches[characteristic.componentName].value !== characteristic.value
        //             )
        //         )
        //     )
        //     return {
        //         ...s,
        //         visible: visible
        //     }
        // }));
    }

    const connectWebSocket = (disconnectedByError?: boolean): WebSocket => {
        setIsWaiting(true);
        const socket = new WebSocket(`${WEBSOCKET_SERVER_HOST}/connect/sessions/${service.port}`);

        socket.onopen = () => {
            console.log('Connected');
            if (disconnectedByError) {
                notifyInfo('Reconnected successfully');
            }
        };

        socket.onclose = (e: CloseEvent) => {
            if (e.code === 1006) {
                console.log('Disconnected with error: ', e.reason);
                dispatch(putError('unable to connect. Try again'));
                // setTimeout(() => {
                //     socketRef.current.socket = connectWebSocket(true);
                // }, 5000);
            } else {
                console.log('Disconnected');
            }
        }

        socket.onmessage = onMessage;

        return socket;
    }

    useEffect(() => {
        if (analyzerPayloadData) {
            analyzerPayloadRef.current = analyzerPayloadData;
        }
    }, [analyzerPayloadData]);

    useEffect(() => {
        const socket = connectWebSocket();
        socketRef.current.socket = socket;

        return () => {
            (async () => {
                setSessions(prev => []);
                setCurrentSession(null);
                currentFilterRef.current = undefined;

                await new Promise<void>((resolve, reject) => {
                    socket.onclose = (e: CloseEvent) => {
                        if (e.code === 1006) {
                            console.log(e.reason);
                            console.log('Disconnected with error');
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
    }, [location.pathname]);

    return (
        <div className="flex w-full flex-col mx-4 my-2">
            <div className="flex flex-row h-1/6 w-full">
                <div className="flex w-full justify-between py-2">
                    <div className="flex flex-row w-1/2">
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
                        <Button
                            onClick={() => {
                                setSessions([]);
                                setCurrentSession(null);
                            }}
                        ><Icon tip="clear sessions" type="negative" name="brush3Line" size={25}/></Button>
                        <Button
                            onClick={() => navigate(`/logs/${service.port}`)}
                        ><Icon tip="watch service logs" type="contrast" name="terminal" size={25}/></Button>
                        <Select 
                            className="relative mr-1"
                            icon="letterW"
                            items={[
                                {
                                    icon: {
                                        value: "mdAdd",
                                        color: theme.accents.positive
                                    },
                                    onItem: () => setCreateWorkspacePanelActive(true)
                                },
                                ...service.workspaces.map(workspace => {
                                    const icon: keyof Icons = workspace.status === "LISTENING" 
                                        ? "headphones"
                                        : workspace.status === "WAITING"
                                        ? "time"
                                        : "done";
                                    return {
                                        text: workspace.name,
                                        icon: {
                                            value: icon,
                                            color: theme.text
                                        },
                                        onItem: () => onClickWorkspace(workspace)
                                    }
                                })
                            ]}
                        />
                        <Select 
                            className="relative"
                            icon="options"
                            items={[
                                {
                                    text: t('update') + "",
                                    onItem: () => setUpdateServicePanelActive(true)
                                },
                                {
                                    text: t('delete') + "",
                                    onItem: () => setConfirmPanel({
                                        active: true,
                                        message: t('delete_service_confirm'),
                                        onConfirm: onDeleteService
                                    })
                                }
                            ]}
                        />
                        <Button
                            className="ml-2 w-3/5 px-1 flex justify-center rounded-md duration-200"
                            backgroundColor={theme.tertiary}
                            onClick={onClickLink} 
                        >
                            <Text className="text-lg font-bold text-start">{service.link}</Text>
                        </Button>
                        <Clipboard text={service.link} size={25}/>
                        {
                            isAnalyzerPayloadLoading ? <Loading size={30} hidden={false} /> : (
                                <AnalyzerPanel 
                                    analyzerPayload={analyzerPayloadRef.current}
                                    onClick={onClickCharacteristic}
                                />
                            )
                        }
                    </div>
                    <FilterPanel 
                        onClick={onClickFilter}
                        selectedFilter={currentFilterRef.current}
                        onUpdate={onRefreshConnection}
                    />
                </div>
            </div>
            <div className="flex flex-row h-[88vh] w-full">
                <div className="flex flex-row w-1/3">
                    <SessionListPanel 
                        hasCheckerMask={analyzerPayloadRef.current.hasChecker}
                        sessions={
                            sessions.filter(s => s.visible).map(s => s.session).length > 100 
                            ? sessions.filter(s => s.visible).map(s => s.session).slice(sessions.length - (paginationIndex + 1) * 100, sessions.length) 
                            : sessions.filter(s => s.visible).map(s => s.session)
                        }
                        onClick={setCurrentSession}
                        onExpand={() => setPaginationIndex(prev => prev + 1)}
                    />
                </div>
                {
                    currentSession && 
                    <div className="w-2/3">
                        <SessionPanel 
                            session={currentSession}
                            filters={currentSession.matchedFilters}
                        />
                    </div>
                }
            </div>
            {
                updateServicePanelActive && (
                    <Overlay>
                        <ServiceUpdatePanel 
                            onClose={() => {setUpdateServicePanelActive(false); onRefreshConnection();}}
                        />
                    </Overlay>
                )
            }
            {
                createWorkspacePanelActive && (
                    <Overlay>
                        <WorkspaceCreationPanel 
                            servicePort={service.port}
                            onClose={() => {setCreateWorkspacePanelActive(false)}}
                        />
                    </Overlay>
                )
            }
            {
                confirmPanel.active && (
                    <Overlay>
                        <Form 
                            label={confirmPanel.message}
                            onCancel={() => setConfirmPanel({
                                ...confirmPanel,
                                active: false
                            })}
                            onSubmit={confirmPanel.onConfirm}
                        />
                    </Overlay>
                )
            }
        </div>
    );
}