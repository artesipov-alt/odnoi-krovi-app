import Button from '@mui/material/Button';
import MainLogo from 'imgs/svg/mainLogo';
import { FC } from 'react';
import { useNavigate } from 'react-router';

import styles from './Final.module.less';

type Props = {
    onBackToStart: () => void;
};

const Final: FC<Props> = ({ onBackToStart }) => {
    const navigate = useNavigate();

    const onBackToSearchClickHandler = () => {
        navigate('/owner#donor');
    };

    return (
        <div className={styles.wrapper}>
            <div className={styles.content}>
                <div className={styles.logo}>
                    <MainLogo />
                </div>
                <h2 className={styles.title}>
                    Начинаем искать кровь
                    <br />
                    Вашему питомцу
                </h2>
                <Button onClick={onBackToSearchClickHandler} className={styles.button} variant='contained'>
                    К поиску
                </Button>
                <p className={styles.link} onClick={onBackToStart}>
                    Добавить еще питомца
                </p>
            </div>
        </div>
    );
};

export default Final;
