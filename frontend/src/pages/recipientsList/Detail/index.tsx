import { Button } from '@mui/material';
import cn from 'classnames';
import { useLocationsQuery } from 'hooks/useDicts';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import Bone from 'imgs/svg/bone';
import Certificates from 'imgs/svg/certificates';
import CrossedEye from 'imgs/svg/crossedEye';
import Eye from 'imgs/svg/eye';
import Location from 'imgs/svg/location';
import Pin from 'imgs/svg/pin';
import PrioritySearch from 'imgs/svg/prioritySearch';
import Taxi from 'imgs/svg/taxi';
import Accordion from 'pages/adding/common/Accordion';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { bloodSearchApply } from 'api/apiServices/bloodSearchApply';
import { getRecipientDetail } from 'api/apiServices/getRecipientDetail';
import { GetRecipientDetailsResponse } from 'api/donor';
import { PetType } from 'api/types';
import { CompensationType } from 'api/user';
import { CircularProgress } from 'components/CircularProgress';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';
import Loading from 'components/Loading';
import Switch from 'components/Switch';

import styles from './RecipientsListDetail.module.less';

type Props = {
    id: string;
    onClose: () => void;
    isBlurByDefault: boolean;
};

const RecipientsListDetail: FC<Props> = ({ id, isBlurByDefault, onClose }) => {
    const navigate = useNavigate();

    const [isBlur, setIsBlur] = useState(isBlurByDefault);
    const [isLoading, setLoading] = useState<boolean>(true);
    const [isTaxi, setIsTaxi] = useState<boolean | null>(null);
    const [isConfirmed, setIsConfirmed] = useState<boolean>(false);
    const [isConditionsOpen, setIsConditionsOpen] = useState(false);
    const [isConfirmCurtainOpen, setIsConfirmCurtainOpen] = useState(false);
    const [reward, setReward] = useState<CompensationType | null>(null);
    const [checkedDonor, setCheckedDonor] = useState<string | null>(null);
    const [recipient, setRecipient] = useState<GetRecipientDetailsResponse | null>(null);

    const { data: locationsDict = [] } = useLocationsQuery();

    const showToast = useCallback(
        (text: string, goToOwner = true) => {
            toast.warn(text, {
                onClose: () => {
                    if (goToOwner) {
                        navigate('/owner');
                    }
                },
            });
        },
        [navigate],
    );

    const fetchDetails = useCallback(async () => {
        const response = await getRecipientDetail(id);

        if (!response) {
            showToast('Ошибка при загрузке данных, попробуйте еще раз');

            return;
        }

        setRecipient(response);
        setLoading(false);
    }, [id, showToast]);

    const isBlurToggle = () => {
        setIsBlur((prevState) => !prevState);
    };

    const onConditionsToggle = () => {
        setIsConditionsOpen((prevState) => !prevState);
    };

    const onRewardItemClickHandler = (type: CompensationType) => () => {
        setReward(type);
    };

    const onChangeSwitchHandler = (_, isChecked: boolean) => {
        setIsTaxi(isChecked);
    };

    const onCheckItemClickHandler = (petId: string) => () => {
        setCheckedDonor((prevState) => {
            if (prevState === petId) {
                return null;
            }

            return petId;
        });
    };

    const onDonorClickHandler = async () => {
        if (!checkedDonor) {
            return;
        }

        setIsConfirmed(true);

        const response = await bloodSearchApply({
            id: checkedDonor,
            compensationType: (reward || recipient?.defaultPrefs?.compensationType)!,
            taxiCompensation: (isTaxi !== null ? isTaxi : recipient?.defaultPrefs?.taxiCompensation)!,
        });

        if (!response) {
            showToast('Не удалось предложить донора, попробуйте еще раз', false);

            setIsConfirmed(false);

            return;
        }

        setIsConfirmCurtainOpen(true);
    };

    const onConfirmCurtainClickHandler = () => {
        navigate('/owner#donor');
    };

    const onCancelCurtainClickHandler = () => {
        navigate('/owner#donorDonations');
    };

    useEffect(() => {
        fetchDetails();
    }, [fetchDetails]);

    if (isLoading) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    if (!recipient) {
        return null;
    }

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.back} onClick={onClose}>
                    <BackAngularArrow />
                </div>
                <h2 className={styles.name}>{recipient.petName.toUpperCase()}</h2>
                {recipient.photoUrls?.[0] && (
                    <div onClick={isBlurToggle} className={styles.blurIconWrapper}>
                        <div className={cn(styles.blurIcon, { [styles.isBlur]: isBlur })}>
                            {isBlur ? <Eye /> : <CrossedEye />}
                        </div>
                    </div>
                )}
            </div>
            <div className={styles.main}>
                <div className={styles.left}>
                    <div className={cn(styles.leftItem, { [styles.blood]: true })}>
                        <div className={styles.subTitle}>
                            <div className={styles.icon}>
                                <Blood />
                            </div>
                            <p className={styles.text}>Ищет</p>
                        </div>
                        <div className={styles.bloodInfo}>
                            <div className={styles.bloodGroup}>{recipient.bloodGroupName}</div>
                            {recipient.searchingBloodNames.some((group) => group !== recipient.bloodGroupName) && ' +'}
                            {recipient.searchingBloodNames.some((group) => group !== recipient.bloodGroupName)
                                ? recipient.searchingBloodNames
                                      .filter((group) => group !== recipient.bloodGroupName)
                                      .map((group) => (
                                          <div key={group} className={cn(styles.bloodGroup, { [styles.needed]: true })}>
                                              {group}
                                          </div>
                                      ))
                                : ''}
                        </div>
                    </div>
                    <div className={cn(styles.leftItem, { [styles.location]: true })}>
                        <div className={styles.icon}>
                            <Location />
                        </div>
                        <p className={styles.text}>
                            {recipient.regions
                                ?.map((lock) => locationsDict.filter(({ value }) => value === lock)[0]?.label)
                                .join(', ')}
                        </p>
                    </div>
                    <div className={cn(styles.leftItem, { [styles.owner]: true })}>
                        <div className={styles.icon}>{recipient.ownerName.charAt(0).toUpperCase()}</div>
                        <div>
                            <p className={styles.descr}>Хозяин</p>
                            <p className={styles.text}>{recipient.ownerName}</p>
                        </div>
                    </div>
                </div>
                <div className={styles.right}>
                    <div className={styles.avatarWrapper}>
                        <img
                            alt={recipient.petName}
                            className={cn(styles.avatar, { [styles.blured]: recipient.photoUrls?.[0] && isBlur })}
                            src={
                                recipient.photoUrls?.[0] ||
                                (recipient.petType === PetType.DOG ? dogRoundStub : catRoundStub)
                            }
                        />
                        <CircularProgress
                            size={180}
                            strokeWidth={15}
                            color='var(--red10, #FF2727)'
                            total={recipient.bloodVolumeNeeded}
                            current={recipient.bloodVolumeReserved}
                        />
                        <div className={styles.neededVolume}>
                            {recipient.bloodVolumeNeeded}
                            <span>мл</span>
                        </div>
                    </div>
                </div>
            </div>
            {(!!recipient.advancedInfo?.description?.length || !!recipient.advancedInfo?.photoUrls?.[0]) && (
                <Accordion icon={<Pin />} title='Дополнительная информация'>
                    <div className={cn(styles.itemText, { [styles.accordionText]: true })}>
                        {recipient.advancedInfo?.description}
                    </div>
                    {recipient.advancedInfo?.photoUrls?.[0] && (
                        <img
                            alt='Фото рецепиента'
                            className={styles.bloodRequestPhoto}
                            src={recipient.advancedInfo.photoUrls[0]}
                        />
                    )}
                </Accordion>
            )}
            <div className={styles.settings}>
                <div className={styles.setting}>
                    <p className={styles.text}>
                        Бонусы
                        <br />
                        Портала
                    </p>
                    <div className={styles.icons}>
                        <div className={styles.bonusIcon}>
                            <PrioritySearch />
                        </div>
                        <div className={styles.bonusIcon}>
                            <Certificates />
                        </div>
                    </div>
                </div>
                <div className={styles.setting} onClick={onConditionsToggle}>
                    <p className={styles.text}>
                        Ваши
                        <br />
                        условия
                    </p>
                    <div className={styles.icons}>
                        <div className={cn(styles.rewardFeedIcon, { [styles.button]: true })}>
                            <div className={styles.icon}>
                                <Bone />
                            </div>
                        </div>
                        <div className={cn(styles.taxiIcon, { [styles.button]: true })}>
                            <div className={styles.icon}>
                                <Taxi />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            {!!recipient.matchingDonors?.length && (
                <div className={styles.donors}>
                    <p className={styles.donorDescr}>Выберите донора</p>
                    <div className={styles.showcase}>
                        {recipient.matchingDonors.map(({ amount, donorBloodGroup, petId, petName, photoUrls }) => (
                            <div key={`${petId}`} className={styles.pet}>
                                <div className={styles.photo}>
                                    <img alt={petName} src={photoUrls[0]} className={styles.img} />
                                    <div className={styles.donorBloodInfo}>
                                        <div className={styles.donorBloodGroup}>{donorBloodGroup}</div>
                                        <div className={styles.bloodVolume}>
                                            <p className={styles.bloodVolumeNumber}>{amount}</p>
                                            <p className={styles.bloodVolumeDescr}>мл</p>
                                        </div>
                                    </div>
                                    <div
                                        onClick={onCheckItemClickHandler(petId)}
                                        className={cn(styles.checkItem, { [styles.checked]: checkedDonor === petId })}
                                    />
                                    <div className={styles.photoFooter}>
                                        <p className={styles.domorName}>{petName.toUpperCase()}</p>
                                    </div>
                                    <div className={styles.gradient} />
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
            <Button
                fullWidth
                onClick={onDonorClickHandler}
                className={cn(styles.confirmButton, { [styles.enabled]: !!checkedDonor && !isConfirmed })}
            >
                Предложить донора
            </Button>
            {isConditionsOpen && (
                <Curtain noRednerButtons title='Ваши условия' shouldCloseByWrapperClick onClose={onConditionsToggle}>
                    <h3 className={styles.rewardTitle}>Вознаграждение от реципиента</h3>
                    <div className={styles.reward}>
                        <div
                            onClick={onRewardItemClickHandler(CompensationType.FREE)}
                            className={cn(styles.rewardItem, {
                                [styles.checked]: reward
                                    ? reward === CompensationType.FREE
                                    : recipient.defaultPrefs?.compensationType === CompensationType.FREE,
                                [styles.disabled]: reward
                                    ? reward !== CompensationType.FREE
                                    : recipient.defaultPrefs?.compensationType !== CompensationType.FREE,
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
                                [styles.checked]: reward
                                    ? reward === CompensationType.PAID
                                    : recipient.defaultPrefs?.compensationType === CompensationType.PAID,
                                [styles.disabled]: reward
                                    ? reward !== CompensationType.PAID
                                    : recipient.defaultPrefs?.compensationType !== CompensationType.PAID,
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
                                [styles.checked]: reward
                                    ? reward === CompensationType.FOOD
                                    : recipient.defaultPrefs?.compensationType === CompensationType.FOOD,
                                [styles.disabled]: reward
                                    ? reward !== CompensationType.FOOD
                                    : recipient.defaultPrefs?.compensationType !== CompensationType.FOOD,
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
                    <div className={styles.taxi}>
                        <div className={styles.taxiTitle}>
                            <div className={styles.taxiIcon}>
                                <div className={styles.icon}>
                                    <Taxi />
                                </div>
                            </div>
                            <p className={styles.taxiDescr}>Нужно компенсировать такси до клиники и обратно</p>
                        </div>
                        <Switch
                            checked={isTaxi === null ? recipient.defaultPrefs?.taxiCompensation : isTaxi}
                            onChange={onChangeSwitchHandler}
                        />
                    </div>
                    <Button
                        fullWidth
                        onClick={onConditionsToggle}
                        className={cn(styles.confirmButton, { [styles.enabled]: !!reward || isTaxi !== null })}
                    >
                        Применить
                    </Button>
                </Curtain>
            )}
            {isConfirmCurtainOpen && (
                <Curtain
                    columnOfButtons
                    cancelButtonTitle='Планируемые донации'
                    title={
                        <>
                            Мы уведомили хозяина
                            <br />о Вашем предложении
                        </>
                    }
                    confirmButtonTitle='К моим питомцам'
                    onConfirm={onConfirmCurtainClickHandler}
                    onCancel={onCancelCurtainClickHandler}
                    subTitle={
                        <div className={styles.curtainSubtitle}>
                            Если в течение часа хозяин реципиента не
                            <br />
                            ответит, Вы сможете выбрать другого
                            <br />
                            реципиента или клинику
                        </div>
                    }
                />
            )}
        </Layout>
    );
};

export default RecipientsListDetail;
