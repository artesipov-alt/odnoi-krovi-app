import { Button, TextField as MuiTextField } from '@mui/material';
import Autocomplete from '@mui/material/Autocomplete';
import cn from 'classnames';
import { BloodAndBreedGroupsDict } from 'hooks/useDicts';
import FemaleIcon from 'imgs/svg/femaleIcon';
import MaleIcon from 'imgs/svg/maleIcon';
import Params from 'imgs/svg/params';
import FormItem from 'pages/adding/common/FormItem';
import { BirthDate } from 'pages/adding/Donor/types';
import { ChangeEvent, FC, useEffect, useState } from 'react';
import { regexInt, regexReal } from 'utils/regexps';
import { getCorrectDeclension, Variants } from 'utils/utils';

import { updatePet } from 'api/apiServices/updatePet';
import { Pet } from 'api/pets';
import { Dict, PetGenderDict, PetTypeDict } from 'api/reference';
import { PetGender, PetType } from 'api/types';
import Alert from 'components/Alert';
import DatePicker from 'components/DatePicker';
import TextField from 'components/TextField';

import Header from '../Header';
import ViewString from '../ViewString';
import styles from './ParamsStep.module.less';

type Props = {
    id: string;
    name: string;
    type: string;
    gender?: string;
    weightKg: number;
    breedId?: string;
    birthDate?: Date;
    bloodGroup: string;
    onClose: () => void;
    chipNumber?: string;
    isEditMode: boolean;
    petTypes: PetTypeDict[];
    livingCondition?: string;
    onErrorUpdate: () => void;
    reproductiveStatus?: string;
    petGenders: PetGenderDict[];
    onSuccessUpdate: () => void;
    livingConditionsDict: Dict[];
    breedsDict: BloodAndBreedGroupsDict;
    reproductiveStatusesDict: Dict[];
    bloodGroupDict: BloodAndBreedGroupsDict;
};

const calculateAge = (birthDate: string | Date) => {
    const birth = new Date(birthDate);
    const today = new Date();

    // Устанавливаем временные переменные для года, месяца и дня
    let years = today.getFullYear() - birth.getFullYear();
    let months = today.getMonth() - birth.getMonth();

    // Корректируем, если текущий день меньше дня рождения — значит, год/месяц ещё не прошёл полностью
    if (today.getDate() < birth.getDate()) {
        months--;
    }

    // Если месяцев получилось отрицательное количество, корректируем годы
    if (months < 0) {
        years--;
        months += 12;
    }

    const yearsDescr = years ? `${years} ${getCorrectDeclension(Variants.YEARS, +years)}` : '';
    const monthDescr = months ? `${months} ${getCorrectDeclension(Variants.MONTHS, +months)}` : '';

    return yearsDescr || monthDescr ? `${yearsDescr} ${monthDescr}` : 'Меньше месяца';
};

const ParamsStep: FC<Props> = ({
    id,
    type,
    name,
    gender,
    breedId,
    onClose,
    weightKg,
    petTypes,
    birthDate,
    breedsDict,
    bloodGroup,
    petGenders,
    chipNumber,
    isEditMode,
    onErrorUpdate,
    bloodGroupDict,
    onSuccessUpdate,
    livingCondition,
    reproductiveStatus,
    livingConditionsDict,
    reproductiveStatusesDict,
}) => {
    const [newExactDate, setNewExactDate] = useState<Date | null>(birthDate || null);
    const [birthDateType, setBirthDateType] = useState<BirthDate | null>(birthDate ? BirthDate.EXACT_DATE : null);
    const [approximateDateYear, setApproximateDateYear] = useState<string>('');
    const [approximateDateMonth, setApproximateDateMonth] = useState<string>('');

    const [newName, setNewName] = useState<string>(name);
    const [newPetType, setNewPetType] = useState<string>(type);
    const [newBreed, setNewBreed] = useState<Dict | null>(
        (breedsDict[type] || []).filter(({ value }) => value === breedId)?.[0] || null,
    );
    const [newGender, setNewGender] = useState<string | undefined>(gender);
    const [newBloodGroup, setNewBloodGroup] = useState<string>(bloodGroup);
    const [newWeight, setNewWeight] = useState<string>(`${weightKg}`);
    const [newChipNumber, setNewChipNumber] = useState<string>(chipNumber || 'none');
    const [autocompleteInputValue, setAutocompleteInputValue] = useState('');
    const [isSaveButtonActive, setIsSaveButtonActive] = useState<boolean>(false);
    const [newLivingCondition, setNewLivingCondition] = useState<string | undefined>(livingCondition);
    const [newReproductiveStatus, setNewReproductiveStatus] = useState<string | undefined>(reproductiveStatus);

    const onChangeNameHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setNewName(value);
    };

    const onPetTypeChangeHandler = (newType: string) => () => {
        setNewPetType(newType);
        setNewBreed(null);
        setNewBloodGroup('');
        setNewWeight('');
        setNewReproductiveStatus(undefined);
    };

    const onPetGenderChangeHandler = (newPetGender: string) => () => {
        setNewGender(newPetGender);
        setNewReproductiveStatus(undefined);
    };

    const onChangeChipNumberHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setNewChipNumber(value.trim());
    };

    const onNoChipClickHandler = () => {
        setNewChipNumber('none');
    };

    const onChangeBloodGroupHandler = (newGroup: string) => () => {
        if (newBloodGroup === newGroup) {
            return;
        }

        setNewBloodGroup(newGroup);
    };

    const onBlurPetWeightHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        if (!value) {
            return;
        }

        const [int, real] = value.split(',');

        if (Number(int) === 0 && Number(real || '0') === 0) {
            setNewWeight('0,05');
        }
    };

    const onChangePetWeightHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        const newValue = value.replaceAll(' ', '');

        if (!newValue) {
            setNewWeight('');

            return;
        }

        if (!newValue.match(regexReal)) {
            return;
        }

        if (newPetType === PetType.CAT && Number(newValue) > 15) {
            return;
        }

        if (newPetType === PetType.DOG && Number(newValue) > 170) {
            return;
        }

        setNewWeight(newValue);
    };

    const onAutocompleteInputChangeHandler = (_, newValue: string) => {
        setAutocompleteInputValue(newValue);
    };

    const onAutocompleteChangeHandler = (_, newValue: Dict | null) => {
        setNewBreed(newValue);
    };

    const onChangeLivingConditionsHandler = (newValue: string) => () => {
        if (newLivingCondition === newValue) {
            return;
        }

        setNewLivingCondition(newValue);
    };

    const onChangerReproductiveStatusHandler = (newValue: string) => () => {
        if (newValue === newReproductiveStatus) {
            setNewReproductiveStatus(undefined);

            return;
        }

        setNewReproductiveStatus(newValue);
    };

    const onBirthDateTypeClickHandler = (newType: BirthDate) => () => {
        if (birthDateType === newType) {
            return;
        }

        setBirthDateType(newType);
        setApproximateDateYear('');
        setApproximateDateMonth('');
    };

    const onChangeExactDateHandler = (newDate: Date | null) => {
        setNewExactDate(newDate);
    };

    const onChangeApproximateDateHandler =
        (view: 'year' | 'months') =>
        ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
            const newValue = value.trim();

            if (newValue === '0') {
                return;
            }

            if (value && !newValue.match(regexInt)) {
                return;
            }

            setNewExactDate(null);

            if (view === 'months' && Number(newValue) > 11) {
                setApproximateDateMonth('11');

                return;
            }

            if (view === 'year' && Number(newValue) > 40) {
                setApproximateDateYear('40');

                return;
            }

            if (view === 'year') {
                setApproximateDateYear(newValue);
            } else {
                setApproximateDateMonth(newValue);
            }
        };

    const onSaveButtonClickHandler = async () => {
        const newData: Partial<Pet> = { id };

        if (name !== newName) {
            newData.name = newName;
        }

        if (type !== newPetType) {
            newData.type = newPetType as PetType;
        }

        if (gender !== newGender) {
            newData.gender = newGender as PetGender;
        }

        if (chipNumber !== newChipNumber && newChipNumber !== 'none') {
            newData.chipNumber = newChipNumber;
        } else if (chipNumber && newChipNumber === 'none') {
            newData.chipNumber = undefined;
        }

        if (weightKg !== Number(newWeight)) {
            newData.weightKg = Number(newWeight);
        }

        if (newBreed?.value !== breedId) {
            newData.breedId = newBreed?.value;
        }

        if (livingCondition !== newLivingCondition) {
            newData.livingCondition = newLivingCondition;
        }

        if (reproductiveStatus !== newReproductiveStatus) {
            newData.reproductiveStatus = newReproductiveStatus;
        }

        if (bloodGroup !== newBloodGroup) {
            newData.bloodGroup = newBloodGroup === 'none' || newBloodGroup === 'Не знаю' ? undefined : newBloodGroup;
        }

        if (birthDateType === BirthDate.EXACT_DATE && newExactDate) {
            if (!birthDate || new Date(birthDate).getTime() !== new Date(newExactDate).getTime()) {
                newData.birthDate = newExactDate;
            }
        } else if (birthDateType === BirthDate.APPROXIMATE_DATE) {
            if (approximateDateYear) {
                newData.ageYears = Number(approximateDateYear);
            }

            if (approximateDateMonth) {
                newData.ageMonths = Number(approximateDateMonth);
            }
        }

        const { success } = await updatePet(newData as Pet);

        if (success) {
            onSuccessUpdate();
        } else {
            onErrorUpdate();
        }

        onClose();
    };

    const renderEditView = () => (
        <>
            <FormItem title='Кличка'>
                <TextField
                    name='name'
                    value={newName}
                    onChange={onChangeNameHandler}
                    placeholder='Как зовут питомца?'
                />
            </FormItem>
            {!!petTypes.length && (
                <FormItem title='Вид'>
                    <div className={styles.buttonsRow}>
                        {petTypes.map(({ label, value }) => (
                            <Button
                                key={value}
                                onClick={onPetTypeChangeHandler(value)}
                                className={cn(styles.buttonsRowItem, { [styles.checked]: newPetType === value })}
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
                            const isChecked = newGender === value;

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
                                    onClick={onPetGenderChangeHandler(value)}
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
                            value={newChipNumber === 'none' ? '' : newChipNumber || ''}
                            inputClass={newChipNumber === 'none' ? styles.input : undefined}
                        />
                    </div>
                    <Button
                        onClick={onNoChipClickHandler}
                        className={cn(styles.buttonsRowItem, { [styles.checked]: newChipNumber === 'none' })}
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
                        <DatePicker
                            onChange={onChangeExactDateHandler}
                            value={newExactDate ? new Date(newExactDate) : newExactDate}
                        />
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
            <FormItem title='Группа крови'>
                <>
                    <div className={cn(styles.buttonsRow, { [styles.dogGroup]: newPetType === PetType.DOG })}>
                        {[...bloodGroupDict[newPetType], { value: 'none', label: 'Не знаю' }].map(
                            ({ label, value }) => (
                                <Button
                                    key={value}
                                    onClick={onChangeBloodGroupHandler(label)}
                                    className={cn(styles.buttonsRowItem, { [styles.checked]: newBloodGroup === label })}
                                >
                                    {label}
                                </Button>
                            ),
                        )}
                    </div>
                    {newPetType === PetType.DOG && (
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
            <FormItem
                title='Вес'
                /* eslint-disable-next-line no-nested-ternary */
                subtitle={!newPetType ? undefined : newPetType === PetType.CAT ? 'до 15 кг' : 'до 170 кг'}
            >
                <TextField
                    name='weight'
                    isDigitInput
                    value={newWeight}
                    onBlur={onBlurPetWeightHandler}
                    onChange={onChangePetWeightHandler}
                    placeholder='Сколько весит питомец?'
                    endAdornment={<div className={styles.endAdornment}>кг</div>}
                />
            </FormItem>
            <FormItem title='Порода'>
                <Autocomplete
                    id='breed'
                    value={newBreed}
                    options={breedsDict[newPetType]}
                    inputValue={autocompleteInputValue}
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
                            color: newBreed && '#8B7069',
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
                            className={cn(styles.buttonsRowItem, {
                                [styles.responsive]: true,
                                [styles.checked]: newLivingCondition === value,
                            })}
                        >
                            {label}
                        </Button>
                    ))}
                </div>
            </FormItem>
            {newGender === PetGender.FEMALE && (
                <FormItem title='Состояние питомца (выберите, если есть)'>
                    <div className={styles.buttonsRow}>
                        {reproductiveStatusesDict.map(({ label, value }) => (
                            <Button
                                key={value}
                                onClick={onChangerReproductiveStatusHandler(value)}
                                className={cn(styles.buttonsRowItem, {
                                    [styles.responsive]: true,
                                    [styles.checked]: newReproductiveStatus === value,
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
                onClick={onSaveButtonClickHandler}
                className={cn(styles.confirm, { [styles.enabled]: isSaveButtonActive })}
            >
                Сохранить изменения
            </Button>
        </>
    );

    const renderView = () => (
        <>
            <ViewString name='Кличка' value={name.toUpperCase()} />
            <ViewString name='Вид' value={petTypes.filter(({ value }) => value === type)[0].label} />
            {gender && <ViewString name='Пол' value={petGenders.filter(({ value }) => value === gender)[0].label} />}
            {chipNumber && <ViewString name='Чип' value={chipNumber} />}
            <ViewString name='Вес' value={`${weightKg} кг`} />
            {birthDate && <ViewString name='Возраст' value={calculateAge(birthDate)} />}
            <ViewString name='Группа крови' value={bloodGroup === 'none' ? 'Не указано' : bloodGroup} />
            {breedId && (
                <ViewString
                    name='Порода'
                    value={(breedsDict[type] || []).filter(({ value }) => value === breedId)[0]?.label}
                />
            )}
            {livingCondition && (
                <ViewString
                    name='Условия содержания'
                    value={livingConditionsDict.filter(({ value }) => value === livingCondition)[0].label}
                />
            )}
            {reproductiveStatus && (
                <ViewString
                    name='Состояние питомца'
                    value={reproductiveStatusesDict.filter(({ value }) => value === reproductiveStatus)[0].label}
                />
            )}
        </>
    );

    useEffect(() => {
        setIsSaveButtonActive(() => {
            if (
                !newName ||
                !newChipNumber ||
                !newWeight ||
                !newBreed ||
                !newLivingCondition ||
                !birthDateType ||
                !newBloodGroup ||
                (birthDateType === BirthDate.APPROXIMATE_DATE && !approximateDateYear && !approximateDateMonth) ||
                (birthDateType === BirthDate.EXACT_DATE && !newExactDate)
            ) {
                return false;
            }

            return (
                name !== newName ||
                type !== newPetType ||
                gender !== newGender ||
                (chipNumber ? chipNumber !== newChipNumber : !!newChipNumber && newChipNumber !== 'none') ||
                bloodGroup !== newBloodGroup ||
                weightKg !== Number(newWeight) ||
                newBreed.value !== breedId ||
                livingCondition !== newLivingCondition ||
                reproductiveStatus !== newReproductiveStatus ||
                birthDate !== newExactDate ||
                (!birthDate && (approximateDateYear || approximateDateMonth))
            );
        });
    }, [
        name,
        type,
        gender,
        newName,
        breedId,
        newBreed,
        weightKg,
        newGender,
        birthDate,
        newWeight,
        newPetType,
        bloodGroup,
        chipNumber,
        newExactDate,
        newBloodGroup,
        newChipNumber,
        birthDateType,
        livingCondition,
        newLivingCondition,
        reproductiveStatus,
        approximateDateYear,
        approximateDateMonth,
        newReproductiveStatus,
    ]);

    return (
        <div className={styles.wrapper}>
            <Header title='Параметры' onClose={onClose} isEditMode={isEditMode} icon={<Params />} />
            {isEditMode ? renderEditView() : renderView()}
        </div>
    );
};

export default ParamsStep;
