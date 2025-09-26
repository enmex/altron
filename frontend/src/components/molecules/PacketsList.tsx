import { useEffect, useRef, useState } from "react"
import { Packet, Session } from "../../app/types/Service";
import { b64ToString, randomKey } from "../../utils/utils";
import { PacketPanel } from "./PacketPanel";
import { Loading } from "../atoms/Loading";
import { Button } from "../atoms/Button";
import { SessionFilter } from "../../app/types/Filter";
import { useAppSelector } from "../../app/store/hooks";

export const PacketsList = (props: {
    session: Session,
    onClickConversion: (idx: number, exportType: string) => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const pageRef = useRef(0);
    const loadingRef = useRef(false);
    const [highlightedPacket, setHighlightedPacket] = useState(-1);
    const fetchNext = (page: number, quantity: number): Promise<Packet[]> => {
        return new Promise((resolve,) => {
            let from = page * quantity;
            const to = Math.min(((page + 1) * quantity), props.session.packets.length);

            const packets = props.session.packets.slice(from, to).map(packet => {
                return {
                    ...packet,
                    payload: b64ToString(packet.payload)
                }
            });
            return resolve(packets);
        });
    }
    const [currentFilterPacket, setCurrentFilterPacket] = useState<{
        filterId: string,
        matchedPacketIdx: number
    }>({
        filterId: "",
        matchedPacketIdx: 0
    });

    const [listData, setListData] = useState<Packet[]>([]);
    const listRef = useRef<HTMLDivElement>(null);

    const next = (): Promise<void> => {
        loadingRef.current = true;
        return fetchNext(pageRef.current++, 10).then(newData => {
            setListData(prev => [...prev, ...newData]);
            loadingRef.current = false;
        });
    }

    const onClickFilter = async (filter: SessionFilter) => {
        const toIdx = filter.id === currentFilterPacket.filterId
            ? (currentFilterPacket.matchedPacketIdx + 1) % filter.matchedPackets.length
            : 0;
        const matchedPacketIdx = filter.matchedPackets[toIdx];
        if (matchedPacketIdx > listData.length - 1) {
            for(let i = 0; i <= (matchedPacketIdx - listData.length - 1) / 10; i++) {
                listRef.current?.children[listRef.current?.children.length - 1].scrollIntoView({
                    behavior: 'smooth'
                })
                await next();
            }
        }
        listRef.current?.children[matchedPacketIdx].scrollIntoView({
            behavior: 'smooth'
        })
        setHighlightedPacket(matchedPacketIdx);
        setTimeout(() => {
            setHighlightedPacket(-1);
        }, 1000);
        setCurrentFilterPacket({
            filterId: filter.id,
            matchedPacketIdx: toIdx
        })
    }

    useEffect(() => {
        fetchNext(0, 10).then(newData => {
            setListData(newData);
        });
        pageRef.current = 1;
        setCurrentFilterPacket({
            filterId: "",
            matchedPacketIdx: 0
        })
        const scroll = listRef.current;
        if (!scroll) {
            return;
        }
        if (props.session.packets.length <= 10) {
            return;
        }
        scroll.scrollTop = 0;
        const handleScroll = (event: Event) => {
            if (!listRef.current) {
                return
            }
            const { scrollTop, scrollHeight, clientHeight } = listRef.current;
            if (scrollTop + clientHeight >= scrollHeight) {
                next();
            }
        }

        scroll.addEventListener('scroll', handleScroll);

        return () => {
            scroll.removeEventListener('scroll', handleScroll);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [props.session.id]);

    return (
        <div className="flex flex-col">
            <div className="flex flex-row">
                {
                    props.session.matchedFilters.filter(f => f.regex).map(matchedFilter => {
                        return (
                            <Button
                                key={randomKey()}
                                onClick={() => onClickFilter(matchedFilter)}
                                className="mt-2 mx-2 flex justify-start"
                            >
                                <span   
                                    className="p-1 rounded duration-200 font-bold text-lg"
                                    style={{ 
                                        backgroundColor: matchedFilter.color,
                                        color: theme.text
                                    }}
                                >
                                    {matchedFilter.name + (matchedFilter.matchesCount > 1 ? ` x${matchedFilter.matchesCount}` : '')}
                                </span>
                            </Button>
                        )
                    })
                }
            </div>
            <div ref={listRef} className="overflow-auto flex flex-col h-[80vh] p-4">
                {
                    listData.map((packet, idx) => {
                        if (idx > props.session.packets.length - 1) {
                            return <div key={idx}></div>
                        }
                        const time = idx == 0 
                            ? 0 
                            : new Date(packet.sentAt).getTime() - new Date(props.session.packets[idx - 1].sentAt).getTime();
                        
                        return (
                            <PacketPanel 
                                packet={packet}
                                key={idx}
                                time={time}
                                idx={idx}
                                filters={props.session.matchedFilters}
                                onClickConversion={props.onClickConversion}
                                highlighted={idx === highlightedPacket}
                            />
                        )
                    })
                }
                <Loading hidden={!loadingRef.current} size={40} />
            </div>
        </div>
    )
}