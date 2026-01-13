import { Button } from '@mui/material';
import cn from 'classnames';
import FemaleIcon from 'imgs/svg/femaleIcon';
import MaleIcon from 'imgs/svg/maleIcon';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useEffect, useState } from 'react';
import { regexInt } from 'utils/regexps';

import { PetGenderDict, PetTypeDict } from 'api/reference';
import { PetGender } from 'api/types';
import Alert from 'components/Alert';
import DatePicker from 'components/DatePicker';
import ImgEditor from 'components/ImgEditor';
import TextField from 'components/TextField';

import { BirthDate } from '../../types';
import styles from './First.module.less';

type Props = {
    name: string;
    petType: string;
    petGender: string;
    photo: File | null;
    chipNumber: string;
    exactDate: Date | null;
    petTypes: PetTypeDict[];
    petGenders: PetGenderDict[];
    approximateDateYear: string;
    approximateDateMonth: string;
    birthDateType: BirthDate | null;
    onChangeName: (name: string) => void;
    onChangePetType: (type: string) => void;
    onLoadPhoto: (photo: File | null) => void;
    onChangePetGender: (type: string) => void;
    onChangeChipNumber: (number: string) => void;
    onConfirmButtonClick: (step: number) => void;
    onChangeExactDate: (date: Date | null) => void;
    onChangeBirthDateType: (type: BirthDate) => void;
    onChangeApproximateDateYear: (newYear: string) => void;
    onChangeApproximateDateMonth: (newMonth: string) => void;
};

const First: FC<Props> = ({
    name,
    photo,
    petType,
    petTypes,
    exactDate,
    petGender,
    chipNumber,
    petGenders,
    onLoadPhoto,
    onChangeName,
    birthDateType,
    onChangePetType,
    onChangePetGender,
    onChangeExactDate,
    onChangeChipNumber,
    approximateDateYear,
    onConfirmButtonClick,
    approximateDateMonth,
    onChangeBirthDateType,
    onChangeApproximateDateYear,
    onChangeApproximateDateMonth,
}) => {
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(false);

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(1);
    };

    const onChangeNameHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        onChangeName(value);
    };

    const onPetTypeButtonClickHandler = (newPetType: string) => () => {
        onChangePetType(newPetType);
    };

    const onPetGenderButtonClickHandler = (newGender: string) => () => {
        onChangePetGender(newGender);
    };

    const onChangeChipNumberHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        onChangeChipNumber(value.trim());
    };

    const onNoChipClickHandler = () => {
        onChangeChipNumber('none');
    };

    const onBirthDateTypeClickHandler = (newType: BirthDate) => () => {
        onChangeBirthDateType(newType);
    };

    const onChangeApproximateDateHandler =
        (type: 'year' | 'months') =>
        ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
            const newValue = value.trim();

            if (newValue === '0') {
                return;
            }

            if (value && !newValue.match(regexInt)) {
                return;
            }

            if (type === 'months' && Number(newValue) > 11) {
                onChangeApproximateDateMonth('11');

                return;
            }

            if (type === 'year' && Number(newValue) > 40) {
                onChangeApproximateDateYear('40');

                return;
            }

            if (type === 'year') {
                onChangeApproximateDateYear(newValue);
            } else {
                onChangeApproximateDateMonth(newValue);
            }
        };

    useEffect(() => {
        setIsConfirmButtonActive(
            !!photo &&
                !!name &&
                !!petType &&
                !!petGender &&
                !!chipNumber &&
                (birthDateType === BirthDate.EXACT_DATE
                    ? !!exactDate
                    : !!approximateDateMonth || !!approximateDateYear),
        );
    }, [
        name,
        photo,
        petType,
        exactDate,
        petGender,
        chipNumber,
        birthDateType,
        approximateDateYear,
        approximateDateMonth,
    ]);

    return (
        <>
            <ImgEditor src={photo} onLoad={onLoadPhoto} className={styles.photo} />
            <FormItem title='Кличка'>
                <TextField name='name' value={name} placeholder='Как зовут питомца?' onChange={onChangeNameHandler} />
            </FormItem>
            {!!petTypes.length && (
                <FormItem title='Вид'>
                    <div className={styles.buttonsRow}>
                        {petTypes.map(({ label, value }) => (
                            <Button
                                key={value}
                                onClick={onPetTypeButtonClickHandler(value)}
                                className={cn(styles.buttonsRowItem, { [styles.checked]: petType === value })}
                            >
                                {label}
                            </Button>
                        ))}
                    </div>
                </FormItem>
            )}
            {!!petGenders.length && (
                <FormItem title='Пол'>
                    <div className={styles.buttonsRow}>
                        {petGenders.map(({ label, value }) => {
                            const isMale = value === PetGender.MALE;
                            const isChecked = petGender === value;

                            return (
                                <Button
                                    key={value}
                                    startIcon={
                                        <div
                                            className={cn(styles.buttonIcon, {
                                                [styles.female]: !isMale,
                                                [styles.isChecked]: isChecked,
                                            })}
                                        >
                                            {isMale ? <MaleIcon /> : <FemaleIcon />}
                                        </div>
                                    }
                                    onClick={onPetGenderButtonClickHandler(value)}
                                    className={cn(styles.buttonsRowItem, { [styles.checked]: isChecked })}
                                >
                                    {label}
                                </Button>
                            );
                        })}
                    </div>
                </FormItem>
            )}
            <FormItem title='Чип'>
                <div className={cn(styles.buttonsRow, { [styles.chip]: true })}>
                    <div className={styles.chipField}>
                        <TextField
                            name='chip'
                            placeholder='Укажите № чипа'
                            onChange={onChangeChipNumberHandler}
                            value={chipNumber === 'none' ? '' : chipNumber}
                            inputClass={chipNumber === 'none' ? styles.input : undefined}
                        />
                    </div>
                    <Button
                        onClick={onNoChipClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: chipNumber === 'none' })}
                    >
                        Отсутвует
                    </Button>
                </div>
            </FormItem>
            <FormItem title='Дата рождения'>
                <div className={styles.buttonsRow}>
                    <Button
                        onClick={onBirthDateTypeClickHandler(BirthDate.EXACT_DATE)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.checked]: birthDateType === BirthDate.EXACT_DATE,
                        })}
                    >
                        Точная дата
                    </Button>
                    <Button
                        onClick={onBirthDateTypeClickHandler(BirthDate.APPROXIMATE_DATE)}
                        className={cn(styles.buttonsRowItem, {
                            [styles.checked]: birthDateType === BirthDate.APPROXIMATE_DATE,
                        })}
                    >
                        Примерный возраст
                    </Button>
                </div>
                {birthDateType === BirthDate.EXACT_DATE && (
                    <div className={styles.picker}>
                        <DatePicker value={exactDate} onChange={onChangeExactDate} />
                    </div>
                )}
                {birthDateType === BirthDate.APPROXIMATE_DATE && (
                    <>
                        <div className={styles.buttonsRow}>
                            <FormItem className={styles.ageItem} title='Полных лет'>
                                <TextField
                                    name='year'
                                    isDigitInput
                                    placeholder='0'
                                    value={approximateDateYear}
                                    onChange={onChangeApproximateDateHandler('year')}
                                />
                            </FormItem>
                            <FormItem className={styles.ageItem} title='Полных месяцев'>
                                <TextField
                                    name='months'
                                    isDigitInput
                                    placeholder='0'
                                    value={approximateDateMonth}
                                    onChange={onChangeApproximateDateHandler('months')}
                                />
                            </FormItem>
                        </div>
                        <Alert text='Должно быть заполнено хотя бы одно поле' />
                    </>
                )}
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

export default First;
