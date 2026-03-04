import { Button } from '@mui/material';
import didNotRecover from 'imgs/didNotRecover.png';
import { FC, useEffect } from 'react';

import Layout from 'components/Layout';

import styles from './DidNotRecover.module.less';

type Props = {
    onClose: () => void;
    onOpenPetProfile: () => void;
};

const DidNotRecover: FC<Props> = ({ onClose, onOpenPetProfile }) => {
    useEffect(() => {
        document.documentElement.classList.add('useStatusBg1');

        return () => {
            document.documentElement.classList.remove('useStatusBg1');
        };
    }, []);

    return (
        <Layout>
            <img className={styles.img} src={didNotRecover} alt='didNotRecover' />
            <h1 className={styles.title}>
                Донор
                <br />
                на восстановлении
            </h1>
            <p className={styles.descr}>Запланируйте донацию через 15 дней</p>
            <div className={styles.buttonWrapper}>
                <Button onClick={onOpenPetProfile} className={styles.profileButton}>
                    В профиль питомца
                </Button>
            </div>
            <p className={styles.back} onClick={onClose}>
                Вернуться
            </p>
        </Layout>
    );
};

export default DidNotRecover;
