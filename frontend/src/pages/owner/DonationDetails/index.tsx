import { Button } from '@mui/material';
import cn from 'classnames';
import { useLocationsQuery, usePetTypesAndBloodGroupsQuery } from 'hooks/useDicts';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import BloodVolume from 'imgs/svg/bloodVolume';
import Bone from 'imgs/svg/bone';
import Cancel from 'imgs/svg/cancel';
import Chat from 'imgs/svg/chat';
import Exclamation from 'imgs/svg/exclamation';
import Info from 'imgs/svg/info';
import Location from 'imgs/svg/location';
import Lock from 'imgs/svg/lock';
import Max from 'imgs/svg/max';
import MiniSinglePaw from 'imgs/svg/miniSinglePaw';
import Phone from 'imgs/svg/phone';
import Pin from 'imgs/svg/pin';
import PrioritySearch from 'imgs/svg/prioritySearch';
import Taxi from 'imgs/svg/taxi';
import Telegram from 'imgs/svg/telegram';
import Accordion from 'pages/adding/common/Accordion';
import { ChangeEvent, FC, MouseEvent, useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { regexReal } from 'utils/regexps';
import { isWithinHours, matchIdentities } from 'utils/utils';

import { cancelDonation } from 'api/apiServices/cancelDonation';
import { completeDonation } from 'api/apiServices/completeDonation';
import { confirmDonation } from 'api/apiServices/confirmDonation';
import { getUserContacts } from 'api/apiServices/getUserContacts';
import { updatePet } from 'api/apiServices/updatePet';
import { BonusType, DonorStatus, PlannedDonation } from 'api/donor';
import { Pet } from 'api/pets';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import { CompensationType, Identities } from 'api/user';
import Alert from 'components/Alert';
import Bonuses from 'components/Bonuses';
import { CircularProgress } from 'components/CircularProgress';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';
import RejectedForm, { RejectView } from 'components/RejectedForm';
import TextField from 'components/TextField';
import Timer from 'components/Timer';

import styles from './DonationDetails.module.less';

type ChatCurtain = {
    isOpen: boolean;
    identities?: Identities[];
};

type RejectedFormType = {
    isOpen: boolean;
    view?: RejectView;
};

type Props = {
    userId: string;
    onClose: () => void;
    identities?: Identities[];
    donation: PlannedDonation;
};

const curtainList = [
    'Не переходите по подозрительным ссылкам',
    'Не передавайте свои паспортные данные',
    'Вы получите бонусы только после подтверждения хозяином реципиента донации на Портале',
];

const DonationDetails: FC<Props> = ({ userId, onClose, donation, identities }) => {
    const [isConditionsOpen, setIsConditionsOpen] = useState(false);
    const [chatCurtain, setChatCurtain] = useState<ChatCurtain>({ isOpen: false });
    const [donatedBloodVolume, setDonatedBloodVolume] = useState<string>('');
    const [isBonusesPageOpen, setIsBonusesPageOpen] = useState<boolean>(false);
    const [isTooltipRegionOpen, setIsTooltipRegionOpen] = useState<boolean>(false);
    const [isPendingTimerExpired, setIsPendingTimerExpired] = useState<boolean>(false);
    const [isDonorConfirmationCurtainOpen, setIsDonorConfirmationCurtainOpen] = useState(false);
    const [rejectDonationFormParams, setRejectDonationFormParams] = useState<RejectedFormType>({ isOpen: false });

    const [bloodGroup, setBloodGroup] = useState<string | null>(null);

    const { data: locationsDict = [], isError: isErrorLocations } = useLocationsQuery();
    const { data: { bloodGroupDict = {} } = {}, isError: isErrorPetTypesAndBloodGroups } =
        usePetTypesAndBloodGroupsQuery();

    const isDonorUnknownBloodGroup = donation.applicationData.donorBloodGroup === 'UNKNOWN';

    const showToast = useCallback((text: string, type = 'warn') => {
        if (type === 'success') {
            toast.success(text);

            return;
        }

        toast.warn(text);
    }, []);

    const onCloseChatCurtainClickHandler = () => {
        setChatCurtain({ isOpen: false });
    };

    const onCallClickHandler = (e) => {
        const storedEnv = localStorage.getItem('environment');

        if (storedEnv === 'tg') {
            e.preventDefault();
            const phone = `tel:+${donation.recipientData.ownerPhone.replace(/[^\d]/g, '')}`;
            const a = document.createElement('a');
            a.href = phone;
            a.target = '_blank';
            a.rel = 'noopener noreferrer';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } else {
            window.location.href = `tel:${donation.recipientData.ownerPhone}`;
        }
    };

    const onMessengerClickHandler = (providerName: string) => async () => {
        const response = await getUserContacts({ id: donation.recipientData.ownerID, provider: providerName });

        if (!response) {
            showToast('Не удалось получить контакт хозяина донора');

            setChatCurtain({ isOpen: false });

            return;
        }

        showToast(response.data.message, 'success');

        setChatCurtain({ isOpen: false });
    };

    const onConfirmDonationToggle = () => {
        if (!donation) {
            return;
        }

        setDonatedBloodVolume(
            donation.applicationData.amount > donation.recipientData.bloodVolumeNeeded
                ? `${donation.recipientData.bloodVolumeNeeded}`.replace('.', ',')
                : `${donation.applicationData.amount}`.replace('.', ','),
        );
        setIsDonorConfirmationCurtainOpen((prevState) => !prevState);
        setBloodGroup(null);
    };

    const onChatOpenHandler = () => {
        setChatCurtain({
            isOpen: true,
            identities: matchIdentities(identities!, donation.recipientData.identities),
        });
    };

    const onConfirmDonationClickHandler = async () => {
        if (bloodGroup) {
            const { success } = await updatePet({ id: donation.applicationData.donorPetID, bloodGroup } as Pet);

            if (!success) {
                showToast('Не удалось сохранить выбранную группу крови, для донора');
            }
        }

        const response = await completeDonation({
            id: donation.applicationData.id,
            amount: Number(donatedBloodVolume.replace(',', '.')),
        });

        if (!response) {
            showToast('Не удалось подтвердить донацию');

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });
        await queryClient.invalidateQueries({ queryKey: ['plannedDonations', userId] });

        onClose();
    };

    const applyCancel = async (reason?: string) => {
        const response = await cancelDonation({ id: donation.applicationData.id, reason: reason || '' });

        if (!response) {
            showToast('Не удалось отменить донацию');

            return false;
        }

        return true;
    };

    const onCancelClickHandler = async () => {
        const isOk = await applyCancel();

        if (!isOk) {
            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });
        await queryClient.invalidateQueries({ queryKey: ['plannedDonations', userId] });

        onClose();
    };

    const onCancelAcceptDonationClickHandler = async (reason: string) => {
        const isOk = await applyCancel(reason);

        if (!isOk) {
            return;
        }

        setRejectDonationFormParams({ isOpen: true, view: RejectView.FINAL });
    };

    const onBlurDonatedBloodVolumeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        if (Number(value.replace(',', '.')) < 10) {
            setDonatedBloodVolume('10');
        }
    };

    const onChangeDonatedBloodVolumeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        const newValue = value.replaceAll(' ', '');

        if (!newValue) {
            setDonatedBloodVolume('');

            return;
        }

        if (!newValue.match(regexReal)) {
            return;
        }

        if (Number(newValue.replace(',', '.')) > (donation.applicationData.amount || 0)) {
            return;
        }

        setDonatedBloodVolume(newValue);
    };

    const onBonusesClickToggle = () => {
        setIsBonusesPageOpen((prevState) => !prevState);
    };

    const onConditionsClickToggle = () => {
        setIsConditionsOpen((prevState) => !prevState);
    };

    const onChangeBloodGroupHandler = (newBloodGroup: string) => () => {
        setBloodGroup(newBloodGroup);
    };

    const onPendingTimerExpired = () => {
        setIsPendingTimerExpired(true);
    };

    const onCompleteTimerExpired = async () => {
        const response = await confirmDonation({
            id: donation.applicationData.id,
            amount: donation.applicationData.amount,
        });

        if (!response) {
            showToast('Не удалось подтвердить донацию');

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });
        await queryClient.invalidateQueries({ queryKey: ['plannedDonations', userId] });

        onClose();
    };

    const onCloseRejectFormHandler = () => {
        setRejectDonationFormParams({ isOpen: false });
    };

    const onRejectDonationClickHandler = () => {
        setRejectDonationFormParams({ isOpen: true, view: RejectView.CANCEL });
    };

    const onTooltipRegionClick = (e: MouseEvent) => {
        e.stopPropagation();

        setIsTooltipRegionOpen(true);
    };

    useEffect(() => {
        const onOutsideClickHandler = () => {
            setIsTooltipRegionOpen(false);
        };

        window.addEventListener('click', onOutsideClickHandler);

        return () => {
            window.removeEventListener('click', onOutsideClickHandler);
        };
    }, []);

    useEffect(() => {
        if (isErrorLocations) {
            showToast('Не удалось загрузить словарь регионов, попробуйте перезагрузить приложение');
        }
    }, [isErrorLocations, showToast]);

    useEffect(() => {
        if (isErrorPetTypesAndBloodGroups) {
            showToast('Не удалось загрузить словарь типов животных и групп крови, попробуйте перезагрузить приложение');
        }
    }, [isErrorPetTypesAndBloodGroups, showToast]);

    if (rejectDonationFormParams.isOpen && rejectDonationFormParams.view) {
        return (
            <Layout>
                <RejectedForm
                    isDonor
                    userId={userId}
                    onBack={onCloseRejectFormHandler}
                    view={rejectDonationFormParams.view}
                    onSubmit={onCancelAcceptDonationClickHandler}
                />
            </Layout>
        );
    }

    if (isBonusesPageOpen) {
        return <Bonuses fromDonationDetails onClose={onBonusesClickToggle} items={donation.applicationData.bonuses} />;
    }

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.back} onClick={onClose}>
                    <BackAngularArrow />
                </div>
                <h2 className={styles.title}>Планируемая донация</h2>
            </div>
            <div className={styles.photos}>
                <div className={styles.avatarWrapper}>
                    <img
                        className={styles.avatar}
                        alt={donation.applicationData.petName}
                        src={donation.applicationData.photoUrls[0]}
                    />
                    <CircularProgress size={156} strokeWidth={10} total={1} current={1} color='var(--red10, #FF2727)' />
                    <p className={styles.name}>{donation.applicationData.petName.toUpperCase()}</p>
                </div>
                <div className={styles.photoDivider} />
                <div className={styles.avatarWrapper}>
                    <img
                        className={styles.avatar}
                        alt={donation.recipientData.petType}
                        src={
                            donation.recipientData.photoUrls?.[0] ||
                            (donation.recipientData.petType === PetType.DOG ? dogRoundStub : catRoundStub)
                        }
                    />
                    <CircularProgress
                        size={156}
                        strokeWidth={10}
                        color='var(--red10, #FF2727)'
                        total={donation.recipientData.bloodVolumeNeeded}
                        current={
                            donation.applicationData.amount > donation.recipientData.bloodVolumeNeeded
                                ? donation.recipientData.bloodVolumeNeeded
                                : donation.applicationData.amount
                        }
                    />
                    <p className={styles.name}>{donation.recipientData.petName.toUpperCase()}</p>
                </div>
                <div className={styles.neededVolume}>
                    {donation.applicationData.amount > donation.recipientData.bloodVolumeNeeded
                        ? donation.recipientData.bloodVolumeNeeded
                        : donation.applicationData.amount}
                    <span>мл</span>
                </div>
            </div>
            {donation.applicationData.status === DonorStatus.ACCEPTED && !!donation.applicationData.rejectedReason && (
                <Alert
                    className={styles.alert}
                    reason={donation.applicationData.rejectedReason}
                    text='Хозяин реципиента не подтвердил донацию. Свяжитесь с ним для обсуждения деталей.'
                />
            )}
            <div
                className={cn(styles.actions, {
                    [styles.noMargin]: donation.applicationData.status === DonorStatus.PENDING,
                })}
            >
                {donation.applicationData.status === DonorStatus.ACCEPTED && (
                    <>
                        <div onClick={onChatOpenHandler} className={styles.chatIcon}>
                            <Chat />
                        </div>
                        <Button className={styles.confirmDonation} onClick={onConfirmDonationToggle}>
                            Донация состоялась
                        </Button>
                        <Button
                            onClick={onRejectDonationClickHandler}
                            className={cn(styles.confirmDonation, { [styles.reject]: true })}
                        >
                            Донация отменилась
                        </Button>
                    </>
                )}
                {donation.applicationData.status === DonorStatus.COMPLETED && (
                    <>
                        <div onClick={onChatOpenHandler} className={styles.chatIcon}>
                            <Chat />
                        </div>
                        {isWithinHours(donation.applicationData.updatedAt, 72) ? (
                            <div className={cn(styles.count, { [styles.completed]: true })}>
                                <Timer
                                    hoursToAdd={72}
                                    className={styles.countTimer}
                                    onTimeEnd={onCompleteTimerExpired}
                                    digitClassName={styles.countDigits}
                                    separatorClassName={styles.countSeparator}
                                    updatedAt={donation.applicationData.updatedAt}
                                />
                                <p className={styles.timerText}>до подтверждения донации</p>
                            </div>
                        ) : (
                            <div className={styles.recipientConfirmation}>Ожидается подтверждение реципиента</div>
                        )}
                    </>
                )}
            </div>
            {donation.applicationData.status === DonorStatus.PENDING &&
                isWithinHours(donation.applicationData.createdAt, 1) &&
                !isPendingTimerExpired && (
                    <div className={styles.count}>
                        <Timer
                            hoursToAdd={1}
                            className={styles.countTimer}
                            onTimeEnd={onPendingTimerExpired}
                            digitClassName={styles.countDigits}
                            separatorClassName={styles.countSeparator}
                            updatedAt={donation.applicationData.createdAt}
                        />
                        <p className={styles.timerText}>до того, как сможете отказаться</p>
                    </div>
                )}
            <div className={styles.recipientInfo}>
                <div className={styles.infoTitle}>
                    <div className={styles.infoTitleIcon}>
                        <MiniSinglePaw />
                    </div>
                    <p className={styles.infoTitleDescr}>О реципиенте</p>
                </div>
                <div className={styles.line}>
                    <div className={styles.lineItem}>
                        <div className={styles.lineTitle}>
                            <div className={styles.lineTitleIcon}>
                                <Blood />
                            </div>
                            <p className={styles.lineTitleText}>Ищет</p>
                        </div>
                        <div className={styles.bloodInfo}>
                            <div className={styles.bloodGroup}>{donation.recipientData.bloodGroup}</div>
                            {(donation.recipientData.searchingBloodNames.some(
                                (group) => group !== donation.recipientData.bloodGroup,
                            ) ||
                                donation.recipientData.includeUnknownBloodGroup) &&
                                ' +'}
                            {donation.recipientData.searchingBloodNames.some(
                                (group) => group !== donation.recipientData.bloodGroup,
                            )
                                ? donation.recipientData.searchingBloodNames
                                      .filter((group) => group !== donation.recipientData.bloodGroup)
                                      .map((group) => (
                                          <div
                                              key={group}
                                              className={cn(styles.bloodGroup, {
                                                  [styles.needed]: true,
                                              })}
                                          >
                                              {group}
                                          </div>
                                      ))
                                : ''}
                            {donation.recipientData.includeUnknownBloodGroup && (
                                <div className={cn(styles.bloodGroup, { [styles.needed]: true })}>?</div>
                            )}
                        </div>
                    </div>
                    <div className={styles.lineItem}>
                        <div className={cn(styles.lineTitle, { [styles.leftMargin]: true })}>
                            <div className={styles.lineTitleIcon}>
                                <BloodVolume />
                            </div>
                            <p className={styles.lineTitleText}>Нужный объем</p>
                        </div>
                        <div className={styles.bloodVolume}>
                            <p className={styles.bloodVolumeNumber}>{donation.recipientData.bloodVolumeNeeded}</p>
                            <p className={styles.bloodVolumeDescr}>мл</p>
                        </div>
                    </div>
                </div>
                <div className={cn(styles.line, { [styles.last]: true })}>
                    <div className={cn(styles.lineItem, { [styles.isAvatar]: true })}>
                        <div className={styles.ownerAvatar}>
                            {donation.recipientData.ownerName.charAt(0).toUpperCase()}
                        </div>
                        <div className={styles.ownerInfo}>
                            <p className={styles.ownerTitle}>Хозяин</p>
                            <p className={styles.ownerName}>{donation.recipientData.ownerName}</p>
                        </div>
                    </div>
                    <div className={cn(styles.lineItem, { [styles.isRegion]: true })}>
                        <div className={cn(styles.lineTitle, { [styles.leftMargin]: true })}>
                            <div className={styles.lineTitleIcon}>
                                <Location />
                            </div>
                            <p className={styles.lineTitleText}>Регион</p>
                        </div>

                        {donation.recipientData.regions.length === 1 ? (
                            <p className={styles.text}>
                                {donation.recipientData.regions.map(
                                    (lock) => locationsDict.filter(({ value }) => value === lock)[0]?.label,
                                )}
                            </p>
                        ) : (
                            <div className={styles.region}>
                                <p className={styles.text}>{donation.recipientData.regions.length} региона</p>
                                <div onClick={onTooltipRegionClick} className={styles.infoIcon}>
                                    <Info />
                                </div>
                                {isTooltipRegionOpen && (
                                    <div className={styles.tooltip}>
                                        {donation.recipientData.regions.map((lock) => (
                                            <p key={lock} className={styles.text}>
                                                {locationsDict.filter(({ value }) => value === lock)[0]?.label}
                                            </p>
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            </div>
            {(!!donation.recipientData?.advancedInfo?.description?.length ||
                !!donation.recipientData?.advancedInfo?.photoUrls?.[0]) && (
                <Accordion className={styles.accordion} icon={<Pin />} title='Дополнительная информация'>
                    <div className={styles.itemText}>{donation.recipientData?.advancedInfo?.description}</div>
                    {donation.recipientData?.advancedInfo?.photoUrls?.[0] && (
                        <img
                            alt='Фото рецепиента'
                            className={styles.bloodRequestPhoto}
                            src={donation.recipientData.advancedInfo?.photoUrls[0]}
                        />
                    )}
                </Accordion>
            )}
            <div className={styles.settings}>
                <div className={styles.setting} onClick={onBonusesClickToggle}>
                    <p className={styles.text}>
                        Бонусы
                        <br />
                        Портала
                    </p>
                    <div className={styles.icons}>
                        <div className={styles.bonusIcon}>
                            <PrioritySearch />
                        </div>
                        {!!donation.applicationData.bonuses?.length &&
                            donation.applicationData.bonuses[0].type === BonusType.LOCK && (
                                <div className={cn(styles.bonusIcon, { [styles.lock]: true })}>
                                    <Lock />
                                </div>
                            )}
                        {!!donation.applicationData.bonuses?.length &&
                            donation.applicationData.bonuses[0].type !== BonusType.LOCK && (
                                <div className={styles.bonusesCount}>
                                    <p className={styles.countPlus}>+ </p>
                                    <div>{donation.applicationData.bonuses?.length}</div>
                                </div>
                            )}
                    </div>
                </div>
                <div className={styles.setting} onClick={onConditionsClickToggle}>
                    <p className={styles.text}>
                        Ваши
                        <br />
                        условия
                    </p>
                    <div className={styles.icons}>
                        {donation.applicationData.compensationType === CompensationType.FOOD && (
                            <div className={cn(styles.rewardFeedIcon, { [styles.button]: true })}>
                                <div className={styles.icon}>
                                    <Bone />
                                </div>
                            </div>
                        )}
                        {donation.applicationData.compensationType === CompensationType.FREE && (
                            <div className={cn(styles.rewardFreeIcon, { [styles.button]: true })}>
                                <p className={styles.sum}>0</p>
                                <p className={styles.descr}>₽</p>
                            </div>
                        )}
                        {donation.applicationData.compensationType === CompensationType.PAID && (
                            <div className={cn(styles.rewardNotFreeIcon, { [styles.button]: true })}>
                                <p className={styles.descr}>₽</p>
                            </div>
                        )}
                        {donation.applicationData.taxiCompensation && (
                            <div className={cn(styles.taxiIcon, { [styles.button]: true })}>
                                <div className={styles.icon}>
                                    <Taxi />
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
            {donation.applicationData.status === DonorStatus.PENDING &&
                (!isWithinHours(donation.applicationData.createdAt, 1) || isPendingTimerExpired) && (
                    <div className={styles.cancel} onClick={onCancelClickHandler}>
                        <div className={styles.cancelIcon}>
                            <Cancel />
                        </div>
                        <p className={styles.cancelDescr}>Отказаться от донации</p>
                    </div>
                )}
            {isDonorConfirmationCurtainOpen && (
                <Curtain
                    columnOfButtons
                    shouldCloseByWrapperClick
                    cancelButtonTitle='Подтвердить'
                    title={
                        <>
                            Укажите параметры
                            <br />
                            проведенной донации
                        </>
                    }
                    onClose={onConfirmDonationToggle}
                    onConfirm={onConfirmDonationToggle}
                    onCancel={onConfirmDonationClickHandler}
                    isDisableCancelButton={
                        !donatedBloodVolume ||
                        Number(donatedBloodVolume) < 10 ||
                        (isDonorUnknownBloodGroup && !bloodGroup)
                    }
                >
                    <p className={styles.confirmParamDescr}>Объем</p>
                    <TextField
                        name='volume'
                        isDigitInput
                        placeholder=''
                        value={donatedBloodVolume}
                        inputClass={styles.volumeInput}
                        htmlInputClass={styles.volumeHtmlInput}
                        onBlur={onBlurDonatedBloodVolumeHandler}
                        onChange={onChangeDonatedBloodVolumeHandler}
                        endAdornment={<div className={styles.endAdornment}>мл</div>}
                    />
                    {isDonorUnknownBloodGroup && !!bloodGroupDict && (
                        <div className={styles.blood}>
                            <p className={styles.confirmParamDescr}>Группа крови донора</p>
                            <div
                                className={cn(styles.bloodGroups, {
                                    [styles.dogGroup]: donation.recipientData.petType === PetType.DOG,
                                })}
                            >
                                {bloodGroupDict[donation.recipientData.petType]
                                    ?.filter((item) => item.value !== 'UNKNOWN')
                                    .map(({ label, value }) => (
                                        <div
                                            key={value}
                                            onClick={onChangeBloodGroupHandler(label)}
                                            className={cn(styles.bloodItem, { [styles.checked]: bloodGroup === label })}
                                        >
                                            {label}
                                        </div>
                                    ))}
                            </div>
                        </div>
                    )}
                </Curtain>
            )}
            {chatCurtain.isOpen && (
                <Curtain noRednerButtons shouldCloseByWrapperClick onClose={onCloseChatCurtainClickHandler}>
                    <div className={styles.exclamation}>
                        <Exclamation />
                    </div>
                    <h3 className={styles.curtainTitle}>Будьте внимательны!</h3>
                    <div className={styles.curtainList}>
                        {curtainList.map((item, i) => (
                            <div key={item} className={styles.curtainListItem}>
                                <div className={styles.curtainListItemNumber}>{i + 1}</div>
                                <p className={styles.curtainListItemText}>{item}</p>
                            </div>
                        ))}
                    </div>
                    <p className={styles.linkDescr}>Связаться с хозяином реципиента</p>
                    <div className={styles.callButtonWrapper}>
                        <Button
                            onClick={onCallClickHandler}
                            className={styles.callButton}
                            startIcon={
                                <div className={styles.phoneIcon}>
                                    <Phone />
                                </div>
                            }
                        >
                            Позвонить
                        </Button>
                    </div>
                    <div className={styles.messengers}>
                        {chatCurtain.identities?.map(({ providerId, providerName }) => (
                            <div
                                key={providerId}
                                className={styles.identity}
                                onClick={onMessengerClickHandler(providerName)}
                            >
                                {providerName === 'telegram_bot' ? <Telegram /> : <Max />}
                            </div>
                        ))}
                    </div>
                    <p className={styles.backLink} onClick={onCloseChatCurtainClickHandler}>
                        Вернуться
                    </p>
                </Curtain>
            )}
            {isConditionsOpen && (
                <Curtain
                    noRednerButtons
                    title='Ваши условия'
                    shouldCloseByWrapperClick
                    onClose={onConditionsClickToggle}
                >
                    <div
                        className={cn(styles.conditions, {
                            [styles.isTaxi]: donation.applicationData.taxiCompensation,
                        })}
                    >
                        <div
                            className={cn(styles.conditionTile, {
                                [styles.isTaxi]: donation.applicationData.taxiCompensation,
                            })}
                        >
                            {donation.applicationData.compensationType === CompensationType.FOOD && (
                                <>
                                    <div className={styles.rewardFeedIcon}>
                                        <div className={styles.icon}>
                                            <Bone />
                                        </div>
                                    </div>
                                    <p className={styles.conditionDescr}>Готов помочь за корм</p>
                                </>
                            )}
                            {donation.applicationData.compensationType === CompensationType.FREE && (
                                <>
                                    <div className={styles.rewardFreeIcon}>
                                        <p className={styles.sum}>0</p>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                    <p className={styles.conditionDescr}>Готов помочь безвозмездно</p>
                                </>
                            )}
                            {donation.applicationData.compensationType === CompensationType.PAID && (
                                <>
                                    <div className={styles.rewardNotFreeIcon}>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                    <p className={styles.conditionDescr}>Не готов помочь безвозмездно</p>
                                </>
                            )}
                        </div>
                        {donation.applicationData.taxiCompensation && (
                            <div
                                className={cn(styles.conditionTile, {
                                    [styles.isTaxi]: donation.applicationData.taxiCompensation,
                                })}
                            >
                                <div className={styles.taxiIcon}>
                                    <div className={styles.icon}>
                                        <Taxi />
                                    </div>
                                </div>
                                <p className={styles.conditionDescr}>Нужно компенсировать такси до клиники и обратно</p>
                            </div>
                        )}
                    </div>
                </Curtain>
            )}
        </Layout>
    );
};

export default DonationDetails;
