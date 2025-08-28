import { useAppSelector } from "../../app/store/hooks";
import { useGetSessionQuery } from "../../app/store/session/session.api";
import { Loading } from "../../components/atoms/Loading";
import { Text } from "../../components/atoms/Text";
import { SessionPanel } from "../../components/organisms/SessionPanel";

export const Session = () => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const sessionId = window.location.pathname.split('/')[2];
    const {isLoading, data: session} = useGetSessionQuery(sessionId);
    
    return (
        <div className="flex flex-col justify-center items-center w-full h-full">
            {
                isLoading ? <Loading hidden={false} size={50}/> : session && (
                    <div className="flex flex-col justify-center w-2/3">
                        <div className="flex flex-row w-1/2 justify-between items-center mb-4">
                            <div className="flex flex-row items-center">
                                <Text 
                                    className="font-bold text-2xl mr-2" 
                                    color={theme.accents.contrast}
                                >{session.protocol.toUpperCase()}</Text>
                                <Text 
                                    className="text-lg mr-2" 
                                    color={theme.text}
                                >{`${session.clientHost} -> [::1]:${session.serverPort}`}</Text>
                                <Text 
                                    className="mr-2 font-semibold text-lg" 
                                    color={theme.accents.neutral}
                                >{new Date(session.sentAt).toLocaleTimeString()}</Text>
                                <Text 
                                    className="text-lg font-bold mr-2" 
                                    color={theme.accents.neutral}
                                >{`TTL ${session.ttl}`}</Text>
                                <Text 
                                    className="text-lg font-bold" 
                                    color={theme.accents.positive}
                                >{session.packetsCount + " packets"}</Text>
                            </div>
                            {
                                session.matchedFilters.length > 0 && <div className="flex flex-row">
                                {
                                    session.matchedFilters.map(filter => {
                                        return (
                                            <Text
                                                className="cursor-default p-1 m-2 font-bold text-lg rounded"
                                                backgroundColor={filter.color}
                                            >{filter.name + (filter.matchesCount > 1 ? ` x${filter.matchesCount}` : '')}</Text>
                                        )
                                    }) 
                                }
                                </div>
                            } 
                        </div>
                        <SessionPanel watchMode session={session} filters={session.matchedFilters}/>
                    </div>
                )
            }
        </div>
    )
}