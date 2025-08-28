import { useState } from "react";
import { Dropdown } from "./Dropdown";
import { useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text"; 
import { Icon } from "../atoms/Icon";
import { Icons } from "../../config/icons";
import { Item } from "../../app/types/Item";

export const Select = (props: {
    items: Item[],
    className?: string;
    icon?: keyof Icons;
    placeholder?: string;
    background?: string;
    remainOnSelect?: boolean;
    onHover?: () => void;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [dropdownActive, setDropdownActive] = useState(false);

    const onHover = () => {
        setDropdownActive(true)
        if (props.onHover) {
            props.onHover();
        }
    }

    return (
        <div 
            className={"cursor-pointer " + props.className}
            onMouseEnter={onHover}    
            onMouseLeave={() => setDropdownActive(false)}
            style={{
                backgroundColor: props.background
            }}
        >
            <div 
                className="flex flex-row duration-200"
                style={{
                    borderColor: dropdownActive ? theme.accents.neutral : theme.text
                }}
            >
                {
                    props.placeholder && (
                        <Text 
                            className="cursor-default text-xl font-bold"
                            color={dropdownActive ? theme.accents.neutral : theme.text}
                        >{props.placeholder}</Text>
                    )
                }
                {
                    props.icon && (
                        <Icon 
                            name={props.icon}
                            type={dropdownActive ? "positive" : "neutral"}
                            size={30}
                        />
                    )
                }
            </div>
            {
                dropdownActive && props.items?.length > 0 && (
                    <Dropdown 
                        items={props.items}
                        onClose={!props.remainOnSelect ? () => setDropdownActive(false) : undefined}
                    />
                )
            }
        </div>
    );
}