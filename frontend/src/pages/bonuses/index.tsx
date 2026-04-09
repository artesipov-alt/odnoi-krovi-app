import cn from 'classnames';
import cardBonus from 'imgs/cardBonus.png';
import emptyBg from 'imgs/emptyBg.png';
import profilePhoto from 'imgs/profilePhoto.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Drugs from 'imgs/svg/drugs';
import Feed from 'imgs/svg/feed';
import Other from 'imgs/svg/other';
import Priority from 'imgs/svg/priority';
import userPreferenseOnboarding from 'imgs/userPreferenseOnboarding.png';
import { FC, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import Layout from 'components/Layout';

import styles from './Bonuses.module.less';

type BonusTab = 'priority' | 'drugs' | 'food' | 'other';

type BonusItem = {
    id: string;
    title: string;
    subtitle: string;
    statusLabel: string;
    image?: string;
    isDisabled?: boolean;
};

const tabs: Array<{ key: BonusTab; title: string; count: number; Icon: FC }> = [
    { key: 'priority', title: 'Приоритет', count: 4, Icon: Priority },
    { key: 'drugs', title: 'Препараты', count: 3, Icon: Drugs },
    { key: 'food', title: 'Корм', count: 3, Icon: Feed },
    { key: 'other', title: 'Другое', count: 0, Icon: Other },
];

const bonusByTab: Record<BonusTab, BonusItem[]> = {
    priority: [],
    drugs: [],
    // drugs: [
    //     {
    //         id: 'drug-1',
    //         title: '“Вемелкам”',
    //         subtitle: 'Открыть QR-код',
    //         statusLabel: 'бессрочно',
    //         image: cardBonus,
    //     },
    //     {
    //         id: 'drug-2',
    //         title: '“Ветом”',
    //         subtitle: 'Открыть QR-код',
    //         statusLabel: 'до 18.05.2028',
    //     },
    //     {
    //         id: 'drug-3',
    //         title: '“Ветом”',
    //         subtitle: 'Открыть QR-код',
    //         statusLabel: 'до 28.05.2028',
    //     },
    //     {
    //         id: 'drug-4',
    //         title: '“Ветом”',
    //         subtitle: 'Барсик - донация от 21.102025',
    //         statusLabel: 'использован',
    //         isDisabled: true,
    //     },
    //     {
    //         id: 'drug-5',
    //         title: '“Ветом”',
    //         subtitle: 'Барсик - донация от 21.102025',
    //         statusLabel: 'истек 28.05.2025',
    //         isDisabled: true,
    //     },
    // ],
    food: [],
    // food: [
    //     {
    //         id: 'food-1',
    //         title: 'ProPlan с индейкой',
    //         subtitle: 'Открыть QR-код',
    //         statusLabel: 'бессрочно',
    //         image: cardBonus,
    //     },
    //     {
    //         id: 'food-2',
    //         title: 'ProPlan с индейкой',
    //         subtitle: 'Открыть промокод',
    //         statusLabel: 'бессрочно',
    //         image: cardBonus,
    //     },
    //     {
    //         id: 'food-3',
    //         title: 'ProPlan для стерилизованных кошек',
    //         subtitle: 'Открыть промокод',
    //         statusLabel: 'бессрочно',
    //         image: cardBonus,
    //     },
    // ],
    other: [],
};

const Bonuses: FC = () => {
    const navigate = useNavigate();
    const [activeTab, setActiveTab] = useState<BonusTab>('drugs');

    const activeItems = useMemo(() => bonusByTab[activeTab], [activeTab]);
    const validItems = useMemo(
        () => activeItems.filter((item) => item.title.trim() && item.subtitle.trim() && item.statusLabel.trim()),
        [activeItems],
    );
    const hasListContent = activeTab !== 'priority' && validItems.length > 0;

    const renderContent = () => {
        if (activeTab === 'priority') {
            return (
                <div className={styles.emptyBlock}>
                    <div className={styles.emptyTitle}>Приоритетный поиск!</div>
                    <div className={styles.emptyText}>
                        Начисляется за донации и&nbsp;помогает&nbsp;быстрее найти кровь
                    </div>
                    <button type='button' className={styles.primaryButton}>
                        Запланировать донацию
                    </button>
                    <img src={profilePhoto} alt='Питомцы' className={styles.emptyImage} />
                </div>
            );
        }

        if (!validItems.length) {
            return (
                <div className={styles.emptyBlock}>
                    <div className={styles.emptyTitle}>Здесь пока пусто...</div>
                    <button type='button' className={styles.primaryButton}>
                        Запланировать донацию
                    </button>
                    <img src={emptyBg} alt='Питомцы' className={styles.emptyImage} />
                </div>
            );
        }

        return (
            <div className={styles.list}>
                {validItems.map((item) => (
                    <div
                        key={item.id}
                        className={cn(styles.card, {
                            [styles.card_disabled]: item.isDisabled,
                            [styles.card_withImage]: !!item.image,
                            [styles.card_noImage]: !item.image,
                        })}
                    >
                        {item.image && <img src={item.image || cardBonus} alt='' className={styles.cardImage} />}
                        <div className={styles.cardTitle}>{item.title}</div>
                        <div className={styles.cardSubtitle}>{item.subtitle}</div>
                        <span className={cn(styles.cardChip, { [styles.cardChip_disabled]: item.isDisabled })}>
                            {item.statusLabel}
                        </span>
                    </div>
                ))}
            </div>
        );
    };

    return (
        <Layout>
            <div className={styles.page} style={{ backgroundImage: `url(${userPreferenseOnboarding})` }}>
                <div className={styles.header}>
                    <button type='button' className={styles.backButton} onClick={() => navigate(-1)}>
                        <BackAngularArrow />
                    </button>
                    <h1 className={styles.title}>Мои бонусы</h1>
                </div>

                <div className={styles.panel}>
                    <div className={styles.tabs}>
                        {tabs.map((tab) => (
                            <button
                                key={tab.key}
                                type='button'
                                className={cn(styles.tab, { [styles.tab_active]: activeTab === tab.key })}
                                onClick={() => setActiveTab(tab.key)}
                            >
                                <span className={styles.tabTop}>
                                    <span className={styles.counterDot}>
                                        <tab.Icon />
                                    </span>
                                    <span className={styles.counterValue}>{tab.count}</span>
                                </span>
                                <span className={styles.tabLabel}>{tab.title}</span>
                            </button>
                        ))}
                    </div>

                    {renderContent()}
                    {hasListContent && (
                        <div className={styles.bottomAction}>
                            <button type='button' className={styles.primaryButton}>
                                Запланировать донацию
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </Layout>
    );
};

export default Bonuses;
