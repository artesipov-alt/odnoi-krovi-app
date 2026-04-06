import Button from '@mui/material/Button';
import cn from 'classnames';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import { FC } from 'react';

import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';

import styles from './DonationComplete.module.less';

type Props = {
    avatar?: string;
    type: PetType;
    onEndSearch: () => void;
    bloodVolumeNeeded: number;
    bloodVolumeDonated: number;
    onBackToSearch: () => void;
};

const DonationComplete: FC<Props> = ({
    type,
    avatar,
    onEndSearch,
    onBackToSearch,
    bloodVolumeNeeded,
    bloodVolumeDonated,
}) => (
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
                    current={bloodVolumeDonated}
                    color='var(--red10, #FF2727)'
                />
                <div className={styles.neededVolume}>
                    {bloodVolumeNeeded}
                    <span>мл</span>
                </div>
            </div>
            <h2 className={styles.title}>Вы получили кровь</h2>
            <div className={styles.volume}>
                <p className={styles.volumeNumber}>{bloodVolumeDonated}</p>
                <p className={styles.volumeDescr}>мл</p>
            </div>
            <p className={styles.descr}>
                Если крови достаточно - отмените поиск, чтобы другие животные получили помощь!
            </p>
            <div className={styles.buttons}>
                <Button onClick={onEndSearch} className={styles.button} variant='contained'>
                    Завершить поиск
                </Button>
                <Button
                    variant='contained'
                    onClick={onBackToSearch}
                    className={cn(styles.button, { [styles.more]: true })}
                >
                    Найти ещё
                </Button>
            </div>
        </div>
    </Layout>
);

export default DonationComplete;
