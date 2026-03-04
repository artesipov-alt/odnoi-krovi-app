import { Button } from '@mui/material';
import notReady from 'imgs/notReady.png';
import StatusTriangle from 'imgs/svg/statusTriangle';
import { FC, useEffect } from 'react';

import { StopFactors } from 'api/pets';
import Layout from 'components/Layout';

import styles from './NotReady.module.less';

type Props = {
    onClose: () => void;
    factors?: StopFactors[];
    onOpenPetProfile: () => void;
};

const NotReady: FC<Props> = ({ onClose, onOpenPetProfile, factors = [] }) => {
    useEffect(() => {
        document.documentElement.classList.add('useOrange2');

        return () => {
            document.documentElement.classList.remove('useOrange2');
        };
    }, []);

    return (
        <Layout>
            <img className={styles.img} src={notReady} alt='notReady' />
            <div className={styles.container}>
                <h1 className={styles.title}>
                    Питомец не может
                    <br />
                    быть донором
                </h1>
                <p className={styles.descr}>Причины:</p>
                <div>
                    {factors.map(({ description }) => (
                        <div key={description} className={styles.item}>
                            <div className={styles.icon}>
                                <StatusTriangle />
                            </div>
                            <p className={styles.itemDescr}>{description}</p>
                        </div>
                    ))}
                </div>
                <Button fullWidth onClick={onOpenPetProfile} className={styles.profileButton}>
                    В профиль питомца
                </Button>
                <p className={styles.back} onClick={onClose}>
                    Вернуться
                </p>
            </div>
        </Layout>
    );
};

export default NotReady;
