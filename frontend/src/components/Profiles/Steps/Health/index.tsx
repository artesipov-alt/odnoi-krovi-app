import { Button } from '@mui/material';
import cn from 'classnames';
import Health from 'imgs/svg/health';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useEffect, useState } from 'react';

import { updatePet } from 'api/apiServices/updatePet';
import { Pet } from 'api/pets';
import { Dict } from 'api/reference';
import Alert, { View } from 'components/Alert';
import Layout from 'components/Layout';
import TextField from 'components/TextField';

import Header from '../Header';
import ViewString from '../ViewString';
import styles from './Health.module.less';

type Props = {
    id: string;
    isEditMode: boolean;
    onClose: () => void;
    transfused?: boolean;
    medications?: string;
    healthStatus?: string;
    isProfileLock?: boolean;
    onErrorUpdate?: () => void;
    onSuccessUpdate?: () => void;
    surgicalInterventions?: string;
    healthStatusesDict: Dict[];
};

const MAX_LETTERS = 250;

const HealthStep: FC<Props> = ({
    id,
    onClose,
    isEditMode,
    transfused,
    medications,
    healthStatus,
    onErrorUpdate,
    isProfileLock,
    onSuccessUpdate,
    healthStatusesDict,
    surgicalInterventions,
}) => {
    const [newMedications, setNewMedications] = useState<string | undefined>(medications);
    const [isSaveButtonActive, setIsSaveButtonActive] = useState<boolean>(false);
    const [newHealthStatus, setNewHealthStatus] = useState<string | undefined>(healthStatus);
    const [newTransfusedValue, setNewTransfusedValue] = useState<boolean | undefined>(transfused);
    const [newSurgicalInterventions, setNewSurgicalInterventions] = useState<string | undefined>(surgicalInterventions);
    const [isTakingMedications, setIsTakingMedications] = useState<boolean | undefined>(
        medications === undefined ? undefined : !!medications,
    );
    const [wasSurgicalInterventions, setWasSurgicalInterventions] = useState<boolean | undefined>(
        surgicalInterventions === undefined ? undefined : !!surgicalInterventions,
    );

    const isAllFieldsEmpty = !medications && !transfused && !healthStatus && !surgicalInterventions;

    const onChangerHealthStatusHandler = (newValue: string) => () => {
        if (newValue === newHealthStatus) {
            return;
        }

        setNewHealthStatus(newValue);
    };

    const onWasBloodTransfusionClickHandler = (newValue: boolean) => () => {
        if (newValue === newTransfusedValue) {
            return;
        }

        setNewTransfusedValue(newValue);
    };

    const onIsTakingMedicationsClickHandler = (newValue: boolean) => () => {
        if (newValue !== null && newValue === isTakingMedications) {
            return;
        }

        if (!newValue) {
            setNewMedications('');
        }

        setIsTakingMedications(newValue);
    };

    const onMedicationsListChangeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setNewMedications(value);
    };

    const onWasSurgicalInterventionsClickHandler = (newValue: boolean) => () => {
        if (newValue !== null && newValue === wasSurgicalInterventions) {
            return;
        }

        if (!newValue) {
            setNewSurgicalInterventions('');
        }

        setWasSurgicalInterventions(newValue);
    };

    const onSurgicalListChangeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setNewSurgicalInterventions(value);
    };

    const onSaveButtonClickHandler = async () => {
        const newData: Partial<Pet> = {
            id,
            health: {
                healthStatus: newHealthStatus!,
                transfused: newTransfusedValue,
                medications: newMedications || undefined,
                surgicalInterventions: newSurgicalInterventions || undefined,
            },
        };

        const { success } = await updatePet(newData as Pet);

        if (success) {
            onSuccessUpdate?.();
        } else {
            onErrorUpdate?.();
        }

        onClose();
    };

    const renderEditView = () => (
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
                                [styles.checked]: newHealthStatus === value,
                            })}
                        >
                            {label}
                        </Button>
                    ))}
                </div>
            </FormItem>
            <FormItem className={styles.formItem} title='Питомцу проводили переливания?'>
                <div className={cn(styles.buttonsRow, { [styles.isLock]: isProfileLock })}>
                    <Button
                        onClick={onWasBloodTransfusionClickHandler(true)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: newTransfusedValue,
                        })}
                    >
                        Да
                    </Button>
                    <Button
                        onClick={onWasBloodTransfusionClickHandler(false)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.yewNo]: true,
                            [styles.checked]: newTransfusedValue === false,
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
                        name='medications'
                        maxLength={MAX_LETTERS}
                        value={newMedications || ''}
                        htmlInputClass={styles.input}
                        placeholder='Укажите препараты'
                        onChange={onMedicationsListChangeHandler}
                    />
                    <div
                        className={cn(styles.counter, {
                            [styles.bigText]: (newMedications || '').length === MAX_LETTERS,
                        })}
                    >
                        {(newMedications || '').length}/{MAX_LETTERS}
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
                        name='surgical'
                        maxLength={MAX_LETTERS}
                        htmlInputClass={styles.input}
                        value={newSurgicalInterventions || ''}
                        onChange={onSurgicalListChangeHandler}
                        placeholder='Укажите, когда и какие хирургические вмешательства проводились питомцу'
                    />
                    <div
                        className={cn(styles.counter, {
                            [styles.bigText]: (newSurgicalInterventions || '').length === MAX_LETTERS,
                        })}
                    >
                        {(newSurgicalInterventions || '').length}/{MAX_LETTERS}
                    </div>
                </div>
            )}
            <Button
                fullWidth
                onClick={onSaveButtonClickHandler}
                className={cn(styles.confirm, { [styles.enabled]: isSaveButtonActive })}
            >
                Сохранить изменения
            </Button>
        </>
    );

    const renderView = () => (
        <>
            {healthStatus && (
                <ViewString
                    noAlignCanter
                    name='Состояние здоровья'
                    tooltip='Есть ли у питомца хронические, инфекционные, аутоиммунные, онкологические заболевания?'
                    value={healthStatusesDict.filter(({ value }) => value === healthStatus)[0].label}
                />
            )}
            <ViewString noAlignCanter name='Питомцу проводили переливания?' value={transfused ? 'Да' : 'Нет'} />
            <ViewString
                noAlignCanter
                value={medications ? '' : 'Нет'}
                descr={medications || undefined}
                name='Питомец принимает препараты'
            />
            <ViewString
                noAlignCanter
                value={surgicalInterventions ? '' : 'Нет'}
                descr={surgicalInterventions || undefined}
                name='Питомцу проводили хирургические вмешательства'
            />
        </>
    );

    useEffect(() => {
        if ((isTakingMedications && !newMedications) || (wasSurgicalInterventions && !newSurgicalInterventions)) {
            setIsSaveButtonActive(false);

            return;
        }

        setIsSaveButtonActive(
            !isAllFieldsEmpty
                ? newHealthStatus !== healthStatus ||
                      newTransfusedValue !== transfused ||
                      (!!newMedications && newMedications !== medications) ||
                      (!!medications && isTakingMedications === false) ||
                      (!!newSurgicalInterventions && newSurgicalInterventions !== surgicalInterventions) ||
                      (!!surgicalInterventions && wasSurgicalInterventions === false)
                : !!newHealthStatus &&
                      newTransfusedValue !== undefined &&
                      (isTakingMedications ? !!newMedications?.length : isTakingMedications === false) &&
                      (wasSurgicalInterventions
                          ? !!newSurgicalInterventions?.length
                          : wasSurgicalInterventions !== undefined),
        );
    }, [
        transfused,
        medications,
        healthStatus,
        newMedications,
        newHealthStatus,
        isAllFieldsEmpty,
        newTransfusedValue,
        isTakingMedications,
        surgicalInterventions,
        newSurgicalInterventions,
        wasSurgicalInterventions,
    ]);

    return (
        <Layout className={styles.wrapper}>
            <Header title='Здоровье' onClose={onClose} isEditMode={isEditMode} icon={<Health />} />
            {isEditMode ? renderEditView() : renderView()}
        </Layout>
    );
};

export default HealthStep;
