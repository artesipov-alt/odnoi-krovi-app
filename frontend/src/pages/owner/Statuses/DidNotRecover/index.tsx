import { Button } from '@mui/material';
import didNotRecover from 'imgs/didNotRecover.png';
import { FC, useEffect } from 'react';
import { getCorrectDeclension, Variants } from 'utils/utils';

import Layout from 'components/Layout';

import styles from './DidNotRecover.module.less';

type Props = {
    onClose: () => void;
    recoveryDays: number;
    onOpenPetProfile: () => void;
};

const DidNotRecover: FC<Props> = ({ onClose, onOpenPetProfile, recoveryDays }) => {
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
            <p className={styles.descr}>
                Запланируйте донацию через {recoveryDays} {getCorrectDeclension(Variants.DAYS, recoveryDays)}
            </p>
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
