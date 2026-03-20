import { Button } from '@mui/material';
import cn from 'classnames';
import { useBloodComponentsQuery, useLocationsQuery } from 'hooks/useDicts';
import catRoundStub from 'imgs/catRoundStub.png';
import dogRoundStub from 'imgs/dogRoundStub.png';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Blood from 'imgs/svg/blood';
import BloodComponents from 'imgs/svg/bloodComponents';
import Cancel from 'imgs/svg/cancel';
import Location from 'imgs/svg/location';
import MiniPaw from 'imgs/svg/miniPaw';
import MiniSinglePaw from 'imgs/svg/miniSinglePaw';
import Pin from 'imgs/svg/pin';
import Accordion from 'pages/adding/common/Accordion';
import { FC, useState } from 'react';
import { getDateFormat } from 'utils/utils';

import { GetPoolRequestResponse, Onboardings } from 'api/bloodRequest';
import { PetType } from 'api/types';
import { CircularProgress } from 'components/CircularProgress';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import SearchOnboarding, { View } from '../Onboarding';
import styles from './SearchCard.module.less';

type Props = GetPoolRequestResponse & {
    name: string;
    petId: string;
    type: PetType;
    avatar?: string;
    isLoading: boolean;
    bloodGroup: string;
    onClose: () => void;
    onBoarding?: Onboardings[];
    onBoardingConfirm: () => void;
};

const tabs = [
    { title: 'Детали', id: 0, disabled: false },
    { title: 'Выбрано', id: 1, disabled: false },
    { title: 'Получено', id: 2, disabled: true },
];

const SearchCard: FC<Props> = ({
    id,
    name,
    type,
    petId,
    avatar,
    onClose,
    regions,
    isLoading,
    photoUrls,
    bloodGroup,
    onBoarding,
    description,
    bloodGroupNames,
    onBoardingConfirm,
    bloodComponentIds,
    bloodVolumeNeeded,
    bloodVolumeReserved,
    smallPetsNotifyAllowed,
    createdAt = '',
}) => {
    const [tab, setTab] = useState<number>(0);

    const { data: locationsDict = [] } = useLocationsQuery();
    const { data: bloodComponentsDict = [] } = useBloodComponentsQuery();

    const onTabClickHandler = (tabId: number) => () => {
        setTab(tabId);
    };

    if (!(onBoarding || []).includes(Onboardings.BLOOD_CARD)) {
        return <SearchOnboarding view={View.CARD} id={id} onSucess={onBoardingConfirm} onBoarding={onBoarding} />;
    }

    return (
        <Layout>
            <div className={styles.wrapper}>
                <div className={styles.header}>
                    <div className={styles.back} onClick={onClose}>
                        <BackAngularArrow />
                    </div>
                    <h2 className={styles.name}>Поиск от {getDateFormat(new Date(createdAt))}</h2>
                </div>
                <div className={styles.tabs}>
                    {tabs.map(({ title, id: tabId, disabled }) => (
                        <div
                            key={title}
                            onClick={onTabClickHandler(tabId)}
                            className={cn(styles.tab, { [styles.active]: tab === tabId, [styles.disabled]: disabled })}
                        >
                            {title}
                        </div>
                    ))}
                </div>
                {tab === 0 && (
                    <>
                        <div className={styles.main}>
                            <div className={styles.left}>
                                <div className={cn(styles.leftItem, { [styles.name]: true })}>
                                    <div className={styles.icon}>
                                        <MiniSinglePaw />
                                    </div>
                                    <p className={styles.text}>{name.toUpperCase()}</p>
                                </div>
                                <div className={cn(styles.leftItem, { [styles.blood]: true })}>
                                    <div className={styles.subTitle}>
                                        <div className={styles.icon}>
                                            <Blood />
                                        </div>
                                        <p className={styles.text}>Ищем</p>
                                    </div>
                                    <div className={styles.bloodInfo}>
                                        <div className={styles.bloodGroup}>{bloodGroup}</div>
                                        {bloodGroupNames.some((group) => group !== bloodGroup) && ' +'}
                                        {bloodGroupNames.some((group) => group !== bloodGroup)
                                            ? bloodGroupNames
                                                  .filter((group) => group !== bloodGroup)
                                                  .map((group) => (
                                                      <div
                                                          key={group}
                                                          className={cn(styles.bloodGroup, { [styles.needed]: true })}
                                                      >
                                                          {group}
                                                      </div>
                                                  ))
                                            : ''}
                                    </div>
                                </div>
                                <div className={cn(styles.leftItem, { [styles.location]: true })}>
                                    <div className={styles.icon}>
                                        <Location />
                                    </div>
                                    <p className={styles.text}>
                                        {regions
                                            ?.map(
                                                (lock) => locationsDict.filter(({ value }) => value === lock)[0]?.label,
                                            )
                                            .join(', ')}
                                    </p>
                                </div>
                            </div>
                            <div className={styles.right}>
                                <div className={styles.avatarWrapper}>
                                    <div className={styles.avatarWrapper}>
                                        <img
                                            alt={name}
                                            className={styles.avatar}
                                            src={avatar || (type === PetType.DOG ? dogRoundStub : catRoundStub)}
                                        />
                                        <CircularProgress
                                            showDot
                                            size={180}
                                            strokeWidth={15}
                                            total={bloodVolumeNeeded}
                                            color='var(--red10, #FF2727)'
                                            current={bloodVolumeReserved || 0}
                                        />
                                        <div className={styles.neededVolume}>
                                            {bloodVolumeNeeded}
                                            <span>мл</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        <div className={styles.infoItem}>
                            <div className={styles.infoItemTitle}>
                                <div className={styles.itemIcon}>
                                    <BloodComponents />
                                </div>
                                <div className={styles.text}>Компоненты крови</div>
                            </div>
                            <div className={styles.bloodComponents}>
                                {bloodComponentIds.map((comp) => (
                                    <div key={comp} className={styles.itemText}>
                                        {bloodComponentsDict.filter(({ value }) => value === comp)[0]?.label}
                                    </div>
                                ))}
                            </div>
                        </div>
                        <div className={styles.infoItem}>
                            <div className={cn(styles.infoItemTitle, { [styles.noMargin]: true })}>
                                <div className={cn(styles.itemIcon, { [styles.miniPaw]: true })}>
                                    <MiniPaw />
                                </div>
                                <div className={styles.text}>
                                    {smallPetsNotifyAllowed
                                        ? 'Уведомлять небольших доноров'
                                        : 'Не уведомлять небольших доноров'}
                                </div>
                            </div>
                        </div>
                        {(!!description?.length || !!photoUrls?.[0]) && (
                            <>
                                <Accordion icon={<Pin />} title='Дополнительная информация'>
                                    <div className={cn(styles.itemText, { [styles.accordionText]: true })}>
                                        {description}
                                    </div>
                                    {photoUrls?.[0] && (
                                        <img
                                            alt='Фото рецепиента'
                                            className={styles.bloodRequestPhoto}
                                            src={photoUrls[0]}
                                        />
                                    )}
                                </Accordion>
                            </>
                        )}
                        <Button
                            fullWidth
                            // onClick={onConfirmButtonClickHandler}
                            className={styles.confirm}
                        >
                            Расширить поиск
                        </Button>
                        <Button
                            fullWidth
                            startIcon={<Cancel />}
                            // onClick={onConfirmButtonClickHandler}
                            className={styles.cancel}
                        >
                            Отменить поиск
                        </Button>
                    </>
                )}
                {tab === 1 && <div>В разработке</div>}
                {isLoading && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
            </div>
        </Layout>
    );
};

export default SearchCard;
