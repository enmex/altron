import { useState } from "react"
import { Form } from "./Form"
import { useAppDispatch } from "../../app/store/hooks"
import { useCreateSessionMutation } from "../../app/store/session/session.api"
import { Session } from "../../app/types/Service"
import { useAppNavigation } from "../../hooks/navigate"
import { putError } from "../../app/store/error/error.slice"
import { Loading } from "../atoms/Loading"

export const ImportJsonPanel = (props:{
    onClose: () => void
}) => {
    const dispatch = useAppDispatch();
    const navigate = useAppNavigation();
    const [createSession] = useCreateSessionMutation();
    const [isLoading, setIsLoading] = useState(false);
    const [json, setJson] = useState("");

    const onSubmit = () => {
        setIsLoading(true);
        const session = JSON.parse(json) as Session
        session.id = "a638a3d5-25f3-488a-bd9c-4452fae59922"
        createSession(session).
            unwrap().
            then(() => {
                navigate(`/sessions/${session.id}`)
                props.onClose();
            }).catch((err) => {
                dispatch(putError(err.data.message));
            }).finally(() => {
                setIsLoading(false);
            })
    }

    return (
        <Form 
            label={"Import JSON session"}
            onSubmit={onSubmit}
            onCancel={props.onClose}
        >
            <textarea onChange={(e) => setJson(e.target.value)}/>
            <Loading size={30} hidden={!isLoading}/>
        </Form>
    )
}