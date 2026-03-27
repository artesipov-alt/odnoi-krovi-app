import cn from 'classnames';
import useBodyScrollLock from 'hooks/useBodyScrollLock';
import {
    BloodAndBreedGroupsDict,
    useGendersQuery,
    useHealthStatusesQuery,
    useLivingConditionsQuery,
    usePetTypesAndBloodGroupsQuery,
    useReproductiveStatusesQuery,
} from 'hooks/useDicts';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import AccordionArrow from 'imgs/svg/accordionArrow';
import Analizes from 'imgs/svg/analizes';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import Bone from 'imgs/svg/bone';
import DonorButton from 'imgs/svg/donorButton';
import Exclamation from 'imgs/svg/exclamation';
import Health from 'imgs/svg/health';
import Max from 'imgs/svg/max';
import Params from 'imgs/svg/params';
import Processing from 'imgs/svg/processing';
import StatusQuestion from 'imgs/svg/statusQuestion';
import Taxi from 'imgs/svg/taxi';
import Telegram from 'imgs/svg/telegram';
import DonationQuestions from 'pages/owner/Statuses/DonationQuestions';
import { FC, useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';

import { applyDonorRespond } from 'api/apiServices/applyDonorRespond';
import { getDonorInfo } from 'api/apiServices/getDonorInfo';
import { getPets } from 'api/apiServices/getPets';
import { getUserIdentities } from 'api/apiServices/getUserIdentities';
import { GetDonorInfoResponse } from 'api/bloodRequest';
import { Pet } from 'api/pets';
import { PetGender, PetType } from 'api/types';
import { CompensationType, Identities, Role } from 'api/user';
import { CircularProgress } from 'components/CircularProgress';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import AnalysesStep from '../Steps/Analyses';
import HealthStep from '../Steps/Health';
import ParamsStep from '../Steps/Params';
import TreatmentsStep from '../Steps/Treatments';
import CheckOtherDonors from './CheckOtherDonors';
import styles from './DonorForRecipient.module.less';

enum TileName {
    PARAMS = 'params',
    HEALTH = 'health',
    ANALYSES = 'analyses',
    DONATIONS = 'donations',
    TREATMENTS = 'treatments',
    CONDITIONS = 'conditions',
}

type OtherDonors = {
    pets?: Pet[];
    isOpen: boolean;
};

type ChatCurtain = {
    isOpen: boolean;
    identities?: Identities[];
};

type Props = {
    donorId: string;
    responseId: string;
    onClose: () => void;
    onBackToSearch: () => void;
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

const DonorForRecipient: FC<Props> = ({ onClose, donorId, responseId, onBackToSearch }) => {
    const [isLoading, setIsLoading] = useState(false);
    const [isGlobalLoading, setIsGlobalLoading] = useState(true);
    const [isWarnFactorsOpen, setIsWarnFactorsOpen] = useState(false);
    const [info, setInfo] = useState<GetDonorInfoResponse | null>(null);
    const [chatCurtain, setChatCurtain] = useState<ChatCurtain>({ isOpen: false });
    const [activeTile, setaActiveTile] = useState<TileName | null>(null);
    const [checkOtherDonors, setCheckOtherDonors] = useState<OtherDonors>({ isOpen: false });

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

    useBodyScrollLock(isGlobalLoading || isLoading);

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

    const fetchInfo = useCallback(async () => {
        const response = await getDonorInfo(donorId);

        if (!response) {
            showToast('Не удалось получить информацию о доноре');

            return;
        }

        setInfo(response.data);
        setIsGlobalLoading(false);
    }, [donorId, showToast]);

    const getOwnerIdentities = async () => {
        const response = await getUserIdentities(info?.ownerId!);

        if (!response) {
            showToast('Не удалось получить данные для формирования ссылок для чата');

            return;
        }

        setChatCurtain({ isOpen: true, identities: response.data.identities });
    };

    const checkIsTrueDonor = async () => {
        setIsLoading(true);

        try {
            const { pets } = await getPets(info?.ownerId!);

            const recoveringPets = pets.filter((pet) => pet.petStatus === Role.RECOVERING);

            if (recoveringPets.length) {
                setCheckOtherDonors({ isOpen: true, pets: recoveringPets });
            }

            getOwnerIdentities();
        } catch (error) {
            showToast('Не удалось получить остальных питомцев хозяина донора, для проверки');
        } finally {
            setIsLoading(false);
        }
    };

    const onTileClickHandler = (tileName: TileName) => () => {
        setaActiveTile(tileName);
    };

    const onTileBackHandler = () => {
        setaActiveTile(null);
    };

    const onStatusItemClickToggle = () => {
        setIsWarnFactorsOpen((prevState) => !prevState);
    };

    const onChatClickHandler = () => {
        checkIsTrueDonor();
    };

    const onCloseCheckOtherDonors = () => {
        setCheckOtherDonors({ isOpen: false });
    };

    const onConfirmCheckOtherDonors = () => {
        setCheckOtherDonors({ isOpen: false });

        getOwnerIdentities();
    };

    const onCloseChatCurtainClickHandler = () => {
        setChatCurtain({ isOpen: false });
    };

    const onMessengerClickHandler = async () => {
        const response = await applyDonorRespond(responseId);

        if (!response) {
            showToast('Не удалось откликнуться на предложение донора');

            setChatCurtain({ isOpen: false });

            return;
        }

        onBackToSearch();
    };

    useEffect(() => {
        fetchInfo();
    }, [fetchInfo]);

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

    if (checkOtherDonors.isOpen && !!checkOtherDonors.pets?.length) {
        return (
            <CheckOtherDonors
                pets={checkOtherDonors.pets}
                avatar={info?.photoUrls?.[0]!}
                onClose={onCloseCheckOtherDonors}
                onSuccess={onConfirmCheckOtherDonors}
            />
        );
    }

    if (isWarnFactorsOpen) {
        return <DonationQuestions factors={info?.donorRestrictions?.warnFactors} onClose={onStatusItemClickToggle} />;
    }

    if (activeTile === TileName.PARAMS && info) {
        return (
            <ParamsStep
                id={info.id}
                name={info.name}
                type={info.type}
                isEditMode={false}
                gender={info.gender}
                breedId={info.breedId}
                petTypes={petTypesDict}
                weightKg={info.weightKg}
                birthDate={info.birthDate}
                onClose={onTileBackHandler}
                petGenders={petGendersDict}
                chipNumber={info.chipNumber}
                bloodGroup={info.bloodGroup}
                livingCondition={info.livingCondition}
                livingConditionsDict={livingConditionsDict}
                reproductiveStatus={info.reproductiveStatus}
                breedsDict={breedsDict as BloodAndBreedGroupsDict}
                reproductiveStatusesDict={reproductiveStatusesDict}
                bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
            />
        );
    }

    if (activeTile === TileName.HEALTH && info) {
        return (
            <HealthStep
                id={info.id}
                isEditMode={false}
                onClose={onTileBackHandler}
                transfused={info.health?.transfused}
                medications={info.health?.medications}
                healthStatusesDict={healthStatusesDict}
                healthStatus={info.health?.healthStatus}
                surgicalInterventions={info.health?.surgicalInterventions}
            />
        );
    }

    if (activeTile === TileName.TREATMENTS && info) {
        return (
            <TreatmentsStep
                id={info.id}
                isEditMode={false}
                onClose={onTileBackHandler}
                dewormingDate={info.treatments?.dewormingDate}
                rabiesVaccinationDate={info.treatments?.rabiesVaccinationDate}
                infectionVaccinationDate={info.treatments?.infectionVaccinationDate}
                ectoparasiteTreatmentDate={info.treatments?.ectoparasiteTreatmentDate}
            />
        );
    }

    if (activeTile === TileName.ANALYSES && info && info.analyses) {
        return (
            <AnalysesStep
                petId={info.id}
                isEditMode={false}
                petType={info.type}
                analyses={info.analyses}
                onClose={onTileBackHandler}
            />
        );
    }

    if (!info || isGlobalLoading) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.back} onClick={onClose}>
                    <BackAngularArrow />
                </div>
                <h2 className={styles.title}>{info.name.toUpperCase()}</h2>
            </div>
            <div className={styles.main}>
                <div className={styles.left}>
                    <div className={styles.leftItem}>
                        <div className={styles.leftItemTitle}>
                            <div className={styles.icon}>
                                <Blood />
                            </div>
                            <p className={styles.leftItemText}>Объем донации</p>
                        </div>
                        <div className={styles.leftItemValue}>
                            <div className={styles.volume}>~ {info.availableBloodAmount}</div>
                            <span className={styles.volumeDescr}>мл</span>
                        </div>
                    </div>
                    {!!info.donorRestrictions?.warnFactors?.length && (
                        <div
                            onClick={onStatusItemClickToggle}
                            className={cn(styles.leftItem, { [styles.status]: true })}
                        >
                            <div className={styles.leftItemTitle}>
                                <div className={styles.icon}>
                                    <StatusQuestion />
                                </div>
                                <p className={styles.leftItemText}>Статус</p>
                            </div>
                            <div className={styles.leftItemValue}>
                                <div className={styles.volume}>Вопросы к донорству</div>
                            </div>
                        </div>
                    )}
                    <div onClick={onChatClickHandler} className={cn(styles.leftItem, { [styles.chat]: true })}>
                        <div className={styles.ownerAvatar}>{info.ownerName.charAt(0).toUpperCase()}</div>
                        <div className={styles.ownerInfo}>
                            <p className={styles.ownerTitle}>Чат с хозяином</p>
                            <p className={styles.ownerName}>{info.ownerName}</p>
                        </div>
                        <div className={styles.arrowTileIcon}>
                            <AccordionArrow />
                        </div>
                    </div>
                </div>
                <div className={styles.right}>
                    <div className={styles.avatarWrapper}>
                        <img
                            alt={info.name}
                            className={styles.avatar}
                            src={info.photoUrls?.[0] || (info.type === PetType.DOG ? dogRoundStub : catRoundStub)}
                        />
                        <CircularProgress
                            showDot
                            size={180}
                            strokeWidth={15}
                            total={10}
                            color='var(--red10, #FF2727)'
                            current={10}
                        />
                        <div className={styles.avatarNameChar}>{info.name.charAt(0).toUpperCase()}</div>
                    </div>
                </div>
            </div>
            <div className={styles.tiles}>
                {tiles.map(({ name: tileName, title, icon }) => (
                    <div
                        key={tileName}
                        onClick={onTileClickHandler(tileName)}
                        className={cn(styles.tile, {
                            [styles.conditionTile]: tileName === TileName.CONDITIONS,
                            [styles.noActive]:
                                tileName === TileName.DONATIONS || (tileName === TileName.ANALYSES && !info.analyses),
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
                                    [styles.noActive]: !info.analyses,
                                })}
                            >
                                {Object.keys(info.analyses || []).length} из{' '}
                                {info.type === PetType.DOG ? dogAnalizesCount : catAnalizesCount}
                            </div>
                        )}
                        {tileName === TileName.CONDITIONS && (
                            <div className={styles.conditions}>
                                {info.compensationType === CompensationType.FOOD && (
                                    <div className={cn(styles.rewardFeedIcon, { [styles.button]: true })}>
                                        <div className={styles.icon}>
                                            <Bone />
                                        </div>
                                    </div>
                                )}
                                {info.compensationType === CompensationType.FREE && (
                                    <div className={cn(styles.rewardFreeIcon, { [styles.button]: true })}>
                                        <p className={styles.sum}>0</p>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                )}
                                {info.compensationType === CompensationType.PAID && (
                                    <div className={cn(styles.rewardNotFreeIcon, { [styles.button]: true })}>
                                        <p className={styles.descr}>₽</p>
                                    </div>
                                )}
                                {info.taxi && (
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
            {isLoading && (
                <div className={styles.loading}>
                    <Loading size={90} thickness={4} />
                </div>
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
                            <div onClick={onMessengerClickHandler} key={providerId} className={styles.identity}>
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

export default DonorForRecipient;
