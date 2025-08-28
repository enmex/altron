import { CirclePicker } from "react-color";
import { useTranslation } from "react-i18next";
import { Text } from "../atoms/Text";

export const ColorPicker = (props: {
    onChange: (color: string) => void
}) => {
    const { t } = useTranslation();
    
    return (
        <div className="flex flex-col items-center">
            <Text className="font-bold text-lg">{t('color')}</Text>
            <CirclePicker className="my-2" color={"white"} onChange={(e) => props.onChange(e.hex)}/>
        </div>
    );
}