import { useEffect } from "react";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { setLanguage } from "../../app/store/language/language.slice";
import i18n from "../../i18";
import { Button } from "../atoms/Button";
import { putError } from "../../app/store/error/error.slice";

export const LanguageSelector = () => {
    const dispatch = useAppDispatch();
    const lng = useAppSelector(state => state.rootReducer.language);

    const changeLanguage = (lng: string) => {
      i18n.changeLanguage(lng).then(() => {
        dispatch(setLanguage(lng));
      }).catch((err) => {
        dispatch(putError(err.data.message));
      });
    };
    const currentLang = lng === "en" ? "gb" : lng;

    useEffect(() => {
        changeLanguage(lng);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    
    return (
      <Button
        onClick={() => changeLanguage(lng === "en" ? "ru" : "en")}
      >
        <span className={`text-2xl rounded-full fi fi-${currentLang} fis`}/>
      </Button>
    );
}