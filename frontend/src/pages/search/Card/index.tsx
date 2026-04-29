import { Button } from '@mui/material';
import cn from 'classnames';
import { useBloodComponentsQuery, useLocationsQuery } from 'hooks/useDicts';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import BloodComponents from 'imgs/svg/bloodComponents';
import Cancel from 'imgs/svg/cancel';
import Location from 'imgs/svg/location';
import MiniPaw from 'imgs/svg/miniPaw';
import MiniSinglePaw from 'imgs/svg/miniSinglePaw';
import Pin from 'imgs/svg/pin';
import Accordion from 'pages/adding/common/Accordion';
import { FC, useCallback, useState } from 'react';
import { toast } from 'react-toastify';
import { getDateFormat } from 'utils/utils';

import { closeSearch } from 'api/apiServices/closeSearch';
import { GetPoolRequestResponse, Onboardings, PoolRequestStatus, RespondingDonorStatus } from 'api/bloodRequest';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import SearchOnboarding, { View } from '../Onboarding';
import CompletedDonations from './CompletedDonations';
import DonationComplete from './DonationComplete';
import DonationDetails from './DonationDetails';
import PlaningDonations from './PlaningDonations';
import styles from './SearchCard.module.less';
import SearchFinish from './SearchFinish';
import PrioritySearch from '../../../imgs/svg/prioritySearch';

type SelectedDonation = {
    id: string;
    status: RespondingDonorStatus;
};

type DonationCompletePage = {
    isOpen: boolean;
    volume?: number;
};

type Props = GetPoolRequestResponse & {
    name: string;
    petId: string;
    type: PetType;
    userId: string;
    avatar?: string;
    isLoading: boolean;
    bloodGroup: string;
    onClose: () => void;
    goToOwner: () => void;
    defaultOpenTab?: number;
    onBoarding?: Onboardings[];
    expireLimitWasShown: boolean;
    poolRequestRefetch: () => void;
};

const tabs = [
    { title: 'Детали', id: 0 },
    { title: 'Найдено', id: 1 },
    { title: 'Получено', id: 2 },
];

const SearchCard: FC<Props> = ({
    id,
    name,
    type,
    petId,
    userId,
    avatar,
    status,
    onClose,
    regions,
    goToOwner,
    isLoading,
    photoUrls,
    bloodGroup,
    onBoarding,
    description,
    prioritySearch,
    acceptedDonors,
    defaultOpenTab = 0,
    bloodGroupNames,
    bloodComponentIds,
    bloodVolumeNeeded,
    bloodVolumeDonated,
    completedDonations,
    poolRequestRefetch,
    bloodVolumeReserved,
    expireLimitWasShown,
    includeUnknownBloodGroup,
    smallPetsNotifyAllowed,
    createdAt = '',
}) => {
    const [tab, setTab] = useState<number>(defaultOpenTab);
    const [isSearchFinishPageOpen, setIsSearchFinishPageOpen] = useState<boolean>(false);
    const [selectedDonation, setSelectedDonation] = useState<SelectedDonation | null>(null);
    const [donationCompletePage, setDonationCompletePage] = useState<DonationCompletePage>({ isOpen: false });

    const { data: locationsDict = [] } = useLocationsQuery();
    const { data: bloodComponentsDict = [] } = useBloodComponentsQuery();

    const showToast = useCallback((text: string) => {
        toast.warn(text, {
            onClose: () => {
                // onClose();
            },
        });
    }, []);

    const onTabClickHandler = (tabId: number) => () => {
        setTab(tabId);
    };

    const onDonationClickHandler = (openDonationId: string, donationStatus: RespondingDonorStatus) => {
        setSelectedDonation({ id: openDonationId, status: donationStatus });
    };

    const onCloseDonationHandler = () => {
        setSelectedDonation(null);
    };

    const onDonationCompleteHandler = (volume: number) => {
        onCloseDonationHandler();

        setDonationCompletePage({ isOpen: true, volume });
    };

    const onIsSearchFinishHandler = async () => {
        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });

        onCloseDonationHandler();

        setIsSearchFinishPageOpen(true);
    };

    const onCloseSearchClickHandler = async () => {
        const response = await closeSearch(id);

        if (!response) {
            showToast('Не удалось завершить поиск');
        }

        setTab(0);
        setDonationCompletePage({ isOpen: false });

        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });

        poolRequestRefetch();
    };

    const onBackToSearchClickHandler = () => {
        setDonationCompletePage({ isOpen: false });

        poolRequestRefetch();
    };

    const onOpenDetailsClickHandler = () => {
        setTab(0);
        setIsSearchFinishPageOpen(false);

        poolRequestRefetch();
    };

    if (selectedDonation) {
        return (
            <DonationDetails
                userId={userId}
                onReject={poolRequestRefetch}
                donationId={selectedDonation.id}
                status={selectedDonation.status}
                onClose={onCloseDonationHandler}
                onIsSearchFinish={onIsSearchFinishHandler}
                onDonationComplete={onDonationCompleteHandler}
            />
        );
    }

    if (!(onBoarding || []).includes(Onboardings.BLOOD_CARD)) {
        return <SearchOnboarding view={View.CARD} id={id} onSucess={poolRequestRefetch} onBoarding={onBoarding} />;
    }

    if (donationCompletePage.isOpen && !!donationCompletePage.volume) {
        return (
            <DonationComplete
                type={type}
                avatar={avatar}
                bloodVolumeNeeded={bloodVolumeNeeded}
                onEndSearch={onCloseSearchClickHandler}
                onBackToSearch={onBackToSearchClickHandler}
                bloodVolumeDonated={donationCompletePage.volume || 44}
            />
        );
    }

    if (isSearchFinishPageOpen) {
        return (
            <SearchFinish
                type={type}
                avatar={avatar}
                onBackToOwner={goToOwner}
                bloodVolumeNeeded={bloodVolumeNeeded}
                bloodVolumeDonated={bloodVolumeNeeded}
                onOpenDetail={onOpenDetailsClickHandler}
            />
        );
    }

    return (
        <Layout>
            <div className={styles.wrapper}>
                <div className={styles.header}>
                    <div
                        className={styles.back}
                        onClick={expireLimitWasShown || status === PoolRequestStatus.CLOSED ? goToOwner : onClose}
                    >
                        <BackAngularArrow />
                    </div>
                    <div className={styles.nameWrapper}>
                        <h2 className={styles.name}>Поиск от {getDateFormat(new Date(createdAt))}</h2>
                        {!!prioritySearch && (
                            <div className={styles.prioritySearch}>
                                <PrioritySearch />
                            </div>
                        )}
                    </div>
                </div>
                <div className={styles.tabs}>
                    {(status !== PoolRequestStatus.CLOSED ? tabs : [tabs[0], tabs[2]]).map(({ title, id: tabId }) => (
                        <div
                            key={title}
                            onClick={onTabClickHandler(tabId)}
                            className={cn(styles.tab, {
                                [styles.active]: tab === tabId,
                                [styles.searchIsClosed]: status === PoolRequestStatus.CLOSED,
                                [styles.disabled]: tabId === 2 && !completedDonations?.length,
                            })}
                        >
                            {title}
                            {tabId === 1 &&
                                !!acceptedDonors?.filter(
                                    (resStatus) =>
                                        resStatus.status !== RespondingDonorStatus.CANCELED &&
                                        resStatus.status !== RespondingDonorStatus.REJECTED,
                                ).length && (
                                    <div className={styles.tabCounter}>
                                        {
                                            acceptedDonors?.filter(
                                                (resStatus) =>
                                                    resStatus.status !== RespondingDonorStatus.CANCELED &&
                                                    resStatus.status !== RespondingDonorStatus.REJECTED,
                                            ).length
                                        }
                                    </div>
                                )}
                            {tabId === 2 && !!completedDonations?.length && (
                                <div className={styles.tabCounter}>{completedDonations.length}</div>
                            )}
                        </div>
                    ))}
                </div>
                {tab === 0 && (
                    <>
                        <div className={styles.main}>
                            <div className={styles.left}>
                                <div className={cn(styles.leftItem, { [styles.name]: true })}>
                                    <div className={styles.icon}>
                                        <MiniSinglePaw />
                                    </div>
                                    <p className={styles.text}>{name.toUpperCase()}</p>
                                </div>
                                <div className={cn(styles.leftItem, { [styles.blood]: true })}>
                                    <div className={styles.subTitle}>
                                        <div className={styles.icon}>
                                            <Blood />
                                        </div>
                                        <p className={styles.text}>
                                            {status === PoolRequestStatus.CLOSED ? 'Искал' : 'Ищем'}
                                        </p>
                                    </div>
                                    <div className={styles.bloodInfo}>
                                        <div className={styles.bloodGroup}>{bloodGroup}</div>
                                        {(bloodGroupNames.some((group) => group !== bloodGroup) ||
                                            includeUnknownBloodGroup) &&
                                            ' +'}
                                        {bloodGroupNames.some((group) => group !== bloodGroup)
                                            ? bloodGroupNames
                                                  .filter((group) => group !== bloodGroup)
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
                                        {includeUnknownBloodGroup && (
                                            <div
                                                className={cn(styles.bloodGroup, {
                                                    [styles.needed]: true,
                                                })}
                                            >
                                                ?
                                            </div>
                                        )}
                                    </div>
                                </div>
                                <div className={cn(styles.leftItem, { [styles.location]: true })}>
                                    <div className={styles.icon}>
                                        <Location />
                                    </div>
                                    <p className={styles.text}>
                                        {regions
                                            ?.map(
                                                (lock) => locationsDict.filter(({ value }) => value === lock)[0]?.label,
                                            )
                                            .join(', ')}
                                    </p>
                                </div>
                            </div>
                            <div className={styles.right}>
                                <div className={styles.avatarWrapper}>
                                    <div className={styles.avatarWrapper}>
                                        <img
                                            alt={name}
                                            className={styles.avatar}
                                            src={avatar || (type === PetType.DOG ? dogRoundStub : catRoundStub)}
                                        />
                                        <CircularProgress
                                            showDot
                                            size={180}
                                            strokeWidth={15}
                                            showWhiteBackStroke
                                            total={bloodVolumeNeeded}
                                            color='var(--red10, #FF2727)'
                                            current={
                                                // eslint-disable-next-line no-nested-ternary
                                                status === PoolRequestStatus.CLOSED
                                                    ? bloodVolumeDonated > bloodVolumeNeeded
                                                        ? bloodVolumeNeeded
                                                        : bloodVolumeDonated
                                                    : bloodVolumeReserved || 0
                                            }
                                        />
                                        <div
                                            className={cn(styles.neededVolume, {
                                                [styles.siClosed]: status === PoolRequestStatus.CLOSED,
                                            })}
                                        >
                                            {bloodVolumeNeeded}
                                            <span>мл</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className={styles.infoItem}>
                            <div className={styles.infoItemTitle}>
                                <div className={styles.itemIcon}>
                                    <BloodComponents />
                                </div>
                                <div className={styles.text}>Компоненты крови</div>
                            </div>
                            <div className={styles.bloodComponents}>
                                {bloodComponentIds.map((comp) => (
                                    <div key={comp} className={styles.itemText}>
                                        {bloodComponentsDict.filter(({ value }) => value === comp)[0]?.label}
                                    </div>
                                ))}
                            </div>
                        </div>
                        <div className={styles.infoItem}>
                            <div className={cn(styles.infoItemTitle, { [styles.noMargin]: true })}>
                                <div className={cn(styles.itemIcon, { [styles.miniPaw]: true })}>
                                    <MiniPaw />
                                </div>
                                <div className={styles.text}>
                                    {smallPetsNotifyAllowed
                                        ? 'Уведомлять небольших доноров'
                                        : 'Не уведомлять небольших доноров'}
                                </div>
                            </div>
                        </div>
                        {(!!description?.length || !!photoUrls?.[0]) && (
                            <>
                                <Accordion icon={<Pin />} title='Дополнительная информация'>
                                    <div className={cn(styles.itemText, { [styles.accordionText]: true })}>
                                        {description}
                                    </div>
                                    {photoUrls?.[0] && (
                                        <img
                                            alt='Фото рецепиента'
                                            className={styles.bloodRequestPhoto}
                                            src={photoUrls[0]}
                                        />
                                    )}
                                </Accordion>
                            </>
                        )}
                        {status !== PoolRequestStatus.CLOSED ? (
                            <>
                                <Button
                                    fullWidth
                                    // onClick={onConfirmButtonClickHandler}
                                    className={cn(styles.confirm, { [styles.disabled]: true })}
                                >
                                    Расширить поиск
                                </Button>
                                <Button
                                    fullWidth
                                    startIcon={<Cancel />}
                                    className={styles.cancel}
                                    onClick={onCloseSearchClickHandler}
                                >
                                    {!!completedDonations?.length || !!acceptedDonors?.length
                                        ? 'Завершить поиск'
                                        : 'Отменить поиск'}
                                </Button>
                            </>
                        ) : (
                            <Button
                                fullWidth
                                disabled
                                // onClick={onConfirmButtonClickHandler}
                                className={cn(styles.confirm, { [styles.disabled]: true })}
                            >
                                Создать новый поиск
                            </Button>
                        )}
                    </>
                )}
                {tab === 1 && !!acceptedDonors?.length && (
                    <div className={styles.selected}>
                        <PlaningDonations donorResponses={acceptedDonors} onDonationClick={onDonationClickHandler} />
                    </div>
                )}
                {tab === 2 && !!completedDonations?.length && (
                    <div className={styles.selected}>
                        <CompletedDonations donations={completedDonations} />
                    </div>
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

export default SearchCard;
