import { Button } from '@mui/material';
import cn from 'classnames';
import { BloodAndBreedGroupsDict } from 'hooks/useDicts';
import Lock from 'imgs/svg/lock';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useState } from 'react';
import { regexReal } from 'utils/regexps';

import { PetTypeDict } from 'api/reference';
import { PetType } from 'api/types';
import Alert from 'components/Alert';
import ImgEditor from 'components/ImgEditor';
import TextField from 'components/TextField';

import styles from './First.module.less';

type Props = {
    name: string;
    weight: string;
    petType: string;
    bloodGroup: string;
    photo: File | null;
    petTypes: PetTypeDict[];
    onChangeName: (name: string) => void;
    onChangePetType: (type: string) => void;
    bloodGroupDict: BloodAndBreedGroupsDict;
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
    bloodGroup,
    onLoadPhoto,
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
            setIsConfirmButtonActive(false);

            onChangeWeight('');

            return;
        }

        if (!newValue.match(regexReal)) {
            setIsConfirmButtonActive(false);

            return;
        }

        if (petType === PetType.CAT && Number(newValue) > 15) {
            setIsConfirmButtonActive(!!name && !!petType && !!bloodGroup && !!weight);

            return;
        }

        if (petType === PetType.DOG && Number(newValue) > 170) {
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
            <div className={styles.photo}>
                <ImgEditor src={photo} onLoad={onLoadPhoto} />
            </div>
            <Alert
                className={styles.alert}
                text='Фото может вызвать эмоциональный отклик у хозяев доноров и увеличить шансы найти помощь'
            />
            <FormItem title='Кличка'>
                <TextField name='name' value={name} placeholder='Как зовут питомца?' onChange={onChangeNameHandler} />
            </FormItem>
            {!!petTypes.length && (
                <FormItem title='Вид'>
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
                </FormItem>
            )}
            {/* eslint-disable-next-line no-nested-ternary */}
            <FormItem title='Вес' subtitle={!petType ? undefined : petType === PetType.CAT ? 'до 15 кг' : 'до 170 кг'}>
                <TextField
                    name='weight'
                    isDigitInput
                    value={weight}
                    disabled={!petType}
                    onBlur={onBlurPetWeightHandler}
                    onChange={onChangePetWeightHandler}
                    inputClass={cn({ [styles.lock]: !petType })}
                    htmlInputClass={cn({ [styles.lock]: !petType })}
                    placeholder={petType ? 'Сколько весит питомец?' : 'Сначала выберите вид'}
                    endAdornment={
                        !petType ? (
                            <div className={styles.lockIcon}>
                                <Lock />
                            </div>
                        ) : (
                            <div className={styles.endAdornment}>кг</div>
                        )
                    }
                />
            </FormItem>
            <FormItem title='Группа крови питомца'>
                {petType ? (
                    <>
                        <div className={cn(styles.bloodGroups, { [styles.dogGroup]: petType === PetType.DOG })}>
                            {bloodGroupDict[petType]
                                .filter(({ value }) => value !== 'UNKNOWN')
                                .map(({ label, value }) => (
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
                        name='bloodGroup'
                        value={bloodGroup}
                        placeholder='Сначала выберите вид'
                        inputClass={cn({ [styles.lock]: !petType })}
                        htmlInputClass={cn({ [styles.lock]: !petType })}
                        endAdornment={
                            !petType ? (
                                <div className={styles.lockIcon}>
                                    <Lock />
                                </div>
                            ) : (
                                <div className={styles.endAdornment}>кг</div>
                            )
                        }
                    />
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
