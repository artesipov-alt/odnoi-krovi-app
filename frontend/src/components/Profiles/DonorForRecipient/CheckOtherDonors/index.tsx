import Button from '@mui/material/Button';
import Exclamation from 'imgs/svg/exclamation';
import { FC } from 'react';
import { getCorrectDeclension, Variants } from 'utils/utils';

import { Pet } from 'api/pets';
import Layout from 'components/Layout';

import styles from './CheckOtherDonors.module.less';

type Props = {
    pets: Pet[];
    avatar: string;
    onClose: () => void;
    onSuccess: () => void;
};

const CheckOtherDonors: FC<Props> = ({ pets, onSuccess, onClose, avatar }) => (
    <Layout className={styles.wrapper}>
        <div className={styles.exclamation}>
            <Exclamation />
        </div>
        <h1 className={styles.title}>Избегайте черного донорства</h1>
        <p className={styles.subTitle}>Проверьте фото других доноров этого хозяина</p>
        <div className={styles.donor}>
            <img className={styles.avatar} alt='avatar' src={avatar} />
            <p className={styles.descr}>
                Если Вам кажется, что донор недавно сдавал кровь - рекомендуем найти другого донора
            </p>
        </div>
        <div className={styles.showcase}>
            {pets.map((pet) => (
                <div key={pet.id} className={styles.pet}>
                    <div className={styles.photo}>
                        <img alt={pet.name} src={pet.photoUrls?.[0]!} className={styles.img} />
                        <div className={styles.photoFooter}>
                            <p className={styles.donorName}>{pet.name.toUpperCase()}</p>
                            <div className={styles.label}>
                                <div className={styles.recover}>
                                    <p className={styles.recoverDays}>pet.recoveryDays</p>
                                    <p className={styles.recoverDescr}>
                                        {getCorrectDeclension(Variants.DAYS, pet.recoveryDays || 1)}
                                    </p>
                                </div>
                                <div className={styles.labelText}>До восстановления</div>
                            </div>
                        </div>
                        <div className={styles.gradient} />
                    </div>
                </div>
            ))}
        </div>
        <div className={styles.buttons}>
            <Button onClick={onSuccess} fullWidth className={styles.button} variant='contained'>
                Продолжить
            </Button>
            <Button onClick={onClose} fullWidth className={styles.button} variant='contained'>
                Отказаться от донора
            </Button>
        </div>
    </Layout>
);

export default CheckOtherDonors;
