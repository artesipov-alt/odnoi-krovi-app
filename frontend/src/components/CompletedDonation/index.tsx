import cn from 'classnames';
import { useLocationsQuery } from 'hooks/useDicts';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import BloodVolume from 'imgs/svg/bloodVolume';
import Bone from 'imgs/svg/bone';
import Certificates from 'imgs/svg/certificates';
import Location from 'imgs/svg/location';
import MiniSinglePaw from 'imgs/svg/miniSinglePaw';
import Pin from 'imgs/svg/pin';
import PrioritySearch from 'imgs/svg/prioritySearch';
import Taxi from 'imgs/svg/taxi';
import Accordion from 'pages/adding/common/Accordion';
import { FC, useCallback, useEffect } from 'react';
import { toast } from 'react-toastify';
import { getDateFormat } from 'utils/utils';

import { PlannedDonation } from 'api/donor';
import { PetType } from 'api/types';
import { CompensationType } from 'api/user';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';

import styles from './CompletedDonation.module.less';

type Props = {
    onClose: () => void;
    donation: PlannedDonation;
};

const CompletedDonation: FC<Props> = ({ onClose, donation }) => {
    const { data: locationsDict = [], isError: isErrorLocations } = useLocationsQuery();

    const showToast = useCallback(
        (text: string) => {
            toast.warn(text, {
                onClose: () => {
                    onClose();
                },
            });
        },
        [onClose],
    );

    useEffect(() => {
        if (isErrorLocations) {
            showToast('Не удалось загрузить словарь регионов, попробуйте перезагрузить приложение');
        }
    }, [isErrorLocations, showToast]);

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.back} onClick={onClose}>
                    <BackAngularArrow />
                </div>
                <h2 className={styles.title}>
                    Донация от {getDateFormat(new Date(donation.applicationData.updatedAt))}
                </h2>
            </div>
            <div className={styles.photos}>
                <div className={styles.avatarWrapper}>
                    <img
                        className={styles.avatar}
                        alt={donation.applicationData.petName}
                        src={donation.applicationData.photoUrls[0]}
                    />
                    <CircularProgress size={156} strokeWidth={10} total={1} current={1} color='var(--red10, #FF2727)' />
                    <p className={styles.name}>{donation.applicationData.petName.toUpperCase()}</p>
                </div>
                <div className={styles.photoDivider} />
                <div className={styles.avatarWrapper}>
                    <img
                        className={styles.avatar}
                        alt={donation.recipientData.petType}
                        src={
                            donation.recipientData.photoUrls?.[0] ||
                            (donation.recipientData.petType === PetType.DOG ? dogRoundStub : catRoundStub)
                        }
                    />
                    <CircularProgress
                        size={156}
                        strokeWidth={10}
                        color='var(--red10, #FF2727)'
                        total={donation.recipientData.bloodVolumeNeeded}
                        current={
                            donation.applicationData.amount > donation.recipientData.bloodVolumeNeeded
                                ? donation.recipientData.bloodVolumeNeeded
                                : donation.applicationData.amount
                        }
                    />
                    <p className={styles.name}>{donation.recipientData.petName.toUpperCase()}</p>
                </div>
                <div className={styles.neededVolume}>
                    {donation.applicationData.amount > donation.recipientData.bloodVolumeNeeded
                        ? donation.recipientData.bloodVolumeNeeded
                        : donation.applicationData.amount}
                    <span>мл</span>
                </div>
            </div>
            <div className={styles.recipientInfo}>
                <div className={styles.infoTitle}>
                    <div className={styles.infoTitleIcon}>
                        <MiniSinglePaw />
                    </div>
                    <p className={styles.infoTitleDescr}>О реципиенте</p>
                </div>
                <div className={styles.line}>
                    <div className={styles.lineItem}>
                        <div className={styles.lineTitle}>
                            <div className={styles.lineTitleIcon}>
                                <Blood />
                            </div>
                            <p className={styles.lineTitleText}>Искал</p>
                        </div>
                        <div className={styles.bloodInfo}>
                            <div className={styles.bloodGroup}>{donation.recipientData.bloodGroup}</div>
                            {(donation.recipientData.searchingBloodNames.some(
                                (group) => group !== donation.recipientData.bloodGroup,
                            ) ||
                                donation.recipientData.includeUnknownBloodGroup) &&
                                ' +'}
                            {donation.recipientData.searchingBloodNames.some(
                                (group) => group !== donation.recipientData.bloodGroup,
                            )
                                ? donation.recipientData.searchingBloodNames
                                      .filter((group) => group !== donation.recipientData.bloodGroup)
                                      .map((group) => (
                                          <div
                                              key={group}
                                              className={cn(styles.bloodGroup, {
                                                  [styles.needed]: true,
                                              })}
                                          >
                                              {group}
                                          </div>
                                      ))
                                : ''}
                            {donation.recipientData.includeUnknownBloodGroup && (
                                <div className={cn(styles.bloodGroup, { [styles.needed]: true })}>?</div>
                            )}
                        </div>
                    </div>
                    <div className={styles.lineItem}>
                        <div className={cn(styles.lineTitle, { [styles.leftMargin]: true })}>
                            <div className={styles.lineTitleIcon}>
                                <BloodVolume />
                            </div>
                            <p className={styles.lineTitleText}>Нужный объем</p>
                        </div>
                        <div className={styles.bloodVolume}>
                            <p className={styles.bloodVolumeNumber}>{donation.recipientData.bloodVolumeNeeded}</p>
                            <p className={styles.bloodVolumeDescr}>мл</p>
                        </div>
                    </div>
                </div>
                <div className={cn(styles.line, { [styles.last]: true })}>
                    <div className={cn(styles.lineItem, { [styles.isAvatar]: true })}>
                        <div className={styles.ownerAvatar}>
                            {donation.recipientData.ownerName.charAt(0).toUpperCase()}
                        </div>
                        <div className={styles.ownerInfo}>
                            <p className={styles.ownerTitle}>Хозяин</p>
                            <p className={styles.ownerName}>{donation.recipientData.ownerName}</p>
                        </div>
                    </div>
                    <div className={styles.lineItem}>
                        <div className={cn(styles.lineTitle, { [styles.leftMargin]: true })}>
                            <div className={styles.lineTitleIcon}>
                                <Location />
                            </div>
                            <p className={styles.lineTitleText}>Регион</p>
                        </div>
                        <p className={styles.location}>
                            {donation.recipientData.regions
                                ?.map((lock) => locationsDict.filter(({ value }) => value === lock)[0]?.label)
                                .join(', ')}
                        </p>
                    </div>
                </div>
            </div>
            {(!!donation.recipientData?.advancedInfo?.description?.length ||
                !!donation.recipientData?.advancedInfo?.photoUrls?.[0]) && (
                <Accordion className={styles.accordion} icon={<Pin />} title='Дополнительная информация'>
                    <div className={styles.itemText}>{donation.recipientData?.advancedInfo?.description}</div>
                    {donation.recipientData?.advancedInfo?.photoUrls?.[0] && (
                        <img
                            alt='Фото рецепиента'
                            className={styles.bloodRequestPhoto}
                            src={donation.recipientData.advancedInfo?.photoUrls[0]}
                        />
                    )}
                </Accordion>
            )}
            <div className={styles.settings}>
                <div className={styles.setting}>
                    <p className={styles.text}>
                        Бонусы
                        <br />
                        Портала
                    </p>
                    <div className={styles.icons}>
                        <div className={styles.bonusIcon}>
                            <PrioritySearch />
                        </div>
                        <div className={styles.bonusIcon}>
                            <Certificates />
                        </div>
                    </div>
                </div>
                <div className={styles.setting}>
                    <p className={styles.text}>
                        Ваши
                        <br />
                        условия
                    </p>
                    <div className={styles.icons}>
                        {donation.applicationData.compensationType === CompensationType.FOOD && (
                            <div className={cn(styles.rewardFeedIcon, { [styles.button]: true })}>
                                <div className={styles.icon}>
                                    <Bone />
                                </div>
                            </div>
                        )}
                        {donation.applicationData.compensationType === CompensationType.FREE && (
                            <div className={cn(styles.rewardFreeIcon, { [styles.button]: true })}>
                                <p className={styles.sum}>0</p>
                                <p className={styles.descr}>₽</p>
                            </div>
                        )}
                        {donation.applicationData.compensationType === CompensationType.PAID && (
                            <div className={cn(styles.rewardNotFreeIcon, { [styles.button]: true })}>
                                <p className={styles.descr}>₽</p>
                            </div>
                        )}
                        {donation.applicationData.taxiCompensation && (
                            <div className={cn(styles.taxiIcon, { [styles.button]: true })}>
                                <div className={styles.icon}>
                                    <Taxi />
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </Layout>
    );
};

export default CompletedDonation;
