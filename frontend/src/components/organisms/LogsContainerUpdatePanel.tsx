import { useTranslation } from "react-i18next";
import { Form } from "../molecules/Form";
import { useScanServicesQuery, useUpdateServiceMutation } from "../../app/store/service/service.api";
import { Loading } from "../atoms/Loading";
import { randomKey } from "../../utils/utils";
import { Button } from "../atoms/Button";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text";
import { useState } from "react";
import { setService } from "../../app/store/service/service.slice";
import { notifyInfo, notifyWarning } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";
import { Icon } from "../atoms/Icon";

export const LogsContainerUpdatePanel = (props:{
    onClose: () => void
}) => {
    const { t } = useTranslation();
    const dispatch = useAppDispatch();
    const service = useAppSelector(state => state.rootReducer.service);
    const {data: dockerizedServices} = useScanServicesQuery({
        scope: "containers"
    });
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [updateService] = useUpdateServiceMutation();
    const [currentService, setCurrentService] = useState<{
        container?: {
            id: string;
            name: string;
            image: string;
        } | undefined;
        port: number;
        name: string;
        isPublic: boolean;
    }>();

    const onSubmit = () => {
        if (!currentService || !currentService.container) {
            notifyWarning('you need to choose container');
            return;
        }
        updateService({
            ...service,
            containerID: currentService.container?.id
        })
        .unwrap()
        .then(() => {
            dispatch(setService({
                ...service,
                containerID: currentService.container?.id
            }));
            notifyInfo(`service '${service.name} updated'`);
            props.onClose();
        })
        .catch((err) => {
            dispatch(putError(err.data.message));
        })
    }

    return (
        <Form
            label={t('update_logs_container')}
            onCancel={props.onClose}
            onSubmit={onSubmit}
        >
            <div
                className="overflow-auto flex flex-col h-[80vh] p-4"
            >
                {
                    dockerizedServices ? dockerizedServices.services.map(service => {
                        if (!service.container) return;

                        return (
                            <Button 
                                key={randomKey()}
                                className="flex flex-col w-full p-2 border-b-2 duration-200 my-2"
                                backgroundColor={
                                    service.port === currentService?.port && 
                                    service.container.id === currentService.container?.id
                                        ? theme.tertiary 
                                        : theme.secondary}
                                borderColor={theme.text}
                                onClick={() => setCurrentService(service)} 
                            >
                                <div className="flex flex-col text-start">
                                    <div className="flex flex-row items-center w-3/4 justify-between">
                                        <Text 
                                            className="text-lg mr-2"
                                            color={theme.text}
                                        >{t('port')}</Text>
                                        <Text 
                                            className="text-lg font-extrabold mr-2"
                                            color={theme.accents.contrast}
                                        >{service.port}</Text>
                                        <Icon 
                                            tip={service.isPublic ? "public port" : "private port"}
                                            type="contrast"
                                            name={service.isPublic ? 'unlock' : 'lock'}
                                            size={15}
                                        />
                                    </div>
                                    <div className="flex flex-row items-center">
                                        <Text 
                                            className="font-bold mr-2"
                                            color={theme.text}
                                        >{service.name}</Text>
                                        <Icon 
                                            tip="dockerized service"
                                            type="neutral"
                                            name="docker"
                                            size={20}
                                        />
                                    </div>
                                    <Text 
                                        className="font-bold"
                                        color={theme.accents.positive}
                                    >{`${service.container.id}`}</Text>
                                    <Text 
                                        className="font-bold"
                                        color={theme.accents.positive}
                                    >{`image: ${service.container.image}`}</Text>
                                </div>
                            </Button>
                        );
                    }) : (
                        <Loading hidden={false} size={30} />
                    )
                }
            </div>
        </Form>
    );
}