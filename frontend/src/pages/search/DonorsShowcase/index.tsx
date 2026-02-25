import { Button } from '@mui/material';
import cn from 'classnames';
import donorShowcaseStart from 'imgs/donorShowcaseStart.png';
import { FC, useState } from 'react';

import { RespondingDonor } from 'api/bloodRequest';
import { PetType } from 'api/types';

import styles from './DonorsShowcase.module.less';

type Props = {
    petType?: PetType;
    list?: RespondingDonor[];
};

const DonorsShowcase: FC<Props> = ({ list, petType }) => {
    const [showStartView, setShowStartView] = useState(true);

    const onConfirmButtonClickHandler = () => {
        setShowStartView(false);
    };

    if (!petType) {
        return null;
    }

    if (showStartView) {
        return (
            <div className={styles.startView}>
                <h1 className={styles.title}>Мы нашли для Вас донора!</h1>
                <p className={styles.descr}>Если не договоритесь - продолжайте поиск</p>
                <div className={styles.confirmButton}>
                    <Button className={styles.confirm} onClick={onConfirmButtonClickHandler}>
                        <div className={styles.labelText}>
                            К донорам <p className={styles.labelArrow}>⟶</p>
                        </div>
                    </Button>
                </div>
                <div className={styles.cancelButton}>Завершить поиск</div>
                <img className={styles.startViewImg} src={donorShowcaseStart} alt='donorShowcaseStart' />
            </div>
        );
    }

    return (
        <div className={styles.showcase}>
            {list?.map((pet) => (
                <div key={`${pet?.id}`} className={cn(styles.pet, { [styles[petType]]: true })}>
                    <div className={styles.photo} onClick={() => {}}>
                        {!!pet?.donorPhotos?.[0] && (
                            <img className={styles.img} src={pet?.donorPhotos?.[0]} alt={pet?.donorName} />
                        )}
                        <div className={styles.bloodGroup}>{pet?.donorBloodGroup || '?'}</div>
                        <div className={styles.photoFooter}>
                            <p className={styles.name}>{pet?.donorName.toUpperCase()}</p>
                        </div>
                        <div className={styles.gradient} />
                    </div>
                </div>
            ))}
        </div>
    );
};

export default DonorsShowcase;
