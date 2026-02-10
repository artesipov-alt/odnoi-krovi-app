import { Button } from '@mui/material';
import cn from 'classnames';
import Analizes from 'imgs/svg/analizes';
import Success from 'imgs/svg/success';
import { Analiz } from 'pages/adding/Donor/types';
import { FC, useEffect, useRef, useState } from 'react';
import { getDateFormat } from 'utils/utils';

import { updatePet } from 'api/apiServices/updatePet';
import { Analyses, Analyses as AnalysesType, AnalysesItem, Pet } from 'api/pets';
import { AnalysesEnum, AnalysesMapping, AnalysesNamesMapping, AnalysesTypes, PetType } from 'api/types';
import DatePicker from 'components/DatePicker';

import Header from '../Header';
import styles from './AnalysesStep.module.less';

type Props = {
    petId: string;
    petType?: string;
    isEditMode: boolean;
    onClose: () => void;
    analyses: AnalysesType;
    onErrorUpdate: () => void;
    onSuccessUpdate: () => void;
};

const getAnalizValue = (itemName: AnalysesTypes, analiz?: AnalysesItem[]) => {
    if (!analiz) {
        return {
            value: null,
            defaultValue: null,
        };
    }

    const itemIndex = analiz.findIndex((item) => item.analysisType === itemName);

    if (itemIndex === -1) {
        return {
            value: null,
            defaultValue: null,
        };
    }

    const result = new Date(analiz[itemIndex].analysisDate);

    return {
        value: result,
        defaultValue: result,
    };
};

const AnalysesStep: FC<Props> = ({ petId, onClose, isEditMode, analyses, petType, onErrorUpdate, onSuccessUpdate }) => {
    const [isSaveButtonActive, setIsSaveButtonActive] = useState<boolean>(false);

    const [leukemia, setLeukemia] = useState<Analiz>({
        type: AnalysesEnum.LEUKEMIA,
        name: AnalysesNamesMapping[AnalysesEnum.LEUKEMIA],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.LEUKEMIA]),
            },
            {
                name: 'ИФА',
                ...getAnalizValue(AnalysesTypes.ELISA, analyses?.[AnalysesEnum.LEUKEMIA]),
            },
            {
                name: 'ИХА',
                ...getAnalizValue(AnalysesTypes.ICA, analyses?.[AnalysesEnum.LEUKEMIA]),
            },
        ],
    });
    const [immunodeficiency, setImmunodeficiency] = useState<Analiz>({
        type: AnalysesEnum.IMMUNODEFICIENCY,
        name: AnalysesNamesMapping[AnalysesEnum.IMMUNODEFICIENCY],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.IMMUNODEFICIENCY]),
            },
            {
                name: 'ИФА',
                ...getAnalizValue(AnalysesTypes.ELISA, analyses?.[AnalysesEnum.IMMUNODEFICIENCY]),
            },
            {
                name: 'ИХА',
                ...getAnalizValue(AnalysesTypes.ICA, analyses?.[AnalysesEnum.IMMUNODEFICIENCY]),
            },
        ],
    });
    const [hemoplasmosis, setHemoplasmosis] = useState<Analiz>({
        type: AnalysesEnum.HEMOPLASMOSIS,
        name: AnalysesNamesMapping[AnalysesEnum.HEMOPLASMOSIS],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.HEMOPLASMOSIS]),
            },
        ],
    });
    const [bartonellosis, setBartonellosis] = useState<Analiz>({
        type: AnalysesEnum.BARTONELLOSIS,
        name: AnalysesNamesMapping[AnalysesEnum.BARTONELLOSIS],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.BARTONELLOSIS]),
            },
        ],
    });
    const [babesiosis, setBabesiosis] = useState<Analiz>({
        type: AnalysesEnum.BABESIOSIS,
        name: AnalysesNamesMapping[AnalysesEnum.BABESIOSIS],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.BABESIOSIS]),
            },
            {
                name: 'Микроскопия мазка',
                ...getAnalizValue(AnalysesTypes.MICROSCOPY, analyses?.[AnalysesEnum.BABESIOSIS]),
            },
        ],
    });
    const [dirofilaria, setDirofilaria] = useState<Analiz>({
        type: AnalysesEnum.DIROFILARIA,
        name: AnalysesNamesMapping[AnalysesEnum.DIROFILARIA],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.DIROFILARIA]),
            },
            {
                name: 'Микроскопия мазка',
                ...getAnalizValue(AnalysesTypes.MICROSCOPY, analyses?.[AnalysesEnum.DIROFILARIA]),
            },
        ],
    });
    const [ehrlichiosis, setEhrlichiosis] = useState<Analiz>({
        type: AnalysesEnum.EHRLICHIOSIS,
        name: AnalysesNamesMapping[AnalysesEnum.EHRLICHIOSIS],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.EHRLICHIOSIS]),
            },
            {
                name: 'Экспресс-тест',
                ...getAnalizValue(AnalysesTypes.EXPRESS, analyses?.[AnalysesEnum.EHRLICHIOSIS]),
            },
        ],
    });
    const [anaplasmosis, setAnaplasmosis] = useState<Analiz>({
        type: AnalysesEnum.ANAPLASMOSIS,
        name: AnalysesNamesMapping[AnalysesEnum.ANAPLASMOSIS],
        items: [
            {
                name: 'ПЦР',
                ...getAnalizValue(AnalysesTypes.PCR, analyses?.[AnalysesEnum.ANAPLASMOSIS]),
            },
            {
                name: 'Экспресс-тест',
                ...getAnalizValue(AnalysesTypes.EXPRESS, analyses?.[AnalysesEnum.ANAPLASMOSIS]),
            },
        ],
    });

    const changesRef = useRef<Record<string, boolean>>({});

    const checkChange = (defaultValue: Date | null, newValue: Date | null, id: string) => {
        const wasSet = defaultValue instanceof Date;
        const nowSet = newValue instanceof Date;

        if (wasSet && !nowSet) {
            changesRef.current[id] = true;
        } else if (!wasSet && nowSet) {
            changesRef.current[id] = true;
        } else if (wasSet && nowSet && newValue?.getTime() !== defaultValue?.getTime()) {
            changesRef.current[id] = true;
        } else {
            delete changesRef.current[id]; // если вернули как было — удаляем из изменённых
        }
    };

    const onDateChangeHandler = (type: string, analizName: string) => (value: Date | null) => {
        const updateState = (setter: React.Dispatch<React.SetStateAction<Analiz>>) => {
            setter((prevState) => {
                const updatedItems = prevState.items.map((item) => {
                    if (item.name === analizName) {
                        const id = `${type}_${analizName}`;
                        checkChange(item.defaultValue || null, value, id);

                        return { ...item, value };
                    }

                    return item;
                });

                return { ...prevState, items: updatedItems };
            });
        };

        switch (type) {
            case AnalysesEnum.LEUKEMIA:
                updateState(setLeukemia);
                break;
            case AnalysesEnum.IMMUNODEFICIENCY:
                updateState(setImmunodeficiency);
                break;
            case AnalysesEnum.HEMOPLASMOSIS:
                updateState(setHemoplasmosis);
                break;
            case AnalysesEnum.BARTONELLOSIS:
                updateState(setBartonellosis);
                break;
            case AnalysesEnum.BABESIOSIS:
                updateState(setBabesiosis);
                break;
            case AnalysesEnum.DIROFILARIA:
                updateState(setDirofilaria);
                break;
            case AnalysesEnum.EHRLICHIOSIS:
                updateState(setEhrlichiosis);
                break;
            case AnalysesEnum.ANAPLASMOSIS:
                updateState(setAnaplasmosis);
                break;
            default:
                break;
        }
    };

    const getAnalizesToRequest = () => {
        const result = (
            petType === PetType.DOG
                ? [babesiosis, dirofilaria, hemoplasmosis, bartonellosis, ehrlichiosis, anaplasmosis]
                : [leukemia, immunodeficiency, hemoplasmosis, bartonellosis]
        ).reduce((acc, analiz) => {
            acc[analiz.type] = analiz.items.reduce(
                (res, item) => {
                    if (item.value) {
                        res.push({
                            analysisDate: item.value,
                            analysisName: analiz.type,
                            analysisType: AnalysesMapping[item.name],
                        });
                    }

                    return res;
                },
                [] as unknown as AnalysesItem[],
            );

            return acc;
        }, {} as Analyses);

        const filteredResult = Object.keys(result).reduce((res, key) => {
            if (result[key].length) {
                res[key] = result[key];
            }

            return res;
        }, {} as Analyses);

        if (Object.keys(filteredResult).length) {
            return filteredResult;
        }

        return {};
    };

    const onSaveButtonClickHandler = async () => {
        const newData: Partial<Pet> = {
            id: petId,
            analyses: getAnalizesToRequest(),
        };

        const { success } = await updatePet(newData as Pet);

        if (success) {
            onSuccessUpdate();
        } else {
            onErrorUpdate();
        }

        onClose();
    };

    // Отслеживаем изменения в состояниях анализов
    useEffect(() => {
        setIsSaveButtonActive(Object.keys(changesRef.current).length > 0);
    }, [leukemia, immunodeficiency, hemoplasmosis, bartonellosis, babesiosis, dirofilaria, ehrlichiosis, anaplasmosis]);

    return (
        <div className={styles.wrapper}>
            <Header title='Здоровье' onClose={onClose} isEditMode={isEditMode} icon={<Analizes />} />
            <>
                {(petType === PetType.DOG
                    ? [babesiosis, dirofilaria, hemoplasmosis, bartonellosis, ehrlichiosis, anaplasmosis]
                    : [leukemia, immunodeficiency, hemoplasmosis, bartonellosis]
                ).map((analiz) => (
                    <div key={analiz.name} className={styles.group}>
                        <p className={styles.analiz}>{analiz.name}</p>
                        {analiz.items.map((item) => (
                            <div key={analiz.type} className={styles.item}>
                                <p className={cn(styles.name, { [styles.isFilled]: !!item.value })}>{item.name}</p>
                                <div className={styles.picker}>
                                    {isEditMode ? (
                                        <DatePicker
                                            value={item.value}
                                            onChange={onDateChangeHandler(analiz.type, item.name)}
                                        />
                                    ) : (
                                        <p className={styles.viewValue}>
                                            {item.value ? getDateFormat(item.value) : 'Отсутствует'}
                                        </p>
                                    )}
                                </div>
                                {item.value && (
                                    <div className={styles.icon}>
                                        <Success />
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                ))}
            </>
            {isEditMode && (
                <Button
                    fullWidth
                    onClick={onSaveButtonClickHandler}
                    className={cn(styles.confirm, { [styles.enabled]: isSaveButtonActive })}
                >
                    Сохранить изменения
                </Button>
            )}
        </div>
    );
};

export default AnalysesStep;
