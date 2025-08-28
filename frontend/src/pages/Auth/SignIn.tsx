import { useState } from "react";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { useSignInMutation } from "../../app/store/auth/auth.api";
import { INDEX_PATH } from "../../config/constants";
import { Input } from "../../components/atoms/Input";
import { useTranslation } from "react-i18next";
import { useAppNavigation } from "../../hooks/navigate";
import { Form } from "../../components/molecules/Form";
import { AnimatedLayout } from "../../layouts/AnimatedLayout";
import { Overlay } from "../../components/atoms/Overlay";
import { putError } from "../../app/store/error/error.slice";
import { setAuth } from "../../app/store/auth/auth.slice";

export const SignIn = () => {
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const theme = useAppSelector(state => state.rootReducer.theme);
    const [signIn] = useSignInMutation();
    const [state, setState] = useState({
        password: "",
    });
    const navigate = useAppNavigation();
    const [greetingAudio] = useState(theme.greetingsSoundPath && new Audio(theme.greetingsSoundPath));

    const onSubmit = () => { 
        signIn({
            ...state,
        }).unwrap().then((auth) => {
            if (!auth) {
                return;
            }
            dispatch(setAuth(auth));
            navigate(INDEX_PATH);
            if (greetingAudio) {
                greetingAudio.volume = 0.3;
                greetingAudio.play().then(() => navigate(INDEX_PATH));
            } else {
                navigate(INDEX_PATH);
            }
        }).catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    return (
        <AnimatedLayout>
            <Overlay>
                <div className="flex w-full h-full justify-center items-center">
                    <Form
                        label={t('sign_in')}
                        onSubmit={onSubmit}
                    >
                        <Input 
                            type="password"
                            label={t('password')}
                            onChange={(e) => setState({
                                password: e.target.value
                            })}    
                            value={state.password}
                        />
                    </Form>
                </div>
            </Overlay>
        </AnimatedLayout>
        
    )
}