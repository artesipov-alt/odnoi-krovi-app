import cn from 'classnames';
import noDonorBg from 'imgs/noDonorBg.png';
import noPacketsBg from 'imgs/noPacketsBg.png';
import { FC } from 'react';

import styles from './NoResults.module.less';

type Props = {
    tab: number;
    suitableDonors?: number;
};

const NoResults: FC<Props> = ({ tab, suitableDonors = 0 }) => {
    const renderTitle = () => {
        if (tab === 1) {
            return 'Раздел в разработке, пока можете посмотреть доноров';
        }

        return suitableDonors < 1
            ? 'Пока нет свободных доноров - но мы активно ведем поиск!'
            : `Уведомили ${suitableDonors} подходящих доноров`;
    };

    const renderSubtitle = () => {
        if (tab === 1 || suitableDonors) {
            return (
                <>
                    {/*Как только найдем -<br />*/}
                    {/*направим уведомление*/}
                </>
            );
        }

        return (
            <>
                Как только они ответят –<br />
                направим уведомление
            </>
        );
    };

    return (
        <>
            <h1 className={styles.title}>{renderTitle()}</h1>
            <h5 className={styles.subtitle}>{renderSubtitle()}</h5>
            <img
                alt='no results'
                src={tab === 1 ? noPacketsBg : noDonorBg}
                className={cn(styles.noPackets, { [styles.donorsTab]: tab === 0 })}
            />
        </>
    );
};

export default NoResults;
