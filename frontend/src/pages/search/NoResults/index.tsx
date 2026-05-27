import cn from 'classnames';
import { FC } from 'react';

import styles from './NoResults.module.less';

type Props = {
    tab: number;
    suitableDonors?: number;
};

const NoResults: FC<Props> = ({ tab, suitableDonors = 0 }) => {
    const renderTitle = () => {
        if (tab === 1) {
            return 'Скоро будет доступно!';
        }

        return suitableDonors < 1
            ? 'Пока нет свободных доноров - но мы активно ведем поиск!'
            : `Уведомили ${suitableDonors} подходящих доноров`;
    };

    const renderSubtitle = () => {
        if (tab === 1 || suitableDonors) {
            return (
                <>
                    Функционал пока в разработке.
                    <br />
                    Но Вы можете найти донора!
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
        <div className={cn(styles.wrapper, { [styles.clinicsTab]: tab === 1 })}>
            <h1 className={styles.title}>{renderTitle()}</h1>
            <h5 className={styles.subtitle}>{renderSubtitle()}</h5>
        </div>
    );
};

export default NoResults;
