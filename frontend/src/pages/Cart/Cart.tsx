import { useEffect, useState } from "react";
import { Session } from "../../app/types/Service";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Overlay } from "../../components/atoms/Overlay";
import { Button } from "../../components/atoms/Button";
import { SessionListPanel } from "../../components/organisms/SessionListPanel";
import { SessionPanel } from "../../components/organisms/SessionPanel";
import { useTranslation } from "react-i18next";
import { Form } from "../../components/molecules/Form";
import { useAppNavigation } from "../../hooks/navigate";
import { Icon } from "../../components/atoms/Icon";
import { Loading } from "../../components/atoms/Loading";
import { getSession } from "../../app/store/session/session.api";
import { getWorkspaceCart, useClearWorkspaceCartMutation, useMergeSessionsMutation } from "../../app/store/cart/cart.api";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const Cart = () => {
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const [currentSession, setCurrentSession] = useState<Session | null>(null);
    const workspace = useAppSelector(state => state.rootReducer.workspace);
    const [getSessionTrigger] = getSession.useLazyQuery();
    const [getCartTrigger] = getWorkspaceCart.useLazyQuery();
    const [clearCart] = useClearWorkspaceCartMutation();
    const [paginationIndex, setPaginationIndex] = useState(0);
    const [sessions, setSessions] = useState<Session[]>([]);
    const [confirmPanel, setConfirmPanel] = useState({
        active: false,
        message: "",
        onConfirm: () => {},
    });
    const [isLoading, setIsLoading] = useState(false);
    const [mergeSessions] = useMergeSessionsMutation();

    const navigate = useAppNavigation();

    const onClickSession = (session: Session) => {
        setIsLoading(true);
        getSessionTrigger(session.id).unwrap().then((data) => {
            setCurrentSession(data);
            setIsLoading(false);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const onDeleteSessionsConfirm = () => {
        setIsLoading(true);
        clearCart(workspace.id)
        .unwrap().then(() => {
            setCurrentSession(null);
            setSessions(prev => []);
            notifyInfo(t('delete_sessions_success'));
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const onClickTrashCan = () => {
        if (sessions.length === 0) {
            return;
        }
        setConfirmPanel({
            active: true,
            message: t('delete_sessions_confirm'),
            onConfirm: onDeleteSessionsConfirm
        });
    }

    useEffect(() => {
        setIsLoading(true);
        getCartTrigger({
            workspaceID: workspace.id,
            pagination: paginationIndex
        }).unwrap().then((res) => {
            setSessions(res.sessions);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [paginationIndex]);

    const onClickMerge = () => {
        if (sessions.length === 0) {
            return;
        }
        setIsLoading(true);
        mergeSessions({
            workspaceID: workspace.id,
            sessions: sessions.map(s => s.id)
        }).unwrap().then((data) => {
            setSessions(prev => []);
            window.open(`http://${window.location.host}/sessions/${data.sessionID}`, '_blank');
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        })
    }

    return (
        <div className="flex flex-col mx-4 my-2">
            <div className="flex flex-row w-full py-2 justify-start">
                <Button
                    onClick={onClickTrashCan}
                >
                    <Icon tip="clear sessions" type="negative" name="trash" size={30}/>
                </Button>
                <Button
                    onClick={onClickMerge}
                >
                    <Icon tip="merge sessions" type="contrast" name="merge" size={30}/>
                </Button>
                <Button
                    onClick={() => navigate(`/workspaces/${workspace.id}`)}
                >
                    <Icon tip="return to workspace" name="workspace" type="contrast" size={30} />
                </Button>
                <Loading className="flex" hidden={!isLoading} size={30} />
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
                                active: false,
                            })}
                            onSubmit={() => {
                                confirmPanel.onConfirm();
                                setConfirmPanel({
                                    ...confirmPanel,
                                    active: false,
                                });
                            }}
                        />
                    </Overlay>
                )
            }
        </div>
    );
}