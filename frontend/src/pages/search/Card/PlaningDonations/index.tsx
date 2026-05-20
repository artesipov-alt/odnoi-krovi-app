import { Button } from '@mui/material';
import cn from 'classnames';
import Exclamation from 'imgs/svg/exclamation';
import Info from 'imgs/svg/info';
import Message from 'imgs/svg/message';
import { FC, useCallback, useState } from 'react';
import { toast } from 'react-toastify';
import { isWithinHours } from 'utils/utils';

import { confirmDonation } from 'api/apiServices/confirmDonation';
import { RespondingDonor, RespondingDonorStatus } from 'api/bloodRequest';
import { queryClient } from 'api/queryClient';
import Curtain from 'components/Curtain';
import Timer from 'components/Timer';

import styles from './PlaningDonations.module.less';

type Props = {
    userId: string;
    poolRequestRefetch: () => void;
    donorResponses: RespondingDonor[];
    onDonationClick: (id: string, status: RespondingDonorStatus, updatedAt: string) => void;
};

const getSortedDonations = (donations: RespondingDonor[]) =>
    [...donations].sort((a, b) => {
        const isLowPriority = (status: RespondingDonorStatus) =>
            status === RespondingDonorStatus.REJECTED || status === RespondingDonorStatus.CANCELED;

        const aLow = isLowPriority(a.status);
        const bLow = isLowPriority(b.status);

        if (aLow && !bLow) return 1;

        if (!aLow && bLow) return -1;

        return 0;
    });

const PlaningDonations: FC<Props> = ({ userId, donorResponses, onDonationClick, poolRequestRefetch }) => {
    const [isInfoCurtainOpen, setIsInfoCurtainOpen] = useState(false);

    const showToast = useCallback((text: string) => {
        toast.warn(text);
    }, []);

    const onCurtainOpenToggle = () => {
        setIsInfoCurtainOpen((prevState) => !prevState);
    };

    const onDonationClickHandler = (id: string, status: RespondingDonorStatus, updatedAt: string) => () => {
        onDonationClick(id, status, updatedAt);
    };

    const onEndTimerClickHandler = (id: string, amount: number) => async () => {
        const response = await confirmDonation({
            id,
            amount,
        });

        if (!response) {
            showToast('Не удалось подтвердить донацию');
        }

        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });

        poolRequestRefetch();
    };

    return (
        <>
            <div className={styles.title}>
                <div className={styles.titleIcon}>
                    <Message />
                </div>
                <p className={styles.titleText}>Планируемые донации</p>
                <div onClick={onCurtainOpenToggle} className={styles.titleInfo}>
                    <Info />
                </div>
            </div>
            <div className={styles.list}>
                {getSortedDonations(donorResponses).map(
                    ({ id, status, donorName, donorBloodGroup, donorPhotos, amount, updatedAt }) => (
                        <div
                            key={id}
                            onClick={onDonationClickHandler(id, status, updatedAt)}
                            className={cn(styles.listItem, {
                                [styles.notActive]:
                                    status === RespondingDonorStatus.REJECTED ||
                                    status === RespondingDonorStatus.CANCELED,
                            })}
                        >
                            <div className={styles.photo}>
                                <img className={styles.photoImg} src={donorPhotos[0]} alt={donorName} />
                                <div className={styles.bloodGroup}>
                                    {donorBloodGroup !== 'UNKNOWN' ? donorBloodGroup : '?'}
                                </div>
                            </div>
                            <div className={styles.info}>
                                <p className={styles.name}>{donorName.toUpperCase()}</p>
                                <div className={styles.status}>
                                    {status === RespondingDonorStatus.ACCEPTED && (
                                        <p className={styles.acceptedText}>Донация состоялась?</p>
                                    )}
                                    {status === RespondingDonorStatus.COMPLETED && isWithinHours(updatedAt, 48) && (
                                        <div className={styles.count}>
                                            <Timer
                                                hoursToAdd={48}
                                                updatedAt={updatedAt}
                                                className={styles.countTimer}
                                                digitClassName={styles.countDigits}
                                                separatorClassName={styles.countSeparator}
                                                onTimeEnd={onEndTimerClickHandler(id, amount)}
                                            />
                                            <p className={styles.timerText}>до подтверждения донации</p>
                                        </div>
                                    )}
                                    {status === RespondingDonorStatus.COMPLETED && !isWithinHours(updatedAt, 48) && (
                                        <p className={styles.acceptedText}>Хозяин донора сообщил о донации</p>
                                    )}
                                    {status === RespondingDonorStatus.REJECTED && (
                                        <p className={styles.acceptedText}>Реципиент отказался от донации</p>
                                    )}
                                    {status === RespondingDonorStatus.CANCELED && (
                                        <p className={styles.acceptedText}>Донор отказался от донации</p>
                                    )}
                                </div>
                            </div>
                            <div className={styles.bloodVolume}>
                                <p className={styles.bloodVolumeNumber}>{amount}</p>
                                <p className={styles.bloodVolumeDescr}>мл</p>
                            </div>
                        </div>
                    ),
                )}
            </div>
            {isInfoCurtainOpen && (
                <Curtain noRednerButtons shouldCloseByWrapperClick onClose={onCurtainOpenToggle}>
                    <div className={styles.exclamation}>
                        <Exclamation />
                    </div>
                    <h3 className={styles.curtainTitle}>
                        Обязательно
                        <br />
                        прочитайте условия
                    </h3>
                    <div className={styles.curtainString}>
                        <div className={styles.curtainStringNumber}>1</div>
                        <p className={styles.curtainStringText}>
                            Если донация состоялась - подтвердите ее, чтобы донор получил бонусы
                        </p>
                    </div>
                    <div className={cn(styles.curtainString, { [styles.last]: true })}>
                        <div className={styles.curtainStringNumber}>2</div>
                        <p className={styles.curtainStringText}>
                            Если донор сообщил о донации - у вас 3 дня на подтверждение, иначе она подтвердится
                            автоматически
                        </p>
                    </div>
                    <Button fullWidth className={styles.curtainConfirm} onClick={onCurtainOpenToggle}>
                        Понятно
                    </Button>
                </Curtain>
            )}
        </>
    );
};

export default PlaningDonations;
