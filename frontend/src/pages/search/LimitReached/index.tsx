import Button from '@mui/material/Button';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import { FC } from 'react';

import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';

import styles from './LimitReached.module.less';

type Props = {
    avatar?: string;
    type: PetType;
    bloodVolumeNeeded: number;
    onBackToSearch?: () => void;
};

const LimitReached: FC<Props> = ({ onBackToSearch, avatar, type, bloodVolumeNeeded }) => (
    <Layout className={styles.wrapper}>
        <div className={styles.content}>
            <div className={styles.avatarWrapper}>
                <img
                    alt='avatar'
                    className={styles.avatar}
                    src={avatar || (type === PetType.DOG ? dogRoundStub : catRoundStub)}
                />
                <CircularProgress
                    showDot
                    size={180}
                    strokeWidth={15}
                    total={bloodVolumeNeeded}
                    current={bloodVolumeNeeded}
                    color='var(--red10, #FF2727)'
                />
                <div className={styles.neededVolume}>
                    {bloodVolumeNeeded}
                    <span>мл</span>
                </div>
            </div>
            <h2 className={styles.title}>Вы достигли лимита</h2>
            <p className={styles.descr}>
                Ваш поиск приостановлен. Договоритесь о переливании или скройте найденные предложения, чтобы продолжить
                поиск.
            </p>
            <Button fullWidth onClick={onBackToSearch} className={styles.button} variant='contained'>
                К найденному
            </Button>
        </div>
    </Layout>
);

export default LimitReached;
