import { Button } from '@mui/material';
import TextField from '@mui/material/TextField';
import cn from 'classnames';
import Lock from 'imgs/svg/lock';
import { ChangeEvent, FC, useState } from 'react';
import { regexReal } from 'utils/regexps';

import { Dict, PetDict } from 'api/reference';
import { PetType } from 'api/types';
import Alert from 'components/Alert';
import ImgEditor from 'components/ImgEditor';

import styles from './First.module.less';

type Props = {
    name: string;
    weight: string;
    petType: string;
    bloodGroup: string;
    photo: File | null;
    petTypes: PetDict[];
    onChangeName: (name: string) => void;
    onChangePetType: (type: string) => void;
    bloodGroupDict: Record<PetType, Dict[]>;
    onChangeWeight: (weight: string) => void;
    onLoadPhoto: (photo: File | null) => void;
    onChangeBloodGroup: (blood: string) => void;
    onConfirmButtonClick: (step: number) => void;
};

const First: FC<Props> = ({
    name,
    photo,
    weight,
    petType,
    petTypes,
    onLoadPhoto,
    bloodGroup,
    onChangeName,
    bloodGroupDict,
    onChangeWeight,
    onChangePetType,
    onChangeBloodGroup,
    onConfirmButtonClick,
}) => {
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(
        !!name && !!petType && !!weight && !!bloodGroup,
    );

    const onChangeNameHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setIsConfirmButtonActive(!!value && !!petType && !!weight && !!bloodGroup);

        onChangeName(value);
    };

    const onChangePetWeightHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        const newValue = value.replaceAll(' ', '');

        if (!newValue) {
            setIsConfirmButtonActive(false);

            onChangeWeight('');

            return;
        }

        if (!newValue.match(regexReal)) {
            setIsConfirmButtonActive(false);

            return;
        }

        if (petType === PetType.CAT && Number(newValue) > 20) {
            setIsConfirmButtonActive(!!name && !!petType && !!bloodGroup && !!weight);

            return;
        }

        if (petType === PetType.DOG && Number(newValue) > 150) {
            setIsConfirmButtonActive(!!name && !!petType && !!bloodGroup && !!weight);

            return;
        }

        setIsConfirmButtonActive(!!name && !!petType && !!bloodGroup);

        onChangeWeight(newValue);
    };

    const onPetTypeButtonClickHandler = (newPetType: string) => () => {
        setIsConfirmButtonActive(false);

        onChangePetType(newPetType);
        onChangeWeight('');
        onChangeBloodGroup('');
    };

    const onChangeBloodGroupHandler = (newBloodGroup: string) => () => {
        if (newBloodGroup === bloodGroup) {
            return;
        }

        setIsConfirmButtonActive(!!name && !!petType && !!weight && !!newBloodGroup);

        onChangeBloodGroup(newBloodGroup);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(1);
    };

    return (
        <>
            <ImgEditor src={photo} onLoad={onLoadPhoto} className={styles.photo} />
            <Alert
                className={styles.alert}
                text='Фото может вызвать эмоциональный отклик у хозяев доноров и увеличить шансы найти помощь'
            />
            <div className={styles.formItem}>
                <p className={styles.label}>Кличка</p>
                <TextField
                    fullWidth
                    name='name'
                    value={name}
                    onChange={onChangeNameHandler}
                    placeholder='Как зовут питомца?'
                    slotProps={{
                        htmlInput: { className: styles.input },
                        input: { className: styles.inputWrapper },
                    }}
                    sx={{
                        '& .MuiOutlinedInput-root': {
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
            </div>
            {!!petTypes.length && (
                <div className={styles.formItem}>
                    <p className={styles.label}>Вид</p>
                    <div className={styles.petType}>
                        {petTypes.map(({ label, value }) => (
                            <Button
                                key={value}
                                onClick={onPetTypeButtonClickHandler(value)}
                                className={cn(styles.petTypeButton, { [styles.checked]: petType === value })}
                            >
                                {label}
                            </Button>
                        ))}
                    </div>
                </div>
            )}
            <div className={styles.formItem}>
                <div className={styles.labelWrapper}>
                    <p className={styles.label}>Вес</p>
                    <span className={styles.subLabel}>
                        {/* eslint-disable-next-line no-nested-ternary */}
                        {!petType ? null : petType === PetType.CAT ? 'до 20 кг' : 'до 150 кг'}
                    </span>
                </div>
                <TextField
                    fullWidth
                    name='weight'
                    value={weight}
                    disabled={!petType}
                    onChange={onChangePetWeightHandler}
                    placeholder={petType ? 'Сколько весит питомец?' : 'Сначала выберите вид'}
                    slotProps={{
                        input: {
                            endAdornment: !petType ? (
                                <div className={styles.lockIcon}>
                                    <Lock />
                                </div>
                            ) : (
                                <div className={styles.endAdornment}>кг</div>
                            ),
                            className: cn(styles.inputWrapper, { [styles.lock]: !petType }),
                        },
                        htmlInput: {
                            className: cn(styles.input, { [styles.lock]: !petType }),
                        },
                    }}
                    sx={{
                        '& .MuiOutlinedInput-root': {
                            '& .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#dee2e9',
                            },
                            '&.Mui-disabled .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#dee2e9',
                            },
                            '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                borderColor: '#dee2e9',
                                borderWidth: '1px',
                            },
                        },
                    }}
                />
            </div>
            <div className={styles.formItem}>
                <p className={styles.label}>Группа крови питомца</p>
                {petType ? (
                    <>
                        <div className={cn(styles.bloodGroups, { [styles.dogGroup]: petType === PetType.DOG })}>
                            {bloodGroupDict[petType].map(({ label, value }) => (
                                <div
                                    key={value}
                                    onClick={onChangeBloodGroupHandler(value)}
                                    className={cn(styles.bloodItem, { [styles.checked]: bloodGroup === value })}
                                >
                                    {label}
                                </div>
                            ))}
                        </div>
                        {petType === PetType.DOG && (
                            <Alert
                                className={styles.alert}
                                text='Сведения вносятся по системе DEA.&nbsp;Не используйте данные из других систем (KAI, DAL)'
                            />
                        )}
                    </>
                ) : (
                    <TextField
                        disabled
                        fullWidth
                        name='bloodGroup'
                        value={bloodGroup}
                        placeholder='Сначала выберите вид'
                        slotProps={{
                            input: {
                                endAdornment: !petType ? (
                                    <div className={styles.lockIcon}>
                                        <Lock />
                                    </div>
                                ) : (
                                    <div className={styles.endAdornment}>кг</div>
                                ),
                                className: cn(styles.inputWrapper, { [styles.lock]: !petType }),
                            },
                            htmlInput: {
                                className: cn(styles.input, { [styles.lock]: !petType }),
                            },
                        }}
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                '& .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#dee2e9',
                                },
                                '&.Mui-disabled .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#dee2e9',
                                },
                            },
                        }}
                    />
                )}
            </div>
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
