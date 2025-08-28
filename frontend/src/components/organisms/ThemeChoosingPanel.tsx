import { useTranslation } from "react-i18next";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks"
import { Form } from "../molecules/Form";
import { useState } from "react";
import { setTheme } from "../../app/store/theme/theme.slice";
import { Button } from "../atoms/Button";
import { darkenColor, negativeColor, randomKey } from "../../utils/utils";
import { THEMES } from "../../config/themes";
import { Theme } from "../../app/types/Theme";
import { Text } from "../atoms/Text";

export const ThemePanelChoosing = (props: {
    onClose: () => void
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const dispatch = useAppDispatch();
    const [currentTheme, setCurrentTheme] = useState<Theme>({
        ...theme
    });
    const { t } = useTranslation();

    const onConfirm = () => {
        dispatch(setTheme(currentTheme));
        props.onClose();
    }

    return (
        <Form 
            label={t('choose_theme')}
            onCancel={props.onClose}
            onSubmit={onConfirm}
        >
            <div className="flex px-4 flex-col h-[70vh] overflow-auto">
            {
                THEMES.map(t => {
                    const colors = [t.primary, t.secondary, t.tertiary, t.accents.positive, t.accents.negative, t.accents.neutral, t.accents.contrast]
                    return (
                        <Button 
                            key={randomKey()}
                            className="flex w-full mb-2 py-2 px-4 rounded justify-center duration-200"
                            backgroundColor={currentTheme.name === t.name 
                                ? theme.accents.contrast : darkenColor(theme.secondary, 0.7)}
                            onClick={() => setCurrentTheme(t)}
                        >
                            <div className="flex flex-col">
                                <Text className="text-xl font-bold" color={currentTheme.name === t.name ? negativeColor(theme.text) : theme.text}>{t.name}</Text>
                                <div className="flex flex-row">
                                    {
                                        colors.map(color => {
                                            return (
                                                <div 
                                                    key={randomKey()}
                                                    className="mx-1 border w-8 h-8 rounded-2xl" 
                                                    style={{ backgroundColor: color}}
                                                />
                                            );
                                        })
                                    }
                                </div>
                            </div>
                        </Button>
                    );
                }) 
            }
            </div>
        </Form>
    );
}