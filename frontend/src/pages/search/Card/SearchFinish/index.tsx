import Button from '@mui/material/Button';
import cn from 'classnames';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import { FC } from 'react';

import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';

import styles from './SearchFinish.module.less';

type Props = {
    type: PetType;
    avatar?: string;
    onOpenDetail: () => void;
    bloodVolumeNeeded: number;
    bloodVolumeDonated: number;
    onBackToOwner: () => void;
};

const SearchFinish: FC<Props> = ({
    type,
    avatar,
    onOpenDetail,
    onBackToOwner,
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
            <h2 className={styles.title}>Ваш поиск выполнен!</h2>
            <p className={styles.descr}>При необходимости создайте новый поиск</p>
            <div className={styles.buttons}>
                <Button onClick={onBackToOwner} className={styles.button} variant='contained'>
                    К питомцам
                </Button>
                <Button
                    variant='contained'
                    onClick={onOpenDetail}
                    className={cn(styles.button, { [styles.more]: true })}
                >
                    Детали поиска
                </Button>
            </div>
        </div>
    </Layout>
);

export default SearchFinish;
