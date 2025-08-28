import { ReactNode, useEffect, useState } from "react"
import { Text } from "./Text";
import { useAppSelector } from "../../app/store/hooks";
import { negativeColor } from "../../utils/utils";

export const Tooltip = (props:{
    tip?: string
    children: ReactNode
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [isVisible, setIsVisible] = useState(false);
    const [timer, setTimer] = useState<NodeJS.Timeout>();

    const onMouseEnter = () => {
        if (!props.tip) {
            return;
        }
        setTimer(setTimeout(() => {
            setIsVisible(true);
        }, 1000));
    }

    const onMouseLeave = () => {
        if (!props.tip) {
            return;
        }
        clearTimeout(timer);
        setTimer(undefined);
        setIsVisible(false);
    }

    useEffect(() => {
        return () => clearTimeout(timer);
    }, []);

    return (
        <div
            onMouseEnter={onMouseEnter}
            onMouseLeave={onMouseLeave}
        >
            {props.children}
            {isVisible && props.tip && (
                <Text
                    className="animate-fade absolute mt-1 font-bold p-2 rounded-md transition-all duration-200"
                    backgroundColor={negativeColor(theme.text)}
                >
                {props.tip}
                </Text>
            )}
        </div>
    );
}