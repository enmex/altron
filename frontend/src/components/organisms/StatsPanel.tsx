import { useAppSelector } from "../../app/store/hooks";
import { Stats } from "../../app/types/Stats";
import { Panel } from "../atoms/Panel";
import { StatsBox } from "../molecules/StatsBox";

export const StatsPanel = (props:{
    stats: Stats[]
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    return (
        <Panel
            withBorder
            className="flex flex-wrap items-center justify-center rounded"
            color={theme.primary}
        >
            {
                props.stats.map(s => {
                    return <StatsBox stats={s}/>
                })
            }
        </Panel>
    );
}