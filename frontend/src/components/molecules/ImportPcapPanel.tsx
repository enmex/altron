import { useTranslation } from "react-i18next";
import { DragDrop } from "../atoms/DragDrop"
import { Form } from "./Form";
import { useState } from "react";
import { useAppDispatch } from "../../app/store/hooks";
import { useUploadPcapMutation } from "../../app/store/pcap/pcap.api";
import { useAppNavigation } from "../../hooks/navigate";
import { Loading } from "../atoms/Loading";
import { setPcapWorkspace } from "../../app/store/pcap/pcap.slice";
import { notifyWarning } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const ImportPcapPanel = (props:{
    onClose: () => void;
}) => {
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const [currentFile, setCurrentFile] = useState<File>();
    const [isLoading, setIsLoading] = useState(false);
    const [uploadPcap] = useUploadPcapMutation();
    const navigate = useAppNavigation();

    const onUpload = () => {
        if (!currentFile) {
            notifyWarning(t('no_pcap_file'));
            return;
        }
        const formData = new FormData();
        formData.append('pcap', new Blob([currentFile], { type: currentFile.type }), currentFile.name);
        setIsLoading(true);
        uploadPcap(formData).unwrap()
            .then((res) => {
                dispatch(setPcapWorkspace(res))
                navigate(`/pcaps/${res.id}`);
                props.onClose();
            })
            .catch((err) => {
                dispatch(putError(err.data.message));
            })
            .finally(() => {
                setIsLoading(false);
            });
    }

    return (
        <Form 
            label={t('import_pcap')}
            onSubmit={onUpload}
            onCancel={props.onClose}
        >
            <DragDrop 
                onDrop={setCurrentFile}
            />
            <Loading hidden={!isLoading} size={40} />
        </Form>
    )
}