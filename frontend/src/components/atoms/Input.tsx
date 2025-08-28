import { useState } from "react";
import { useAppSelector } from "../../app/store/hooks";
import { randomKey } from "../../utils/utils";

export const Input = (props: {
    type?: string;
    label: string,
    value?: any | null,
    notRequired?: boolean
    disabled?: boolean
    onChange: (value: any) => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [key] = useState(randomKey());
    const [hovered, setHovered] = useState(false);
    const [focused, setFocused] = useState(false);

    return (
        <div className="flex flex-col items-start w-full my-6 relative group">
            <input 
                id={key}
                className="w-full rounded border-2 outline-none py-2 duration-200 peer"
                style={{
                    color: theme.text,
                    backgroundColor: theme.secondary,
                    borderColor: focused ? theme.accents.contrast : hovered ? theme.accents.neutral : theme.tertiary
                }}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
                onFocus={() => setFocused(!focused)}
                type={props.type ?? "text"}
                onChange={props.onChange} 
                value={props.value} 
                required={!props.notRequired}
                disabled={props.disabled}
            />
            <label 
                htmlFor={key}
                className="
                    cursor-text
                    font-extralight text-lg transform transition-all absolute 
                    top-0 left-0 h-full flex items-center pl-2 group-focus-within:text-lg 
                    group-focus-within:font-bold group-focus-within:h-1/2 peer-valid:h-1/2 
                    group-focus-within:-translate-y-full peer-valid:font-bold peer-valid:-translate-y-full 
                    group-focus-within:pl-0 peer-valid:pl-0"
                style={{
                    color: theme.text
                }}
            >{props.label}</label>
        </div>
    );
}