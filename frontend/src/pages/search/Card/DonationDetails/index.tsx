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
import Processing from 'imgs/svg/processing';
import Taxi from 'imgs/svg/taxi';
import Telegram from 'imgs/svg/telegram';
import { ChangeEvent, FC, useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { regexReal } from 'utils/regexps';

import { confirmDonation } from 'api/apiServices/confirmDonation';
import { getDonationForRecipientById } from 'api/apiServices/getDonationForRecipientById';
import { getUserContacts } from 'api/apiServices/getUserContacts';
import { rejectDonation } from 'api/apiServices/rejectDonation';
import { GetDonationForRecipientByIdResponse, RespondingDonorStatus } from 'api/bloodRequest';
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
import TextField from 'components/TextField';

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

type Props = {
    userId: string;
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
    donationId,
    onIsSearchFinish,
    onDonationComplete,
}) => {
    const [isLoading, setIsLoading] = useState(true);
    const [activeTile, setaActiveTile] = useState<TileName | null>(null);
    const [chatCurtain, setChatCurtain] = useState<ChatCurtain>({ isOpen: false });
    const [donatedBloodVolume, setDonatedBloodVolume] = useState<string>('');
    const [donation, setDonation] = useState<GetDonationForRecipientByIdResponse | null>(null);
    const [isDonorConfirmationCurtainOpen, setIsDonorConfirmationCurtainOpen] = useState(false);
    const [isRecipientConfirmationCurtainOpen, setIsRecipientConfirmationCurtainOpen] = useState(false);

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

    const onTileClickHandler = (tileName: TileName) => () => {
        setaActiveTile(tileName);
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
    };

    const onDonorConfirmDonationToggle = () => {
        setIsDonorConfirmationCurtainOpen((prevState) => !prevState);
    };

    const onRejectDonationClickHandler = async () => {
        const response = await rejectDonation(donationId);

        if (!response) {
            showToast('Не удалось отменить донацию');
        }

        onClose();
        onReject();
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

        if (Number(newValue.replace(',', '.')) > (donation?.recipientData.bloodVolumeNeeded || 0)) {
            return;
        }

        setDonatedBloodVolume(newValue);
    };

    const onCloseChatCurtainClickHandler = () => {
        setChatCurtain({ isOpen: false });
    };

    const onMessengerClickHandler = (providerName: string) => async () => {
        const response = await getUserContacts({ id: userId, provider: providerName });

        if (!response) {
            showToast('Не удалось получить контакт хозяина донора');

            setChatCurtain({ isOpen: false });

            return;
        }

        showToast(response.data.message, 'success');

        setChatCurtain({ isOpen: false });
    };

    const onChatOpenHandler = () => {
        setChatCurtain({ isOpen: true, identities: userData?.identities });
    };

    const onConfirmDonationClickHandler =
        (fromRecipient = true) =>
        async () => {
            if (!donation) {
                return;
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

    if (activeTile === TileName.ANALYSES && donation.donorData.analyses) {
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
                            onClick={onRejectDonationClickHandler}
                        >
                            Донация отменилась
                        </Button>
                    </>
                )}
                {status === RespondingDonorStatus.COMPLETED && (
                    <>
                        <div className={styles.completedTile}>Хозяин донора сообщил о донации</div>
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
                            [styles.conditionTile]: tileName === TileName.CONDITIONS,
                            [styles.noActive]:
                                tileName === TileName.DONATIONS ||
                                (tileName === TileName.ANALYSES && !donation.donorData.analyses),
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
                            <div
                                className={cn(styles.analizesCount, {
                                    [styles.noActive]: !donation.donorData.analyses,
                                })}
                            >
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
                    cancelButtonTitle='Потвердить'
                    title={
                        <>
                            Укажите объем
                            <br />
                            проведенной донации
                        </>
                    }
                    onClose={onConfirmDonationToggle}
                    onConfirm={onConfirmDonationToggle}
                    onCancel={onConfirmDonationClickHandler()}
                    isDisableCancelButton={!donatedBloodVolume}
                >
                    <TextField
                        name='volume'
                        isDigitInput
                        placeholder=''
                        value={donatedBloodVolume}
                        inputClass={styles.volumeInput}
                        htmlInputClass={styles.volumeHtmlInput}
                        onChange={onChangeDonatedBloodVolumeHandler}
                        endAdornment={<div className={styles.endAdornment}>мл</div>}
                    />
                </Curtain>
            )}
            {isDonorConfirmationCurtainOpen && (
                <Curtain
                    shouldCloseByWrapperClick
                    cancelButtonTitle='Подтверждаю'
                    confirmButtonTitle='Не подтверждаю'
                    onClose={onDonorConfirmDonationToggle}
                    onConfirm={onRejectDonationClickHandler}
                    subTitleClassName={styles.donorConfirmationSubTitle}
                    confirmButtonClassName={styles.donorConfirmationConfirm}
                    onCancel={onConfirmDonationClickHandler(false)}
                    subTitle='Подтвердите, чтобы донор смог получить бонусы'
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
                    <p className={styles.linkDescr}>Пришлем контакт донора в мессенджер</p>
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
        </Layout>
    );
};

export default DonationDetails;
