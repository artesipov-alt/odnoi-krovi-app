import cn from 'classnames';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import BloodNo from 'imgs/svg/bloodNo';
import BloodOk from 'imgs/svg/bloodOk';
import BloodVolume from 'imgs/svg/bloodVolume';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';
import { getDateFormat } from 'utils/utils';

import { getCompletedDonations } from 'api/apiServices/getCompletedDonations';
import { DonorStatus, GetCompletedDonationsResponse, PlannedDonation } from 'api/donor';
import { PetType } from 'api/types';
import CompletedDonation from 'components/CompletedDonation';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import styles from './DonationsHistory.module.less';

type DetailsPage = {
    isOpen: boolean;
    donation?: PlannedDonation;
};

type Props = {
    id: string;
};

const tiles = ['количество донаций', 'объем донаций'];

const sortDonationsByStatusAndDate = (items: PlannedDonation[]) => {
    return items.sort((a, b) => {
        const statusA = a.applicationData.status;
        const statusB = b.applicationData.status;

        const isCancelledOrRejectedA = statusA === DonorStatus.CANCELLED || statusA === DonorStatus.REJECTED;
        const isCancelledOrRejectedB = statusB === DonorStatus.CANCELLED || statusB === DonorStatus.REJECTED;

        // Если один из статусов — отменён/отклонён, а другой — нет, сортируем по приоритету: активные вперед
        if (isCancelledOrRejectedA !== isCancelledOrRejectedB) {
            return isCancelledOrRejectedA ? 1 : -1;
        }

        // Если оба статуса одинаковы (оба отменены или оба нет), сортируем по дате updated_at (свежие — первыми)
        const dateA = new Date(a.applicationData.updatedAt).getTime();
        const dateB = new Date(b.applicationData.updatedAt).getTime();

        return dateB - dateA; // Свежие даты в начало
    });
};

const getDefaultPhoto = (petType: PetType) => (petType === PetType.DOG ? dogRoundStub : catRoundStub);

const DonationsHistory: FC<Props> = ({ id }) => {
    const navigate = useNavigate();

    const [isLoading, setIsLoading] = useState(true);
    const [detailsPage, setDetailsPage] = useState<DetailsPage>({ isOpen: false });
    const [history, setHistory] = useState<GetCompletedDonationsResponse | null>(null);

    const showToast = useCallback(
        (text: string) => {
            toast.warn(text, {
                onClose: () => {
                    navigate('/owner#donor');
                },
            });
        },
        [navigate],
    );

    const fetchHistory = useCallback(async () => {
        const response = await getCompletedDonations(id);

        if (!response) {
            showToast('Не удалось получить историю донаций');

            return;
        }

        setIsLoading(false);
        setHistory(response.data);
    }, [id, showToast]);

    const onCloseClickHandler = () => {
        navigate('/owner#donor');
    };

    const onDonationClickHandler = (donation: PlannedDonation) => () => {
        setDetailsPage({ isOpen: true, donation });
    };

    const onDetailsCloseClickHandler = () => {
        setDetailsPage({ isOpen: false });
    };

    useEffect(() => {
        fetchHistory();
    }, [fetchHistory]);

    if (isLoading) {
        return (
            <div className={styles.loading}>
                <Loading size={90} thickness={4} />
            </div>
        );
    }

    if (detailsPage.isOpen && detailsPage.donation) {
        return <CompletedDonation donation={detailsPage.donation} onClose={onDetailsCloseClickHandler} />;
    }

    return (
        <Layout>
            <div className={styles.header}>
                <div className={styles.back} onClick={onCloseClickHandler}>
                    <BackAngularArrow />
                </div>
                <h2 className={styles.title}>История донаций</h2>
            </div>
            <div className={styles.tiles}>
                {tiles.map((tile, i) => (
                    <div className={styles.tile} key={tile}>
                        <div className={styles.tileInfo}>
                            <div className={styles.tileIcon}>{i === 0 ? <Blood /> : <BloodVolume />}</div>
                            <p className={styles.tileValue}>
                                {i === 0
                                    ? history?.totalCompletedDonations || 0
                                    : history?.totalDonatedVolume.toString().replace('.', ',')}
                            </p>
                            {i === 1 && <p className={styles.tileUnit}>мл</p>}
                        </div>
                        <p className={styles.tileDescr}>{tile}</p>
                    </div>
                ))}
            </div>
            {!history?.items.length ? (
                <div className={styles.noItems}>
                    <h2 className={styles.noItemsTitle}>
                        У Вас пока
                        <br />
                        не было донаций
                    </h2>
                </div>
            ) : (
                <div className={styles.list}>
                    {sortDonationsByStatusAndDate(history.items).map((donation) => {
                        const isCanceled =
                            donation.applicationData.status === DonorStatus.CANCELLED ||
                            donation.applicationData.status === DonorStatus.REJECTED;

                        return (
                            <div
                                className={styles.donationTile}
                                key={donation.applicationData.id}
                                onClick={isCanceled ? undefined : onDonationClickHandler(donation)}
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
                                    <p className={styles.infoTitle}>
                                        Донация от {getDateFormat(new Date(donation.applicationData.updatedAt))}
                                    </p>
                                    <div className={cn(styles.status, { [styles.isCanceled]: isCanceled })}>
                                        <div className={cn(styles.statusIcon, { [styles.isCanceled]: isCanceled })}>
                                            {isCanceled ? <BloodNo /> : <BloodOk />}
                                        </div>
                                        <p className={styles.statusText}>{isCanceled ? 'отменилась' : 'состоялась'}</p>
                                    </div>
                                </div>
                                {!isCanceled && (
                                    <div className={styles.bloodVolume}>
                                        <p className={styles.bloodVolumeNumber}>
                                            {donation.applicationData.amount.toString().replace('.', ',')}
                                        </p>
                                        <p className={styles.bloodVolumeDescr}>мл</p>
                                    </div>
                                )}
                            </div>
                        );
                    })}
                </div>
            )}
        </Layout>
    );
};

export default DonationsHistory;
