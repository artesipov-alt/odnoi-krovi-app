import { Button, TextField as MuiTextField } from '@mui/material';
import Autocomplete from '@mui/material/Autocomplete';
import cn from 'classnames';
import { BloodAndBreedGroupsDict } from 'hooks/useDicts';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useEffect, useState } from 'react';
import { regexReal } from 'utils/regexps';

import { Dict } from 'api/reference';
import { PetGender, PetType } from 'api/types';
import Alert from 'components/Alert';
import TextField from 'components/TextField';

import styles from './Second.module.less';

type Props = {
    weight: string;
    petType: string;
    petGender: string;
    breed: Dict | null;
    bloodGroup: string;
    livingCondition: string;
    reproductiveStatus: string;
    livingConditionsDict: Dict[];
    breedsDict: BloodAndBreedGroupsDict;
    reproductiveStatusesDict: Dict[];
    bloodGroupDict: BloodAndBreedGroupsDict;
    onChangeWeight: (weight: string) => void;
    onChangeBreed: (breed: Dict | null) => void;
    onChangeBloodGroup: (blood: string) => void;
    onConfirmButtonClick: (step: number) => void;
    onChangeReproductiveStatus: (status: string) => void;
    onChangeLivingCondition: (condition: string) => void;
};

const Second: FC<Props> = ({
    breed,
    weight,
    petType,
    petGender,
    breedsDict,
    bloodGroup,
    onChangeBreed,
    bloodGroupDict,
    onChangeWeight,
    livingCondition,
    reproductiveStatus,
    onChangeBloodGroup,
    livingConditionsDict,
    onConfirmButtonClick,
    onChangeLivingCondition,
    reproductiveStatusesDict,
    onChangeReproductiveStatus,
}) => {
    const [inputValue, setInputValue] = useState('');
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(false);

    const onChangeBloodGroupHandler = (newBloodGroup: string) => () => {
        if (newBloodGroup === bloodGroup) {
            return;
        }

        onChangeBloodGroup(newBloodGroup);
    };

    const onChangeLivingConditionsHandler = (newCondition: string) => () => {
        if (newCondition === livingCondition) {
            return;
        }

        onChangeLivingCondition(newCondition);
    };

    const onChangerReproductiveStatusHandler = (newStatus: string) => () => {
        if (newStatus === reproductiveStatus) {
            onChangeReproductiveStatus('none');

            return;
        }

        onChangeReproductiveStatus(newStatus);
    };

    const onBlurPetWeightHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        if (!value) {
            return;
        }

        const [int, real] = value.split(',');

        if (Number(int) === 0 && Number(real || '0') === 0) {
            onChangeWeight('0,05');
        }
    };

    const onChangePetWeightHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        const newValue = value.replaceAll(' ', '');

        if (!newValue) {
            onChangeWeight('');

            return;
        }

        if (!newValue.match(regexReal)) {
            return;
        }

        if (petType === PetType.CAT && Number(newValue) > 15) {
            return;
        }

        if (petType === PetType.DOG && Number(newValue) > 170) {
            return;
        }

        onChangeWeight(newValue);
    };

    const onAutocompleteInputChangeHandler = (_, newValue: string) => {
        setInputValue(newValue);
    };

    const onAutocompleteChangeHandler = (_, newValue: Dict | null) => {
        onChangeBreed(newValue);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(2);
    };

    useEffect(() => {
        setIsConfirmButtonActive(!!bloodGroup && !!weight && !!breed && !!livingCondition);
    }, [bloodGroup, breed, livingCondition, weight]);

    return (
        <>
            <FormItem title='Группа крови питомца'>
                <>
                    <div className={cn(styles.buttonsRow, { [styles.dogGroup]: petType === PetType.DOG })}>
                        {bloodGroupDict[petType].map(({ label, value }) => (
                            <Button
                                key={value}
                                onClick={onChangeBloodGroupHandler(value)}
                                className={cn(styles.buttonsRowItem, { [styles.checked]: bloodGroup === value })}
                            >
                                {label !== 'UNKNOWN' ? label : 'Не знаю'}
                            </Button>
                        ))}
                    </div>
                    {petType === PetType.DOG && (
                        <Alert
                            text={
                                <>
                                    Сведения вносятся по системе DEA.
                                    <br />
                                    Не используйте данные из других систем (KAI, DAL)
                                </>
                            }
                        />
                    )}
                </>
            </FormItem>
            {/* eslint-disable-next-line no-nested-ternary */}
            <FormItem title='Вес' subtitle={!petType ? undefined : petType === PetType.CAT ? 'до 15 кг' : 'до 170 кг'}>
                <TextField
                    name='weight'
                    value={weight}
                    isDigitInput
                    onBlur={onBlurPetWeightHandler}
                    onChange={onChangePetWeightHandler}
                    placeholder='Сколько весит питомец?'
                    endAdornment={<div className={styles.endAdornment}>кг</div>}
                />
            </FormItem>
            <FormItem title='Порода'>
                <Autocomplete
                    id='breed'
                    value={breed}
                    inputValue={inputValue}
                    options={breedsDict[petType]}
                    onChange={onAutocompleteChangeHandler}
                    noOptionsText='Нет подходящих вариантов'
                    onInputChange={onAutocompleteInputChangeHandler}
                    renderInput={(params) => <MuiTextField {...params} placeholder='Выберите из списка' />}
                    sx={{
                        '& .MuiOutlinedInput-root': {
                            backgroundColor: 'white',
                            height: '50px',
                            padding: '0 9px',
                            borderRadius: '16px',
                            color: breed && '#8B7069',
                            '& .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#dee2e9',
                            },
                            '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#dee2e9',
                                borderWidth: '1px',
                            },
                        },
                    }}
                />
            </FormItem>
            <FormItem title='Условия содержания'>
                <div className={styles.buttonsRow}>
                    {livingConditionsDict.map(({ label, value }) => (
                        <Button
                            key={value}
                            onClick={onChangeLivingConditionsHandler(value)}
                            className={cn(styles.buttonsRowItem, { [styles.checked]: livingCondition === value })}
                        >
                            {label}
                        </Button>
                    ))}
                </div>
            </FormItem>
            {petGender === PetGender.FEMALE && (
                <FormItem title='Состояние питомца (выберите, если есть)'>
                    <div className={styles.buttonsRow}>
                        {reproductiveStatusesDict.map(({ label, value }) => (
                            <Button
                                key={value}
                                onClick={onChangerReproductiveStatusHandler(value)}
                                className={cn(styles.buttonsRowItem, {
                                    [styles.checked]: reproductiveStatus === value,
                                })}
                            >
                                {label}
                            </Button>
                        ))}
                    </div>
                </FormItem>
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

export default Second;
