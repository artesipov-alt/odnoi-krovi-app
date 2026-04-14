import Message from 'imgs/svg/message';
import { FC } from 'react';
import { getDateFormat } from 'utils/utils';

import { RespondingDonor } from 'api/bloodRequest';

import styles from './CompletedDonations.module.less';

type Props = {
    donations: RespondingDonor[];
};

const CompletedDonations: FC<Props> = ({ donations }) => (
    <>
        <div className={styles.title}>
            <div className={styles.titleIcon}>
                <Message />
            </div>
            <p className={styles.titleText}>Проведенные донации</p>
        </div>
        <div className={styles.list}>
            {donations.map(({ id, donorName, donorBloodGroup, donorPhotos, amount, updatedAt }) => (
                <div key={id} className={styles.listItem}>
                    <div className={styles.photo}>
                        <img className={styles.photoImg} src={donorPhotos[0]} alt={donorName} />
                        <div className={styles.bloodGroup}>{donorBloodGroup !== 'UNKNOWN' ? donorBloodGroup : '?'}</div>
                    </div>
                    <div className={styles.info}>
                        <p className={styles.name}>{donorName.toUpperCase()}</p>
                        <div className={styles.status}>
                            <p className={styles.acceptedText}>Подтверждена {getDateFormat(new Date(updatedAt))}</p>
                        </div>
                    </div>
                    <div className={styles.bloodVolume}>
                        <p className={styles.bloodVolumeNumber}>{amount}</p>
                        <p className={styles.bloodVolumeDescr}>мл</p>
                    </div>
                </div>
            ))}
        </div>
    </>
);

export default CompletedDonations;
