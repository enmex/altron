import { getAllFilters, useDeleteFilterMutation } from "../../app/store/filter/filter.api";
import { Filter } from "../../app/types/Filter";
import { Button } from "../atoms/Button";
import { Overlay } from "../atoms/Overlay";
import { FilterCreationPanel } from "../organisms/FilterCreationPanel";
import { useEffect, useState } from "react";
import { DropdownButton } from "./DropdownButton";
import { useTranslation } from "react-i18next";
import { useAppDispatch, useAppSelector } from "../../app/store/hooks";
import { FilterUpdatePanel } from "../organisms/FilterUpdatePanel";
import { randomKey } from "../../utils/utils";
import { Loading } from "../atoms/Loading";
import { Text } from "../atoms/Text";
import { Icon } from "../atoms/Icon";
import { notifyInfo } from "../../app/notifications/notifier";
import { putError } from "../../app/store/error/error.slice";

export const FilterPanel = (props: {
    onClick: (filter: Filter) => void
    onUpdate: () => void
    selectedFilter?: Filter
}) => {
    const dispatch = useAppDispatch();
    const theme = useAppSelector(state => state.rootReducer.theme);
    const servicePort = useAppSelector(state => state.rootReducer.service.port);
    const [filters, setNewFilters] = useState<Filter[]>([]);
    const { t } = useTranslation();
    const [getAllFiltersTrigger] = getAllFilters.useLazyQuery();
    const [deleteFilter] = useDeleteFilterMutation();
    const [createFilterPanelActive, setCreateFilterPanelActive] = useState(false);
    const [filterToUpdate, setFilterToUpdate] = useState<Filter>({
        id: "",
        name: "",
        inRequest: false,
        inResponse: false,
        color: "",
        isBlocking: false,
    });
    const [isLoading, setIsLoading] = useState(true);

    const onDeleteFilter = (filter: Filter) => {
        deleteFilter(filter.id)
            .unwrap()
            .then(() => {
                const newFilters = filters.filter(f => f.id !== filter.id);
                setNewFilters(prev => newFilters);
                onGetAllFilters();
                props.onUpdate();
                notifyInfo(t('delete_filter_success'));
            })
            .catch((err) => {
                dispatch(putError(err.data.message));
            });
    }

    const onGetAllFilters = () => {
        setIsLoading(true);
        getAllFiltersTrigger(servicePort)
            .unwrap()
            .then((data) => {
                setNewFilters(prev => data.filters);
                setIsLoading(false);
            })
            .catch((err) => {
                dispatch(putError(err.data.message));
            });
    }


    const onCreateFilterPanelClose = () => {
        setCreateFilterPanelActive(false);
        onGetAllFilters();
        props.onUpdate();
    }

    const onUpdateFilter = () => {
        setFilterToUpdate({
            id: "",
            name: "",
            inRequest: false,
            inResponse: false,
            color: "",
            isBlocking: false,
        }); 
        onGetAllFilters();
        props.onUpdate();
    }

    useEffect(() => {
        onGetAllFilters();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [servicePort]);

    return (
        <div className="flex flex-row px-4">
            <div className="overflow-x-auto w-full flex items-center">
                <div className="flex flex-row">
                    {
                        isLoading ? (
                            <div className="flex w-full h-full justify-center items-center">
                                <Loading hidden={false} size={25} />
                            </div>
                        ) : filters.map(filter => {
                            return (
                                <DropdownButton 
                                    key={randomKey()}
                                    text={
                                        <div 
                                            className="duration-200 rounded px-1 border-2"
                                            style={{
                                                borderColor: props.selectedFilter && props.selectedFilter.name === filter.name ? theme.accents.contrast : "transparent",
                                                backgroundColor: filter.color
                                            }}
                                        >
                                            <Text
                                                className="px-1 rounded duration-200 font-bold text-lg"
                                            >{filter.name.length > 5 ? filter.name.slice(0, 5) + "..." : filter.name}</Text>
                                        </div>
                                        
                                    }
                                    icon={ filter.isBlocking && <Icon name="filterAlt" size={25} /> }
                                    onClick={() => {
                                        if (!filter.isBlocking) {
                                            props.onClick(filter);
                                        }
                                    }}
                                    dropdownItems={[
                                        {
                                            text: t('update') + "",
                                            onItem: () => setFilterToUpdate(filter)
                                        },
                                        {
                                            text: t('delete') + "",
                                            onItem: () => onDeleteFilter(filter)
                                        }
                                    ]}
                                />
                            )
                        })
                    }
                </div>
            </div>
            <Button 
                onClick={() => setCreateFilterPanelActive(true)}
            ><Icon tip="new filter" name="mdAdd" size={27}/></Button>
            <Text className="cursor-default text-lg font-bold px-2 flex items-center">{t('filters')}</Text>
            {
                createFilterPanelActive && (
                    <Overlay>
                        <FilterCreationPanel 
                            onClose={onCreateFilterPanelClose}
                        />
                    </Overlay>
                )
            }
            {
                filterToUpdate.id && (
                    <Overlay>
                        <FilterUpdatePanel 
                            filter={filterToUpdate}
                            onClose={onUpdateFilter}
                        />
                    </Overlay>
                )
            }
        </div>
    );
}