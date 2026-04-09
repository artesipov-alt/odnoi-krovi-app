import cn from 'classnames';
import { usePlannedDonations } from 'hooks/usePlannedDonations';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import AccordionArrow from 'imgs/svg/accordionArrow';
import { FC, useCallback, useEffect } from 'react';
import { toast } from 'react-toastify';

import { DonorStatus, PlannedDonation } from 'api/donor';
import { PetType } from 'api/types';
import Loading from 'components/Loading';

import styles from './PlannedDonations.module.less';

type Props = {
    id: string;
    onDonationClick: (donation: PlannedDonation) => void;
};

const PlannedDonations: FC<Props> = ({ id, onDonationClick }) => {
    const { data: donations, isLoading: isDonationsLoading, refetch, isError } = usePlannedDonations(id);

    const showToast = useCallback((text: string) => {
        toast.warn(text);
    }, []);

    const onDonationClickHandler = (donation: PlannedDonation) => () => {
        onDonationClick(donation);
    };

    const getDefaultPhoto = (petType: PetType) => (petType === PetType.DOG ? dogRoundStub : catRoundStub);

    useEffect(() => {
        if (isError) {
            showToast('Не удалось получить список запланированных донаций');
        }
    }, [isError, showToast]);

    if (isDonationsLoading) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    return (
        <div className={styles.wrapper}>
            {donations?.map((donation) => (
                <div
                    className={styles.tile}
                    key={donation.applicationData.id}
                    onClick={onDonationClickHandler(donation)}
                >
                    <div className={styles.avatars}>
                        <div className={styles.pet}>
                            <img
                                className={styles.photo}
                                alt={donation.applicationData.petName}
                                src={donation.applicationData.photoUrls[0]}
                            />
                            <p className={styles.name}>{donation.applicationData.petName.toUpperCase()}</p>
                        </div>
                        <div className={cn(styles.pet, { [styles.recipient]: true })}>
                            <img
                                alt={donation.recipientData.petName}
                                src={
                                    donation.recipientData.photoUrls?.[0]
                                        ? donation.recipientData.photoUrls[0]
                                        : getDefaultPhoto(donation.recipientData.petType)
                                }
                                className={cn(styles.photo, { [styles.isRecipient]: true })}
                            />
                            <p className={styles.name}>{donation.recipientData.petName.toUpperCase()}</p>
                        </div>
                    </div>
                    <div className={styles.info}>
                        <p className={styles.infoText}>
                            {donation.applicationData.status === DonorStatus.PENDING &&
                                'Можно отказаться и запланировать новую донацию'}
                            {donation.applicationData.status === DonorStatus.ACCEPTED &&
                                !donation.applicationData.rejectedReason &&
                                'Проведите донацию или откажитесь'}
                            {donation.applicationData.status === DonorStatus.ACCEPTED &&
                                !!donation.applicationData.rejectedReason &&
                                'Хозяин реципиента не подтвердил донацию'}
                            {donation.applicationData.status === DonorStatus.COMPLETED &&
                                'Ожидается подтверждение реципиента'}
                        </p>
                        <div className={styles.infoIcon}>
                            <AccordionArrow />
                        </div>
                    </div>
                    <div className={styles.bloodVolume}>
                        <p className={styles.bloodVolumeNumber}>{donation.recipientData.bloodVolumeNeeded}</p>
                        <p className={styles.bloodVolumeDescr}>мл</p>
                    </div>
                </div>
            ))}
        </div>
    );
};

export default PlannedDonations;
