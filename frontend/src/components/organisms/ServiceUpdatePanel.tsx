import { useEffect, useState } from "react";
import { Input } from "../atoms/Input";
import { Service } from "../../app/types/Service";
import { getAllPlugins } from "../../app/store/plugin/plugin.api";
import { useUpdateServiceMutation } from "../../app/store/service/service.api";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { setContainerID, setService } from "../../app/store/service/service.slice";
import { Checkbox } from "../atoms/Checkbox";
import { Plugin } from "../../app/types/Plugin";
import { useTranslation } from "react-i18next";
import { Form } from "../molecules/Form";
import { Text } from "../atoms/Text";
import { randomKey, negativeColor } from "../../utils/utils";
import { CounterInput } from "../atoms/CounterInput";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { Loading } from "../atoms/Loading";

export const ServiceUpdatePanel = (props: {
    onClose: () => void
}) => {
    const dispatch = useAppDispatch();
    const theme = useAppSelector(state => state.rootReducer.theme);
    const service = useAppSelector(state => state.rootReducer.service);
    const { t } = useTranslation();
    const [state, setState] = useState<Service>(service);
    const [getAllPluginsTrigger] = getAllPlugins.useLazyQuery();
    const [plugins, setPlugins] = useState<Plugin[]>([]);
    const [updateService] = useUpdateServiceMutation();
    const [isLoading, setIsLoading] = useState(false);

    const onSubmit = () => {
        console.log(service.id);
        const selectedPlugins = plugins.filter(p => p.checked).map(p => p.name);
        updateService({
            ...state,
            containerID: null,
            plugins: selectedPlugins
        }).unwrap().then(() => {
            const name = state.name;
            dispatch(setService({
                ...state,
                plugins: selectedPlugins
            }));
            if (state.containerID) {
                dispatch(setContainerID(state.containerID));
            }
            notifyInfo(`service '${name} updated'`);
            props.onClose();
        }).catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    const onCheckboxChange = (plugin: Plugin) => {
        setPlugins(prev => prev.map(p => {
            if (p.name !== plugin.name) {
                return {
                    ...p,
                }
            }
            return {
                ...p,
                checked: !p.checked,
            }
        }));
    }

    useEffect(() => {
        setIsLoading(true);
        getAllPluginsTrigger().
            unwrap().
            then((res) => {
                setPlugins(res.plugins.map(p => {
                    return {
                        name: p,
                        checked: service.plugins.includes(p)
                    }
                }));
            }).catch((err) => {
                dispatch(putError(err.data.message));
            }).finally(() => {
                setIsLoading(false);
            });
    }, []);

    return (
        <Form
            label={t('update_service')}
            onCancel={props.onClose}
            onSubmit={onSubmit}
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
                disabled
            />
            {
                plugins && (
                    <div className="flex flex-col mb-4 items-center">
                        <Text className="font-semibold mb-2">{t('plugins')}</Text>
                        <div className="overflow-y-auto items-center w-1/2">
                            <div className="flex flex-wrap justify-center">
                                {
                                    isLoading ? <Loading hidden={false} size={20}/> : (
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
                                    )
                                }
                            </div>
                        </div>
                    </div>
                )
            } 
        </Form>
    );
}