import { useAppSelector } from "../../app/store/hooks";

export const Range = (props:{
    min: number;
    max: number;
    value?: number;
    label: string;
    onChange: (value: number) => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);

    return (
        <div
            className="flex flex-col items-start"
            style={{
                color: theme.text
            }}
        >
            <span>{props.label}</span>
            <input
                type="range"
                className="w-full h-2 mt-2 rounded-lg appearance-none cursor-pointer"
                style={{
                    backgroundColor: theme.accents.neutral
                }}
                min={props.min}
                max={props.max}
                onChange={(e) => props.onChange(Number(e.target.value))}
                value={props.value}
            />
        </div>
    )
}