import { Button } from '@mui/material';
import cn from 'classnames';
import {
    BloodAndBreedGroupsDict,
    useGendersQuery,
    useHealthStatusesQuery,
    useLivingConditionsQuery,
    usePetTypesAndBloodGroupsQuery,
    useReproductiveStatusesQuery,
} from 'hooks/useDicts';
import { useGetUserById } from 'hooks/useGetUserById';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import AccordionArrow from 'imgs/svg/accordionArrow';
import Analizes from 'imgs/svg/analizes';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Bone from 'imgs/svg/bone';
import Chat from 'imgs/svg/chat';
import DonatedVolume from 'imgs/svg/donatedVolume';
import DonorButton from 'imgs/svg/donorButton';
import Exclamation from 'imgs/svg/exclamation';
import Health from 'imgs/svg/health';
import Max from 'imgs/svg/max';
import Params from 'imgs/svg/params';
import Phone from 'imgs/svg/phone';
import Processing from 'imgs/svg/processing';
import Taxi from 'imgs/svg/taxi';
import Telegram from 'imgs/svg/telegram';
import DonationQuestions from 'pages/owner/Statuses/DonationQuestions';
import { ChangeEvent, FC, useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { regexReal } from 'utils/regexps';
import { isWithinHours, matchIdentities } from 'utils/utils';

import { confirmDonation } from 'api/apiServices/confirmDonation';
import { getDonationForRecipientById } from 'api/apiServices/getDonationForRecipientById';
import { getUserContacts } from 'api/apiServices/getUserContacts';
import { rejectDonation } from 'api/apiServices/rejectDonation';
import { updatePet } from 'api/apiServices/updatePet';
import { GetDonationForRecipientByIdResponse, RespondingDonorStatus } from 'api/bloodRequest';
import { Pet } from 'api/pets';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import { CompensationType, Identities } from 'api/user';
import { CircularProgress } from 'components/CircularProgress';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';
import Loading from 'components/Loading';
import AnalysesStep from 'components/Profiles/Steps/Analyses';
import HealthStep from 'components/Profiles/Steps/Health';
import ParamsStep from 'components/Profiles/Steps/Params';
import TreatmentsStep from 'components/Profiles/Steps/Treatments';
import RejectedForm, { RejectView } from 'components/RejectedForm';
import TextField from 'components/TextField';
import Timer from 'components/Timer';

import styles from './DonationDetails.module.less';

enum TileName {
    PARAMS = 'params',
    HEALTH = 'health',
    ANALYSES = 'analyses',
    DONATIONS = 'donations',
    TREATMENTS = 'treatments',
    CONDITIONS = 'conditions',
}

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
    updatedAt: string;
    donationId: string;
    onClose: () => void;
    onReject: () => void;
    onIsSearchFinish: () => void;
    status: RespondingDonorStatus;
    onDonationComplete: (donatedVolume: number) => void;
};

const dogAnalizesCount = 6;
const catAnalizesCount = 4;

const tiles = [
    { name: TileName.PARAMS, title: 'Параметры', icon: <Params /> },
    { name: TileName.HEALTH, title: 'Здоровье', icon: <Health /> },
    { name: TileName.TREATMENTS, title: 'Обработки', icon: <Processing /> },
    { name: TileName.ANALYSES, title: 'Анализы', icon: <Analizes /> },
    { name: TileName.DONATIONS, title: 'История донаций', icon: <DonorButton /> },
    { name: TileName.CONDITIONS, title: 'Условия донора' },
];

const curtainList = [
    'Не передавайте вознаграждение до проведения донации',
    'Не переходите по подозрительным ссылкам',
    'Не передавайте свои паспортные данные',
];

const DonationDetails: FC<Props> = ({
    userId,
    status,
    onClose,
    onReject,
    updatedAt,
    donationId,
    onIsSearchFinish,
    onDonationComplete,
}) => {
    const [isLoading, setIsLoading] = useState(true);
    const [isConditionsOpen, setIsConditionsOpen] = useState(false);
    const [activeTile, setaActiveTile] = useState<TileName | null>(null);
    const [chatCurtain, setChatCurtain] = useState<ChatCurtain>({ isOpen: false });
    const [donatedBloodVolume, setDonatedBloodVolume] = useState<string>('');
    const [isDonorWarnFactorsOpen, setIsDonorWarnFactorsOpen] = useState(false);
    const [donation, setDonation] = useState<GetDonationForRecipientByIdResponse | null>(null);
    const [isDonorConfirmationCurtainOpen, setIsDonorConfirmationCurtainOpen] = useState(false);
    const [rejectDonationFormParams, setRejectDonationFormParams] = useState<RejectedFormType>({ isOpen: false });
    const [isRecipientConfirmationCurtainOpen, setIsRecipientConfirmationCurtainOpen] = useState(false);

    const [donorBloodGroup, setDonorBloodGroup] = useState<string | null>(null);

    // dicts
    const { data: petGendersDict = [], isError: isErrorGenders } = useGendersQuery();
    const { data: healthStatusesDict = [], isError: isErrorHealthStatuses } = useHealthStatusesQuery();
    const { data: livingConditionsDict = [], isError: isErrorLivingConditions } = useLivingConditionsQuery();
    const { data: reproductiveStatusesDict = [], isError: isErrorReproductiveStatusesDict } =
        useReproductiveStatusesQuery();
    const {
        data: { petTypesDict = [], bloodGroupDict = {}, breedsDict = {} } = {},
        isError: isErrorPetTypesAndBloodGroups,
    } = usePetTypesAndBloodGroupsQuery();

    const { data: userData } = useGetUserById(userId);

    const isDonorUnknownBloodGroup = donation?.donorData.bloodGroup === 'UNKNOWN';

    const showToast = useCallback(
        (text: string, type = 'warn') => {
            if (type === 'success') {
                toast.success(text);

                return;
            }

            toast.warn(text, {
                onClose: () => {
                    onClose();
                },
            });
        },
        [onClose],
    );

    const fetchDetails = useCallback(async () => {
        const response = await getDonationForRecipientById(donationId);

        if (!response) {
            showToast('Не удалось получить детали донации');

            return;
        }

        setIsLoading(false);
        setDonation(response.data);
    }, [donationId, showToast]);

    const onConditionsClickToggle = () => {
        setIsConditionsOpen((prevState) => !prevState);
    };

    const onTileClickHandler = (tileName: TileName) => () => {
        setaActiveTile(tileName);

        if (tileName === TileName.CONDITIONS) {
            onConditionsClickToggle();
        }
    };

    const onTileBackHandler = () => {
        setaActiveTile(null);
    };

    const onConfirmDonationToggle = () => {
        if (!donation) {
            return;
        }

        setDonatedBloodVolume(
            donation.donorData.availableBloodAmount > donation.recipientData.bloodVolumeNeeded
                ? `${donation.recipientData.bloodVolumeNeeded}`.replace('.', ',')
                : `${donation.donorData.availableBloodAmount}`.replace('.', ','),
        );
        setIsRecipientConfirmationCurtainOpen((prevState) => !prevState);
        setDonorBloodGroup(null);
    };

    const onDonorConfirmDonationToggle = () => {
        setIsDonorConfirmationCurtainOpen((prevState) => !prevState);
    };

    const onRejectDonationClickHandler = (view: RejectView) => () => {
        setRejectDonationFormParams({ isOpen: true, view });
    };

    const onCallClickHandler = (e) => {
        const storedEnv = localStorage.getItem('environment');

        if (!donation) {
            return;
        }

        if (storedEnv === 'tg') {
            e.preventDefault();
            const phone = `tel:+${donation.donorData.ownerPhone.replace(/[^\d]/g, '')}`;
            const a = document.createElement('a');
            a.href = phone;
            a.target = '_blank';
            a.rel = 'noopener noreferrer';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } else {
            window.location.href = `tel:${donation?.donorData.ownerPhone}`;
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

        if (Number(newValue.replace(',', '.')) > (donation?.donorData.availableBloodAmount || 0)) {
            return;
        }

        setDonatedBloodVolume(newValue);
    };

    const onBlurDonatedBloodVolumeHandler = ({
        target: { value },
    }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        if (Number(value.replace(',', '.')) < 10) {
            setDonatedBloodVolume('10');
        }
    };

    const onCloseChatCurtainClickHandler = () => {
        setChatCurtain({ isOpen: false });
    };

    const onMessengerClickHandler = (providerName: string) => async () => {
        const response = await getUserContacts({ id: donation?.donorData.ownerId!, provider: providerName });

        if (!response) {
            showToast('Не удалось получить контакт хозяина донора');

            setChatCurtain({ isOpen: false });

            return;
        }

        showToast(response.data.message, 'success');

        setChatCurtain({ isOpen: false });
    };

    const onChatOpenHandler = () => {
        setChatCurtain({
            isOpen: true,
            identities: matchIdentities(userData?.identities!, donation?.donorData.identities),
        });
    };

    const onEndTimerClickHandler = async () => {
        if (!donation) {
            return;
        }

        const response = await confirmDonation({
            id: donationId,
            amount: donation.donorData.application.amount,
        });

        if (!response) {
            showToast('Не удалось подтвердить донацию');
        }

        const isFinish = donation.recipientData.bloodVolumeNeeded <= donation.donorData.application.amount;

        if (isFinish) {
            onIsSearchFinish();
        } else {
            onDonationComplete(donation.donorData.application.amount);
        }
    };

    const onConfirmDonationClickHandler =
        (fromRecipient = true) =>
        async () => {
            if (!donation) {
                return;
            }

            if (donorBloodGroup) {
                const { success } = await updatePet({ id: donation.donorData.id, bloodGroup: donorBloodGroup } as Pet);

                if (!success) {
                    showToast('Не удалось сохранить выбранную группу крови, для донора');
                }
            }

            const response = await confirmDonation({
                id: donationId,
                amount: fromRecipient
                    ? Number(donatedBloodVolume.replace(',', '.'))
                    : donation.donorData.application.amount,
            });

            if (!response) {
                showToast('Не удалось подтвердить донацию');

                return;
            }

            const donated = fromRecipient
                ? Number(donatedBloodVolume.replace(',', '.'))
                : donation.donorData.application.amount;
            const isFinish =
                donation.recipientData.bloodVolumeNeeded <= donation.recipientData.bloodVolumeDonated + donated;

            if (isFinish) {
                onIsSearchFinish();
            } else {
                onDonationComplete(donated);
            }
        };

    const onChangeBloodGroupHandler = (newBloodGroup: string) => () => {
        setDonorBloodGroup(newBloodGroup);
    };

    const onCloseRejectFormHandler = () => {
        setRejectDonationFormParams({ isOpen: false });
    };

    const onOpenWarnFactorsToggle = () => {
        setIsDonorWarnFactorsOpen((prevState) => !prevState);
    };

    const onSubmitRejectFormHandler = async (reason: string) => {
        const response = await rejectDonation({ id: donationId, reason });

        if (!response) {
            showToast('Не удалось отменить донацию');
        } else {
            await queryClient.invalidateQueries({ queryKey: ['pets', userId] });
        }

        if (rejectDonationFormParams.view === RejectView.CANCEL) {
            setRejectDonationFormParams({ isOpen: true, view: RejectView.FINAL });
        } else {
            onClose();
            onReject();
        }
    };

    useEffect(() => {
        fetchDetails();
    }, [fetchDetails]);

    useEffect(() => {
        if (isErrorPetTypesAndBloodGroups) {
            showToast('Не удалось загрузить словарь типов животных и групп крови, попробуйте перезагрузить приложение');
        }
    }, [isErrorPetTypesAndBloodGroups, showToast]);

    useEffect(() => {
        if (isErrorGenders) {
            showToast('Не удалось загрузить словарь полов, попробуйте перезагрузить приложение');
        }
    }, [isErrorGenders, showToast]);

    useEffect(() => {
        if (isErrorHealthStatuses) {
            showToast('Не удалось загрузить словарь статусов здоровья, попробуйте перезагрузить приложение');
        }
    }, [isErrorHealthStatuses, showToast]);

    useEffect(() => {
        if (isErrorLivingConditions) {
            showToast('Не удалось загрузить словарь условий проживания, попробуйте перезагрузить приложение');
        }
    }, [isErrorLivingConditions, showToast]);

    useEffect(() => {
        if (isErrorReproductiveStatusesDict) {
            showToast('Не удалось загрузить словарь репродуктивных состояний, попробуйте перезагрузить приложение');
        }
    }, [isErrorReproductiveStatusesDict, showToast]);

    if (isLoading || !donation) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    if (isDonorWarnFactorsOpen) {
        return (
            <DonationQuestions
                isRecipientOpen
                onClose={onOpenWarnFactorsToggle}
                factors={donation.donorData.donorRestrictions?.warnFactors}
            />
        );
    }

    if (rejectDonationFormParams.isOpen && rejectDonationFormParams.view) {
        return (
            <Layout>
                <RejectedForm
                    userId={userId}
                    onBack={onCloseRejectFormHandler}
                    view={rejectDonationFormParams.view}
                    onSubmit={onSubmitRejectFormHandler}
                />
            </Layout>
        );
    }

    if (activeTile === TileName.PARAMS) {
        return (
            <Layout>
                <ParamsStep
                    isEditMode={false}
                    petTypes={petTypesDict}
                    id={donation.donorData.id}
                    onClose={onTileBackHandler}
                    petGenders={petGendersDict}
                    type={donation.donorData.type}
                    name={donation.donorData.name}
                    gender={donation.donorData.gender}
                    breedId={donation.donorData.breedId}
                    weightKg={donation.donorData.weightKg}
                    birthDate={donation.donorData.birthDate}
                    chipNumber={donation.donorData.chipNumber}
                    bloodGroup={donation.donorData.bloodGroup}
                    livingConditionsDict={livingConditionsDict}
                    breedsDict={breedsDict as BloodAndBreedGroupsDict}
                    reproductiveStatusesDict={reproductiveStatusesDict}
                    livingCondition={donation.donorData.livingCondition}
                    reproductiveStatus={donation.donorData.reproductiveStatus}
                    bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
                />
            </Layout>
        );
    }

    if (activeTile === TileName.HEALTH) {
        return (
            <Layout>
                <HealthStep
                    isEditMode={false}
                    id={donation.donorData.id}
                    onClose={onTileBackHandler}
                    healthStatusesDict={healthStatusesDict}
                    transfused={donation.donorData.health?.transfused}
                    medications={donation.donorData.health?.medications}
                    healthStatus={donation.donorData.health?.healthStatus}
                    surgicalInterventions={donation.donorData.health?.surgicalInterventions}
                />
            </Layout>
        );
    }

    if (activeTile === TileName.TREATMENTS) {
        return (
            <Layout>
                <TreatmentsStep
                    isEditMode={false}
                    id={donation.donorData.id}
                    onClose={onTileBackHandler}
                    dewormingDate={donation.donorData.treatments?.dewormingDate}
                    rabiesVaccinationDate={donation.donorData.treatments?.rabiesVaccinationDate}
                    infectionVaccinationDate={donation.donorData.treatments?.infectionVaccinationDate}
                    ectoparasiteTreatmentDate={donation.donorData.treatments?.ectoparasiteTreatmentDate}
                />
            </Layout>
        );
    }

    if (activeTile === TileName.ANALYSES) {
        return (
            <Layout>
                <AnalysesStep
                    isEditMode={false}
                    onClose={onTileBackHandler}
                    petId={donation.donorData.id}
                    petType={donation.donorData.type}
                    analyses={donation.donorData.analyses}
                />
            </Layout>
        );
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
                        alt={donation.donorData.name}
                        src={
                            donation.donorData.photoUrls[0] ||
                            (donation.donorData.type === PetType.DOG ? dogRoundStub : catRoundStub)
                        }
                    />
                    <CircularProgress size={156} strokeWidth={10} total={1} current={1} color='var(--red10, #FF2727)' />
                    <p className={styles.name}>{donation.donorData.name.toUpperCase()}</p>
                    {!!donation.donorData.donorRestrictions?.warnFactors?.length && (
                        <span onClick={onOpenWarnFactorsToggle} className={styles.warnFactors}>
                            ?
                        </span>
                    )}
                </div>
                <div className={styles.photoDivider} />
                <div className={styles.avatarWrapper}>
                    <img
                        className={styles.avatar}
                        alt={donation.donorData.name}
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
                            donation.recipientData.bloodVolumeDonated + donation.donorData.application.amount >
                            donation.recipientData.bloodVolumeNeeded
                                ? donation.recipientData.bloodVolumeNeeded
                                : donation.recipientData.bloodVolumeDonated + donation.donorData.application.amount
                        }
                    />
                    <p className={styles.name}>{donation.recipientData.petName.toUpperCase()}</p>
                </div>
                <div className={styles.neededVolume}>
                    {donation.donorData.availableBloodAmount}
                    <span>мл</span>
                </div>
            </div>
            <div className={cn(styles.actions, { [styles.completed]: status === RespondingDonorStatus.COMPLETED })}>
                {status === RespondingDonorStatus.ACCEPTED && (
                    <>
                        <div onClick={onChatOpenHandler} className={styles.chatIcon}>
                            <Chat />
                        </div>
                        <Button className={styles.confirmDonation} onClick={onConfirmDonationToggle}>
                            Донация состоялась
                        </Button>
                        <Button
                            className={cn(styles.confirmDonation, { [styles.reject]: true })}
                            onClick={onRejectDonationClickHandler(RejectView.CANCEL)}
                        >
                            Донация отменилась
                        </Button>
                    </>
                )}
                {status === RespondingDonorStatus.COMPLETED && (
                    <>
                        {isWithinHours(updatedAt, 72) && (
                            <div className={styles.count}>
                                <p className={styles.timerText}>Хозяин донора сообщил о донации</p>
                                <div className={styles.timerWrapper}>
                                    <Timer
                                        hoursToAdd={72}
                                        updatedAt={updatedAt}
                                        className={styles.countTimer}
                                        onTimeEnd={onEndTimerClickHandler}
                                        digitClassName={styles.countDigits}
                                        separatorClassName={styles.countSeparator}
                                    />
                                    <p className={styles.timerCaption}>до автоматического подтверждения</p>
                                </div>
                            </div>
                        )}
                        <div onClick={onChatOpenHandler} className={styles.chatIcon}>
                            <Chat />
                        </div>
                        <Button className={styles.confirmDonation} onClick={onDonorConfirmDonationToggle}>
                            Подробнее
                        </Button>
                    </>
                )}
            </div>
            <div className={styles.tiles}>
                {tiles.map(({ name: tileName, title, icon }) => (
                    <div
                        key={tileName}
                        onClick={onTileClickHandler(tileName)}
                        className={cn(styles.tile, {
                            [styles.hide]: tileName === TileName.DONATIONS,
                            [styles.conditionTile]: tileName === TileName.CONDITIONS,
                        })}
                    >
                        {icon && (
                            <div
                                className={cn(styles.tileIcon, {
                                    [styles.needFill]: tileName === TileName.DONATIONS,
                                })}
                            >
                                {icon}
                            </div>
                        )}
                        <p className={styles.tileTitle}>{title}</p>
                        {tileName !== TileName.CONDITIONS && (
                            <div className={styles.arrowTileIcon}>
                                <AccordionArrow />
                            </div>
                        )}
                        {tileName === TileName.ANALYSES && (
                            <div className={styles.analizesCount}>
                                {Object.keys(donation.donorData.analyses || []).length} из{' '}
                                {donation.donorData.type === PetType.DOG ? dogAnalizesCount : catAnalizesCount}
                            </div>
                        )}
                        {tileName === TileName.CONDITIONS && (
                            <div className={styles.conditions}>
                                {donation.donorData.application.compensationType === CompensationType.FOOD && (
                                    <div className={cn(styles.rewardFeedIcon, { [styles.button]: true })}>
                                        <div className={styles.icon}>
                                            <Bone />
                                        </div>
                                    </div>
                                )}
                                {donation.donorData.application.compensationType === CompensationType.FREE && (
                                    <div className={cn(styles.rewardFreeIcon, { [styles.button]: true })}>
                                        <p className={styles.sum}>0</p>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                )}
                                {donation.donorData.application.compensationType === CompensationType.PAID && (
                                    <div className={cn(styles.rewardNotFreeIcon, { [styles.button]: true })}>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                )}
                                {donation.donorData.application.taxiCompensation && (
                                    <div className={cn(styles.taxiIcon, { [styles.button]: true })}>
                                        <div className={styles.icon}>
                                            <Taxi />
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                ))}
            </div>
            {isRecipientConfirmationCurtainOpen && (
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
                    onCancel={onConfirmDonationClickHandler()}
                    isDisableCancelButton={
                        !donatedBloodVolume ||
                        Number(donatedBloodVolume) < 10 ||
                        (isDonorUnknownBloodGroup && !donorBloodGroup)
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
                                            className={cn(styles.bloodItem, {
                                                [styles.checked]: donorBloodGroup === label,
                                            })}
                                        >
                                            {label}
                                        </div>
                                    ))}
                            </div>
                        </div>
                    )}
                </Curtain>
            )}
            {isDonorConfirmationCurtainOpen && (
                <Curtain
                    shouldCloseByWrapperClick
                    cancelButtonTitle='Подтверждаю'
                    confirmButtonTitle='Не подтверждаю'
                    onClose={onDonorConfirmDonationToggle}
                    subTitleClassName={styles.donorConfirmationSubTitle}
                    subTitle='Подтвердите, чтобы донор смог получить бонусы'
                    confirmButtonClassName={styles.donorConfirmationConfirm}
                    onCancel={onConfirmDonationClickHandler(false)}
                    onConfirm={onRejectDonationClickHandler(RejectView.NOT_CONFIRM)}
                    title={
                        <>
                            Хозяин донора
                            <br />
                            сообщил о донации
                        </>
                    }
                >
                    <div className={styles.donorConfirmationCurtain}>
                        <div className={styles.donorConfirmationIcon}>
                            <DonatedVolume />
                        </div>
                        <div className={styles.donorConfirmationVolume}>
                            <p className={styles.donorConfirmationNumber}>{donation.donorData.application.amount}</p>
                            <p className={styles.donorConfirmationDescr}>мл</p>
                        </div>
                    </div>
                </Curtain>
            )}
            {chatCurtain.isOpen && (
                <Curtain noRednerButtons shouldCloseByWrapperClick onClose={onCloseChatCurtainClickHandler}>
                    <div className={styles.exclamation}>
                        <Exclamation />
                    </div>
                    <h3 className={styles.curtainTitle}>Будьте внимательны!</h3>
                    <p className={styles.curtainDecr}>
                        Обсудите условия и встретьтесь в клинике для получения помощи. Если не договоритесь - отмените
                        донацию, чтобы освободить лимит поиска.
                    </p>
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
                    title='Условия донора'
                    shouldCloseByWrapperClick
                    onClose={onConditionsClickToggle}
                >
                    <div
                        className={cn(styles.donorConditions, {
                            [styles.isTaxi]: donation.donorData.application.taxiCompensation,
                        })}
                    >
                        <div
                            className={cn(styles.donorConditionTile, {
                                [styles.isTaxi]: donation.donorData.application.taxiCompensation,
                            })}
                        >
                            {donation.donorData.application.compensationType === CompensationType.FOOD && (
                                <>
                                    <div className={styles.rewardFeedIcon}>
                                        <div className={styles.icon}>
                                            <Bone />
                                        </div>
                                    </div>
                                    <p className={styles.donorConditionDescr}>Готов помочь за корм</p>
                                </>
                            )}
                            {donation.donorData.application.compensationType === CompensationType.FREE && (
                                <>
                                    <div className={styles.rewardFreeIcon}>
                                        <p className={styles.sum}>0</p>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                    <p className={styles.donorConditionDescr}>Готов помочь безвозмездно</p>
                                </>
                            )}
                            {donation.donorData.application.compensationType === CompensationType.PAID && (
                                <>
                                    <div className={styles.rewardNotFreeIcon}>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                    <p className={styles.donorConditionDescr}>Не готов помочь безвозмездно</p>
                                </>
                            )}
                        </div>
                        {donation.donorData.application.taxiCompensation && (
                            <div
                                className={cn(styles.donorConditionTile, {
                                    [styles.isTaxi]: donation.donorData.application.taxiCompensation,
                                })}
                            >
                                <div className={styles.taxiIcon}>
                                    <div className={styles.icon}>
                                        <Taxi />
                                    </div>
                                </div>
                                <p className={styles.donorConditionDescr}>
                                    Нужно компенсировать такси до клиники и обратно
                                </p>
                            </div>
                        )}
                    </div>
                </Curtain>
            )}
        </Layout>
    );
};

export default DonationDetails;
