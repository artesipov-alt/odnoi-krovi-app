import { Button } from '@mui/material';
import cn from 'classnames';
import { usePetsQuery } from 'hooks/usePetsQuery';
import bonusBg from 'imgs/bonusBg.png';
import emptyBg from 'imgs/emptyBg.png';
import fourPaws from 'imgs/fourPaws.png';
import ozon from 'imgs/ozon.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Copied from 'imgs/svg/copied';
import Drugs from 'imgs/svg/drugs';
import Feed from 'imgs/svg/feed';
import Other from 'imgs/svg/other';
import Priority from 'imgs/svg/priority';
import taily from 'imgs/taily.png';
import userPreferenseOnboarding from 'imgs/userPreferenseOnboarding.png';
import wb from 'imgs/wb.png';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import { getDateFormat, isExpiredDate } from 'utils/utils';

import { getAllBonuses } from 'api/apiServices/getAllBonuses';
import { BonusCategory, BonusType, GetAllBonusesResponse, PlatformName } from 'api/donor';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import styles from './Bonuses.module.less';

type BonusTab = BonusType | 'priority';

const tabs = [
    { key: 'priority', title: 'Приоритет', Icon: Priority },
    { key: BonusType.PREPARATION, title: 'Препараты', Icon: Drugs },
    { key: BonusType.FOOD, title: 'Корм', Icon: Feed },
    { key: BonusType.OTHER, title: 'Другое', Icon: Other },
];

type Props = {
    userId: string;
};

const Bonuses: FC<Props> = ({ userId }) => {
    const navigate = useNavigate();

    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [activeTab, setActiveTab] = useState<BonusTab>(BonusType.PREPARATION);
    const [isCurtainOPen, setIsCurtainOpen] = useState<boolean>(false);
    const [isCopiedPromo, setIsCopiedPromo] = useState<boolean>(false);
    const [selectedBonus, setSelectedBonus] = useState<BonusCategory | null>(null);
    const [bonuses, setBonuses] = useState<GetAllBonusesResponse | undefined>(undefined);

    const { data: pets, refetch } = usePetsQuery(userId);

    const showToast = useCallback(
        (text: string) => {
            toast.warn(text, {
                onClose: () => {
                    navigate('/owner#donor');
                },
            });
        },
        [navigate],
    );

    const fetchBonuses = useCallback(async () => {
        const response = await getAllBonuses(userId);

        if (!response || !response.data) {
            showToast('Не удалось получить бонусы');

            return;
        }

        setBonuses(response.data);
        setIsLoading(false);
    }, [showToast, userId]);

    const goToBackClickHandler = () => {
        navigate('/owner#donor');
    };

    const onCloseCurtainHandler = () => {
        setSelectedBonus(null);
        setIsCurtainOpen(false);
        setIsCopiedPromo(false);
    };

    const onToShopClickHandler = () => {
        if (!selectedBonus) {
            return;
        }

        window.open(selectedBonus.platformUrl);

        onCloseCurtainHandler();
    };

    const onCopiClickHandler = async () => {
        if (!selectedBonus?.promoCode) {
            return;
        }

        try {
            await navigator.clipboard.writeText(selectedBonus?.promoCode);

            setIsCopiedPromo(true);
        } catch (err) {
            // error
        }
    };

    const onBonusClickHandler = (bonus: BonusCategory) => async () => {
        try {
            await navigator.clipboard.writeText(bonus.promoCode);
        } catch (err) {
            // error
        }

        setIsCurtainOpen(true);
        setSelectedBonus(bonus);
    };

    const renderContent = () => {
        // if (activeTab === 'priority') {
        //     return (
        //         <div className={styles.emptyBlock}>
        //             <div className={styles.emptyTitle}>Приоритетный поиск!</div>
        //             <div className={styles.emptyText}>
        //                 Начисляется за донации и&nbsp;помогает&nbsp;быстрее найти кровь
        //             </div>
        //             <button type='button' className={styles.primaryButton}>
        //                 Запланировать донацию
        //             </button>
        //             <img src={profilePhoto} alt='Питомцы' className={styles.emptyImage} />
        //         </div>
        //     );
        // }

        if (!bonuses?.[activeTab].length) {
            return (
                <div className={styles.emptyBlock}>
                    <div className={styles.emptyTitle}>Здесь пока пусто...</div>
                    <button onClick={goToBackClickHandler} type='button' className={styles.primaryButton}>
                        Запланировать донацию
                    </button>
                    <img src={emptyBg} alt='Питомцы' className={styles.emptyImage} />
                </div>
            );
        }

        return (
            <div className={styles.list}>
                {bonuses[activeTab]
                    .sort((a, b) => new Date(b.expiresAt).getTime() - new Date(a.expiresAt).getTime())
                    .map((item) => {
                        const isExpired = isExpiredDate(item.expiresAt);

                        return (
                            <div
                                key={item.id}
                                onClick={onBonusClickHandler(item)}
                                className={cn(styles.card, { [styles.isExpired]: isExpired })}
                            >
                                <div className={styles.cardHeader}>
                                    <p className={styles.cardTitle}>&#34;{item.partnerName}&#34;</p>
                                    <p className={cn(styles.expirationDate, { [styles.isExpired]: isExpired })}>
                                        {isExpired ? 'истек' : 'до'} {getDateFormat(item.expiresAt)}
                                    </p>
                                </div>
                                <p className={styles.cardDescr}>{item.description}</p>
                            </div>
                        );
                    })}
            </div>
        );
    };

    useEffect(() => {
        fetchBonuses();
    }, [fetchBonuses]);

    if (isLoading) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.page} style={{ backgroundImage: `url(${userPreferenseOnboarding})` }}>
                <div className={styles.innerWrapper}>
                    <div className={styles.header}>
                        <button type='button' className={styles.backButton} onClick={goToBackClickHandler}>
                            <BackAngularArrow />
                        </button>
                        <h1 className={styles.title}>Мои бонусы</h1>
                    </div>
                    <div className={styles.panel}>
                        <div className={styles.tabs}>
                            {tabs.map(({ title, Icon, key }, i) => (
                                <button
                                    key={key}
                                    type='button'
                                    className={cn(styles.tab, {
                                        [styles.isFirst]: i === 0,
                                        [styles.tab_active]: activeTab === key,
                                    })}
                                    onClick={() => setActiveTab(key as BonusTab)}
                                >
                                    <span className={styles.tabTop}>
                                        <span className={styles.counterDot}>
                                            <Icon />
                                        </span>
                                        <span className={styles.counterValue}>
                                            {i === 0 ? pets?.totalPrioritySearch || 0 : bonuses?.[key]?.length || 0}
                                        </span>
                                    </span>
                                    <span className={styles.tabLabel}>{title}</span>
                                </button>
                            ))}
                        </div>
                    </div>
                </div>
                <div className={styles.content}>
                    {renderContent()}
                    {!!bonuses?.[activeTab].length && (
                        <div className={styles.bottomAction} onClick={goToBackClickHandler}>
                            <button type='button' className={styles.primaryButton}>
                                Запланировать донацию
                            </button>
                        </div>
                    )}
                </div>
            </div>
            {isCurtainOPen && (
                <Curtain
                    noRednerButtons
                    contentBorderRadius={0}
                    contentOverflow='visible'
                    shouldCloseByWrapperClick
                    backgroundImage={bonusBg}
                    onClose={onCloseCurtainHandler}
                    title={<span style={{ display: 'none' }} />}
                >
                    <div className={styles.popupShell}>
                        <div className={styles.popupCard}>
                            <h3 className={styles.popupTitle}>Бонус за донацию</h3>
                            <p className={styles.popupText}>{selectedBonus?.description}</p>
                            <div className={styles.popupDivider} />
                            <p className={styles.promoText}>
                                {isCopiedPromo ? 'Промокод скопирован' : 'Промокод уже скопирован'}
                            </p>
                            <div className={styles.promo}>
                                <p className={styles.promoCode}>{selectedBonus?.promoCode}</p>
                                <div onClick={onCopiClickHandler} className={styles.copiedIcon}>
                                    <Copied />
                                </div>
                            </div>
                            <Button
                                onClick={onToShopClickHandler}
                                className={cn(styles.toShop, {
                                    [styles.isIcon]:
                                        selectedBonus?.platformName.toLowerCase() === PlatformName.OZON ||
                                        selectedBonus?.platformName.toLowerCase() === PlatformName.TAILY_PLATFORM ||
                                        selectedBonus?.platformName.toLowerCase() === PlatformName.FOUR_PAWS ||
                                        selectedBonus?.platformName.toLowerCase() === PlatformName.WB,
                                })}
                            >
                                Применить в магазине
                                {selectedBonus?.platformName.toLowerCase() === PlatformName.OZON && (
                                    <img className={styles.logo} alt={PlatformName.OZON} src={ozon} />
                                )}
                                {selectedBonus?.platformName.toLowerCase() === PlatformName.TAILY_PLATFORM && (
                                    <img className={styles.logo} alt={PlatformName.TAILY_PLATFORM} src={taily} />
                                )}
                                {selectedBonus?.platformName.toLowerCase() === PlatformName.FOUR_PAWS && (
                                    <img className={styles.logo} alt={PlatformName.FOUR_PAWS} src={fourPaws} />
                                )}
                                {selectedBonus?.platformName.toLowerCase() === PlatformName.WB && (
                                    <img className={styles.logo} alt={PlatformName.WB} src={wb} />
                                )}
                            </Button>
                        </div>
                    </div>
                </Curtain>
            )}
        </Layout>
    );
};

export default Bonuses;
