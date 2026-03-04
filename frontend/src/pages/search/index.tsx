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
import { TelegramUser } from 'types';

import { Onboardings } from 'api/bloodRequest';
import { Pet } from 'api/pets';
import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import SearchCard from './Card';
import DonorsShowcase from './DonorsShowcase';
import NoResults from './NoResults';
import SearchOnboarding, { View } from './Onboarding';
import styles from './Search.module.less';

type Props = {
    user: TelegramUser;
};

const tabs = [
    { title: 'Доноры', ind: 0 },
    { title: 'Пакеты крови', ind: 1 },
];

const Search: FC<Props> = ({ user }) => {
    const [searchParams] = useSearchParams();
    const { id } = useParams<{ id?: string; bloodFound?: string }>();

    const navigate = useNavigate();

    const isBloodFound = searchParams.get('bloodFound');
    const isPacketsBloodFound = false;

    const {
        data: pets = [],
        refetch: petsRefetch,
        isError: petsIsError,
        isLoading: petsIsLoading,
    } = usePetsQuery(user.id); // ?

    const [tab, setTab] = useState(0);
    const [selectedPet, setSelectedPet] = useState<Pet | null>(null);
    const [isSearchCardOpen, setIsSearchCardOpen] = useState<boolean>(false);

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

    const renderTabCounter = (count?: number) => {
        if (!count) {
            return null;
        }

        return <div className={styles.tabCounter}>{count}</div>;
    };

    useEffect(() => {
        setSelectedPet(pets?.find((pet) => pet.id === id) || null);
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

    // карточка поиска
    if (isSearchCardOpen && selectedPet && !poolRequestIsLoading && poolRequest) {
        return (
            <SearchCard
                {...poolRequest}
                petId={selectedPet.id}
                name={selectedPet.name}
                type={selectedPet.type}
                onClose={onCardOpenToggle}
                isLoading={poolRequestIsLoading}
                avatar={selectedPet.photoUrls?.[0]}
                bloodGroup={selectedPet.bloodGroup}
                onBoarding={poolRequest?.onBoarding}
                onBoardingConfirm={poolRequestRefetch}
            />
        );
    }

    // Онбординг
    if (!poolRequestIsLoading && !(poolRequest?.onBoarding || []).includes(Onboardings.SEARCH)) {
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
                                current={(poolRequest?.bloodVolumeNeeded || 30) / 2}
                            />
                        </div>
                        <div className={styles.info}>
                            <div className={styles.name}>{selectedPet?.name}</div>
                            <div className={styles.searchResult}>Найдено: 0 из {poolRequest?.bloodVolumeNeeded} мл</div>
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
                            {ind === 0 && isBloodFound && renderTabCounter(poolRequest?.responses?.length)}
                            {ind === 1 && isPacketsBloodFound && renderTabCounter()}
                        </div>
                    ))}
                </div>
                {tab === 0 && !isLoading && !isBloodFound && (
                    <NoResults tab={tab} suitableDonors={poolRequest?.suitableDonors} />
                )}
                {tab === 1 && !isLoading && !isPacketsBloodFound && <NoResults tab={tab} />}
                {tab === 0 && !isLoading && isBloodFound && (
                    <DonorsShowcase petType={selectedPet?.type} list={poolRequest?.responses} />
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
