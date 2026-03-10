import Button from '@mui/material/Button';
import { FC } from 'react';
import { useNavigate } from 'react-router';

import ImgEditor from 'components/ImgEditor';

import styles from './Final.module.less';

type Props = {
    photo: File | null;
    onBackToStart: () => void;
};

const Final: FC<Props> = ({ onBackToStart, photo }) => {
    const navigate = useNavigate();

    const onBackToSearchClickHandler = () => {
        navigate('/owner#donor');
    };

    return (
        <div className={styles.wrapper}>
            <div className={styles.content}>
                <div className={styles.photo}>
                    <ImgEditor showStub src={photo} className={styles.photo} isMiniView />
                </div>
                <h2 className={styles.title}>
                    Профиль питомца
                    <br />
                    создан!
                </h2>
                <Button onClick={onBackToSearchClickHandler} className={styles.button} variant='contained'>
                    К питомцам
                </Button>
                <p className={styles.link} onClick={onBackToStart}>
                    Добавить еще питомца
                </p>
            </div>
        </div>
    );
};

export default Final;
