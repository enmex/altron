import { useEffect, useState } from "react";
import { Session } from "../../app/types/Service";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Overlay } from "../../components/atoms/Overlay";
import { Button } from "../../components/atoms/Button";
import { SessionListPanel } from "../../components/organisms/SessionListPanel";
import { SessionPanel } from "../../components/organisms/SessionPanel";
import { getPaginatedSessions, useClearWorkspaceSessionsMutation, useDeleteWorkspaceMutation, useSearchWorkspaceSessionsMutation } from "../../app/store/workspace/workspace.api";
import { WorkspaceUpdatePanel } from "../../components/organisms/WorkspaceUpdatePanel";
import { useTranslation } from "react-i18next";
import { FilterPanel } from "../../components/molecules/FilterPanel";
import { Filter } from "../../app/types/Filter";
import { Form } from "../../components/molecules/Form";
import { useAppNavigation } from "../../hooks/navigate";
import { Icon } from "../../components/atoms/Icon";
import { Select } from "../../components/molecules/Select";
import { SearchBar } from "../../components/atoms/SearchBar";
import { Loading } from "../../components/atoms/Loading";
import { AnalyzerPanel } from "../../components/molecules/AnalyzerPanel";
import { getSession } from "../../app/store/session/session.api";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { useGetWorkspaceAnalyzerPayloadQuery } from "../../app/store/analyzer/analyzer.api";
import { useAddWorkspaceSessionsToCartMutation } from "../../app/store/cart/cart.api";

export const Workspace = () => {
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const [currentSession, setCurrentSession] = useState<Session | null>(null);
    const [currentFilter, setCurrentFilter] = useState<Filter | null>(null);
    const workspace = useAppSelector(state => state.rootReducer.workspace);
    const [getPaginatedSessionsTrigger] = getPaginatedSessions.useLazyQuery();
    const [searchWorkspaceSessionsTrigger] = useSearchWorkspaceSessionsMutation();
    const [getSessionTrigger] = getSession.useLazyQuery();
    const [paginationIndex, setPaginationIndex] = useState(0);
    const [sessions, setSessions] = useState<Session[]>([]);
    const [confirmPanel, setConfirmPanel] = useState({
        active: false,
        message: "",
        onConfirm: () => {},
    });
    const [updateWorkspacePanelActive, setUpdateWorkspacePanelActive] = useState(false);
    const [clearWorkspaceSessions] = useClearWorkspaceSessionsMutation();
    const [deleteWorkspace] = useDeleteWorkspaceMutation();
    const [isLoading, setIsLoading] = useState(false);
    const [searchValue, setSearchValue] = useState("");
    const {data: analyzerPayload} = useGetWorkspaceAnalyzerPayloadQuery({
        workspaceId: workspace.id
    });
    const [currentCharacteristics, setCurrentCharacteristics] = useState<{
        [componentName: string]: {
            value: string;
            selected: boolean;
            blocked: boolean;
        }[]
    }>(analyzerPayload ? Object.keys(analyzerPayload.analyzer).reduce(
        (acc, componentName) => ({
            ...acc,
            [componentName]: analyzerPayload?.analyzer[componentName].map((characteristic) => ({
                value: characteristic.value, 
                selected: false,
                blocked: false, 
            })),
        }),
        {}
    ) : {});
    const [addSessionsToCart] = useAddWorkspaceSessionsToCartMutation();

    const navigate = useAppNavigation();

    const onClickFilter = (filter: Filter) => {
        setPaginationIndex(prev => 0);
        if (currentFilter && currentFilter.id === filter.id) {
            setCurrentFilter(null);
            return;
        }
        setCurrentFilter(prev => filter);
    }

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

    const onDeleteSessionsConfirm = () => {
        clearWorkspaceSessions({
            workspaceId: workspace.id,
        }).unwrap().then(() => {
            setCurrentSession(null);
            setSessions(prev => []);
            notifyInfo(t('delete_sessions_success'));
        }).catch((err) => {
            dispatch(putError(err.data.message));
        });
        return;
    }

    const onClickTrashCan = () => {
        setConfirmPanel({
            active: true,
            message: t('delete_sessions_confirm'),
            onConfirm: onDeleteSessionsConfirm
        });
    }

    const onDeleteWorkspace = () => {
        setConfirmPanel({
            ...confirmPanel,
            active: false,
        });
        deleteWorkspace(workspace.id).unwrap().then(() => {
            notifyInfo(t('delete_workspace_success', {workspaceName: workspace.name}));
            navigate('/home');
        }).catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    useEffect(() => {
        onSearch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [searchValue, currentFilter, currentCharacteristics, paginationIndex]);

    const onSearch = () => {
        setIsLoading(true);
        if (!searchValue && !currentFilter && Object.values(currentCharacteristics).every((chs) => chs.length === 0)) {
            getPaginatedSessionsTrigger({
                workspaceId: workspace.id,
                paginationIndex: paginationIndex
            }).unwrap().then((data) => {
                setSessions(prev => paginationIndex > 0 ? [...data.sessions, ...prev] : [...data.sessions]);
            }).catch((err) => {
                dispatch(putError(err.data.message));
            }).finally(() => {
                setIsLoading(false);
            });
            return;
        }
        searchWorkspaceSessionsTrigger({
            workspaceId: workspace.id,
            searchValue: searchValue,
            pagination: paginationIndex,
            filterId: currentFilter?.id,
            selectedCharacteristics: currentCharacteristics
        }).unwrap().then((data) => {
            setSessions(prev => [...data.sessions]);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const onClickMerge = () => {
        setIsLoading(true);
        addSessionsToCart({
            workspaceId: workspace.id,
            sessions: sessions.map(s => s.id)
        }).unwrap().then(() => {
            notifyInfo('sessions added to cart');
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        })
    }

    return (
        <div className="flex flex-col mx-4 my-2">
            <div className="flex flex-row w-full py-2 justify-between">
                <Button
                    onClick={onClickTrashCan}
                >
                    <Icon tip="clear" type="negative" name="trash" size={30}/>
                </Button>
                <Button
                    onClick={onClickMerge}
                >
                    <Icon tip="add sessions to cart" type="contrast" name="merge" size={30}/>
                </Button>
                <Button
                    onClick={() => navigate(`/carts/${workspace.id}`)}
                >
                    <Icon tip="cart" name="basket" type="contrast" size={30} />
                </Button>
                <Select 
                    className="relative"
                    icon="options"
                    items={[
                        {
                            text: t('update') + "",
                            onItem: () => setUpdateWorkspacePanelActive(true)
                        },
                        {
                            text: t('delete') + "",
                            onItem: () => setConfirmPanel({
                                active: true,
                                message: t('delete_workspace_confirm'),
                                onConfirm: onDeleteWorkspace
                            })
                        }
                    ]}
                />
                <SearchBar 
                    onChange={setSearchValue}
                />
                <AnalyzerPanel 
                    analyzerPayload={analyzerPayload}
                    onUpdate={setCurrentCharacteristics}
                />
                <Loading className="flex" hidden={!isLoading} size={30} />
                <FilterPanel 
                    onClick={onClickFilter}
                    selectedFilter={currentFilter ?? undefined}
                    onUpdate={() => {}}
                />
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
            {
                updateWorkspacePanelActive && (
                    <Overlay>
                        <WorkspaceUpdatePanel 
                            workspace={workspace}
                            onClose={() => setUpdateWorkspacePanelActive(false)}
                        />
                    </Overlay>
                )
            }
        </div>
    );
}