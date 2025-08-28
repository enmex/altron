import { useState } from "react";
import { Button } from "../atoms/Button";
import { Icon } from "../atoms/Icon";
import { Overlay } from "../atoms/Overlay";
import { ThemePanelChoosing } from "../organisms/ThemeChoosingPanel";

export const ThemeSwitcher = () => {
    const [chooseThemePanelActive, setChooseThemePanelActive] = useState(false);

    return (
        <>
        <Button
            onClick={() => setChooseThemePanelActive(true)}
        >
            <Icon tip="switch theme" type="neutral" name="colorSwatch" size={30}/>
        </Button> 
        {
            chooseThemePanelActive && (
                <Overlay>
                    <ThemePanelChoosing 
                        onClose={() => setChooseThemePanelActive(false)}
                    />
                </Overlay>
            )
        }
        </>
    );
}