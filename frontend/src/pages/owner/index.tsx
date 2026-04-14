import Button from '@mui/material/Button';
import cn from 'classnames';
import { useGetUserById } from 'hooks/useGetUserById';
import { usePetsQuery } from 'hooks/usePetsQuery';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import BloodFound from 'imgs/svg/bloodFound';
import BloodSearch from 'imgs/svg/bloodSearch';
import Bonus from 'imgs/svg/bonus';
import DonorButton from 'imgs/svg/donorButton';
import Pause from 'imgs/svg/pause';
import Paw from 'imgs/svg/paw';
import RecipientButton from 'imgs/svg/recipientButton';
import RoundCancel from 'imgs/svg/roundCancel';
import RoundQuestion from 'imgs/svg/roundQuestion';
import Settings from 'imgs/svg/settings';
import { FC, MouseEvent, useEffect, useLayoutEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { getCorrectDeclension, Variants } from 'utils/utils';

import { PlannedDonation } from 'api/donor';
import { DonorRestrictions, Pet } from 'api/pets';
import { queryClient } from 'api/queryClient';
import { Onboarding, Role } from 'api/user';
import Layout from 'components/Layout';
import Loading from 'components/Loading';
import PetProfile from 'components/Profiles/Pet';

import DonationDetails from './DonationDetails';
import DonorPreference, { View as DonorPreferenceView } from './DonorPreference';
import RecipientOnboarding from './Onboardings/Recipient';
import styles from './Owner.module.less';
import PlannedDonations from './PlannedDonations';
import DidNotRecover from './Statuses/DidNotRecover';
import DonationQuestions from './Statuses/DonationQuestions';
import NotReady from './Statuses/NotReady';

type Props = {
    userId: string;
};

type View = Role.DONOR | Role.RECIPIENT | Role.BLOOD_FOUND | Role.NONE;

type DonorStatus = {
    isOpen: boolean;
    status?: 'didNotRecover' | 'donationQuestions' | 'notReady';
};

type DonationDetailsType = {
    isOpen: boolean;
    donation?: PlannedDonation;
};

const tabs = [
    { title: 'Мои питомцы', ind: 0 },
    { title: 'Планируемые донации', ind: 1 },
];

const Owner: FC<Props> = ({ userId }) => {
    const navigate = useNavigate();

    const [tab, setTab] = useState(0);
    const [view, setView] = useState<View>(Role.RECIPIENT);
    const [selectedPet, setSelectedPet] = useState<Pet | null>(null);
    const [donorStatus, setDonorStatus] = useState<DonorStatus>({ isOpen: false });
    const [isPetProfileOpen, setIsPetProfileOpen] = useState<boolean>(false);
    const [isDonorPreferenceOpen, setIsDonorPreferenceOpen] = useState<boolean>(false);
    const [donationDetails, setDonationDetails] = useState<DonationDetailsType>({ isOpen: false });
    const [isDonorPreferenceOnboardingWasShown, setIsDonorPreferenceOnboardingWasShown] = useState<boolean>(false);

    const { data: pets, isLoading, refetch } = usePetsQuery(userId);
    const { data: userData, isLoading: isUserDataLoading, refetch: refetchUserData } = useGetUserById(userId);
    const userAvatarUrl = userData?.photoUrls?.[0];
    const userInitial = userData?.fullName?.charAt(0).toUpperCase() || '?';

    const onButtonClickHandler = (newView: View) => () => {
        if (newView === view) {
            return;
        }

        window.location.hash = `#${newView}`;

        setView(newView);
        setTab(0);

        if (newView === Role.DONOR && !userData?.donorPreference) {
            setIsDonorPreferenceOnboardingWasShown(false);
        }
    };

    const onPetProfileToggleHandler = (petData: Pet | null) => () => {
        setSelectedPet(petData);

        setIsPetProfileOpen((prevState) => !prevState);
    };

    const onAddPetClickHandler = () => {
        navigate('/adding');
    };

    const onTabClick = (tabId: number) => async () => {
        if (tabId === tab) {
            return;
        }

        if (tabId === 1) {
            await queryClient.invalidateQueries({ queryKey: ['plannedDonations', userId] });
        }

        setTab(tabId);
    };

    const onDonorPreferenceOnboardingCloseHandler = () => {
        setIsDonorPreferenceOnboardingWasShown(true);
    };

    const onDonorPreferenceCloseHandler = () => {
        setIsDonorPreferenceOpen(false);
    };

    const onRecipientLabelClickHandler =
        (label: 'startSearch' | 'activeSearch' | 'bloodFound', petId: string) =>
        async (e: MouseEvent<HTMLDivElement>) => {
            e.stopPropagation();

            if (label === 'startSearch') {
                navigate(`/adding#startSearch_${petId}`);
            }

            if (label === 'activeSearch') {
                navigate(`/search/${petId}`);
            }

            if (label === 'bloodFound') {
                await queryClient.invalidateQueries({ queryKey: ['poolRequestByPetId', petId] });

                navigate(`/search/${petId}?bloodFound=true`);
            }
        };

    const renderTabCounter = (count: number) => <div className={styles.tabCounter}>{count}</div>;

    const renderRecipientLabel = (petStatus: Role, petId: string) => {
        switch (true) {
            case petStatus === Role.RECIPIENT: {
                return (
                    <div
                        onClick={onRecipientLabelClickHandler('activeSearch', petId)}
                        className={cn(styles.label, { [styles.activeSearch]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <BloodSearch />
                        </div>
                        <div className={styles.labelText}>
                            Идет поиск <p className={styles.labelArrow}>⟶</p>
                        </div>
                    </div>
                );
            }
            case petStatus === Role.BLOOD_FOUND: {
                return (
                    <div
                        onClick={onRecipientLabelClickHandler('bloodFound', petId)}
                        className={cn(styles.label, { [styles.bloodFound]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <BloodFound />
                        </div>
                        <div className={styles.labelText}>Нашли кровь</div>
                    </div>
                );
            }
            case petStatus === Role.PLANNED_DONATION: {
                return (
                    <>
                        <div className={cn(styles.label, { [styles.pause]: true })}>
                            <div className={styles.statusLabelIcon}>
                                <Pause />
                            </div>
                            <div className={styles.labelText}>Планируется донация</div>
                        </div>
                    </>
                );
            }
            case petStatus === Role.NONE || petStatus === Role.DONOR || petStatus === Role.RECOVERING: {
                return (
                    <div
                        onClick={onRecipientLabelClickHandler('startSearch', petId)}
                        className={cn(styles.label, { [styles.startSearch]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <RecipientButton />
                        </div>
                        Начать поиск
                    </div>
                );
            }
            default: {
                return null;
            }
        }
    };

    const onDonorStatusCloseHandler = () => {
        setDonorStatus({ isOpen: false });
    };

    const onOpenPetProfileFromStatus = () => {
        setIsPetProfileOpen(true);
        onDonorStatusCloseHandler();
    };

    const onDonorLabelClickHandler =
        (label: 'didNotRecover' | 'notReady' | 'donationQuestions', petData: Pet) =>
        (e: MouseEvent<HTMLDivElement>) => {
            e.stopPropagation();

            setSelectedPet(petData);
            setDonorStatus({ isOpen: true, status: label });
        };

    const onNotPreferenceClickHandler = () => {
        setIsDonorPreferenceOpen(true);
    };

    const refetchUserDataClickHandler = () => {
        refetchUserData();

        setIsDonorPreferenceOpen(false);
        setIsDonorPreferenceOpen(false);
    };

    const onDonateBloodClickHandler = async (e: MouseEvent) => {
        e.stopPropagation();

        await queryClient.invalidateQueries({ queryKey: ['recipientsList', userId] });

        navigate('/recipientsList');
    };

    const onDonationClickHandler = (donation: PlannedDonation) => {
        setDonationDetails({ isOpen: true, donation });
    };

    const onDonationDetailsCloseHandler = () => {
        setDonationDetails({ isOpen: false });
    };

    const renderDonorLabel = (petData: Pet, donorRestrictions?: DonorRestrictions) => {
        switch (true) {
            case petData.petStatus === Role.RECOVERING: {
                return (
                    <div
                        onClick={onDonorLabelClickHandler('didNotRecover', petData)}
                        className={cn(styles.label, { [styles.didNotRecover]: true })}
                    >
                        <div className={styles.recover}>
                            <p className={styles.recoverDays}>{petData?.recoveryDays}</p>
                            <p className={styles.recoverDescr}>
                                {getCorrectDeclension(Variants.DAYS, petData.recoveryDays || 1)}
                            </p>
                        </div>
                        <div className={styles.labelText}>До восстановления</div>
                    </div>
                );
            }
            case petData.petStatus === Role.PLANNED_DONATION: {
                return (
                    <>
                        <div className={cn(styles.label, { [styles.pause]: true })}>
                            <div className={styles.statusLabelIcon}>
                                <Pause />
                            </div>
                            <div className={styles.labelText}>Планируется донация</div>
                        </div>
                    </>
                );
            }
            case !!donorRestrictions?.stopFactors?.length: {
                return (
                    <div
                        onClick={onDonorLabelClickHandler('notReady', petData)}
                        className={cn(styles.label, { [styles.notReady]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <RoundCancel />
                        </div>
                        <div className={styles.labelText}>Не готов к донации</div>
                    </div>
                );
            }
            case !!donorRestrictions?.warnFactors?.length: {
                return (
                    <div
                        onClick={onDonorLabelClickHandler('donationQuestions', petData)}
                        className={cn(styles.label, { [styles.donationQuestions]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <RoundQuestion />
                        </div>
                        <div className={styles.labelText}>Вопросы к донорству</div>
                    </div>
                );
            }
            default: {
                return null;
            }
        }
    };

    const renderSettingsTab = () => (
        <div
            className={cn(styles.pet, {
                [styles.notPreference]: !userData?.donorPreference,
                [styles.notCandidats]:
                    userData?.donorPreference &&
                    pets?.pets.every(({ donorRestrictions }) => donorRestrictions?.stopFactors?.length),
                [styles.isCandidats]:
                    userData?.donorPreference &&
                    pets?.pets.some(({ donorRestrictions }) => !donorRestrictions?.stopFactors?.length),
            })}
        >
            {!userData?.donorPreference && (
                <>
                    <p className={styles.settingTabText}>
                        Задайте
                        <br />
                        параметры
                        <br />
                        донорства
                    </p>
                    <Button
                        fullWidth
                        variant='contained'
                        className={styles.settingsButton}
                        onClick={onNotPreferenceClickHandler}
                    >
                        <div className={styles.settingsIcon}>
                            <Settings />
                        </div>
                        Настроить
                    </Button>
                </>
            )}
            {userData?.donorPreference &&
                pets?.pets.every(({ donorRestrictions }) => donorRestrictions?.stopFactors?.length) && (
                    <>
                        <div className={styles.notCandidatsButton} onClick={onNotPreferenceClickHandler}>
                            <div className={styles.preferencesettings}>
                                <Settings />
                            </div>
                        </div>
                        <h3 className={styles.notCandidatsTitle}>У вас нет доноров, готовых к донации</h3>
                        <p className={styles.notCandidatsDescr}>проверьте статус питомцев или добавьте новых</p>
                    </>
                )}
            {userData?.donorPreference &&
                pets?.pets.some(
                    ({ donorRestrictions, petStatus }) =>
                        !donorRestrictions?.stopFactors?.length || petStatus === Role.PLANNED_DONATION,
                ) && (
                    <>
                        <div className={styles.notCandidatsButton} onClick={onNotPreferenceClickHandler}>
                            <div className={styles.preferencesettings}>
                                <Settings />
                            </div>
                        </div>
                        <div>
                            <h3 className={styles.isCandidatsTitle}>Спасайте жизни - получайте награды </h3>
                            <Button
                                fullWidth
                                variant='contained'
                                onClick={onDonateBloodClickHandler}
                                className={styles.isCandidatsButton}
                            >
                                Сдать кровь
                                <div className={styles.candidatsIcon}>
                                    <BackAngularArrow />
                                </div>
                            </Button>
                        </div>
                    </>
                )}
        </div>
    );

    useLayoutEffect(() => {
        if (window.location.hash === '#recipient') {
            setView(Role.RECIPIENT);

            return;
        }

        if (window.location.hash === '#donor') {
            setView(Role.DONOR);

            return;
        }

        if (window.location.hash === '#donorDonations') {
            setTab(1);
            setView(Role.DONOR);

            return;
        }

        window.location.hash = '#recipient';
    }, []);

    useEffect(() => {
        if (!pets) {
            setSelectedPet(null);

            return;
        }

        setSelectedPet((prevState) => {
            if (prevState) {
                return pets?.pets.find((pet) => pet.id === prevState.id) || null;
            }

            return prevState;
        });
    }, [pets]);

    if (donorStatus.isOpen) {
        if (donorStatus.status === 'didNotRecover') {
            return <DidNotRecover onClose={onDonorStatusCloseHandler} onOpenPetProfile={onOpenPetProfileFromStatus} />;
        }

        if (donorStatus.status === 'notReady') {
            return (
                <NotReady
                    onClose={onDonorStatusCloseHandler}
                    onOpenPetProfile={onOpenPetProfileFromStatus}
                    factors={selectedPet?.donorRestrictions?.stopFactors}
                />
            );
        }

        if (donorStatus.status === 'donationQuestions') {
            return (
                <DonationQuestions
                    onClose={onDonorStatusCloseHandler}
                    onOpenPetProfile={onOpenPetProfileFromStatus}
                    factors={selectedPet?.donorRestrictions?.warnFactors}
                />
            );
        }
    }

    if (
        view === Role.DONOR &&
        !isUserDataLoading &&
        !userData?.donorPreference &&
        !isDonorPreferenceOnboardingWasShown
    ) {
        return (
            <DonorPreference
                id={userData?.id}
                view={DonorPreferenceView.ONBOARDING}
                refetchUserData={refetchUserDataClickHandler}
                onClose={onDonorPreferenceOnboardingCloseHandler}
            />
        );
    }

    if (view === Role.DONOR && isDonorPreferenceOpen) {
        return (
            <DonorPreference
                id={userData?.id}
                preference={userData?.donorPreference}
                view={DonorPreferenceView.PREFERENCE}
                onClose={onDonorPreferenceCloseHandler}
                refetchUserData={refetchUserDataClickHandler}
            />
        );
    }

    if (donationDetails.isOpen && donationDetails.donation) {
        return (
            <DonationDetails
                userId={userId}
                identities={userData?.identities}
                donation={donationDetails.donation}
                onClose={onDonationDetailsCloseHandler}
            />
        );
    }

    if (
        (view === Role.RECIPIENT || view === Role.BLOOD_FOUND) &&
        (!userData?.onBoarding || !userData?.onBoarding?.includes(Onboarding.FIND_BLOOD))
    ) {
        return (
            <RecipientOnboarding
                id={userId}
                onboardings={userData?.onBoarding}
                onConfirmButtonClick={refetchUserData}
            />
        );
    }

    if (isPetProfileOpen && selectedPet) {
        return <PetProfile updatePets={refetch} onClose={onPetProfileToggleHandler(null)} {...selectedPet} />;
    }

    return (
        <Layout>
            <div className={cn(styles.wrapper, { [styles.isPets]: !!pets?.pets.length })}>
                <div className={styles.header}>
                    <div className={styles.avatar} onClick={() => navigate('/profile')} role='button'>
                        {userAvatarUrl ? (
                            <img src={userAvatarUrl} alt='Фото профиля' className={styles.avatarImage} />
                        ) : (
                            userInitial
                        )}
                    </div>
                    <h1 className={styles.fullName}>{userData?.fullName || ''}</h1>
                    {view === 'donor' && (
                        <button type='button' className={styles.bonusCounter} onClick={() => navigate('/bonuses')}>
                            <span className={styles.bonusCounterIcon}>
                                <Bonus />
                            </span>
                            <span className={styles.bonusCounterValue}>0</span>
                        </button>
                    )}
                </div>
                {(isLoading || isUserDataLoading) && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
                {!(isLoading || isUserDataLoading) && !pets?.pets.length && (
                    <div className={styles.button} onClick={onAddPetClickHandler}>
                        <div className={styles.pawIcon}>
                            <Paw />
                        </div>
                        <p className={styles.pawButtonText}>Добавить питомца</p>
                    </div>
                )}
                {/* При переключении между вкладками перезапрашивать ли запросы на поиск крови??? */}
                {!(isLoading || isUserDataLoading) && !!pets?.pets.length && (
                    <>
                        {view === 'donor' && (
                            <div className={styles.tabs}>
                                {tabs.map(({ title, ind }) => (
                                    <div
                                        key={title}
                                        onClick={onTabClick(ind)}
                                        className={cn(styles.tab, { [styles.active]: tab === ind })}
                                    >
                                        {title}
                                        {ind === 0 && renderTabCounter(pets.totalPets)}
                                        {ind === 1 &&
                                            !!pets.totalPlannedDonations &&
                                            renderTabCounter(pets.totalPlannedDonations)}
                                    </div>
                                ))}
                            </div>
                        )}
                        {view === 'donor' && tab === 1 && (
                            <div className={styles.plannedDonations}>
                                <PlannedDonations onDonationClick={onDonationClickHandler} id={userId} />
                            </div>
                        )}
                        {tab === 0 && (
                            <div className={cn(styles.showcase, { [styles.donorView]: view === 'donor' })}>
                                {view === 'donor' && renderSettingsTab()}
                                {pets.pets.map((pet) => (
                                    <div key={`${pet.id}`} className={cn(styles.pet, { [styles[pet.type]]: true })}>
                                        <div className={styles.photo} onClick={onPetProfileToggleHandler(pet)}>
                                            {!!pet.photoUrls?.[0] && (
                                                <img className={styles.img} src={pet.photoUrls?.[0]} alt={pet.name} />
                                            )}
                                            <div className={styles.bloodGroup}>
                                                {pet.bloodGroup !== 'UNKNOWN' ? pet.bloodGroup : '?'}
                                            </div>
                                            <div className={styles.photoFooter}>
                                                <p className={styles.name}>{pet.name.toUpperCase()}</p>
                                                {(view === Role.RECIPIENT || view === Role.BLOOD_FOUND) &&
                                                    renderRecipientLabel(pet.petStatus, pet.id)}
                                                {(view === Role.DONOR || view === Role.NONE) &&
                                                    renderDonorLabel(pet, pet.donorRestrictions)}
                                            </div>
                                            <div className={styles.gradient} />
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                        <div className={styles.footer}>
                            <div
                                onClick={onButtonClickHandler(Role.DONOR)}
                                className={cn(styles.viewButton, { [styles.checked]: view === 'donor' })}
                            >
                                <div className={cn(styles.buttonIcon, { [styles.checked]: view === 'donor' })}>
                                    <DonorButton />
                                </div>
                                <p className={styles.buttonText}>Стать донором</p>
                            </div>
                            <div onClick={onAddPetClickHandler} className={styles.pawButton}>
                                <Paw />
                            </div>
                            <div
                                onClick={onButtonClickHandler(Role.RECIPIENT)}
                                className={cn(styles.viewButton, { [styles.checked]: view === 'recipient' })}
                            >
                                <div className={cn(styles.buttonIcon, { [styles.checked]: view === 'recipient' })}>
                                    <RecipientButton />
                                </div>
                                <p className={styles.buttonText}>Найти кровь</p>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </Layout>
    );
};

export default Owner;
