import { Button } from '@mui/material';
import cn from 'classnames';
import FormItem from 'pages/adding/common/FormItem';
import { FC, useEffect, useState } from 'react';

import DatePicker from 'components/DatePicker';

import styles from './Fourth.module.less';

type Props = {
    isNoDeworming: boolean;
    dewormingDate: Date | null;
    isNoRabiesVaccination: boolean;
    rabiesVaccinationDate: Date | null;
    isNoInfectionsVaccination: boolean;
    isNoEctoparasitesTreatment: boolean;
    infectionsVaccinationDate: Date | null;
    ectoparasitesTreatmentDate: Date | null;
    onConfirmButtonClick: (step: number) => void;
    onChangeIsNoDeworming: (newValue: boolean) => void;
    onChangeDewormingDate: (newDate: Date | null) => void;
    onChangeIsNoRabiesVaccination: (newValue: boolean) => void;
    onChangeRabiesVaccinationDate: (newDate: Date | null) => void;
    onChangeIsNoInfectionsVaccination: (newValue: boolean) => void;
    onChangeIsNoEctoparasitesTreatment: (newValue: boolean) => void;
    onChangeInfectionsVaccinationDate: (newDate: Date | null) => void;
    onChangeEctoparasitesTreatmentDate: (newDate: Date | null) => void;
};

const Fourth: FC<Props> = ({
    isNoDeworming,
    dewormingDate,
    onConfirmButtonClick,
    rabiesVaccinationDate,
    isNoRabiesVaccination,
    onChangeDewormingDate,
    onChangeIsNoDeworming,
    isNoInfectionsVaccination,
    infectionsVaccinationDate,
    isNoEctoparasitesTreatment,
    ectoparasitesTreatmentDate,
    onChangeIsNoRabiesVaccination,
    onChangeRabiesVaccinationDate,
    onChangeIsNoInfectionsVaccination,
    onChangeInfectionsVaccinationDate,
    onChangeIsNoEctoparasitesTreatment,
    onChangeEctoparasitesTreatmentDate,
}) => {
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(false);

    const onIsNoRabiesVaccinationClickHandler = () => {
        if (isNoRabiesVaccination) {
            return;
        }

        onChangeIsNoRabiesVaccination(true);
    };

    const onIsNoInfectionsVaccinationClickHandler = () => {
        if (isNoInfectionsVaccination) {
            return;
        }

        onChangeIsNoInfectionsVaccination(true);
    };

    const onIsNoEctoparasitesTreatmentClickHandler = () => {
        if (isNoEctoparasitesTreatment) {
            return;
        }

        onChangeIsNoEctoparasitesTreatment(true);
    };

    const onIsNoDewormingClickHandler = () => {
        if (isNoDeworming) {
            return;
        }

        onChangeIsNoDeworming(true);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(4);
    };

    useEffect(() => {
        setIsConfirmButtonActive(
            (!!rabiesVaccinationDate || isNoRabiesVaccination) &&
                (!!infectionsVaccinationDate || isNoInfectionsVaccination) &&
                (!!ectoparasitesTreatmentDate || isNoEctoparasitesTreatment) &&
                (!!dewormingDate || isNoDeworming),
        );
    }, [
        isNoDeworming,
        dewormingDate,
        rabiesVaccinationDate,
        isNoRabiesVaccination,
        isNoInfectionsVaccination,
        infectionsVaccinationDate,
        isNoEctoparasitesTreatment,
        ectoparasitesTreatmentDate,
    ]);

    return (
        <>
            <FormItem title='Последняя вакцинация от бешенства'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            value={rabiesVaccinationDate}
                            onChange={onChangeRabiesVaccinationDate}
                            backgroundColor={isNoRabiesVaccination ? '#EFF1F6' : undefined}
                        />
                    </div>
                    <Button
                        onClick={onIsNoRabiesVaccinationClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: isNoRabiesVaccination })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Последняя вакцинация от инфекций'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            value={infectionsVaccinationDate}
                            onChange={onChangeInfectionsVaccinationDate}
                            backgroundColor={isNoInfectionsVaccination ? '#EFF1F6' : undefined}
                        />
                    </div>
                    <Button
                        onClick={onIsNoInfectionsVaccinationClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: isNoInfectionsVaccination })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Последняя обработка от эктопаразитов'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            value={ectoparasitesTreatmentDate}
                            onChange={onChangeEctoparasitesTreatmentDate}
                            backgroundColor={isNoEctoparasitesTreatment ? '#EFF1F6' : undefined}
                        />
                    </div>
                    <Button
                        onClick={onIsNoEctoparasitesTreatmentClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: isNoEctoparasitesTreatment })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Последняя дегельминтизация'>
                <div className={styles.buttonsRow}>
                    <div className={styles.buttonsRowItem}>
                        <DatePicker
                            value={dewormingDate}
                            onChange={onChangeDewormingDate}
                            backgroundColor={isNoDeworming ? '#EFF1F6' : undefined}
                        />
                    </div>
                    <Button
                        onClick={onIsNoDewormingClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: isNoDeworming })}
                    >
                        Отсутствует
                    </Button>
                </div>
            </FormItem>
            <Button
                fullWidth
                onClick={onConfirmButtonClickHandler}
                className={cn(styles.confirm, { [styles.enabled]: isConfirmButtonActive })}
            >
                Далее
            </Button>
        </>
    );
};

export default Fourth;
