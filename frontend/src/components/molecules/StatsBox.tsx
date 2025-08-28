import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { setContainerID } from "../../app/store/service/service.slice";
import { Stats } from "../../app/types/Stats";
import { useAppNavigation } from "../../hooks/navigate";
import { convertBytesToLargerUnit } from "../../utils/utils";
import { Button } from "../atoms/Button";
import { Icon } from "../atoms/Icon";
import { Text } from "../atoms/Text";

export const StatsBox = (props:{
    stats: Stats;
}) => {
    const dispatch = useAppDispatch();
    const theme = useAppSelector(state => state.rootReducer.theme);
    const navigate = useAppNavigation();

    const onClickStats = () => {
        dispatch(setContainerID(props.stats.id));
        if (props.stats.status === "down") {
            return;
        } 
        navigate('/logs/altron/' + props.stats.id);
    }

    return (
        <Button
            className="animate-fade flex flex-col w-64 h-56 p-6 m-4 border-2 rounded-md duration-200"
            borderColor={props.stats.status === "up" 
                ? theme.accents.positive 
                : props.stats.status === "mumble"
                ? theme.accents.contrast
                : theme.accents.negative
            }
            onClick={onClickStats}
        >
            <Text
                className="font-bold mb-2"
                color={theme.accents.contrast}
            >{`${props.stats.name.toUpperCase()}`}</Text>
            <Icon 
                color={
                    props.stats.status === "up" 
                    ? theme.accents.positive 
                    : props.stats.status === "mumble"
                    ? theme.accents.contrast
                    : theme.accents.negative
                } 
                name={
                    props.stats.status === "up" 
                    ? "online" 
                    : props.stats.status === "mumble"
                    ? "mumble"
                    : "offline"
                } 
                size={20}
            />
            <div className="flex flex-col">
                <Text
                    className="flex justify-start"
                >{`PIDS ${props.stats.stats?.pidsStats?.current ?? 0}`}</Text>
                <Text>
                    <Icon type="contrast" size={20} name="cpu" />
                    <div className="flex flex-row items-center">
                        {`${props.stats.stats?.cpuUsage.toFixed(2) ?? '--.--'}%`}
                        <Icon 
                            color={props.stats.stats ? props.stats.stats.growthDynamics.cpu > 0 
                                ? theme.accents.positive 
                                : props.stats.stats.growthDynamics.cpu < 0 
                                ? theme.accents.negative 
                                : theme.accents.contrast : theme.accents.contrast
                            } 
                            name={props.stats.stats ? props.stats.stats.growthDynamics.cpu > 0 
                                ? "increase" 
                                : props.stats.stats.growthDynamics.cpu < 0 
                                ? "decrease" 
                                : "stable" : "stable"
                            } 
                            size={15}
                        />
                    </div>
                </Text>
                <Text>
                    <Icon type="contrast" size={20} name="memory" />
                    <div className="flex flex-row items-center">
                        {`${props.stats.stats && props.stats.status !== "down" ? convertBytesToLargerUnit(props.stats.stats.memoryStats.usage) : '--.--'} / ${props.stats.stats && props.stats.status !== "down" ? convertBytesToLargerUnit(props.stats.stats.memoryStats.maxUsage) : '--.--'}`}           
                        <Icon 
                            color={props.stats.stats ? props.stats.stats.growthDynamics.memory > 0 
                                ? theme.accents.positive 
                                : props.stats.stats.growthDynamics.memory < 0 
                                ? theme.accents.negative 
                                : theme.accents.contrast : theme.accents.contrast
                            } 
                            name={props.stats.stats ? props.stats.stats.growthDynamics.memory > 0 
                                ? "increase" 
                                : props.stats.stats.growthDynamics.memory < 0 
                                ? "decrease" 
                                : "stable" : "stable"
                            } 
                            size={15}
                        />
                    </div>
                </Text>
            </div>
        </Button>
    );
}