import { useAppSelector } from "../../app/store/hooks";
import { SessionFilter } from "../../app/types/Filter";
import { Packet } from "../../app/types/Service"
import { Text } from "../atoms/Text"
import { HighlightedText } from "../atoms/HighlightedText";
import { Select } from "./Select";
import { conversionTypes } from "../../config/constants";
import { useState } from "react";
import { Button } from "../atoms/Button";
import { Icon } from "../atoms/Icon";
import { textToHexDump } from "../../utils/utils";

export const PacketPanel = (props: {
    onClickConversion: (idx: number, exportType: string) => void,
    packet: Packet, 
    filters?: SessionFilter[],
    time: number,
    idx: number,
    highlighted?: boolean
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [isHexMode, setIsHexMode] = useState(false);

    return (
        <div 
            className={`flex flex-col duration-200 rounded ${props.packet.isRequest ? "justify-start items-start" : "justify-end items-end"}`}
            style={{
                backgroundColor: props.highlighted ? theme.secondary : "transparent"
            }}
        >
            <div
                className="flex flex-col w-2/3 my-2 p-2 whitespace-pre-line font-mono font-normal text-lg mb-4 border-2 cursor-default rounded-lg"
                style={{
                    backgroundColor: props.packet.isRequest ? theme.request : theme.response
                }}
            >
                <div className="flex flex-row">
                    <div 
                        className="px-2 min-w-1/4 max-w-1/3 flex justify-around text-center items-center font-mono rounded-md"
                        style={{
                            backgroundColor: theme.tertiary
                        }}
                    >
                        <Text className="font-bold mr-4">{`${props.idx + 1}. ${props.packet.isRequest ? "Client" : "Server"}`}</Text>
                        <Text>{`+${props.time} ms`}</Text>
                    </div>
                    <Button
                        onClick={() => setIsHexMode(!isHexMode)}
                    ><Icon tip={isHexMode ? "normal mode" : "hex mode"} type="neutral" name="eye2Line" size={25}/></Button>
                    <Select 
                        icon="outlineExport"
                        className="relative ml-2"
                        items={
                            conversionTypes.filter(ct => props.packet.isRequest ? ct.inRequest : ct.inResponse).map(conversionType => {
                                return {
                                    text: conversionType.text,
                                    onItem: () => props.onClickConversion(props.idx, conversionType.type)
                                }
                            })
                        }
                    />
                </div>
                <div className="m-2 w-full max-h-96 overflow-auto">
                    <HighlightedText filters={props.filters} text={isHexMode ? textToHexDump(props.packet.payload) : props.packet.payload} />
                </div>
            </div>    
        </div>
    )
}