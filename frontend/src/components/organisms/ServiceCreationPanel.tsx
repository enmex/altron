import { useEffect, useState } from "react";
import { scanService, scanServices, useCreateServiceMutation } from "../../app/store/service/service.api";
import { getAllPlugins } from "../../app/store/plugin/plugin.api";
import { Service } from "../../app/types/Service";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Input } from "../atoms/Input";
import { Button } from "../atoms/Button";
import { Checkbox } from "../atoms/Checkbox";
import { Plugin } from "../../app/types/Plugin";
import { useTranslation } from "react-i18next";
import { Form } from "../molecules/Form";
import { Panel } from "../atoms/Panel";
import { negativeColor, randomKey } from "../../utils/utils";
import { Text } from "../atoms/Text";
import { Icon } from "../atoms/Icon";
import { Loading } from "../atoms/Loading";
import { CounterInput } from "../atoms/CounterInput";
import { Select } from "../molecules/Select";
import { useAppNavigation } from "../../hooks/navigate";
import { setService } from "../../app/store/service/service.slice";
import { notifyInfo, notifyWarning } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const ServiceCreationPanel = (props: {
    onClose: () => void,
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const [state, setState] = useState<Service>({
        id: "",
        name: "",
        link: "",
        port: 0,
        plugins: [],
        workspaces: [],
    });
    const [isServicesScanning, setIsServicesScanning] = useState(false);
    const [scanningServiceIdx, setScanningServiceIdx] = useState(-1);
    const [hostServices, setHostServices] = useState<{
        container?: {
            id: string;
            name: string;
            image: string;
        };
        port: number;
        name: string;
        isPublic: boolean;
    }[]>([]);
    const [getAllPluginsTrigger] = getAllPlugins.useLazyQuery();
    const [plugins, setPlugins] = useState<Plugin[]>([]);
    const [pluginsLoading, setPluginsLoading] = useState(false);
    const [createService] = useCreateServiceMutation();
    const [scanServicesTrigger] = scanServices.useLazyQuery();
    const [scanServiceTrigger] = scanService.useLazyQuery();

    const navigate = useAppNavigation();

    const onSubmit = () => {
        createService({
            ...state,
        }).unwrap().then((data) => {
            notifyInfo(t('create_service_success', {serviceName: data.name}));
            dispatch(setService({
                ...data,
                workspaces: []
            }));
            navigate(`/services/${data.port}`);
            props.onClose();
        }).catch((err) => { 
            dispatch(putError(err.data.message));
        });
    }

    const onCheckboxChange = (plugin: Plugin) => {
        setPlugins(prev => plugins.map(p => {
            if (p.name !== plugin.name) {
                return {
                    ...p,
                }
            }
            return {
                ...p,
                checked: !p.checked,
            }
        }))
        setState({
            ...state,
            plugins: state.plugins.includes(plugin.name) 
                ? state.plugins.filter(p => p !== plugin.name) 
                : state.plugins.concat(plugin.name),
        });
    }

    const onScanServices = () => {
        setIsServicesScanning(true);
        scanServicesTrigger({})
        .unwrap()
        .then((data) => {
            setHostServices(data.services);
        })
        .catch((err) => {
            if (err.data.message.includes('agent')) {
                notifyInfo(err.data.message);
                return
            }
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsServicesScanning(false);
        });
    }

    const onScannedService = (idx: number, hostService: {
        container?: {
            id: string;
            name: string;
            image: string;
        };
        port: number;
        name: string;
        isPublic: boolean;
    }) => {
        setScanningServiceIdx(idx);
        scanServiceTrigger(hostService.port)
        .unwrap()
        .then((data) => {
            if (data) {
                setState({
                    ...state,
                    name: hostService.name,
                    link: data.link,
                    port: hostService.port,
                    containerID: hostService.container?.id
                })
            }
        })
        .catch((err) => {
            if (err.data.message.includes('agent')) {
                notifyWarning(err.data.message);
                return
            }
            dispatch(putError(err.data.message));
        })
        .finally(() => {
            setScanningServiceIdx(-1);
        });
    }

    useEffect(() => {
        onScanServices();
        setPluginsLoading(true);
        getAllPluginsTrigger()
        .unwrap()
        .then((data) => {
            setPlugins(prev => data.plugins.map((plugin) => {
                return {
                    name: plugin,
                    checked: false,
                }
            }));
        })
        .catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setPluginsLoading(false);
        });
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return (
        <div className="flex z-10 m-2">
            <Form 
                label={t('create_service')}
                onSubmit={onSubmit}
                onCancel={props.onClose}
            >
                <Input 
                    label={t('name')}
                    onChange={(e) => setState({
                        ...state,
                        name: e.target.value
                    })}
                    value={state.name}
                />
                <Input 
                    label={t('link')}
                    onChange={(e) => setState({
                        ...state,
                        link: e.target.value
                    })}    
                    value={state.link}
                />
                <CounterInput 
                    label={t('port')}
                    onChange={(value) => setState({
                        ...state,
                        port: value,
                    })}    
                    value={state.port}
                />
                <Button 
                    className="flex justify-center border-2 mb-2 p-2 rounded-md duration-200"
                    borderColor={theme.text}
                    onClick={onScanServices}
                >
                    <Icon name="scanEye" size={30}/>
                    <Text className="text-lg font-bold ml-1">{t('scan_localhost')}</Text>
                </Button>
                <Select 
                    items={hostServices.filter(service => service.container).map(service => {
                        return {
                            text: `${service.container?.name} (${service.container?.id})`,
                            onItem: () => setState({
                                ...state,
                                containerID: service.container?.id
                            })
                        }
                    })}
                    className="flex w-full justify-center border-2 my-2 py-2 rounded relative"
                    icon="arrowDropDown"
                    placeholder={state.containerID && state.containerID.length > 0 ? state.containerID : t('logs container').toString()}
                />
                {
                    !pluginsLoading ? (
                        <div className="flex flex-col mb-4 items-center">
                            <Text className="font-semibold mb-2">{t('plugins')}</Text>
                            <div className="overflow-y-auto items-center">
                                <div className="flex justify-center">
                                    {
                                        plugins.map((plugin) => {
                                            return (
                                                <Checkbox
                                                    key={randomKey()}
                                                    onChange={() => onCheckboxChange(plugin)}
                                                    checked={plugin.checked}
                                                >
                                                    <Text 
                                                        className="font-bold" 
                                                        color={plugin.checked ? negativeColor(theme.text) : theme.text}
                                                    >{plugin.name}</Text>
                                                </Checkbox>
                                            );
                                        })
                                    }
                                </div>
                            </div>
                        </div>
                    ) : (
                        <div className="flex w-full h-full justify-center items-center">
                            <Loading hidden={false} size={50} />
                        </div>
                    )
                }
            </Form>
            <Panel withBorder>
                <Text className="cursor-default font-bold text-xl">{t('localhost_services')}</Text>
                <div className="flex flex-grow w-full h-2 overflow-auto">
                    {
                        isServicesScanning ? (
                            <div className="flex w-full h-full justify-center items-center">
                                <Loading hidden={!isServicesScanning} size={50} />
                            </div>
                        ) : 
                        (
                            <div className="flex flex-col">
                                {
                                    hostServices.slice().sort((a, b) => a.port - b.port).map((hostService, idx) => {
                                        return (
                                            <Button 
                                                key={randomKey()}
                                                className="flex flex-col w-full p-2 border-b-2 duration-200 my-2"
                                                backgroundColor={theme.secondary}
                                                borderColor={theme.text}
                                                onClick={() => onScannedService(idx, hostService)} 
                                            >
                                                <div className="flex flex-col text-start">
                                                    <div className="flex flex-row items-center w-3/4 justify-between">
                                                        <Text 
                                                            className="text-lg mr-2"
                                                            color={theme.text}
                                                        >{t('port')}</Text>
                                                        <Text 
                                                            className="text-lg font-extrabold mr-2"
                                                            color={theme.accents.contrast}
                                                        >{hostService.port}</Text>
                                                        <Icon 
                                                            tip={hostService.isPublic ? "public port" : "private port"}
                                                            type="contrast"
                                                            name={hostService.isPublic ? 'unlock' : 'lock'}
                                                            size={15}
                                                        />
                                                        <Loading size={20} hidden={idx !== scanningServiceIdx}/>
                                                    </div>
                                                    <div className="flex flex-row items-center">
                                                        <Text 
                                                            className="font-bold mr-2"
                                                            color={theme.text}
                                                        >{hostService.name}</Text>
                                                        {
                                                            hostService.container && (
                                                                <Icon 
                                                                    tip="dockerized service"
                                                                    type="neutral"
                                                                    name="docker"
                                                                    size={20}
                                                                />
                                                            )
                                                        }
                                                    </div>
                                                    
                                                    {
                                                        hostService.container && (
                                                            <>
                                                            <Text 
                                                                className="font-bold"
                                                                color={theme.accents.positive}
                                                            >{`${hostService.container.id}`}</Text>
                                                            <Text 
                                                                className="font-bold"
                                                                color={theme.accents.positive}
                                                            >{`image: ${hostService.container.image}`}</Text>
                                                            </>
                                                        )
                                                    }
                                                </div>
                                            </Button>
                                        );
                                    })
                                }
                            </div>
                        )
                    }
                </div>
            </Panel>
        </div>
    );
}