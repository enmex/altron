import { useState, useRef, useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { WEBSOCKET_SERVER_HOST } from "../../config";
import { Stats } from "../../app/types/Stats";
import { StatsPanel } from "../../components/organisms/StatsPanel";
import { CONTAINERS } from "../../config/constants";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { Panel } from "../../components/atoms/Panel";
import { Text } from "../../components/atoms/Text";
import { Icon } from "../../components/atoms/Icon";
import { convertBytesToLargerUnit } from "../../utils/utils";

export const Health = () => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const dispatch = useAppDispatch();
    const [stats, setStats] = useState<Stats[]>(CONTAINERS.map((c) => {
        return {
            id: "",
            name: c,
            image: "",
            status: "down"
        }
    }));
    const [totalStats, setTotalStats] = useState<{
        status: string;
        pids: number;
        cpuUsage: number;
        memoryUsage: number;
        growthDynamics: {
            cpu: number;
            memory: number;
        };
    }>({
        status: "up",
        pids: 0,
        cpuUsage: 0,
        memoryUsage: 0,
        growthDynamics: {
            cpu: 0,
            memory: 0
        }
    });
    const socketRef = useRef<{
        socket: WebSocket | null
    }>({
        socket: null
    });

    const onMessage = (e: MessageEvent) => {
        const message = e.data as string;
        if (message === "connect") {
            return;
        }
        const containerStats = JSON.parse(message) as Stats;
        setStats(prev => prev.map(s => s.name === containerStats.name ? containerStats : s));
    }

    const connectWebSocket = (disconnectedByError?: boolean): WebSocket => {
        const socket = new WebSocket(`${WEBSOCKET_SERVER_HOST}/connect/health`);

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

    useEffect(() => {
        const socket = connectWebSocket();
        socketRef.current.socket = socket;

        return () => {
            (async () => {
                await new Promise<void>((resolve, ) => {
                    socket.onclose = (e: CloseEvent) => {
                        if (e.code === 1006) {
                            console.log('Disconnected with error:', e.reason);
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
    }, []);

    useEffect(() => {
        const altronStats = stats.filter(s => s.name !== "agent");
        const pids = altronStats.reduce((acc, current) => acc + (current.stats?.pidsStats?.current ?? 0), 0);
        const cpu = altronStats.reduce((acc, current) => acc + (current.stats ? current.stats.cpuUsage : 0), 0);
        const memory = altronStats.reduce((acc, current) => acc + (current.stats ? current.stats.memoryStats.usage : 0), 0);

        setTotalStats(prev => {
            return {
                status: cpu > 50 ? "mumble" : "up",
                pids: pids,
                cpuUsage: cpu,
                memoryUsage: memory,
                growthDynamics: {
                    cpu: cpu > prev.cpuUsage ? 1 : cpu == prev.cpuUsage ? 0 : -1,
                    memory: memory > prev.memoryUsage ? 1 : memory == prev.memoryUsage ? 0 : -1
                }
            }
        });
    }, [stats]);

    return (
        <div className="flex flex-col h-11/12 w-full justify-center items-center">
            <Panel
                className="flex flex-row my-8 items-center"
                color={theme.primary}
            >
                <Icon 
                    color={
                        totalStats.status === "up" 
                        ? theme.accents.positive 
                        : totalStats.status === "mumble"
                        ? theme.accents.contrast
                        : theme.accents.negative
                    } 
                    name={
                        totalStats.status === "up" 
                        ? "online" 
                        : totalStats.status === "mumble"
                        ? "mumble"
                        : "offline"
                    } 
                    size={20}
                />
                <Text
                    className="font-bold text-xl mx-4"
                >{`PIDS ${totalStats.pids}`}</Text>
                <Text
                    className="flex flex-row items-center font-bold text-xl mr-4"
                >
                    <Icon type="contrast" size={20} name="cpu" />
                    <div className="flex flex-row items-center ml-1">
                        {`${totalStats.cpuUsage.toFixed(2) ?? '--.--'}%`}
                        <Icon 
                            color={totalStats.growthDynamics.cpu > 0 
                                ? theme.accents.positive 
                                : totalStats.growthDynamics.cpu < 0 
                                ? theme.accents.negative 
                                : theme.accents.contrast
                            } 
                            name={totalStats.growthDynamics.cpu > 0 
                                ? "increase" 
                                : totalStats.growthDynamics.cpu < 0 
                                ? "decrease" 
                                : "stable"
                            } 
                            size={15}
                        />
                    </div>
                </Text>
                <Text
                    className="flex flex-row items-center font-bold text-xl mr-4"
                >
                    <Icon type="contrast" size={20} name="memory" />
                    <div className="flex flex-row items-center ml-1">
                        {`${convertBytesToLargerUnit(totalStats.memoryUsage)}`}           
                        <Icon 
                            color={totalStats.growthDynamics.cpu > 0 
                                ? theme.accents.positive 
                                : totalStats.growthDynamics.cpu < 0 
                                ? theme.accents.negative 
                                : theme.accents.contrast
                            } 
                            name={totalStats.growthDynamics.cpu > 0 
                                ? "increase" 
                                : totalStats.growthDynamics.cpu < 0 
                                ? "decrease" 
                                : "stable"
                            } 
                            size={15}
                        />
                    </div>
                </Text>
            </Panel>
            <div className="w-10/12 items-center">
                <StatsPanel stats={stats}/>
            </div>
        </div>
    )
}