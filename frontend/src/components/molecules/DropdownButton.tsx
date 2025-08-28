import { ReactNode, useEffect, useState } from "react";
import { Button } from "../atoms/Button";
import { Dropdown } from "./Dropdown";
import { Item } from "../../app/types/Item";

export const DropdownButton = (props: {
    onClick?: () => void;
    icon?: ReactNode;
    text?: ReactNode;
    dropdownItems: Item[];
}) => {
    const [mouseEnter, setMouseEnter] = useState(false);
    const [dropdownVisible, setDropdownVisible] = useState(false);

    useEffect(() => {
        if (mouseEnter) {
            const timeout = setTimeout(() => {
                setDropdownVisible(true);
            }, 500);

            return () => clearTimeout(timeout);
        } else {
            setDropdownVisible(false);
        }
    }, [mouseEnter]);

    return (
        <div
            onMouseEnter={() => setMouseEnter(true)}    
            onMouseLeave={() => setMouseEnter(false)}
        >
            <Button 
                onClick={props.onClick}
            >{props.icon}{props.text}</Button>
            <div className="absolute w-24">
            {
                dropdownVisible && <Dropdown items={props.dropdownItems} />
            }
            </div>
        </div>
    );
}