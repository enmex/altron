import { useAppSelector } from "../../app/store/hooks"
import { Panel } from "../../components/atoms/Panel"
import { Text } from "../../components/atoms/Text";
import { Button } from "../../components/atoms/Button"
import { useAppNavigation } from "../../hooks/navigate";
import { INDEX_PATH } from "../../config/constants";

export const NotFound = () => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const navigate = useAppNavigation();

    return (
        <Panel
            className="flex flex-col w-full h-full justify-center items-center"
        >
            <Text
                className="font-extrabold text-9xl mb-16"
                color={theme.text}
            >404</Text>
            <Text
                className="font-bold text-6xl mb-6"
                color={theme.accents.contrast}
            >Page not found</Text>
            <Text
                className="font-bold text-xl mb-6"
            >The page you are looking for does not exist or an error occurred</Text>
            <Button 
                className="flex justify-center border-2 mb-2 p-2 rounded-md duration-200"
                borderColor={theme.text}
                onClick={() => navigate(INDEX_PATH)}
            >
                <Text className="text-lg font-bold p-2">Go Home</Text>
            </Button>
        </Panel>
    )
}