import { Button } from '@mui/material';
import cn from 'classnames';
import { useGetRecipientsList } from 'hooks/useGetRecipientsList';
import { useGetUserById } from 'hooks/useGetUserById';
import clinicsBg from 'imgs/clinicsBg.png';
import noDonorBg from 'imgs/noDonorBg.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import CrossedEye from 'imgs/svg/crossedEye';
import Eye from 'imgs/svg/eye';
import PrioritySearch from 'imgs/svg/prioritySearch';
import { FC, useCallback, useState } from 'react';
import { useNavigate } from 'react-router';

import { Onboarding } from 'api/user';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import RecipientsListDetail from './Detail';
import RecipientsListOnboarding from './Onboarding';
import styles from './RecipientsList.module.less';

type Props = {
    userId: string;
};

const tabs = [
    { title: 'Реципиенты', ind: 0 },
    { title: 'Клиники', ind: 1 },
];

const RecipientsList: FC<Props> = ({ userId }) => {
    const navigate = useNavigate();

    const [tab, setTab] = useState(0);
    const [isBlur, setIsBlur] = useState(true);
    const [checkedBloadSearchId, setCheckedBloadSearchId] = useState<string | null>(null);

    const { data: user, isLoading: isUserLoading, refetch: refetchUser } = useGetUserById(userId);
    const { data: list, isLoading: isListLoading, refetch: refetchList } = useGetRecipientsList(userId);

    const goToOwner = useCallback(() => {
        navigate('/owner#donor');
    }, [navigate]);

    const onTabClick = (tabId: number) => () => {
        setTab(tabId);
    };

    const onPetClickHandler = (petId: string) => () => {
        setCheckedBloadSearchId(petId);
    };

    const onSettingsButtonClickHandler = (type: 'blur' | 'show') => () => {
        if ((type === 'blur' && isBlur) || (type === 'show' && !isBlur)) {
            return;
        }

        setIsBlur(type === 'blur');
    };

    const onDetailCloseHandler = () => {
        setCheckedBloadSearchId(null);
    };

    if (isUserLoading || isListLoading) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    if (!user?.onBoarding?.includes(Onboarding.RECIPIENT_LIST)) {
        return <RecipientsListOnboarding onBoarding={user?.onBoarding} id={userId} refetch={refetchUser} />;
    }

    if (checkedBloadSearchId) {
        return (
            <RecipientsListDetail
                userId={userId}
                isBlurByDefault={isBlur}
                id={checkedBloadSearchId}
                onClose={onDetailCloseHandler}
            />
        );
    }

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.back} onClick={goToOwner}>
                    <BackAngularArrow />
                </div>
                <div className={styles.tabs}>
                    {tabs.map(({ title, ind }) => (
                        <div
                            key={title}
                            onClick={onTabClick(ind)}
                            className={cn(styles.tab, { [styles.active]: tab === ind })}
                        >
                            {title}
                        </div>
                    ))}
                </div>
            </div>
            {tab === 0 && !!list?.items.length && (
                <div className={styles.settings}>
                    <p className={styles.settingsText}>Осторожно! Возможен травмирующий контент</p>
                    <div className={styles.settingsButtons}>
                        <div
                            onClick={onSettingsButtonClickHandler('blur')}
                            className={cn(styles.settingsButton, { [styles.blur]: true, [styles.active]: isBlur })}
                        >
                            <div className={styles.settingsIcon}>
                                <CrossedEye />
                            </div>
                            <p>Размыть</p>
                        </div>
                        <div
                            onClick={onSettingsButtonClickHandler('show')}
                            className={cn(styles.settingsButton, { [styles.show]: true, [styles.active]: !isBlur })}
                        >
                            <div className={styles.settingsIcon}>
                                <Eye />
                            </div>
                            <p>Показать</p>
                        </div>
                    </div>
                </div>
            )}
            <div
                className={cn(styles.content, {
                    [styles.clinic]: tab === 1,
                    [styles.short]: tab === 0 && !list?.items.length,
                })}
            >
                {tab === 0 && !list?.items.length && (
                    <div>
                        <h1 className={styles.noListTitle}>Пока нет подходящих реципиентов</h1>
                        <p className={styles.noListDescr}>Когда найдем - направим уведомление</p>
                        <img className={styles.noListImg} src={noDonorBg} alt='no-recipients' />
                    </div>
                )}
                {tab === 1 && (
                    <div>
                        <h1 className={cn(styles.noListTitle, { [styles.isClinics]: true })}>Скоро будет доступно!</h1>
                        <p className={cn(styles.noListDescr, { [styles.isClinics]: true })}>
                            Функционал пока в разработке.
                            <br />
                            Но Вы можете найти и помочь питомцу в беде!
                        </p>
                        <div className={styles.clinicButton}>
                            <Button fullWidth onClick={onTabClick(0)} className={styles.confirm}>
                                Искать реципиентов
                            </Button>
                        </div>
                        <img className={styles.clinkcsImg} src={clinicsBg} alt='no-recipients' />
                    </div>
                )}
                <div className={styles.showcase}>
                    {tab === 0 &&
                        !!list?.items.length &&
                        list.items.map((pet) => (
                            <div
                                key={`${pet.petId}`}
                                className={cn(styles.pet, { [styles[pet.petType]]: !pet?.photoUrls?.[0] })}
                            >
                                <div className={styles.photo} onClick={onPetClickHandler(pet.id)}>
                                    {!!pet?.photoUrls?.[0] && (
                                        <img
                                            alt={pet.petName}
                                            src={pet.photoUrls?.[0]}
                                            className={cn(styles.img, { [styles.blured]: isBlur })}
                                        />
                                    )}
                                    <div className={styles.info}>
                                        <div className={styles.bloodGroup}>{pet.bloodGroupName}</div>
                                        {!!pet.prioritySearch && (
                                            <div className={styles.prioritySearch}>
                                                <PrioritySearch />
                                            </div>
                                        )}
                                    </div>
                                    <div className={styles.bloodVolume}>
                                        <p className={styles.bloodVolumeNumber}>{pet.bloodVolumeRemaining}</p>
                                        <p className={styles.bloodVolumeDescr}>мл</p>
                                    </div>
                                    <div className={styles.photoFooter}>
                                        <p className={styles.name}>{pet.petName.toUpperCase()}</p>
                                        {!!pet.matchingDonors?.length && (
                                            <div className={styles.matchingDonors}>{pet?.matchingDonors?.length}</div>
                                        )}
                                    </div>
                                    <div className={styles.gradient} />
                                </div>
                            </div>
                        ))}
                </div>
            </div>
        </Layout>
    );
};

export default RecipientsList;
