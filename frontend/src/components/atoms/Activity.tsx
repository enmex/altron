import UseAnimations from "react-useanimations";
import activity from "react-useanimations/lib/activity";

export const Activity = (props:{
    size: number;
    color: string;
    speed: number;
}) => {
    return (
        <UseAnimations 
            animation={activity} 
            size={props.size} 
            strokeColor={props.color}
            speed={props.speed}
        />
    )
}