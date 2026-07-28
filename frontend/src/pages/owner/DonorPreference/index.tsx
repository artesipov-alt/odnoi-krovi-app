import { Button } from '@mui/material';
import cn from 'classnames';
import { useLocationsQuery } from 'hooks/useDicts';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Bone from 'imgs/svg/bone';
import Settings from 'imgs/svg/settings';
import Taxi from 'imgs/svg/taxi';
import FormItem from 'pages/adding/common/FormItem';
import { FC, useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';

import { updateUser } from 'api/apiServices/updateUser';
import { queryClient } from 'api/queryClient';
import { Dict } from 'api/reference';
import { CompensationType, DonorPreference as DonorPreferenceType, NotificationFrequency } from 'api/user';
import Layout from 'components/Layout';
import MultiAutocomplete from 'components/MultiAutocomplete';
import Switch from 'components/Switch';

import styles from './DonorPreference.module.less';

export enum View {
    ONBOARDING = 'onboarding',
    PREFERENCE = 'preference',
}

type Props = {
    view: View;
    id?: string;
    onClose: () => void;
    refetchUserData: () => void;
    preference?: DonorPreferenceType;
};

const DonorPreference: FC<Props> = ({ id, view, onClose, refetchUserData, preference }) => {
    const [localView, setLocalView] = useState<View>(view);
    const [isConfirmButtonActive, setIsConfirmButtonActive] = useState<boolean>(false);

    const [locations, setLocations] = useState<Dict[]>([]);
    const [isTaxi, setIsTaxi] = useState<boolean>(preference ? preference.taxiCompensation : false);
    const [recovery, setRecovery] = useState<number>(preference ? preference.recoveryPeriodMonths : 2);
    const [reward, setReward] = useState<CompensationType | null>(preference ? preference.compensationType : null);
    const [isOpenForContact, setIsOpenForContact] = useState<boolean>(preference ? preference.openForContact : true);
    const [notifications, setNotifications] = useState<NotificationFrequency>(
        preference ? (preference.notificationFrequency as NotificationFrequency) : NotificationFrequency.IMMEDIATELY,
    );

    const {
        data: locationsDict = [],
        isError: isErrorLocations,
        isLoading: isLoadingLocations,
        isSuccess: isLocationsSuccess,
    } = useLocationsQuery();

    const showToast = useCallback(
        (text: string) => {
            toast.warn(text, {
                onClose: () => {
                    onClose();
                },
            });
        },
        [onClose],
    );

    const onSettingsClickHandler = () => {
        setLocalView(View.PREFERENCE);
    };

    const onChangeLocationsHandler = (newLocations: Dict[]) => {
        setLocations(newLocations);
    };

    const onRecoveryButtonClickHandler = (value: number) => () => {
        setRecovery(value);
    };

    const onRewardItemClickHandler = (type: CompensationType) => () => {
        setReward(type);
    };

    const onChangeSwitchHandler = (_, isChecked: boolean) => {
        setIsTaxi(isChecked);
    };

    const onNotificationButtonClickHandler = (value: NotificationFrequency) => () => {
        setNotifications(value);
    };

    const onOpenForContactToggle = (value: boolean) => () => {
        setIsOpenForContact((prevState) => {
            if (value === prevState) {
                return prevState;
            }

            return value;
        });
    };

    const onConfirmButtonClickHandler = async () => {
        if (!id) {
            return;
        }

        const { error } = await updateUser({
            id,
            donorPreference: {
                taxiCompensation: isTaxi,
                compensationType: reward!,
                recoveryPeriodMonths: recovery,
                openForContact: isOpenForContact,
                notificationFrequency: notifications,
                preferredLocationIds: locations.map(({ value }) => value),
            },
        });

        if (error) {
            showToast('Не удалось сохранить настройки, попробуйте еще раз');

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', id] });

        refetchUserData();
    };

    useEffect(() => {
        if (isErrorLocations) {
            showToast('Не удалось загрузить словарь регионов, попробуйте eще раз');
        }
    }, [isErrorLocations, showToast]);

    useEffect(() => {
        if (isLocationsSuccess) {
            setLocations(locationsDict.filter((location) => preference?.preferredLocationIds.includes(location.value)));
        }
    }, [isLocationsSuccess, locationsDict, preference?.preferredLocationIds]);

    useEffect(() => {
        setIsConfirmButtonActive(
            preference
                ? preference.taxiCompensation !== isTaxi ||
                      preference.compensationType !== reward ||
                      preference.recoveryPeriodMonths !== recovery ||
                      preference.openForContact !== isOpenForContact ||
                      preference.notificationFrequency !== notifications ||
                      (!!locations.length &&
                          preference.preferredLocationIds.join(',') !== locations.map(({ value }) => value).join(','))
                : !!locations.length && !!reward,
        );
    }, [isTaxi, locations, notifications, preference, recovery, reward, isOpenForContact]);

    return (
        <Layout>
            {localView === View.ONBOARDING && (
                <div className={styles.onboarding}>
                    <div className={styles.container}>
                        <h1 className={styles.title}>
                            Помогайте
                            <br />
                            спасать жизни!
                        </h1>
                        <p className={styles.subtitle}>Но сначала задайте параметры донорства</p>
                        <div className={styles.buttons}>
                            <Button
                                variant='contained'
                                onClick={onSettingsClickHandler}
                                className={cn(styles.button, styles.settings)}
                            >
                                <div className={styles.settingsIcon}>
                                    <Settings />
                                </div>
                                Настроить
                            </Button>
                            <Button className={cn(styles.button, styles.close)} onClick={onClose} variant='contained'>
                                Не сейчас
                            </Button>
                        </div>
                    </div>
                </div>
            )}
            {localView === View.PREFERENCE && (
                <div className={styles.preference}>
                    <div className={styles.header}>
                        <div className={styles.back} onClick={onClose}>
                            <BackAngularArrow />
                        </div>
                        <h2 className={styles.title}>Параметры донорства</h2>
                    </div>
                    {!isLoadingLocations && (
                        <FormItem title='Где хотите помогать?'>
                            <MultiAutocomplete
                                value={locations}
                                dict={locationsDict}
                                onChange={onChangeLocationsHandler}
                            />
                        </FormItem>
                    )}
                    <FormItem
                        title='Период восстановления'
                        className={styles.recoveryItem}
                        tooltip='По стандартам Национальной Ветеринарной Палаты отдых после донации составляет не менее 2-х месяцев, но Вы можете указать и больший период'
                    >
                        <div className={styles.recovery}>
                            {[2, 3, 4, 5, 6].map((item) => (
                                <Button
                                    key={item}
                                    onClick={onRecoveryButtonClickHandler(item)}
                                    className={cn(styles.button, { [styles.checked]: item === recovery })}
                                >
                                    {item} мес
                                </Button>
                            ))}
                        </div>
                    </FormItem>
                    <FormItem title='Вознаграждение от реципиента' tooltip='Реципиент - питомец, нуждающийся в помощи'>
                        <div className={styles.reward}>
                            <div
                                onClick={onRewardItemClickHandler(CompensationType.FREE)}
                                className={cn(styles.rewardItem, {
                                    [styles.checked]: reward === CompensationType.FREE,
                                    [styles.disabled]: !!reward && reward !== CompensationType.FREE,
                                })}
                            >
                                <div className={styles.rewardFreeIcon}>
                                    <p className={styles.sum}>0</p>
                                    <p className={styles.descr}>₽</p>
                                </div>
                                <p className={styles.rewardItemText}>
                                    Готов помочь
                                    <br />
                                    безвозмездно
                                </p>
                            </div>
                            <div
                                onClick={onRewardItemClickHandler(CompensationType.PAID)}
                                className={cn(styles.rewardItem, {
                                    [styles.checked]: reward === CompensationType.PAID,
                                    [styles.disabled]: !!reward && reward !== CompensationType.PAID,
                                })}
                            >
                                <div className={styles.rewardNotFreeIcon}>
                                    <p className={styles.descr}>₽</p>
                                </div>
                                <p className={styles.rewardItemText}>
                                    Не готов помочь
                                    <br />
                                    бесплатно
                                </p>
                            </div>
                            <div
                                onClick={onRewardItemClickHandler(CompensationType.FOOD)}
                                className={cn(styles.rewardItem, {
                                    [styles.checked]: reward === CompensationType.FOOD,
                                    [styles.disabled]: !!reward && reward !== CompensationType.FOOD,
                                })}
                            >
                                <div className={styles.rewardFeedIcon}>
                                    <div className={styles.icon}>
                                        <Bone />
                                    </div>
                                </div>
                                <p className={styles.rewardItemText}>
                                    Готов помочь за
                                    <br />
                                    корм
                                </p>
                            </div>
                        </div>
                    </FormItem>
                    <div className={styles.taxi}>
                        <div className={styles.taxiTitle}>
                            <div className={styles.taxiIcon}>
                                <div className={styles.icon}>
                                    <Taxi />
                                </div>
                            </div>
                            <p className={styles.taxiDescr}>Нужно компенсировать такси до клиники и обратно</p>
                        </div>
                        <Switch checked={isTaxi} onChange={onChangeSwitchHandler} />
                    </div>
                    <FormItem title='Можно связаться с Вами, если нужна помощь?' className={styles.openForContact}>
                        <div className={styles.buttons}>
                            <Button
                                onClick={onOpenForContactToggle(true)}
                                className={cn(styles.button, {
                                    [styles.checked]: isOpenForContact,
                                })}
                            >
                                Да
                            </Button>
                            <Button
                                onClick={onOpenForContactToggle(false)}
                                className={cn(styles.button, {
                                    [styles.checked]: !isOpenForContact,
                                })}
                            >
                                Выберу реципиента самостоятельно
                            </Button>
                        </div>
                    </FormItem>
                    <FormItem title='Как часто уведомлять о новых реципиентах?' className={styles.notifications}>
                        <div className={styles.notificationsFirst}>
                            <Button
                                onClick={onNotificationButtonClickHandler(NotificationFrequency.IMMEDIATELY)}
                                className={cn(styles.button, {
                                    [styles.checked]: notifications === NotificationFrequency.IMMEDIATELY,
                                })}
                            >
                                Сразу
                            </Button>
                            <Button
                                onClick={onNotificationButtonClickHandler(NotificationFrequency.DAILY)}
                                className={cn(styles.button, {
                                    [styles.checked]: notifications === NotificationFrequency.DAILY,
                                })}
                            >
                                Раз в день
                            </Button>
                            <Button
                                onClick={onNotificationButtonClickHandler(NotificationFrequency.WEEKLY)}
                                className={cn(styles.button, {
                                    [styles.checked]: notifications === NotificationFrequency.WEEKLY,
                                })}
                            >
                                Раз в неделю
                            </Button>
                        </div>
                        <Button
                            onClick={onNotificationButtonClickHandler(NotificationFrequency.NEVER)}
                            className={cn(styles.button, {
                                [styles.fullWidth]: true,
                                [styles.checked]: notifications === NotificationFrequency.NEVER,
                            })}
                        >
                            Никогда
                        </Button>
                    </FormItem>
                    <Button
                        fullWidth
                        onClick={onConfirmButtonClickHandler}
                        className={cn(styles.confirm, { [styles.enabled]: isConfirmButtonActive })}
                    >
                        Применить
                    </Button>
                </div>
            )}
        </Layout>
    );
};

export default DonorPreference;
