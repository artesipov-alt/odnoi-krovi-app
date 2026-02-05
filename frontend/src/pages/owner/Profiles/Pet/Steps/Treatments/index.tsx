import { Button } from '@mui/material';
import cn from 'classnames';
import Processing from 'imgs/svg/processing';
import FormItem from 'pages/adding/common/FormItem';
import { FC, useEffect, useState } from 'react';
import { getDateFormat } from 'utils/utils';

import { updatePet } from 'api/apiServices/updatePet';
import { Pet } from 'api/pets';
import DatePicker from 'components/DatePicker';

import Header from '../Header';
import ViewString from '../ViewString';
import styles from './Treatments.module.less';

type Props = {
    id: string;
    isEditMode: boolean;
    onClose: () => void;
    dewormingDate?: Date;
    onErrorUpdate: () => void;
    onSuccessUpdate: () => void;
    rabiesVaccinationDate?: Date;
    infectionVaccinationDate?: Date;
    ectoparasiteTreatmentDate?: Date;
};

const TreatmentsStep: FC<Props> = ({
    id,
    onClose,
    isEditMode,
    dewormingDate,
    onErrorUpdate,
    onSuccessUpdate,
    rabiesVaccinationDate,
    infectionVaccinationDate,
    ectoparasiteTreatmentDate,
}) => {
    const [isSaveButtonActive, setIsSaveButtonActive] = useState<boolean>(false);

    const [newRabiesVaccinationDate, setNewRabiesVaccinationDate] = useState<Date | null>(
        rabiesVaccinationDate ? new Date(rabiesVaccinationDate) : null,
    );
    const [newInfectionVaccinationDate, setNewInfectionVaccinationDate] = useState<Date | null>(
        infectionVaccinationDate ? new Date(infectionVaccinationDate) : null,
    );
    const [newEctoparasiteTreatmentDate, setNewEctoparasiteTreatmentDate] = useState<Date | null>(
        ectoparasiteTreatmentDate ? new Date(ectoparasiteTreatmentDate) : null,
    );
    const [newDewormingDate, setNewDewormingDate] = useState<Date | null>(
        dewormingDate ? new Date(dewormingDate) : null,
    );
    // const [isNotDeworming, setIsNotDeworming] = useState<boolean>(false);
    // const [isNotRabiesVaccination, setIsNotRabiesVaccination] = useState<boolean>(false);
    // const [isNotInfectionVaccination, setIsNotInfectionVaccination] = useState<boolean>(false);
    // const [isNotEctoparasiteTreatment, setIsNotEctoparasiteTreatment] = useState<boolean>(false);

    const isAllFieldsEmpty =
        !rabiesVaccinationDate && !infectionVaccinationDate && !ectoparasiteTreatmentDate && !dewormingDate;

    const onChangeRabiesVaccinationDateHandler = (date: Date | null) => {
        setNewRabiesVaccinationDate(date);
        // setIsNotRabiesVaccination(false);
    };

    const onIsNoRabiesVaccinationClickHandler = () => {
        setNewRabiesVaccinationDate(null);
        // setIsNotRabiesVaccination(true);
    };

    const onChangeInfectionsVaccinationDateHandler = (date: Date | null) => {
        setNewInfectionVaccinationDate(date);
        // setIsNotInfectionVaccination(false);
    };

    const onIsNoInfectionsVaccinationClickHandler = () => {
        setNewInfectionVaccinationDate(null);
        // setIsNotInfectionVaccination(true);
    };

    const onChangeEctoparasitesTreatmentDateHandler = (date: Date | null) => {
        setNewEctoparasiteTreatmentDate(date);
        // setIsNotEctoparasiteTreatment(false);
    };

    const onIsNoEctoparasitesTreatmentClickHandler = () => {
        setNewEctoparasiteTreatmentDate(null);
        // setIsNotEctoparasiteTreatment(true);
    };

    const onChangeDewormingDateHandler = (date: Date | null) => {
        setNewDewormingDate(date);
        // setIsNotDeworming(false);
    };

    const onIsNoDewormingClickHandler = () => {
        setNewDewormingDate(null);
        // setIsNotDeworming(true);
    };

    const onSaveButtonClickHandler = async () => {
        const newData: Partial<Pet> = {
            id,
            treatments: {
                dewormingDate: newDewormingDate || undefined,
                rabiesVaccinationDate: newRabiesVaccinationDate || undefined,
                infectionVaccinationDate: newInfectionVaccinationDate || undefined,
                ectoparasiteTreatmentDate: newEctoparasiteTreatmentDate || undefined,
            },
        };

        const { success } = await updatePet(newData as Pet);

        if (success) {
            onSuccessUpdate();
        } else {
            onErrorUpdate();
        }

        onClose();
    };

    const renderView = () => (
        <>
            <ViewString
                noAlignCanter
                name='Последняя вакцинация от бешенства'
                value={rabiesVaccinationDate ? getDateFormat(rabiesVaccinationDate) : 'Отсутствует'}
            />
            <ViewString
                noAlignCanter
                name='Последняя вакцинация от инфекций'
                value={infectionVaccinationDate ? getDateFormat(infectionVaccinationDate) : 'Отсутствует'}
            />
            <ViewString
                noAlignCanter
                name='Последняя обработка от эктопаразитов'
                value={ectoparasiteTreatmentDate ? getDateFormat(ectoparasiteTreatmentDate) : 'Отсутствует'}
            />
            <ViewString
                noAlignCanter
                name='Последняя дегельминтизация'
                value={dewormingDate ? getDateFormat(dewormingDate) : 'Отсутствует'}
            />
        </>
    );

    const renderEditView = () => (
        <>
            <FormItem title='Последняя вакцинация от бешенства'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            onChange={onChangeRabiesVaccinationDateHandler}
                            backgroundColor={
                                (!newRabiesVaccinationDate)
                                    ? '#EFF1F6'
                                    : undefined
                            }
                            value={newRabiesVaccinationDate ? new Date(newRabiesVaccinationDate) : null}
                        />
                    </div>
                    <Button
                        onClick={onIsNoRabiesVaccinationClickHandler}
                        className={cn(styles.buttonsRowItem, {
                            [styles.checked]: !newRabiesVaccinationDate,
                        })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Последняя вакцинация от инфекций'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            onChange={onChangeInfectionsVaccinationDateHandler}
                            backgroundColor={
                                (!newInfectionVaccinationDate)
                                    ? '#EFF1F6'
                                    : undefined
                            }
                            value={newInfectionVaccinationDate ? new Date(newInfectionVaccinationDate) : null}
                        />
                    </div>
                    <Button
                        onClick={onIsNoInfectionsVaccinationClickHandler}
                        className={cn(styles.buttonsRowItem, {
                            [styles.checked]: !newInfectionVaccinationDate,
                        })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Последняя обработка от эктопаразитов'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            value={newEctoparasiteTreatmentDate ? new Date(newEctoparasiteTreatmentDate) : null}
                            onChange={onChangeEctoparasitesTreatmentDateHandler}
                            backgroundColor={
                                (!newEctoparasiteTreatmentDate)
                                    ? '#EFF1F6'
                                    : undefined
                            }
                        />
                    </div>
                    <Button
                        onClick={onIsNoEctoparasitesTreatmentClickHandler}
                        className={cn(styles.buttonsRowItem, {
                            [styles.checked]: !newEctoparasiteTreatmentDate,
                        })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Последняя дегельминтизация'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            onChange={onChangeDewormingDateHandler}
                            backgroundColor={
                                (!newDewormingDate) ? '#EFF1F6' : undefined
                            }
                            value={newDewormingDate ? new Date(newDewormingDate) : null}
                        />
                    </div>
                    <Button
                        onClick={onIsNoDewormingClickHandler}
                        className={cn(styles.buttonsRowItem, {
                            [styles.checked]: !newDewormingDate,
                        })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <Button
                fullWidth
                onClick={onSaveButtonClickHandler}
                className={cn(styles.confirm, { [styles.enabled]: isSaveButtonActive })}
            >
                Сохранить изменения
            </Button>
        </>
    );

    useEffect(() => {
        const result =
            (!rabiesVaccinationDate && !!newRabiesVaccinationDate) ||
              (!!rabiesVaccinationDate && !newRabiesVaccinationDate) ||
              (!!rabiesVaccinationDate &&
                  !!newRabiesVaccinationDate &&
                  new Date(rabiesVaccinationDate).getTime() !== newRabiesVaccinationDate.getTime()) ||
              (!infectionVaccinationDate && !!newInfectionVaccinationDate) ||
              (!!infectionVaccinationDate && !newInfectionVaccinationDate) ||
              (!!infectionVaccinationDate &&
                  !!newInfectionVaccinationDate &&
                  new Date(infectionVaccinationDate).getTime() !== newInfectionVaccinationDate.getTime()) ||
              (!ectoparasiteTreatmentDate && !!newEctoparasiteTreatmentDate) ||
              (!!ectoparasiteTreatmentDate && !newEctoparasiteTreatmentDate) ||
              (!!ectoparasiteTreatmentDate &&
                  !!newEctoparasiteTreatmentDate &&
                  new Date(ectoparasiteTreatmentDate).getTime() !== newEctoparasiteTreatmentDate.getTime()) ||
              (!dewormingDate && !!newDewormingDate) ||
              (!!dewormingDate && !newDewormingDate) ||
              (!!dewormingDate &&
                  !!newDewormingDate &&
                  new Date(dewormingDate).getTime() !== newDewormingDate.getTime());
            // : (isNotRabiesVaccination || !!newRabiesVaccinationDate) &&
            //   (isNotInfectionVaccination || !!newInfectionVaccinationDate) &&
            //   (isNotEctoparasiteTreatment || !!newEctoparasiteTreatmentDate) &&
            //   (isNotDeworming || !!newDewormingDate);

        setIsSaveButtonActive(result);
    }, [
        dewormingDate,
        // isNotDeworming,
        newDewormingDate,
        isAllFieldsEmpty,
        rabiesVaccinationDate,
        // isNotRabiesVaccination,
        infectionVaccinationDate,
        newRabiesVaccinationDate,
        ectoparasiteTreatmentDate,
        // isNotInfectionVaccination,
        // isNotEctoparasiteTreatment,
        newInfectionVaccinationDate,
        newEctoparasiteTreatmentDate,
    ]);

    return (
        <div className={styles.wrapper}>
            <Header title='Обработки' onClose={onClose} isEditMode={isEditMode} icon={<Processing />} />
            {isEditMode ? renderEditView() : renderView()}
        </div>
    );
};

export default TreatmentsStep;
