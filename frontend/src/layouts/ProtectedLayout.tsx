import { ReactNode, useEffect } from "react";
import { Header } from "../components/organisms/Header";
import { SIGN_IN_PATH } from "../config/constants";
import { useAppNavigation } from "../hooks/navigate";
import { useAppDispatch, useAppSelector } from "../app/store/hooks";
import { setAuth } from "../app/store/auth/auth.slice";

export const ProtectedLayout = (props: {
    children: ReactNode
  }) => {
    const dispatch = useAppDispatch();
    const auth = useAppSelector(state => state.rootReducer.auth);
    const navigate = useAppNavigation();
    
    useEffect(() => {
      const t = localStorage.getItem("auth");
      if (!t) {
        navigate(SIGN_IN_PATH);
      } 
      dispatch(setAuth({
        token: t
      }));
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [auth]);
  
    return (
      <>
      <Header />
      {props.children}
      </>
    )
  }