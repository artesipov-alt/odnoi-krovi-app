import { Button } from '@mui/material';
import cn from 'classnames';
import donorShowcaseStart from 'imgs/donorShowcaseStart.png';
import BoneBig from 'imgs/svg/boneBig';
import NotPaid from 'imgs/svg/notPaid';
import Paid from 'imgs/svg/paid';
import RoundQuestion from 'imgs/svg/roundQuestion';
import TaxiBig from 'imgs/svg/taxiBig';
import { FC, MouseEvent, useCallback } from 'react';
import { toast } from 'react-toastify';

import { closeSearch } from 'api/apiServices/closeSearch';
import { RespondingDonor } from 'api/bloodRequest';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import { CompensationType } from 'api/user';

import styles from './DonorsShowcase.module.less';

type Props = {
    userId: string;
    searchId: string;
    petType?: PetType;
    goToOwner: () => void;
    showStartView: boolean;
    list?: RespondingDonor[];
    setIsStartViewShown: () => void;
    onDonorClick: (donorId: string) => void;
    onOpenWarnFactors: (donorId: string) => void;
};

const DonorsShowcase: FC<Props> = ({
    list,
    userId,
    petType,
    searchId,
    goToOwner,
    onDonorClick,
    showStartView,
    onOpenWarnFactors,
    setIsStartViewShown,
}) => {
    const showToast = useCallback((text: string) => {
        toast.warn(text);
    }, []);

    const onConfirmButtonClickHandler = () => {
        setIsStartViewShown();
    };

    const onCLoseSearchClickHandler = async () => {
        const response = await closeSearch(searchId);

        if (!response) {
            showToast('Не удалось завершить поиск');
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });

        goToOwner();
    };

    const onShowWarnFactorsToggle = (donorId: string) => (e: MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();

        onOpenWarnFactors(donorId);
    };

    const onDonorClickHandler = (donorId: string) => () => {
        onDonorClick(donorId);
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
                <div role='button' onClick={onCLoseSearchClickHandler} className={styles.cancelButton}>
                    Завершить поиск
                </div>
                <img className={styles.startViewImg} src={donorShowcaseStart} alt='donorShowcaseStart' />
            </div>
        );
    }

    return (
        <div className={styles.showcase}>
            {list?.map((pet) => (
                <div
                    key={`${pet?.id}`}
                    onClick={onDonorClickHandler(pet.donorId)}
                    className={cn(styles.pet, { [styles[petType]]: true })}
                >
                    <div className={styles.photo} onClick={() => {}}>
                        {!!pet?.donorPhotos?.[0] && (
                            <img className={styles.img} src={pet?.donorPhotos?.[0]} alt={pet?.donorName} />
                        )}
                        <div className={styles.info}>
                            <div className={styles.bloodGroup}>
                                {pet?.donorBloodGroup !== 'UNKNOWN' ? pet?.donorBloodGroup : '?'}
                            </div>
                            <div className={styles.icon}>
                                {pet.compensationType === CompensationType.FREE && <NotPaid />}
                                {pet.compensationType === CompensationType.PAID && <Paid />}
                                {pet.compensationType === CompensationType.FOOD && <BoneBig />}
                            </div>
                            {pet.taxiCompensation && (
                                <div className={styles.icon}>
                                    <TaxiBig />
                                </div>
                            )}
                        </div>
                        <div className={styles.bloodVolume}>
                            <p className={styles.bloodVolumeNumber}>{pet.amount}</p>
                            <p className={styles.bloodVolumeDescr}>мл</p>
                        </div>
                        <div className={styles.photoFooter}>
                            <p className={styles.name}>{pet?.donorName.toUpperCase()}</p>
                            {!!pet.warnFactors?.length && (
                                <div
                                    onClick={onShowWarnFactorsToggle(pet.donorId)}
                                    className={cn(styles.label, { [styles.donationQuestions]: true })}
                                >
                                    <div className={styles.searchIcon}>
                                        <RoundQuestion />
                                    </div>
                                    <div>Вопросы к донорству</div>
                                </div>
                            )}
                        </div>
                        <div className={styles.gradient} />
                    </div>
                </div>
            ))}
        </div>
    );
};

export default DonorsShowcase;
