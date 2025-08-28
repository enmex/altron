import { useEffect, useRef } from "react";
import { Session } from "../../app/types/Service"
import { Button } from "../atoms/Button";
import { MdOutlineExpandMore } from "react-icons/md";
import { darkenColor, randomKey } from "../../utils/utils";
import { Panel } from "../atoms/Panel";
import { useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text";
import { DropdownButton } from "../molecules/DropdownButton";
import { Icon } from "../atoms/Icon";

export const SessionListPanel = (props: {
    sessions: Session[]
    hasCheckerMask?: boolean
    onClick: (session: Session) => void,
    onExpand: () => void,
}) => {
    const listRef = useRef<HTMLDivElement>(null);
    const theme = useAppSelector(state => state.rootReducer.theme);
    const selectedRef = useRef<{
        index: number
    }>({
        index: -1
    });

    const scrollFunction = (event: KeyboardEvent) => {
        if (!listRef?.current) {
            return;
        }
        if (event.key === "k" || event.key === "л") {
            listRef.current.scrollTo({
                top: listRef.current.scrollTop - 150,
                behavior: "smooth",
            });
            if (selectedRef.current.index < 0) {
                selectedRef.current.index = props.sessions.length - 1;
            }
            const idx = selectedRef.current.index + 1 > props.sessions.length - 1 
                ? selectedRef.current.index 
                : selectedRef.current.index + 1;
            onClickSession(idx, props.sessions[idx]);
        }
        if (event.key === "j" || event.key === "о") {
            if (selectedRef.current.index < 0) {
                selectedRef.current.index = props.sessions.length;
            } else {
                listRef.current.scrollTo({
                    top: listRef.current.scrollTop + 150,
                    behavior: "smooth",
                });
            }
            const idx = selectedRef.current.index - 1 < 0 ? 0 : selectedRef.current.index - 1;
            onClickSession(idx, props.sessions[idx]);
        }
    }

    const onClickSession = (idx: number, session: Session) => {
        selectedRef.current.index = idx;
        props.onClick(session);
    }

    useEffect(() => {
        document.addEventListener("keydown", scrollFunction, false);

        if (listRef?.current) {
            listRef.current.scrollTop = -listRef.current.scrollHeight;
        }

        return () => document.removeEventListener("keydown", scrollFunction, false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [props.sessions.length]);
    
    return (
        <Panel className="p-1 h-full w-full rounded-lg">
            <div ref={listRef} className="w-full overflow-auto rounded-md h-full">
                <div className="flex flex-col-reverse justify-center h-auto max-h-auto">
                    {
                        props.sessions.length >= 100 && (
                            <Button 
                                onClick={() => props.onExpand()}
                                className="flex justify-center duration-200 w-full"
                            >
                                <MdOutlineExpandMore 
                                    size={30}
                                    className="text-white hover:text-cyan-300 duration-200"
                                />
                            </Button>
                        )
                    }
                    {
                        props.sessions.map((session, idx) => {
                            return (
                                <Button 
                                    key={randomKey()}
                                    className="flex justify-start mb-1 border-b-2 duration-200"
                                    borderColor={theme.tertiary}
                                    backgroundColor={selectedRef.current.index >= 0 && selectedRef.current.index === idx ? darkenColor(theme.secondary, 0.5) : theme.secondary}
                                    onClick={() => onClickSession(idx, session)}
                                >
                                    <div className="m-3 p-3">
                                        <div className="flex flex-col text-start font-bold">
                                            <div className="flex flex-row items-center mb-2">
                                                <Text 
                                                    className="font-bold text-lg" 
                                                    color={theme.accents.contrast}
                                                >{session.protocol.toUpperCase()}</Text>
                                                {
                                                    props.hasCheckerMask && !session.isSafe && (
                                                        <DropdownButton 
                                                            icon={
                                                                <Icon 
                                                                    color={theme.accents.negative}
                                                                    name="warning"
                                                                    size={15}
                                                                />
                                                            }
                                                            dropdownItems={
                                                                []
                                                            }
                                                        />
                                                    )
                                                }
                                            </div>
                                            <Text 
                                                className="text-md" 
                                                color={theme.text}
                                            >{`${session.clientHost} -> ${session.serverPort}`}</Text>
                                            <div className="flex flex-row justify-between">
                                                <div className="flex flex-row">
                                                    <Text 
                                                        className="pr-2 font-semibold" 
                                                        color={theme.accents.neutral}
                                                    >{new Date(session.sentAt).toLocaleTimeString()}</Text>
                                                    <Text 
                                                        className="text-md font-bold pr-2" 
                                                        color={theme.accents.neutral}
                                                    >{`TTL ${session.ttl}`}</Text>
                                                    <Text 
                                                        color={theme.accents.positive}
                                                    >{session.packetsCount + " packets"}</Text>
                                                </div>
                                            </div>
                                            <div>
                                                {
                                                    session.matchedFilters.map(filter => {
                                                        return (
                                                            <span   
                                                                key={randomKey()}
                                                                className="p-1 mx-1 rounded duration-200 font-bold text-lg"
                                                                style={{ 
                                                                    backgroundColor: filter.color,
                                                                    color: theme.text
                                                                }}
                                                            >
                                                                {filter.name + (filter.matchesCount > 1 ? ` x${filter.matchesCount}` : '')}
                                                            </span>
                                                        )
                                                    })
                                                }
                                            </div>
                                        </div>
                                    </div>
                                </Button>
                            )
                        })
                    }
                </div>
            </div>
        </Panel>
    );
}