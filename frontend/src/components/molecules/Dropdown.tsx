import { useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text";
import { Button } from "../atoms/Button";
import { Icon } from "../atoms/Icon";
import { Item } from "../../app/types/Item";
import { useState } from "react";
import { randomKey } from "../../utils/utils";

export const Dropdown = (props: {
    items: Item[]
    className?: string
    onClose?: () => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [hovered, setHovered] = useState(-1);

    const onClick = (item: Item) => {
        item.onItem();
        
        if (props.onClose) {
            props.onClose();
        }
    }

    return (
        <div 
            className={props.className ?? "animate-fade absolute z-[100] divide-y overflow-auto max-h-96 rounded-md"}
            onMouseLeave={() => setHovered(-1)}
        >
            {
                props.items.map((item, idx) => {
                    return (
                        <div onContextMenu={item.onContextMenu} key={randomKey()}>
                        <Button 
                            onClick={() => onClick(item)}
                            className="rounded-sm font-bold text-lg flex w-full items-center justify-center py-2 px-4 duration-200"
                            onMouseEnter={() => setHovered(idx)}
                            backgroundColor={item.color ?? theme.secondary}
                        >
                            {
                                item.icon && (
                                    <Icon 
                                        name={item.icon.value} 
                                        color={item.icon.color}
                                        size={20} 
                                    />
                                )
                            }
                            {
                                item.text && <Text 
                                    className="text-lg font-bold px-2"
                                >{item.text}</Text> 
                            }
                        </Button>
                        {
                            item.children && hovered === idx && (
                                <Dropdown 
                                    className="fixed top-40 ml-28"
                                    items={item.children}
                                />
                            )
                        }
                        </div>
                    )
                })
            }
        </div>
    );
}