import { Toaster, resolveValue, toast as toastFunction } from "react-hot-toast";
import { useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text";
import { Button } from "../atoms/Button";
import { Icon } from "../atoms/Icon";
import { useTranslation } from "react-i18next";

export const NotificationToaster = () => {
    const { t } = useTranslation();
    const theme = useAppSelector(state => state.rootReducer.theme);

    return (
        <Toaster
            position="bottom-right"
            reverseOrder={true}
            gutter={8}
        >
            {(toast) => {
                let type = toast.type.toString();
                if (type === "blank") {
                    type = "warning";
                } else if (type === "success") {
                    type = "info";
                }

                return (
                    <div
                        className="w-1/5 rounded-lg font-bold text-lg"
                        style={{
                            color: theme.text,
                            backgroundColor: theme.secondary
                        }}
                    >
                    <div 
                        className="animate-progress top-0 left-0 h-2 w-full"
                        style={{
                            backgroundColor: theme.tertiary,
                        }}
                    ></div>
                    <div className="px-1 py-2">
                        <div className="flex flex-row justify-between w-full">
                            <Text
                                className="cursor-default flex flex-col ml-2 font-extrabold text-xl"
                                color={type === "error" 
                                    ? theme.accents.negative
                                    : type === "info"
                                        ? theme.accents.positive
                                        : theme.accents.contrast
                                    }
                            >{t(type)}</Text>
                            <Button
                                onClick={() => toastFunction.dismiss(toast.id)}
                            >
                                <Icon tip="close" name="close" type="negative" size={20}/>
                            </Button>
                        </div>
                        <Text className="cursor-default font-bold text-lg ml-2">{resolveValue(toast.message, toast)}</Text>
                    </div>
                    </div>
                )
            }}
        </Toaster>
    );
}