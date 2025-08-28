import { DragEvent, useRef, useState } from "react";
import { useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text";
import { useTranslation } from "react-i18next";

export const DragDrop = (props:{
    onDrop: (file: File) => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const { t } = useTranslation();
    const [dragActive, setDragActive] = useState(false);
    const [filename, setFilename] = useState("");
    const fileUploader = useRef<HTMLInputElement>(null);

    const onDrop = (e: DragEvent) => {
        if (!e.dataTransfer) {
            return;
        }
        if (e.dataTransfer.files && e.dataTransfer.files[0]) {
            setDragActive(false);
            setFilename(e.dataTransfer.files[0].name);
            props.onDrop(e.dataTransfer.files[0]);
        }
    }

    const onClick = () => {
        if (!fileUploader.current) {
            return;
        }
        fileUploader.current.click()
    }

    return (
        <div 
            draggable
            className="flex flex-col border-2 border-dashed px-12 py-28 m-2 rounded-lg duration-200 cursor-pointer"
            style={{
                borderColor: dragActive ? theme.accents.contrast : theme.text,
                backgroundColor: theme.secondary
            }}
            onDragEnter={() => setDragActive(true)}
            onDragLeave={() => setDragActive(false)}
            onMouseEnter={() => setDragActive(true)}
            onMouseLeave={() => setDragActive(false)}
            onDrop={onDrop}
            onClick={onClick}
        >
            <Text
                className="text-lg cursor-default"
            >{
                filename ? filename : t('drag_drop_title')
            }</Text>
            <input className="hidden" type="file" multiple ref={fileUploader} onChange={(e) => {
                if (e.target.files) {
                    setFilename(e.target.files[0].name);
                    props.onDrop(e.target.files[0]);
                }
            }}/>
        </div>
    );
}