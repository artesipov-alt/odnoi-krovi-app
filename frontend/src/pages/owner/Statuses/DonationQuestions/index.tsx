import { Button } from '@mui/material';
import donationQuestions from 'imgs/donationQuestions.png';
import Info from 'imgs/svg/info';
import StatusQuestion from 'imgs/svg/statusQuestion';
import { FC, MouseEvent, useEffect, useState } from 'react';

import { WarnFactors } from 'api/pets';
import Layout from 'components/Layout';

import styles from './DonationQuestions.module.less';

type Props = {
    onClose?: () => void;
    factors?: WarnFactors[];
    isRecipientOpen?: boolean;
    onOpenPetProfile?: () => void;
};

const DonationQuestions: FC<Props> = ({ onClose, onOpenPetProfile, isRecipientOpen, factors = [] }) => {
    const [openTooltipId, setOpenTooltipId] = useState<number | null>(null);

    const onTooltipIconClick = (i: number) => (e: MouseEvent) => {
        e.stopPropagation();

        setOpenTooltipId(i);
    };

    useEffect(() => {
        const onOutsideClickHandler = () => {
            setOpenTooltipId(null);
        };

        window.addEventListener('click', onOutsideClickHandler);

        return () => {
            window.removeEventListener('click', onOutsideClickHandler);
        };
    }, []);

    return (
        <Layout>
            <img className={styles.img} src={donationQuestions} alt='donationQuestions' />
            <div className={styles.container}>
                <h1 className={styles.title}>Вопросы к донорству</h1>
                <p className={styles.descr}>{isRecipientOpen ? '' : 'Перед донацией обсудите с врачом следующее:'}</p>
                <div>
                    {factors.map(({ description, subDescription }, i) => (
                        <div key={description} className={styles.item}>
                            <div className={styles.icon}>
                                <StatusQuestion />
                            </div>
                            <p className={styles.itemDescr}>{description}</p>
                            {!!subDescription && (
                                <div onClick={onTooltipIconClick(i)} className={styles.infoIcon}>
                                    <Info />
                                </div>
                            )}
                            {openTooltipId === i && <div className={styles.tooltip}>{subDescription}</div>}
                        </div>
                    ))}
                </div>
                {!!onOpenPetProfile && (
                    <Button fullWidth onClick={onOpenPetProfile} className={styles.profileButton}>
                        В профиль питомца
                    </Button>
                )}
                {!!onClose && (
                    <p className={styles.back} onClick={onClose}>
                        Вернуться
                    </p>
                )}
            </div>
        </Layout>
    );
};

export default DonationQuestions;
