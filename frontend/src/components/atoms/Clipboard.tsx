import { useTranslation } from "react-i18next";
import { Button } from "./Button";
import { Icon } from "./Icon";
import { notifyInfo } from "../../app/notifications/notifier";

export const Clipboard = (props: {
    text: string;
    size?: number;
    onClick?: () => void;
}) => {
    const { t } = useTranslation();

    const onClick = () => {
        const tempTextarea = document.createElement('textarea');
        tempTextarea.value = props.text;
        tempTextarea.style.width = 'auto';
        tempTextarea.style.height = 'auto';
        tempTextarea.style.resize = 'auto';
        tempTextarea.style.overflow = 'hidden';      
        document.body.appendChild(tempTextarea);
        tempTextarea.select();
        document.execCommand('copy');
        document.body.removeChild(tempTextarea);
        notifyInfo(t('copied_to_clipboard'));
        if (props.onClick) {
            props.onClick();
        }
    }
    
    return (
        <Button
            onClick={onClick}
        >
            <Icon 
                tip="copy to clipboard"
                type="neutral"
                name="contentCopy"
                size={props.size ?? 30}
            />
        </Button>
    )
}