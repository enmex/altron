import { Select } from "./Select";
import { useAppSelector } from "../../app/store/hooks";
import { useEffect, useState } from "react";
import { Characteristic } from "../../app/types/Analyzer";
import { randomKey } from "../../utils/utils";

export const AnalyzerPanel = (props:{
    analyzerPayload?: {
        hasChecker: boolean,
        analyzer: {
            [componentName: string]: Characteristic[]
        },
    },
    onClick?: (characteristic: {
        componentName: string,
        value: string
    }, action: string, reset: boolean) => void;
    onUpdate?: (currentCharacteristics: {
        [componentName: string]: {
            value: string;
            selected: boolean;
            blocked: boolean;
        }[]
    }) => void;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [currentCharacteristics, setCurrentCharacteristics] = useState<{
        [componentName: string]: {
            value: string;
            selected: boolean;
            blocked: boolean;
        }[]
    }>(props.analyzerPayload ? Object.keys(props.analyzerPayload.analyzer).reduce(
        (acc, componentName) => ({
            ...acc,
            [componentName]: [],
        }),
        {}
    ) : {});

    const onSelectCharacteristic = (characteristic: {
        componentName: string,
        value: string
    }) => {
        if (props.onClick) {
            props.onClick(
                characteristic, 
                'analyzer-pass', 
                !!(currentCharacteristics[characteristic.componentName].find(ch => ch.selected && ch.value === characteristic.value)),
            );
        }
        
        const component = currentCharacteristics[characteristic.componentName];
        if (!component.find(ch => ch.value === characteristic.value)) {
            component.push({
                value: characteristic.value,
                blocked: false,
                selected: false
            });
        }

        component.forEach(ch => {
            if (ch.value === characteristic.value) {
                ch.selected = !ch.selected;
                ch.blocked = false;
            }
        });

        currentCharacteristics[characteristic.componentName] = component;

        setCurrentCharacteristics({
            ...currentCharacteristics
        });
        if (props.onUpdate) {
            props.onUpdate(currentCharacteristics);
        }
    }

    const onBlockCharacteristic = (characteristic: {
        componentName: string,
        value: string
    }) => {
        if (props.onClick) {
            props.onClick(
                characteristic, 
                'analyzer-block',
                !!(currentCharacteristics[characteristic.componentName].find(ch => ch.blocked && ch.value === characteristic.value))
            );
        }
        
        const component = currentCharacteristics[characteristic.componentName];
        if (!component.find(ch => ch.value === characteristic.value)) {
            component.push({
                value: characteristic.value,
                blocked: false,
                selected: false
            });
        }

        component.forEach(ch => {
            if (ch.value === characteristic.value) {
                ch.blocked = !ch.blocked;
                ch.selected = false;
            }
        });

        currentCharacteristics[characteristic.componentName] = component;

        setCurrentCharacteristics({
            ...currentCharacteristics
        });
        if (props.onUpdate) {
            props.onUpdate(currentCharacteristics);
        }
    }

    useEffect(() => {
        if (props.onUpdate) {
            props.onUpdate(currentCharacteristics);
        }
    }, [currentCharacteristics]);

    useEffect(() => {
        setCurrentCharacteristics(props.analyzerPayload ? Object.keys(props.analyzerPayload.analyzer).reduce(
            (acc, componentName) => ({
                ...acc,
                [componentName]: [],
            }),
            {}
        ) : {});
    }, [props.analyzerPayload]);

    return (
        <div 
            onContextMenu={(e) => e.preventDefault()}
            className="flex flex-row justify-start w-full items-center ml-4"
        >
            {
                props.analyzerPayload ? Object.keys(props.analyzerPayload.analyzer).map(component => {
                    return (
                        <Select 
                            remainOnSelect
                            key={randomKey()}
                            items={props.analyzerPayload ? props.analyzerPayload.analyzer[component].map(characteristic => {
                                return {
                                    text: `${characteristic.value} (${characteristic.number})`,
                                    icon: props.analyzerPayload?.hasChecker ? {
                                        value: "warning",
                                        color: characteristic.isSafe ? theme.accents.positive : theme.accents.negative
                                    } : undefined,
                                    onItem: () => onSelectCharacteristic({
                                        componentName: component,
                                        value: characteristic.value
                                    }),
                                    onContextMenu: () => onBlockCharacteristic({
                                        componentName: component,
                                        value: characteristic.value
                                    }),
                                    color: currentCharacteristics[component]?.find(ch => ch.value === characteristic.value && ch.selected)
                                        ? theme.accents.positive 
                                        : currentCharacteristics[component]?.find(ch => ch.value === characteristic.value && ch.blocked)
                                        ? theme.accents.negative 
                                        : theme.secondary
                                }
                            }) : []}
                            background={theme.secondary}
                            className="mr-2 px-2 rounded"
                            placeholder={component.toUpperCase()}
                        />
                    )
                }) : <></>
            }
        </div>
    );
}