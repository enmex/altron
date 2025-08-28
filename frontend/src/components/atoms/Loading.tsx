import UseAnimations from "react-useanimations";
import loading from "react-useanimations/lib/loading";
import { useAppSelector } from "../../app/store/hooks";

export const Loading = (props:{
    size: number
    hidden: boolean
    className?: string
}) => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    return (
        <div className={props.className ?? "flex w-full h-full justify-center items-center"}>
            <UseAnimations 
                animation={loading} 
                size={props.size} 
                strokeColor={props.hidden ? 'transparent' : theme.text}
            />
        </div>
    );
}