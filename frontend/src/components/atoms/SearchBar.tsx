import { useEffect, useState } from "react";
import { useAppSelector } from "../../app/store/hooks";
import { Icon } from "./Icon";

export const SearchBar = (props: {
    onChange: (value: string) => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [hovered, setHovered] = useState(false);
    const [focused, setFocused] = useState(false);
    const [searchValue, setSearchValue] = useState('');

    useEffect(() => {
        const timeout = setTimeout(() => {
            props.onChange(searchValue);
        }, 500);
        return () => clearTimeout(timeout); 
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [searchValue]);

    return (
        <div className="flex flex-row px-4 items-center w-1/4">
            <input 
                className="text-md w-full border-2 outline-none mr-1 rounded duration-200"
                style={{
                    backgroundColor: theme.secondary,
                    color: theme.text,
                    borderColor: focused ? theme.accents.contrast : hovered ? theme.accents.neutral : theme.tertiary
                }}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
                onFocus={() => setFocused(!focused)}
                onChange={(e) => setSearchValue(e.target.value)}
                placeholder="regex or substring"
            />
            <Icon name="search" size={30}/>
        </div>
    );
}