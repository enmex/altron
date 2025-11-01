import { getDashboard } from "../../app/store/dashboard/dashboard.api";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { Service } from "../../app/types/Service";
import { setService, unsetService } from "../../app/store/service/service.slice";
import { useEffect, useState } from "react";
import { Button } from "../atoms/Button";
import { Select } from "../molecules/Select";
import { useTranslation } from "react-i18next";
import { useAppNavigation } from "../../hooks/navigate";
import { Icon } from "../atoms/Icon";
import { ServiceCreationPanel } from "./ServiceCreationPanel";
import { Overlay } from "../atoms/Overlay";
import { Text } from "../atoms/Text";
import { unsetWorkspace } from "../../app/store/workspace/workspace.slice";
import { Activity } from "../atoms/Activity";
import { LanguageSelector } from "../molecules/LanguageSwitcher";
import { ThemeSwitcher } from "../molecules/ThemeSwitcher";
import { SIGN_IN_PATH } from "../../config/constants";
import { ImportPcapPanel } from "../molecules/ImportPcapPanel";
import { PcapWorkspace } from "../../app/types/PcapWorkspace";
import { DropdownButton } from "../molecules/DropdownButton";
import { unsetPcapWorkspace } from "../../app/store/pcap/pcap.slice";
import { putError } from "../../app/store/error/error.slice";
import { ImportJsonPanel } from "../molecules/ImportJsonPanel";
import { getService } from "../../app/store/service/service.api";

export const Header = () => {
    const theme = useAppSelector(state => state.rootReducer.theme);
    const service = useAppSelector(state => state.rootReducer.service);
    const workspace = useAppSelector(state => state.rootReducer.workspace);
    const pcap = useAppSelector(state => state.rootReducer.pcap);
    const currentRpm = useAppSelector(state => state.rootReducer.metrics);
    const dispatch = useAppDispatch();
    const { t } = useTranslation();
    const navigate = useAppNavigation();
    const [getServiceTrigger] = getService.useLazyQuery();
    const [serviceCreationPanelActive, setServiceCreationPanelActive] = useState(false);
    const [pcapWorkspaces, setPcapWorkspaces] = useState<PcapWorkspace[]>([]);
    const [importPcapPanelActive, setImportPcapPanelActive] = useState(false);
    const [importJsonPanelActive, setImportJsonPanelActive] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [services, setServices] = useState<Service[]>([]);
    const [getDashboardTrigger] = getDashboard.useLazyQuery();
    const [agentStatus, setAgentStatus] = useState("down");

    const onService = (service: Service) => {
        dispatch(unsetWorkspace());
        dispatch(unsetPcapWorkspace());
        refetchDashboard();
        getServiceTrigger(service.port).unwrap().then((res) => {
            dispatch(setService(res));
            navigate("/services/" + service.port);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        });
    }

    const refetchDashboard = () => {
        setIsLoading(true);
        getDashboardTrigger().unwrap().then((dashboard) => {
            setServices(dashboard.services);
            setPcapWorkspaces(dashboard.pcapWorkspaces);
        }).catch((err) => {
            dispatch(putError(err.data.message));
        }).finally(() => {
            setIsLoading(false);
        });
    }

    const checkAgentStatus = () => {
        fetch('/connect/agent/status').
            then((res) => res.json()).
            then((data) => setAgentStatus(data.status)).
            catch((err) => {
                setAgentStatus("down");
                dispatch(putError(err.data.message));
            });
    }

    useEffect(() => {
        checkAgentStatus();
        const interval = setInterval(() => {
            checkAgentStatus();
        }, 5000);

        refetchDashboard();
        return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [navigate]);

    const onHome = () => {
        dispatch(unsetService());
        dispatch(unsetWorkspace());
        dispatch(unsetPcapWorkspace());
        navigate("/home");
    }

    const getRpmColor = (rpm: number): string => {
        const rps = rpm / 60;
        return rps >= 0 && rps < 100 
            ? theme.accents.positive
            : rps >= 100 && rps < 500
            ? theme.accents.contrast
            : theme.accents.negative; 
    }

    const onLogout = () => {
        localStorage.removeItem("auth");
        navigate(SIGN_IN_PATH);
    }

    const onPcapWorkspace = (pcapWorkspace: PcapWorkspace) => {
        dispatch(unsetWorkspace());
        dispatch(unsetService());
        navigate(`/pcaps/${pcapWorkspace.id}`);
    }

    return (
        <div 
            className="flex flex-row justify-between items-center w-screen mx-auto border-b-4 p-2"
            style={{
                backgroundColor: theme.secondary,
                borderColor: theme.tertiary
            }}
        >
            <div className="flex flex-row items-center">
                <Button 
                    onClick={onHome}
                ><Icon tip="Altron" type="contrast" name="burningEye" size={40}/></Button>
                <Button 
                    onClick={() => setServiceCreationPanelActive(true)}
                >
                    <Icon 
                        type="neutral"
                        size={27}
                        name="mdAdd"
                    />
                </Button>
                <Select 
                    items={services && !isLoading ? services.map(service => {
                        return {
                            text: service.name,
                            onItem: () => onService(service)
                        }
                    }) : []}
                    className="relative"
                    icon="arrowDropDown"
                    placeholder={t('services') ?? ""}
                />
                {
                    workspace.servicePort > 0 || service.port > 0 && (
                        <Text
                            className="ml-2 font-bold text-xl border-2 px-2 cursor-default"
                        >{workspace.servicePort > 0 ? workspace.servicePort : service.port}</Text>
                    )
                }
                {
                    workspace.name.length > 0 && (
                        <>
                        <Icon name="arrowRight" size={30}/>
                        <Text
                            className="ml-2 font-bold text-xl border-2 px-2 cursor-default"
                        >{workspace.name}</Text>
                        </>
                    )
                }
                {
                    service.port > 0 && workspace.name.length === 0 && (
                        <div className="flex flex-row">
                            <Text 
                                className="ml-2 font-bold text-xl px-2 cursor-default"
                            >{`${currentRpm} RPM`}</Text>
                            <Text 
                                className="ml-2 font-bold text-xl px-2 cursor-default"
                            >{`${(currentRpm / 60).toFixed(1)} RPS`}</Text>
                            <Activity 
                                size={30}
                                color={getRpmColor(currentRpm)}
                                speed={Math.max(5 * currentRpm/(60 * 1000), 0.5)}
                            />
                        </div>
                    )
                }
                {
                    pcap.id.length > 0 && (
                        <Text
                            className="ml-2 font-bold text-xl border-2 px-2 cursor-default"
                        >{pcap.fileName.length > 12 ? pcap.fileName.slice(0, 12) + '...' : pcap.fileName}</Text>
                    )
                }
            </div>
            <div className="flex flex-row items-center">
                <DropdownButton
                    onClick={() => setImportPcapPanelActive(true)}
                    icon={
                        <Icon tip="import pcap dump" type='contrast' size={35} name='import' />
                    }
                    dropdownItems={pcapWorkspaces.map(pcapWorkspace => {
                        return {
                            text: pcapWorkspace.fileName,
                            onItem: () => onPcapWorkspace(pcapWorkspace)
                        }
                    })}
                />
                <Button
                    onClick={() => setImportJsonPanelActive(true)}
                >
                    <Icon tip="import session in json" type='contrast' size={25} name='json'/>
                </Button>
                <Icon 
                    tip={agentStatus === "up" ? "agent online" : "agent offline"} 
                    type='contrast' size={30} 
                    name={agentStatus === "up" ? 'online' : 'offline'} 
                />
                <ThemeSwitcher />
                <LanguageSelector />
                <Button
                    onClick={onLogout}
                >
                    <Icon name="logout" type="negative" size={35} />
                </Button>
            </div>
            {
                serviceCreationPanelActive && (
                    <Overlay>
                        <ServiceCreationPanel onClose={() => {setServiceCreationPanelActive(false); refetchDashboard();}} />
                    </Overlay>
                )
            }
            {
                importPcapPanelActive && (
                    <Overlay>
                        <ImportPcapPanel onClose={() => {setImportPcapPanelActive(false); refetchDashboard();}}/>
                    </Overlay>
                )
            }
            {
                importJsonPanelActive && (
                    <Overlay>
                        <ImportJsonPanel onClose={() => setImportJsonPanelActive(false)}/>
                    </Overlay>
                )
            }
        </div>
    );
}