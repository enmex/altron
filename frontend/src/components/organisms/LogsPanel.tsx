import { useEffect, useRef } from "react";
import { useAppSelector } from "../../app/store/hooks";
import { Log } from "../../app/types/Log";
import { b64ToString, negativeColor } from "../../utils/utils";
import { Panel } from "../atoms/Panel";
import { Text } from "../atoms/Text";
import Ansi from 'ansi-to-react';

export const LogsPanel = (props: {
    logs: Log[];
    onExpand: () => void;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const listRef = useRef<HTMLDivElement>(null);

    const scrollToBottom = () => {
        if (listRef.current) {
            listRef.current.scrollTop = listRef.current.scrollHeight;
        }
    }

    useEffect(() => {
        scrollToBottom();
    }, [props.logs]);

    return (
        <Panel 
            color={negativeColor(theme.text)} 
            withBorder
            className="justify-center items-center my-1 p-2 w-4/5 overflow-auto rounded-md"
        >
            <div 
                className="m-4 flex flex-col overflow-auto justify-center items-start"
                ref={listRef}
            >
                {
                    props.logs.map((log) => {
                        return (
                            <Text
                                className="whitespace-pre-wrap font-bold text-lg font-mono ml-1"
                            ><Ansi>{b64ToString(log.message)}</Ansi></Text>
                        )
                    })
                }
            </div>
        </Panel>
    );
}