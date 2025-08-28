import { useEffect, useState } from "react"
import { useAppSelector } from "../../app/store/hooks"
import { Button } from "./Button"
import { Text } from "./Text"

export const CounterInput = (props: {
    label: string
    disabled?: boolean
    onChange: (value: number) => void
    value?: number
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [state, setState] = useState<number>(props.value ?? 0);

    const onIncrement = () => {
        const newState = state + 1;
        onChange(newState);
    }

    const onDecrement = () => {
        const newState = state - 1 < 0 ? 0 : state - 1;
        onChange(newState);
    }

    const onChange = (value: number) => {
        setState(value);
        props.onChange(value);
    }

    useEffect(() => {
        setState(props.value ?? 0);
    }, [props.value]);

    return (
        <div className="flex w-full justify-center">
            <div 
                className="h-18 w-44 mb-4"
            >
                <Text
                    className="flex justify-start"
                >{props.label}</Text>
                <div 
                    className="flex flex-row rounded"
                    style={{
                        backgroundColor: theme.tertiary
                    }}
                >
                    <Button 
                        className="h-full w-20 rounded-l cursor-pointer outline-none"
                        onClick={props.disabled ? () => {} : onDecrement}
                    >
                        <Text className="m-auto text-3xl font-thin">-</Text>
                    </Button>
                    <input 
                        type="number" 
                        className="text-xl bg-transparent outline-none text-center w-1/2 font-semibold text-md md:text-basecursor-default flex items-center"
                        value={state?.toString()}
                        onChange={(e) => onChange(Number(e.target.value))}
                        disabled={props.disabled}
                        style={{
                            color: theme.text
                        }}
                        required
                        min={0}
                    />
                    <Button 
                        className="h-full w-20 rounded-l cursor-pointer outline-none"
                        onClick={props.disabled ? () => {} : onIncrement}
                    >
                        <Text className="m-auto text-3xl font-thin">+</Text>
                    </Button>
                </div>
            </div>
        </div>
    )
}