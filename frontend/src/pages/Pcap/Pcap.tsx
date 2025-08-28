import { useEffect, useState } from "react";
import { Session } from "../../app/types/Service";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { SessionListPanel } from "../../components/organisms/SessionListPanel";
import { SessionPanel } from "../../components/organisms/SessionPanel";
import { Loading } from "../../components/atoms/Loading";
import { getSession } from "../../app/store/session/session.api";
import { putError } from "../../app/store/error/error.slice";
import { getPaginatedPcapSessions, useDeletePcapWorkspaceMutation } from "../../app/store/pcapworkspace/pcap-workspace.api";
import { Icon } from "../../components/atoms/Icon";
import { Select } from "../../components/molecules/Select";
import { useTranslation } from "react-i18next";
import { notifyInfo } from "../../app/notifications/notifier";
import { unsetPcapWorkspace } from "../../app/store/pcap/pcap.slice";
import { useAppNavigation } from "../../hooks/navigate";
import { Form } from "../../components/molecules/Form";
import { Overlay } from "../../components/atoms/Overlay";

export const Pcap = () => {
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const navigate = useAppNavigation();
    const [currentSession, setCurrentSession] = useState<Session | null>(null);
    const pcap = useAppSelector(state => state.rootReducer.pcap);
    const [getPaginatedPcapSessionsTrigger] = getPaginatedPcapSessions.useLazyQuery();
    const [deletePcapWorkspace] = useDeletePcapWorkspaceMutation();
    const [getSessionTrigger] = getSession.useLazyQuery();
    const [paginationIndex, setPaginationIndex] = useState(0);
    const [sessions, setSessions] = useState<Session[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [confirmPanel, setConfirmPanel] = useState({
        active: false,
        message: "",
        onConfirm: () => {},
    });

    const onClickSession = (session: Session) => {
        setIsLoading(true);
        getSessionTrigger(session.id).unwrap().then((data) => {
            setCurrentSession(data);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const onDeletePcapWorkspace = () => {
        deletePcapWorkspace(pcap.id)
        .unwrap()
        .then(() => {
            notifyInfo("pcap workspace has been deleted");
            dispatch(unsetPcapWorkspace());
            navigate('/home');
        }).catch((err) => {
            dispatch(putError(err.data.message));
        })
    }

    useEffect(() => {
        setIsLoading(true);
        getPaginatedPcapSessionsTrigger({
            workspaceId: pcap.id,
            paginationIndex: paginationIndex
        }).unwrap().then((res) => {
            setSessions(prev => paginationIndex > 0 ? [...res.sessions, ...prev] : [...res.sessions]);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        })
    }, [paginationIndex]);

    return (
        <div className="flex flex-col mx-4 my-2">
            <div className="flex flex-row w-full py-2 justify-start">
                <Icon 
                    size={30}
                    tip={pcap.status === "LISTENING" 
                        ? "pcap workspace is processing file..."
                        : "pcap workspace has processed file"}
                    name={pcap.status === "LISTENING" 
                        ? "headphones"
                        : "done"}
                />
                <Select 
                    className="relative mx-4"
                    icon="options"
                    items={[
                        {
                            text: t('delete') + "",
                            onItem: () => setConfirmPanel({
                                active: true,
                                message: 'Are you sure you want to delete this pcap workspace?',
                                onConfirm: onDeletePcapWorkspace
                            })
                        }
                    ]}
                />
                <Loading hidden={!isLoading} size={30} />
            </div>
            <div className="flex flex-row h-[88vh] w-full">
                <div className="w-1/3">
                    <SessionListPanel 
                        sessions={sessions}
                        onClick={onClickSession}
                        onExpand={() => setPaginationIndex(paginationIndex + 1)}
                    />
                </div>
                {
                    currentSession && (
                        <div className="w-2/3">
                            <SessionPanel 
                                session={currentSession}
                                filters={currentSession.matchedFilters}
                            />
                        </div>
                    )
                }
            </div>
            {
                confirmPanel.active && (
                    <Overlay>
                        <Form 
                            label={confirmPanel.message}
                            onCancel={() => setConfirmPanel({
                                ...confirmPanel,
                                active: false
                            })}
                            onSubmit={confirmPanel.onConfirm}
                        />
                    </Overlay>
                )
            }
        </div>
    );
}