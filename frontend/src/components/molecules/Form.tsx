import { t } from "i18next";
import { FormEvent, ReactNode, useCallback, useEffect } from "react";
import { Button } from "../atoms/Button";
import { Panel } from "../atoms/Panel";
import { useAppSelector } from "../../app/store/hooks";
import { Text } from "../atoms/Text";

export const Form = (props: {
    label: string;
    onCancel?: () => void;
    onSubmit: () => void;
    children?: ReactNode;
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const escFunction = useCallback((event: KeyboardEvent) => {
        if (props.onCancel && event.key === "Escape") {
            props.onCancel();
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    
    useEffect(() => {
        document.addEventListener("keydown", escFunction, false);

        return () => {
            document.removeEventListener("keydown", escFunction, false);
        };
    }, [escFunction]);

    const onSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        props.onSubmit();
    }

    return (
        <Panel withBorder>
            <form className="flex flex-col w-full max-w-xs mx-2 justify-center text-center" onSubmit={onSubmit}>
                <Text className="font-bold text-xl cursor-default">{props.label}</Text>
                {props.children}
                <div className="flex flex-row justify-center mt-2">
                    {
                        props.onCancel && (
                            <Button 
                                className="font-semibold px-4 p-2 mr-4 rounded-lg duration-200 border-2"
                                borderColor={theme.accents.negative}
                                onClick={props.onCancel}
                            ><Text className="font-bold text-lg">{t('cancel')}</Text></Button>
                        )
                    }
                    <Button 
                        type="submit"
                        className="font-semibold px-4 p-2 rounded-lg duration-200 border-2"
                        borderColor={theme.accents.positive}
                    ><Text 
                        className="font-bold text-lg"
                    >OK</Text></Button>
                </div>
            </form>
        </Panel>
    );
}