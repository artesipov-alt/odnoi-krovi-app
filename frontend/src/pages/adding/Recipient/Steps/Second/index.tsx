import { Button } from '@mui/material';
import { SelectChangeEvent } from '@mui/material/Select';
import cn from 'classnames';
import { BloodAndBreedGroupsDict } from 'hooks/useDicts';
import Lock from 'imgs/svg/lock';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useEffect, useState } from 'react';
import { regexReal } from 'utils/regexps';

import { Dict } from 'api/reference';
import { PetType } from 'api/types';
import Alert from 'components/Alert';
import Multiselect from 'components/Multiselect';
import Switch from 'components/Switch';
import TextField from 'components/TextField';

import styles from './Second.module.less';

type Props = {
    weight: string;
    petType: string;
    bloodGroup: string;
    locations: string[];
    bloodVolume: string;
    locationsDict: Dict[];
    bloodComponents: string[];
    bloodComponentsDict: Dict[];
    desiredBloodGroups: string[];
    notifyOfSmallDonors: boolean;
    bloodGroupDict: BloodAndBreedGroupsDict;
    onConfirmButtonClick: (step: number) => void;
    onChangeBloodVolume: (volume: string) => void;
    onChangeLocations: (locations: string[]) => void;
    onChangeNotifyOfSmallDonors: (isChecked: boolean) => void;
    onChangeDesiredBloodGroups: (bloodGroups: string[]) => void;
    onChangeBloodComponents: (bloodComponents: string[]) => void;
};

const Second: FC<Props> = ({
    weight,
    petType,
    locations,
    bloodGroup,
    bloodVolume,
    locationsDict,
    bloodGroupDict,
    bloodComponents,
    onChangeLocations,
    desiredBloodGroups,
    notifyOfSmallDonors,
    bloodComponentsDict,
    onChangeBloodVolume,
    onConfirmButtonClick,
    onChangeBloodComponents,
    onChangeDesiredBloodGroups,
    onChangeNotifyOfSmallDonors,
}) => {
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(false);

    const onChangeDesiredBloodGroupHandler = (newBloodGroup: string) => () => {
        if (newBloodGroup === bloodGroup) {
            return;
        }

        const index = desiredBloodGroups.findIndex((group) => group === newBloodGroup);

        if (index === -1) {
            onChangeDesiredBloodGroups([...desiredBloodGroups, newBloodGroup]);
        } else {
            const clonedArr = [...desiredBloodGroups];

            clonedArr.splice(index, 1);

            onChangeDesiredBloodGroups(clonedArr);
        }
    };

    const onChangeBloodComponentsHandler = ({ target: { value } }: SelectChangeEvent<typeof bloodComponents>) => {
        const newComponents = typeof value === 'string' ? value.split(',') : value;

        if (newComponents.length > 3) {
            return;
        }

        onChangeBloodComponents(newComponents);
    };

    const onChangeBloodVolumeHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        const newValue = value.replaceAll(' ', '');

        if (!newValue) {
            onChangeBloodVolume('');

            return;
        }

        if (!newValue.match(regexReal)) {
            return;
        }

        if (petType === PetType.CAT && Number(newValue) > Number(weight) * 0.07 * 1000) {
            return;
        }

        if (petType === PetType.DOG && Number(newValue) > Number(weight) * 0.1 * 1000) {
            return;
        }

        onChangeBloodVolume(newValue);
    };

    const onBlurBloodVolumeHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        if (!value) {
            return;
        }

        if (Number(value) < 10) {
            onChangeBloodVolume('10');
        }
    };

    const onChangeLocationsHandler = ({ target: { value } }: SelectChangeEvent<typeof locations>) => {
        const newLocations = typeof value === 'string' ? value.split(',') : value;

        onChangeLocations(newLocations);
    };

    const onChangeSwitchHandler = (_, isChecked) => {
        onChangeNotifyOfSmallDonors(isChecked);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(2);
    };

    useEffect(() => {
        setIsConfirmButtonActive(!!bloodComponents.length && !!bloodVolume && !!locations.length);
    }, [bloodComponents, bloodVolume, locations]);

    return (
        <>
            <FormItem title='Какую группу ищете?'>
                <div className={styles.bloodGroups}>
                    {bloodGroupDict[petType]?.map(({ label, value }) => {
                        const isGroupChecked = desiredBloodGroups.includes(value);

                        return (
                            <div
                                key={value}
                                onClick={onChangeDesiredBloodGroupHandler(value)}
                                className={cn(styles.bloodItem, { [styles.checked]: isGroupChecked })}
                            >
                                <span>{label}</span>
                                {isGroupChecked && value === bloodGroup && (
                                    <div className={styles.lockIcon}>
                                        <Lock />
                                    </div>
                                )}
                            </div>
                        );
                    })}
                </div>
                {((petType === PetType.CAT && desiredBloodGroups.length > 1) ||
                    (petType === PetType.DOG && desiredBloodGroups.length > 1 && `${bloodGroup}` === 'BLG-2')) && (
                    <Alert
                        className={cn(styles.alert, { [styles.isTopMargin]: true })}
                        text='Переливание неподходящей группы крови может быть ОПАСНО! Проконсультируйтесь с врачом!'
                    />
                )}
                {petType === PetType.DOG && desiredBloodGroups.length > 1 && `${bloodGroup}` === 'BLG-1' && (
                    <Alert
                        className={cn(styles.alert, { [styles.isTopMargin]: true })}
                        text='Питомцу подходят обе группы крови.&nbsp;При поиске рекомендуем выбирать родную группу (DEA 1 +), чтобы не создавать дефицит для других собак.'
                    />
                )}
            </FormItem>
            <FormItem title='Какие компоненты нужны?' subtitle='до 3 компонентов'>
                <Multiselect
                    selectValue={bloodComponents}
                    dict={bloodComponentsDict}
                    onChange={onChangeBloodComponentsHandler}
                />
            </FormItem>
            <FormItem
                title='Какой объем требуется?'
                // subtitle={`до ${petType === PetType.CAT ? Big(Number(weight)).times(0.07).times(1000) : Big(Number(weight)).times(0.1).times(1000)} мл`}
                subtitle={`до ${Number((Number(weight) * (petType === PetType.DOG ? 17.6 : 13.2)).toFixed(2))} мл`}
            >
                <TextField
                    name='volume'
                    isDigitInput
                    value={bloodVolume}
                    onBlur={onBlurBloodVolumeHandler}
                    placeholder='Укажите нужный объем'
                    onChange={onChangeBloodVolumeHandler}
                    endAdornment={<div className={styles.endAdornment}>мл</div>}
                />
                <Alert
                    className={cn(styles.alert, { [styles.firstOfFew]: true })}
                    text='Чем меньше объем - тем выше шансы найти кровь'
                />
                <Alert className={styles.alert} text='Могут быть показаны предложения меньшего объема' />
            </FormItem>
            <FormItem title='В каком регионе искать?'>
                <Multiselect dict={locationsDict} selectValue={locations} onChange={onChangeLocationsHandler} />
            </FormItem>
            <div className={styles.formItem}>
                <div className={cn(styles.labelWrapper, { [styles.noMargin]: true })}>
                    <p className={cn(styles.label, { [styles.noMargin]: true })}>Уведомлять о небольших донорах</p>
                    <Switch checked={notifyOfSmallDonors} onChange={onChangeSwitchHandler} />
                </div>
                <p className={styles.donorDescr}>
                    Покажем доноров с меньшим объемом -<br />
                    лучше перелить меньше, чем не перелить совсем.
                </p>
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

export default Second;
