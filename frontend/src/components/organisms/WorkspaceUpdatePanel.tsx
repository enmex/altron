import { useState } from "react";
import { Workspace } from "../../app/types/Workspace";
import { useUpdateWorkspaceMutation } from "../../app/store/workspace/workspace.api";
import { useAppDispatch } from "../../app/store/hooks";
import { Input } from "../atoms/Input";
import { setWorkspace } from "../../app/store/workspace/workspace.slice";
import { useTranslation } from "react-i18next";
import { Form } from "../molecules/Form";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const WorkspaceUpdatePanel = (props: {
    workspace: Workspace,
    onClose: () => void
}) => {
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const [state, setState] = useState(props.workspace.name);
    const [updateWorkspace] = useUpdateWorkspaceMutation();

    const onSubmit = () => {
        updateWorkspace({
            id: props.workspace.id,
            name: state
        }).unwrap().then(() => {
            notifyInfo(t('update_workspace_success', {workspaceName: props.workspace.name}));
            dispatch(setWorkspace({
                ...props.workspace,
                name: state
            }));  
            props.onClose();
        })
        .catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    return (
        <Form
            label={t('update_workspace')}
            onCancel={props.onClose}
            onSubmit={onSubmit}
        >
            <Input 
                label="name"
                onChange={(e) => setState(e.target.value)}
                value={state}
            />
        </Form>
    );
}