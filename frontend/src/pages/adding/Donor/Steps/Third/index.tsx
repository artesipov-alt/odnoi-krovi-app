import { Button } from '@mui/material';
import cn from 'classnames';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useEffect, useState } from 'react';

import { StringDict } from 'api/reference';
import Alert, { View } from 'components/Alert';
import DatePicker from 'components/DatePicker';
import TextField from 'components/TextField';

import styles from './Third.module.less';

type Props = {
    healthStatus: string;
    surgicalList: string;
    medicationsList: string;
    lastDonation: Date | null;
    isNoLastDonation: boolean;
    healthStatusesDict: StringDict[];
    wasBloodTransfusion: boolean | null;
    isTakingMedications: boolean | null;
    wasSurgicalInterventions: boolean | null;
    onConfirmButtonClick: (step: number) => void;
    onChangeHealthStatus: (status: string) => void;
    onChangeSurgicalList: (status: string) => void;
    onChangeMedicationsList: (status: string) => void;
    onChangeLastDonationDate: (date: Date | null) => void;
    onChangeIsNoLastDonation: (newValue: boolean) => void;
    onChangeWasBloodTransfusion: (newValue: boolean) => void;
    onChangeIsTakingMedications: (newValue: boolean) => void;
    onChangeWasSurgicalInterventions: (newValue: boolean) => void;
};

const MAX_LETTERS = 250;

const Third: FC<Props> = ({
    lastDonation,
    healthStatus,
    surgicalList,
    medicationsList,
    isNoLastDonation,
    healthStatusesDict,
    wasBloodTransfusion,
    isTakingMedications,
    onConfirmButtonClick,
    onChangeHealthStatus,
    onChangeSurgicalList,
    onChangeMedicationsList,
    onChangeIsNoLastDonation,
    onChangeLastDonationDate,
    wasSurgicalInterventions,
    onChangeIsTakingMedications,
    onChangeWasBloodTransfusion,
    onChangeWasSurgicalInterventions,
}) => {
    const [surgicalListValue, setSurgicalListValue] = useState<string>(surgicalList);
    const [medicationsListValue, setMedicationsListValue] = useState<string>(medicationsList);
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(false);

    const onChangerHealthStatusHandler = (newStatus: string) => () => {
        if (newStatus === healthStatus) {
            return;
        }

        onChangeHealthStatus(newStatus);
    };

    const onIsNoLastDonationClickHandler = () => {
        if (isNoLastDonation) {
            return;
        }

        onChangeIsNoLastDonation(true);
    };

    const onWasBloodTransfusionClickHandler = (newValue: boolean) => () => {
        if (newValue === wasBloodTransfusion) {
            return;
        }

        onChangeWasBloodTransfusion(newValue);
    };

    const onIsTakingMedicationsClickHandler = (newValue: boolean) => () => {
        if (newValue !== null && newValue === isTakingMedications) {
            return;
        }

        if (!newValue) {
            setMedicationsListValue('');

            onChangeMedicationsList('');
        }

        onChangeIsTakingMedications(newValue);
    };

    const onMedicationsListChangeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setMedicationsListValue(value);

        if (!medicationsList.length && value) {
            onChangeMedicationsList(value);
        }

        if (medicationsList.length && !value) {
            onChangeMedicationsList('');
        }
    };

    const onMedicationsListBlurHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        onChangeMedicationsList(value);
    };

    const onWasSurgicalInterventionsClickHandler = (newValue: boolean) => () => {
        if (newValue !== null && newValue === wasSurgicalInterventions) {
            return;
        }

        if (!newValue) {
            setSurgicalListValue('');

            onChangeSurgicalList('');
        }

        onChangeWasSurgicalInterventions(newValue);
    };

    const onSurgicalListChangeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setSurgicalListValue(value);

        if (!surgicalList.length && value) {
            onChangeSurgicalList(value);
        }

        if (surgicalList.length && !value) {
            onChangeSurgicalList('');
        }
    };

    const onSurgicalListBlurHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        onChangeSurgicalList(value);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(3);
    };

    useEffect(() => {
        setIsConfirmButtonActive(
            !!healthStatus &&
                (isNoLastDonation || !!lastDonation) &&
                wasBloodTransfusion !== null &&
                (isTakingMedications ? !!medicationsList.length : isTakingMedications === false) &&
                (wasSurgicalInterventions ? !!surgicalList.length : wasSurgicalInterventions === false),
        );
    }, [
        lastDonation,
        healthStatus,
        surgicalList,
        medicationsList,
        isNoLastDonation,
        isTakingMedications,
        wasBloodTransfusion,
        wasSurgicalInterventions,
    ]);

    return (
        <>
            <FormItem title='Состояние здоровья'>
                <Alert
                    className={styles.alert}
                    view={View.INFO_WITHOUT_ICON}
                    text='Есть ли у питомца хронические, инфекционные, аутоиммунные, онкологические заболевания?'
                />
                <div className={styles.buttonsRow}>
                    {healthStatusesDict.map(({ label, value }) => (
                        <Button
                            key={value}
                            onClick={onChangerHealthStatusHandler(value)}
                            className={cn(styles.buttonsRowItem, {
                                [styles.checked]: healthStatus === value,
                            })}
                        >
                            {label}
                        </Button>
                    ))}
                </div>
            </FormItem>
            <FormItem title='Последняя донация'>
                <div className={styles.buttonsRow}>
                    <div className={styles.datePicker}>
                        <DatePicker
                            value={lastDonation}
                            onChange={onChangeLastDonationDate}
                            backgroundColor={isNoLastDonation ? '#EFF1F6' : undefined}
                        />
                    </div>
                    <Button
                        onClick={onIsNoLastDonationClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: isNoLastDonation })}
                    >
                        Не был донором
                    </Button>
                </div>
            </FormItem>
            <FormItem className={styles.formItem} title='Питомцу проводили переливания?'>
                <div className={styles.buttonsRow}>
                    <Button
                        onClick={onWasBloodTransfusionClickHandler(true)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: wasBloodTransfusion,
                        })}
                    >
                        Да
                    </Button>
                    <Button
                        onClick={onWasBloodTransfusionClickHandler(false)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: wasBloodTransfusion === false,
                        })}
                    >
                        Нет
                    </Button>
                </div>
            </FormItem>
            <FormItem
                title='Питомец принимает препараты?'
                className={cn(styles.formItem, { [styles.smallMargin]: isTakingMedications })}
            >
                <div className={styles.buttonsRow}>
                    <Button
                        onClick={onIsTakingMedicationsClickHandler(true)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: isTakingMedications,
                        })}
                    >
                        Да
                    </Button>
                    <Button
                        onClick={onIsTakingMedicationsClickHandler(false)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: isTakingMedications === false,
                        })}
                    >
                        Нет
                    </Button>
                </div>
            </FormItem>
            {isTakingMedications && (
                <div className={styles.textarea}>
                    <TextField
                        multiline
                        maxLength={250}
                        name='medications'
                        value={medicationsListValue}
                        htmlInputClass={styles.input}
                        placeholder='Укажите препараты'
                        onBlur={onMedicationsListBlurHandler}
                        onChange={onMedicationsListChangeHandler}
                    />
                    <div
                        className={cn(styles.counter, {
                            [styles.bigText]: medicationsListValue.length === MAX_LETTERS,
                        })}
                    >
                        {medicationsListValue.length}/{MAX_LETTERS}
                    </div>
                </div>
            )}

            <FormItem
                title='Питомцу недавно проводили хирургические вмешательства?'
                className={cn(styles.formItem, { [styles.smallMargin]: wasSurgicalInterventions })}
            >
                <div className={styles.buttonsRow}>
                    <Button
                        onClick={onWasSurgicalInterventionsClickHandler(true)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: wasSurgicalInterventions,
                        })}
                    >
                        Да
                    </Button>
                    <Button
                        onClick={onWasSurgicalInterventionsClickHandler(false)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: wasSurgicalInterventions === false,
                        })}
                    >
                        Нет
                    </Button>
                </div>
            </FormItem>
            {wasSurgicalInterventions && (
                <div className={styles.textarea}>
                    <TextField
                        multiline
                        maxLength={250}
                        name='surgical'
                        value={surgicalListValue}
                        htmlInputClass={styles.input}
                        onBlur={onSurgicalListBlurHandler}
                        onChange={onSurgicalListChangeHandler}
                        placeholder='Укажите, когда и какие хирургические вмешательства проводились питомцу'
                    />
                    <div
                        className={cn(styles.counter, {
                            [styles.bigText]: surgicalListValue.length === MAX_LETTERS,
                        })}
                    >
                        {surgicalListValue.length}/{MAX_LETTERS}
                    </div>
                </div>
            )}

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

export default Third;
