import cn from 'classnames';
import { usePlannedDonations } from 'hooks/usePlannedDonations';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import emptyBg from 'imgs/emptyBg.png';
import AccordionArrow from 'imgs/svg/accordionArrow';
import { FC, useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { isWithinHours } from 'utils/utils';

import { confirmDonation } from 'api/apiServices/confirmDonation';
import { DonorStatus, PlannedDonation } from 'api/donor';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import Loading from 'components/Loading';
import Timer from 'components/Timer';

import styles from './PlannedDonations.module.less';

type Props = {
    id: string;
    onDonationClick: (donation: PlannedDonation) => void;
};

const getDefaultPhoto = (petType: PetType) => (petType === PetType.DOG ? dogRoundStub : catRoundStub);

const PlannedDonations: FC<Props> = ({ id, onDonationClick }) => {
    const [expiredPendingTimers, setExpiredPendingTimers] = useState<string[]>([]);

    const { data: donations, isLoading: isDonationsLoading, refetch, isError } = usePlannedDonations(id);

    const showToast = useCallback((text: string) => {
        toast.warn(text);
    }, []);

    const onDonationClickHandler = (donation: PlannedDonation) => () => {
        onDonationClick(donation);
    };

    const onPendingTimerExpired = (donationId: string) => () => {
        setExpiredPendingTimers((prevState) => [...prevState, donationId]);
    };

    const onCompleteTimerExpired = (donationId: string) => async () => {
        const donation = donations?.find((item) => item.applicationData.id === donationId);

        if (!donation) {
            return;
        }

        const response = await confirmDonation({
            id: donation.applicationData.id,
            amount: donation.applicationData.amount,
        });

        if (!response) {
            showToast('Не удалось подтвердить донацию');

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', id] });

        refetch();
    };

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
        <div className={cn(styles.wrapper, { [styles.noItems]: !donations?.length })}>
            {!!donations?.length &&
                donations.map((donation) => {
                    const showPendingTimer =
                        donation.applicationData.status === DonorStatus.PENDING &&
                        isWithinHours(donation.applicationData.createdAt, 1) &&
                        !expiredPendingTimers.includes(donation.applicationData.id);

                    const showCompletedTimer =
                        donation.applicationData.status === DonorStatus.COMPLETED &&
                        isWithinHours(donation.applicationData.updatedAt, 72);

                    const isPendingTimerExpired =
                        donation.applicationData.status === DonorStatus.PENDING &&
                        (!isWithinHours(donation.applicationData.createdAt, 1) ||
                            expiredPendingTimers.includes(donation.applicationData.id));

                    return (
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
                            <div
                                className={cn(styles.info, {
                                    [styles.isTimer]: showPendingTimer || showCompletedTimer,
                                })}
                            >
                                {showPendingTimer && (
                                    <div className={styles.count}>
                                        <p className={styles.timerText}>Можно отказаться через</p>
                                        <div className={styles.timer}>
                                            <Timer
                                                hoursToAdd={1}
                                                className={styles.countTimer}
                                                digitClassName={styles.countDigits}
                                                separatorClassName={styles.countSeparator}
                                                updatedAt={donation.applicationData.createdAt}
                                                onTimeEnd={onPendingTimerExpired(donation.applicationData.id)}
                                            />
                                            <div className={cn(styles.infoIcon, { [styles.inTimer]: true })}>
                                                <AccordionArrow />
                                            </div>
                                        </div>
                                    </div>
                                )}
                                {showCompletedTimer && (
                                    <div className={styles.count}>
                                        <p className={styles.timerText}>до подтверждения донации</p>
                                        <div className={styles.timer}>
                                            <Timer
                                                hoursToAdd={72}
                                                className={styles.countTimer}
                                                digitClassName={styles.countDigits}
                                                separatorClassName={styles.countSeparator}
                                                updatedAt={donation.applicationData.updatedAt}
                                                onTimeEnd={onCompleteTimerExpired(donation.applicationData.id)}
                                            />
                                            <div className={cn(styles.infoIcon, { [styles.inTimer]: true })}>
                                                <AccordionArrow />
                                            </div>
                                        </div>
                                    </div>
                                )}
                                <p className={styles.infoText}>
                                    {isPendingTimerExpired && 'Можно отказаться и запланировать новую донацию'}
                                    {donation.applicationData.status === DonorStatus.ACCEPTED &&
                                        !donation.applicationData.rejectedReason &&
                                        'Проведите донацию или откажитесь'}
                                    {donation.applicationData.status === DonorStatus.ACCEPTED &&
                                        !!donation.applicationData.rejectedReason &&
                                        'Хозяин реципиента не подтвердил донацию'}
                                    {donation.applicationData.status === DonorStatus.COMPLETED &&
                                        !isWithinHours(donation.recipientData.updatedAt, 72) &&
                                        'Ожидается подтверждение реципиента'}
                                </p>
                                <div
                                    className={cn(styles.infoIcon, {
                                        [styles.hide]: showPendingTimer || showCompletedTimer,
                                    })}
                                >
                                    <AccordionArrow />
                                </div>
                            </div>
                            <div className={styles.bloodVolume}>
                                <p className={styles.bloodVolumeNumber}>
                                    {donation.applicationData.amount > donation.recipientData.bloodVolumeNeeded
                                        ? donation.recipientData.bloodVolumeNeeded
                                        : donation.applicationData.amount}
                                </p>
                                <p className={styles.bloodVolumeDescr}>мл</p>
                            </div>
                        </div>
                    );
                })}
            {!donations?.length && (
                <div className={styles.emptyBlock}>
                    <div className={styles.emptyTitle}>Пока нет новых донаций...</div>
                    <img src={emptyBg} alt='Питомцы' className={styles.emptyImage} />
                </div>
            )}
        </div>
    );
};

export default PlannedDonations;
