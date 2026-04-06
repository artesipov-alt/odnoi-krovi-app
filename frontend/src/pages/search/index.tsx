import cn from 'classnames';
import { useGetPoolRequestByPetId } from 'hooks/useGetPoolRequest';
import { usePetsQuery } from 'hooks/usePetsQuery';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { useParams, useSearchParams } from 'react-router-dom';
import { toast } from 'react-toastify';

import { Onboardings, PoolRequestStatus } from 'api/bloodRequest';
import { Pet } from 'api/pets';
import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';
import Loading from 'components/Loading';
import DonorForRecipient from 'components/Profiles/DonorForRecipient';

import DonationQuestions from '../owner/Statuses/DonationQuestions';
import SearchCard from './Card';
import DonorsShowcase from './DonorsShowcase';
import LimitReached from './LimitReached';
import NoResults from './NoResults';
import SearchOnboarding, { View } from './Onboarding';
import styles from './Search.module.less';

type DonorDetails = {
    id?: string;
    isOpen: boolean;
};

type Props = {
    userId: string;
};

const tabs = [
    { title: 'Доноры', ind: 0 },
    { title: 'Пакеты крови', ind: 1 },
];

const Search: FC<Props> = ({ userId }) => {
    const [searchParams] = useSearchParams();
    const { id } = useParams<{ id?: string; bloodFound?: string }>();

    const navigate = useNavigate();

    const isBloodFound = searchParams.get('bloodFound');
    const isPacketsBloodFound = false;

    const { data: pets, refetch: petsRefetch, isError: petsIsError, isLoading: petsIsLoading } = usePetsQuery(userId); // ?

    const [tab, setTab] = useState(0);
    const [cardDefaultTab, setCardDefaultTab] = useState(0);
    const [showStartView, setShowStartView] = useState(true);
    const [selectedPet, setSelectedPet] = useState<Pet | null>(null);
    const [donorDetails, setDonorDetails] = useState<DonorDetails>({ isOpen: false });
    const [isSearchCardOpen, setIsSearchCardOpen] = useState<boolean>(false);
    const [expireLimitWasShown, setExpireLimitWasShown] = useState<boolean>(false);
    const [donorWarnFactors, setDonorWarnFactors] = useState<DonorDetails>({ isOpen: false });

    const {
        data: poolRequest,
        isError: poolRequestIsError,
        refetch: poolRequestRefetch,
        isLoading: poolRequestIsLoading,
    } = useGetPoolRequestByPetId(selectedPet?.id);

    const isLoading = petsIsLoading || poolRequestIsLoading;

    const goToOwner = useCallback(() => {
        navigate('/owner');
    }, [navigate]);

    const showToast = useCallback(
        (text: string) => {
            toast.warn(text, {
                onClose: () => {
                    setIsSearchCardOpen((prevState) => !prevState);

                    goToOwner();
                },
            });
        },
        [goToOwner],
    );

    const onTabClick = (tabId: number) => () => {
        setTab(tabId);
    };

    const onCardOpenToggle = () => {
        setIsSearchCardOpen((prevState) => !prevState);

        poolRequestRefetch();
    };

    const onOpenCardFromDonorRespond = () => {
        setCardDefaultTab(1);
        setDonorDetails({ isOpen: false });

        onCardOpenToggle();
    };

    const onOpenCardFromLimitReached = () => {
        setExpireLimitWasShown(true);

        onOpenCardFromDonorRespond();
    };

    const onOpenWarnFactorsToggle = (donorId?: string) => {
        setDonorWarnFactors((prevState) => ({
            id: donorId,
            isOpen: !prevState.isOpen,
        }));
    };

    const onDonorToggle = (donorId?: string) => {
        setDonorDetails((prevState) => ({
            id: donorId,
            isOpen: !prevState.isOpen,
        }));
    };

    const setIsStartViewShownHandler = () => {
        setShowStartView(false);
    };

    const renderTabCounter = (count?: number) => {
        if (!count) {
            return null;
        }

        return <div className={styles.tabCounter}>{count}</div>;
    };

    useEffect(() => {
        setSelectedPet(pets?.pets.find((pet) => pet.id === id) || null);
    }, [id, pets]);

    useEffect(() => {
        if (petsIsError) {
            showToast('Не удалось загрузить данные питомца');
        }
    }, [petsIsError, showToast]);

    useEffect(() => {
        if (poolRequestIsError) {
            showToast('Не удалось загрузить данные по запросу крови');
        }
    }, [poolRequestIsError, showToast]);

    if (donorDetails.isOpen && donorDetails.id) {
        return (
            <DonorForRecipient
                userId={userId}
                onClose={onDonorToggle}
                donorId={donorDetails.id}
                onBackToSearch={onOpenCardFromDonorRespond}
                responseId={poolRequest?.responses?.find(({ donorId }) => donorId === donorDetails.id)?.id!}
            />
        );
    }

    if (donorWarnFactors.isOpen && !!poolRequest) {
        return (
            <DonationQuestions
                onClose={onOpenWarnFactorsToggle}
                factors={poolRequest.responses?.find(({ donorId }) => donorId === donorWarnFactors?.id)?.warnFactors}
            />
        );
    }

    // карточка поиска
    if (isSearchCardOpen && selectedPet && !poolRequestIsLoading && poolRequest) {
        return (
            <SearchCard
                {...poolRequest}
                userId={userId}
                goToOwner={goToOwner}
                petId={selectedPet.id}
                name={selectedPet.name}
                type={selectedPet.type}
                onClose={onCardOpenToggle}
                defaultOpenTab={cardDefaultTab}
                isLoading={poolRequestIsLoading}
                avatar={selectedPet.photoUrls?.[0]}
                bloodGroup={selectedPet.bloodGroup}
                onBoarding={poolRequest?.onBoarding}
                poolRequestRefetch={poolRequestRefetch}
                expireLimitWasShown={expireLimitWasShown}
            />
        );
    }

    if (
        !expireLimitWasShown &&
        !!poolRequest &&
        (poolRequest.bloodVolumeReserved || 0) >= (poolRequest.bloodVolumeNeeded || 0)
    ) {
        return (
            <LimitReached
                type={selectedPet?.type!}
                avatar={selectedPet?.photoUrls?.[0]}
                onBackToSearch={onOpenCardFromLimitReached}
                bloodVolumeNeeded={poolRequest?.bloodVolumeNeeded!}
            />
        );
    }

    // Онбординг
    if (!!poolRequest && !(poolRequest.onBoarding || []).includes(Onboardings.SEARCH)) {
        return (
            <SearchOnboarding
                view={View.SEARCH}
                id={poolRequest?.id}
                onSucess={poolRequestRefetch}
                onBoarding={poolRequest?.onBoarding}
            />
        );
    }

    return (
        <Layout>
            <div className={styles.wrapper}>
                <div className={styles.header}>
                    <div className={styles.back} onClick={goToOwner}>
                        <BackAngularArrow />
                    </div>
                    <div className={styles.main}>
                        <div className={styles.avatarWrapper}>
                            <img
                                alt='no results'
                                className={styles.avatar}
                                src={
                                    selectedPet?.photoUrls?.[0] ||
                                    (selectedPet?.type === PetType.DOG ? dogRoundStub : catRoundStub)
                                }
                            />
                            <CircularProgress
                                size={45}
                                strokeWidth={2}
                                color='var(--red10, #FF2727)'
                                total={poolRequest?.bloodVolumeNeeded || 0}
                                current={poolRequest?.bloodVolumeReserved || 0}
                            />
                        </div>
                        <div className={styles.info}>
                            <div className={styles.name}>{selectedPet?.name}</div>
                            <div className={styles.searchResult}>
                                Найдено: {poolRequest?.bloodVolumeReserved || 0} из {poolRequest?.bloodVolumeNeeded} мл
                            </div>
                        </div>
                        <div className={styles.searchDetails} onClick={onCardOpenToggle}>
                            Детали поиска
                        </div>
                    </div>
                </div>
                {/* При переключении между вкладками перезапрашивать ли запросы на поиск крови??? */}
                <div className={styles.tabs}>
                    {tabs.map(({ title, ind }) => (
                        <div
                            key={title}
                            onClick={onTabClick(ind)}
                            className={cn(styles.tab, { [styles.active]: tab === ind })}
                        >
                            {title}
                            {ind === 0 &&
                                isBloodFound &&
                                poolRequest?.responses?.length &&
                                renderTabCounter(poolRequest.responses.length)}
                            {ind === 1 && isPacketsBloodFound && renderTabCounter()}
                        </div>
                    ))}
                </div>
                {tab === 0 && !isLoading && !poolRequest?.responses && (
                    <NoResults tab={tab} suitableDonors={poolRequest?.suitableDonors} />
                )}
                {tab === 1 && !isLoading && !isPacketsBloodFound && <NoResults tab={tab} />}
                {tab === 0 && !isLoading && !!poolRequest?.responses && (
                    <DonorsShowcase
                        userId={userId}
                        goToOwner={goToOwner}
                        searchId={poolRequest.id}
                        petType={selectedPet?.type}
                        onDonorClick={onDonorToggle}
                        list={poolRequest.responses}
                        showStartView={showStartView}
                        onOpenWarnFactors={onOpenWarnFactorsToggle}
                        setIsStartViewShown={setIsStartViewShownHandler}
                    />
                )}
                {isLoading && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default Search;
