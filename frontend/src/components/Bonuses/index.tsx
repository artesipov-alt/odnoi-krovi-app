import { Button } from '@mui/material';
import cn from 'classnames';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Drugs from 'imgs/svg/drugs';
import Feed from 'imgs/svg/feed';
import Other from 'imgs/svg/other';
import PrioritySearch from 'imgs/svg/prioritySearch';
import { FC } from 'react';
import { useNavigate } from 'react-router';

import { Bonus, BonusType } from 'api/donor';
import Layout from 'components/Layout';

import Alert from '../Alert';
import styles from './Bonuses.module.less';

type Props = {
    items?: Bonus[];
    onClose: () => void;
    isReceivedBonuses?: boolean;
    fromDonationDetails?: boolean;
};

const Bonuses: FC<Props> = ({ items, onClose, isReceivedBonuses, fromDonationDetails }) => {
    const navigate = useNavigate();

    const onGoToBonusesClickHandler = () => {
        navigate('/bonuses');
    };

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.innerWrapper}>
                <div className={styles.header}>
                    <div className={styles.back} onClick={onClose}>
                        <BackAngularArrow />
                    </div>
                    <h2 className={styles.title}>Бонусы Портала</h2>
                </div>
                <div className={styles.priority}>
                    <div className={styles.icon}>
                        <PrioritySearch />
                    </div>
                    <p className={styles.iconTitle}>Приоритетный поиск</p>
                </div>
                <p className={cn(styles.priorityDescr, { [styles.withMb]: fromDonationDetails })}>
                    Начисляется за проведенные донации, сможете найти кровь одним из первых при необходимости
                </p>
                {!isReceivedBonuses ? (
                    <Alert
                        className={cn(styles.alert, {
                            [styles.none]: fromDonationDetails && items?.[0]?.type !== BonusType.LOCK,
                        })}
                        text={
                            items?.[0]?.type === BonusType.LOCK
                                ? 'остальные бонусы станут доступны через 2 месяца, следите за обновлениями – набор доступных бонусов может меняться'
                                : 'Бонусы могут закончиться - успейте предложить донора!'
                        }
                    />
                ) : (
                    <div className={styles.buttonWrapper}>
                        <Button onClick={onGoToBonusesClickHandler} className={styles.button}>
                            Применить бонусы
                        </Button>
                    </div>
                )}
            </div>
            {items?.[0]?.type !== BonusType.LOCK && (
                <div className={styles.bonuses}>
                    {items
                        ?.sort((a, b) => {
                            const order = [BonusType.PREPARATION, BonusType.FOOD, BonusType.OTHER];

                            return order.indexOf(a.type) - order.indexOf(b.type);
                        })
                        .map(({ description, type, partner }, i) => (
                            // eslint-disable-next-line react/no-array-index-key
                            <div key={`${partner}_${i}`} className={styles.bonus}>
                                <div className={styles.bonusHeader}>
                                    <div className={styles.icon}>
                                        {type === BonusType.FOOD && <Feed />}
                                        {type === BonusType.OTHER && <Other />}
                                        {type === BonusType.PREPARATION && <Drugs />}
                                    </div>
                                    <p className={styles.iconTitle}>{partner}</p>
                                </div>
                                <p className={styles.bonusDescr}>{description}</p>
                            </div>
                        ))}
                </div>
            )}
        </Layout>
    );
};

export default Bonuses;
