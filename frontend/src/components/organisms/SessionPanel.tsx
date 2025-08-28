import { useState } from "react";
import { useConvertSessionToExploitMutation, useExtractFilesFromPacketMutation } from "../../app/store/conversion/conversion.api";
import { Session } from "../../app/types/Service"
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { ExploitPanel } from "../molecules/ExploitPanel";
import { Overlay } from "../atoms/Overlay";
import { useTranslation } from "react-i18next";
import { Filter } from "../../app/types/Filter";
import { useConvertPacketToExploitMutation } from "../../app/store/conversion/conversion.api";
import { b64ToPythonBytes } from "../../utils/utils";
import { Select } from "../molecules/Select";
import { Loading } from "../atoms/Loading";
import { Button } from "../atoms/Button";
import { Icon } from "../atoms/Icon";
import { useCreateSessionMutation } from "../../app/store/session/session.api";
import { PacketsList } from "../molecules/PacketsList";
import { useAddWorkspaceSessionsToCartMutation } from "../../app/store/cart/cart.api";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const SessionPanel = (props: {
    session: Session,
    watchMode?: boolean
    filters?: Filter[]
}) => {
    const dispatch = useAppDispatch();
    const isWorkspaceMode = useAppSelector(state => state.rootReducer.workspace.id.length > 0);
    const workspace = useAppSelector(state => state.rootReducer.workspace);
    const { t } = useTranslation();
    const [convertSessionToExploit] = useConvertSessionToExploitMutation();
    const [convertPacketToExploit] = useConvertPacketToExploitMutation();
    const [addSessionToCart] = useAddWorkspaceSessionsToCartMutation();
    const [createSession] = useCreateSessionMutation();
    const [extractFiles] = useExtractFilesFromPacketMutation();
    const [exploitPanelActive, setExploitPanelActive] = useState(false);
    const [exploit, setExploit] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    const onClickSessionConversion = (exportType: string) => {
        setIsLoading(true);
        convertSessionToExploit({
            session: props.session,
            exportType: exportType
        })
        .unwrap()
        .then((data) => {
            notifyInfo(t('happy_hacking'));
            setExploit(data.exploit);
            setExploitPanelActive(true);
        })
        .catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        })
    }

    const onClickPacketConversion = (idx: number, exportType: string) => {
        const packet = props.session.packets[idx];
        if (exportType === "python_bytes") {
            setExploit(b64ToPythonBytes(packet.payload));
            setExploitPanelActive(true);
            return;
        }
        if (exportType === "files") {
            setIsLoading(true);
            extractFiles({
                sessionID: props.session.id,
                packetNumber: idx,
                packet: packet
            }).unwrap().then((res) => {
                const fileName = `${props.session.id}_${idx}_output_file.extracted.zip`;
                fetch(`data:application/octet-stream;base64,${res.data}`).then((res) => {
                    res.blob().then((blob) => {
                        const link = document.createElement('a');
                        link.href = window.URL.createObjectURL(blob);

                        link.download = fileName;

                        document.body.appendChild(link);
                        link.click();

                        document.body.removeChild(link);
                    }).catch((err) => {
                        dispatch(putError(err.data.message));
                    })
                }).catch((err) => {
                    dispatch(putError(err.data.message));
                });
            }).catch((err) => {
                dispatch(putError(err.data.message));
            }).finally(() => {
                setIsLoading(false);
            })
            return;
        }
        setIsLoading(true);
        convertPacketToExploit({
            packet: packet,
            exportType: exportType,
            servicePort: props.session.serverPort,
        }).unwrap().then((data) => {
            setIsLoading(false);
            setExploit(data.exploit);
            setExploitPanelActive(true);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    const onClickShare = () => {
        if (isWorkspaceMode) {
            onOpenSessionPage();
            return;
        }
        setIsLoading(true);
        createSession({
            ...props.session,
        }).unwrap().then(() => {
            setIsLoading(false);
            onOpenSessionPage();
        }).catch((err) => {
            if (err.data.message.includes("exists")) {
                onOpenSessionPage();
                return
            }
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const onClickMerge = () => {
        setIsLoading(true);
        addSessionToCart({
            workspaceId: workspace.id,
            sessions: [props.session.id]
        }).unwrap().then(() => {
            notifyInfo("session added to cart");
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const onOpenSessionPage = () => {
        window.open(`http://${window.location.host}/sessions/${props.session.id}`, '_blank');
    }

    return (
        <>
        <div className="flex w-full justify-center">
            <div className="mx-4 px-4 w-full h-[90%] rounded">
                <div className="flex w-full">
                    <div className="flex flex-row w-1/4 items-start">
                        <Select
                            items={[
                                {
                                    text: "as pwntools",
                                    onItem: () => {},
                                    children: [
                                        {
                                            text: "with recvuntil",
                                            onItem: () => onClickSessionConversion("pwntools_recvuntil")
                                        },
                                        {
                                            text: "with recvrepeat",
                                            onItem: () => onClickSessionConversion("pwntools_recvrepeat")
                                        }
                                    ]
                                },
                                {
                                    text: "as requests",
                                    onItem: () => onClickSessionConversion("requests")
                                },
                            ]} 
                            className="relative"
                            icon="outlineExport"
                        />
                        {
                            !props.watchMode && (
                                <>
                                <Button
                                    className="ml-2"
                                    onClick={onClickShare}
                                >
                                    <Icon 
                                        tip="share session"
                                        type="positive"
                                        name="share"
                                        size={30}
                                    />
                                </Button>
                                {
                                    workspace.id.length > 0 && (
                                        <Button
                                            onClick={onClickMerge}
                                        >
                                            <Icon 
                                                tip="add to cart"
                                                type="contrast"
                                                name="merge"
                                                size={30}
                                            />
                                        </Button>
                                    )
                                }
                                </>
                            )
                        }
                        <Loading className="flex" hidden={!isLoading} size={30} />
                    </div>
                </div>
                <PacketsList 
                    session={props.session}
                    onClickConversion={onClickPacketConversion}
                />
            </div>
        </div>
        {
            exploitPanelActive && (
                <Overlay>
                    <ExploitPanel 
                        exploit={exploit}
                        onClose={() => setExploitPanelActive(false)}
                    />
                </Overlay>
            )
        }
        </>
    );
}