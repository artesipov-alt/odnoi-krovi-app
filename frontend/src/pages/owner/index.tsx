import Bell from 'imgs/bell';
import Paw from 'imgs/paw';
import { FC } from 'react';
import { TelegramUser } from 'types';

import Layout from 'components/Layout';

import styles from './Owner.module.less';

type Props = {
    user: TelegramUser;
};

const Owner: FC<Props> = ({ user }) => (
    <Layout>
        <div className={styles.wrapper}>
            <div className={styles.header}>
                <h2 className={styles.title}>Мои питомцы</h2>
                <div className={styles.services}>
                    <div className={styles.bellIcon}>
                        <Bell />
                    </div>
                    <div className={styles.avatar}>{user.fullName.charAt(0).toUpperCase()}</div>
                </div>
            </div>
            <div className={styles.button}>
                <div className={styles.pawIcon}>
                    <Paw />
                </div>
                <p className={styles.buttonText}>Добавить питомца</p>
            </div>
        </div>
    </Layout>
);

export default Owner;
