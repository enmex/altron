import { useState } from "react";
import { useCreateWorkspaceMutation } from "../../app/store/workspace/workspace.api";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Input } from "../atoms/Input";
import { useTranslation } from "react-i18next";
import { Form } from "../molecules/Form";
import { Text } from "../atoms/Text";
import { CreateWorkspaceRequest } from "../../app/store/workspace/workspace.types";
import { getNearestHour, negativeColor, randomKey } from "../../utils/utils";
import { Checkbox } from "../atoms/Checkbox";
import { useAppNavigation } from "../../hooks/navigate";
import { Loading } from "../atoms/Loading";
import { Panel } from "../atoms/Panel";
import { setWorkspace } from "../../app/store/workspace/workspace.slice";
import { Range } from "../atoms/Range";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { putWorkspace } from "../../app/store/service/service.slice";

export const WorkspaceCreationPanel = (props:{
    servicePort: number,
    onClose: () => void,
}) => {
    const intervals = [0, 15, 30, 45, 60];
    const presetIntervals = [3, 5, 8, 10];
    const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
    const service = useAppSelector(state => state.rootReducer.service);
    const theme = useAppSelector(state => state.rootReducer.theme);
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const [state, setState] = useState<CreateWorkspaceRequest>({
        name: "",
        timeout: "8m",
        servicePort: props.servicePort,
        startTime: null
    });
    const [createWorkspace] = useCreateWorkspaceMutation();
    const [inputDisabled, setInputDisabled] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useAppNavigation();

    const onSubmit = () => {
        if (!inputDisabled && state.name === 'checker') {
            dispatch(putError(t('checker_workspace_only_by_checkbox')));
            return;
        }
        createWorkspace({
            ...state
        }).unwrap()
        .then((data) => {
            if (data) {
                notifyInfo(t('create_workspace_success', {workspaceName: data.name}));
                dispatch(setWorkspace(data));
                props.onClose();
                navigate(`/workspaces/${data.id}`);
            }
        })
        .catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    const onClickCheckerCheckbox = () => {
        setInputDisabled(prev => !prev);
        setState({
            ...state,
            name: state.name === 'checker' ? '' : 'checker'
        })
    }

    const onClickPreset = (workspaceName: string) => {
        workspaceName = workspaceName.toLowerCase().replaceAll(' ', '_') + `_${service.workspaces.length}`;
        const timeout = workspaceName.split('_')[2];
        setIsLoading(true);
        createWorkspace({
            name: workspaceName,
            servicePort: props.servicePort,
            timeout: timeout,
            startTime: null
        }).unwrap()
        .then((data) => {
            notifyInfo(t('create_workspace_success', {workspaceName: data.name}));
            props.onClose();
            dispatch(putWorkspace(data));
            dispatch(setWorkspace(data));
            navigate(`/workspaces/${data.id}`);
        })
        .catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    return (
        <div className="flex z-10 m-2">
            <Form
                label={t('create_workspace')}
                onCancel={props.onClose}
                onSubmit={onSubmit}
            >
                <Input 
                    label={t('name')}
                    onChange={(e) => setState({
                        ...state,
                        name: e.target.value
                    })}
                    disabled={inputDisabled}
                    value={state.name}
                />
                <Checkbox
                    key={randomKey()}
                    onChange={onClickCheckerCheckbox}
                    checked={inputDisabled}
                >
                    <Text 
                        className="font-bold" 
                        color={!inputDisabled ? theme.text : negativeColor(theme.text)}
                    >CHECKER</Text>
                </Checkbox>
                <Text className="flex justify-start">{t('listening_start_time')}</Text>
                <div className="flex justify-center">
                    <div className="flex flex-wrap w-2/3 items-center justify-center">
                        <Checkbox
                            key={randomKey()}
                            onChange={() => setState({
                                ...state,
                                startTime: null
                            })}
                            checked={state.startTime == null}
                        >
                            <Text 
                                className="font-bold" 
                                color={state.startTime ? theme.text : negativeColor(theme.text)}
                            >NOW</Text>
                        </Checkbox>
                        {
                            intervals.map((interval) => {
                                const time = getNearestHour(interval);
                                return (
                                    <Checkbox
                                        key={randomKey()}
                                        onChange={() => setState({
                                            ...state,
                                            startTime: time.getTime()
                                        })}
                                        checked={time.getTime() === state.startTime}
                                    >
                                        <Text 
                                            className="font-bold" 
                                            color={state.startTime === time.getTime() ? negativeColor(theme.text) : theme.text}
                                        >
                                            {time.toLocaleTimeString([], {
                                                timeZone: timezone,
                                                hour: '2-digit', 
                                                minute: '2-digit'
                                            })}
                                        </Text>
                                    </Checkbox>
                                )
                            })
                        }
                    </div>
                </div>
                <Range 
                    min={1} 
                    max={15} 
                    value={Number(state.timeout.slice(0, -1))} 
                    label={t('listening_timeout')}
                    onChange={(minutes) => setState({
                        ...state,
                        timeout: `${minutes}m`
                    })}
                />
                <Text className="font-bold text-xl">{`${state.timeout}in`}</Text>
                <Loading hidden={!isLoading} size={30} />
            </Form>
            <Panel withBorder>
                <Text className="font-bold text-xl mb-2">{t('presets')}</Text>
                <div className="flex flex-col justify-center">
                    {
                        presetIntervals.map((presetInterval) => {
                            const workspaceName = `Find exploit ${presetInterval}m`;
                            return (
                                <Checkbox
                                    onChange={() => onClickPreset(workspaceName)}
                                    checked={state.name === workspaceName}
                                >
                                    <Text 
                                        className="font-bold" 
                                        color={state.name === workspaceName ? negativeColor(theme.text) : theme.text}
                                    >{workspaceName}</Text>
                                </Checkbox>
                            )
                        })
                    }
                </div>
            </Panel>
        </div>
    );
}