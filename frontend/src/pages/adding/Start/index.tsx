import Button from '@mui/material/Button';
import MainLogo from 'imgs/svg/mainLogo';
import { FC } from 'react';
import { Link } from 'react-router';

import styles from './Start.module.less';

type Props = {
    onAddDonorClick: () => void;
    onAddRecipientClick: () => void;
};

const Start: FC<Props> = ({ onAddDonorClick, onAddRecipientClick }) => (
    <div className={styles.wrapper}>
        <div className={styles.content}>
            <div className={styles.logo}>
                <MainLogo />
            </div>
            <h2 className={styles.title}>
                Для чего
                <br />
                Вы создаете
                <br />
                профиль питомца?
            </h2>
            <div>
                <Button onClick={onAddRecipientClick} className={styles.button} variant='contained'>
                    Ищу кровь
                </Button>
                <Button onClick={onAddDonorClick} className={styles.button} variant='contained'>
                    Хочу помочь
                </Button>
            </div>
            <Link className={styles.link} to='/owner'>
                На главную
            </Link>
        </div>
    </div>
);

export default Start;
