import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useCreateFilterMutation } from "../../app/store/filter/filter.api";
import { Input } from "../atoms/Input";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Form } from "../molecules/Form";
import { ColorPicker } from "../molecules/ColorPicker";
import { Select } from "../molecules/Select";
import { Checkbox } from "../atoms/Checkbox";
import { Text } from "../atoms/Text";
import { negativeColor } from "../../utils/utils";
import { Filter } from "../../app/types/Filter";
import { CounterInput } from "../atoms/CounterInput";
import { useGetAnalyzerPayloadQuery } from "../../app/store/analyzer/analyzer.api";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const FilterCreationPanel = (props: {
    onClose: (filter?: Filter) => void;
}) => {
    const dispatch = useAppDispatch();
    const theme = useAppSelector(state => state.rootReducer.theme);
    const service = useAppSelector(state => state.rootReducer.service);
    const workspace = useAppSelector(state => state.rootReducer.workspace);
    const {data: totalRequestsPayload} = useGetAnalyzerPayloadQuery({
        servicePort: workspace.id.length > 0 ? undefined : service.port,
        workspaceId: workspace.id.length > 0 ? workspace.id : undefined
    });
    const { t } = useTranslation();
    const [state, setState] = useState<{
        name: string;
        serviceId?: string;
        regex?: string;
        ttl?: number;
        totalPackets?: number;
        inRequest: boolean;
        inResponse: boolean;
        color: string;
        isBlocking: boolean;
    }>({
        name: "",
        inRequest: false,
        inResponse: false,
        color: "#365F88",
        isBlocking: false,
    });
    const [createFilter] = useCreateFilterMutation();

    const onSubmit = () => {
        if (!(state.regex || state.ttl || state.totalPackets)) {
            dispatch(putError('one of the field is required'));
            return;
        }
        if (!state.inRequest && !state.inResponse && !(!state.regex && !state.ttl && state.totalPackets)) {
            dispatch(putError('choose search type'));
            return;
        }
        createFilter({
            ...state,
            serviceID: state.serviceId,
        }).unwrap().then((data) => {
            notifyInfo(t('create_filter_success', {filterName: state.name}))
            props.onClose({
                id: "",
                ...state,
            });
        }).catch((err) => {
            dispatch(putError(err.data.message));
        })
    }

    const onSelect = (item: {
        request: boolean,
        response: boolean
    }) => {
        setState({
            ...state,
            inRequest: item.request,
            inResponse: item.response
        });
    }

    return (
        <Form
            label={t('create_filter')}
            onCancel={() => props.onClose()}
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
                label={t('regex')}
                onChange={(e) => setState({
                    ...state,
                    regex: e.target.value
                })}    
                value={state.regex}
                notRequired
            />
            <CounterInput 
                label="TTL"
                onChange={(value) => setState({
                    ...state,
                    ttl: value
                })}
            />
            <div className="flex flex-row items-center">
                <CounterInput 
                    label="Total packets"
                    onChange={(value) => setState({
                        ...state,
                        totalPackets: value
                    })}
                    value={state.totalPackets}
                />
                <Select 
                    className="relative py-1 rounded"
                    items={totalRequestsPayload ? totalRequestsPayload.analyzer['requests'].map((characteristic) => {
                        return {
                            text: `${characteristic.value} (${characteristic.number})`,
                            onItem: () => setState({
                                ...state,
                                totalPackets: Number(characteristic.value)
                            })
                        }
                    }) : []}
                    icon="arrowDropDown"
                />
            </div>
            <ColorPicker onChange={(color: string) => setState({
                ...state,
                color: color
            })}/>
            <Checkbox
                onChange={() => setState(({
                    ...state,
                    isBlocking: !state.isBlocking,
                }))}
                checked={state.isBlocking}
            >
                <Text 
                    className="font-bold" 
                    color={state.isBlocking ? negativeColor(theme.text) : theme.text}
                >{t('block_content')}</Text>
            </Checkbox>
            <Checkbox
                onChange={() => setState(({
                    ...state,
                    serviceId: state.serviceId ? undefined : service.id,
                }))}
                checked={!(state.serviceId)}
            >
                <Text 
                    className="font-bold" 
                    color={state.serviceId ? theme.text : negativeColor(theme.text)}
                >{t('global')}</Text>
            </Checkbox>
            <div className="flex w-full justify-center">
                <Select 
                    className="relative my-4"
                    items={[
                        {
                            text: t('both') + "",
                            onItem: () => onSelect({
                                request: true,
                                response: true
                            })
                        },
                        {
                            text: t('request') + "",
                            onItem: () => onSelect({
                                request: true,
                                response: false
                            })
                        },
                        {
                            text: t('response') + "",
                            onItem: () => onSelect({
                                request: false,
                                response: true
                            })
                        }
                    ]}
                    placeholder={state.inRequest && state.inResponse 
                        ? t('both') ?? ""
                        : state.inRequest 
                        ? t('request') ?? ""
                        : state.inResponse 
                        ? t('response') ?? ""
                        : t('search_type') ?? ""
                    }
                    icon="arrowDropDown"
                />
            </div>
        </Form>
    );
}