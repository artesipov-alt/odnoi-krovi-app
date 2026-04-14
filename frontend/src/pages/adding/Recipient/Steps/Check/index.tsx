import { Button } from '@mui/material';
import cn from 'classnames';
import { BloodAndBreedGroupsDict } from 'hooks/useDicts';
import BackArrow from 'imgs/svg/backArrow';
import Blood from 'imgs/svg/blood';
import BloodComponents from 'imgs/svg/bloodComponents';
import BloodVolume from 'imgs/svg/bloodVolume';
import Location from 'imgs/svg/location';
import MiniPaw from 'imgs/svg/miniPaw';
import Pin from 'imgs/svg/pin';
import Accordion from 'pages/adding/common/Accordion';
import { FC } from 'react';

import { Dict } from 'api/reference';
import ImgEditor from 'components/ImgEditor';
import Loading from 'components/Loading';

import styles from './Check.module.less';

type Props = {
    name: string;
    weight: string;
    petType: string;
    photoUrl?: string;
    bloodGroup: string;
    photo: File | null;
    isLoading: boolean;
    locations: string[];
    bloodVolume: string;
    description: string;
    locationsDict: Dict[];
    bloodComponents: string[];
    bloodComponentsDict: Dict[];
    desiredBloodGroups: string[];
    notifyOfSmallDonors: boolean;
    bloodRequestPhoto: File | null;
    includeUnknownBloodGroup: boolean;
    bloodGroupDict: BloodAndBreedGroupsDict;
    onConfirmButtonClick: (step: number) => void;
};

const Check: FC<Props> = ({
    name,
    photo,
    weight,
    petType,
    photoUrl,
    isLoading,
    locations,
    bloodGroup,
    description,
    bloodVolume,
    locationsDict,
    bloodGroupDict,
    bloodComponents,
    bloodRequestPhoto,
    desiredBloodGroups,
    bloodComponentsDict,
    notifyOfSmallDonors,
    onConfirmButtonClick,
    includeUnknownBloodGroup,
}) => {
    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(4);
    };

    const onBackClickHandler = () => {
        onConfirmButtonClick(0);
    };

    return (
        <>
            <div className={styles.header} />
            <div className={styles.container}>
                <div className={styles.photo}>
                    <ImgEditor
                        showStub
                        src={photo}
                        name={name}
                        weight={weight}
                        petType={petType}
                        serverSrc={photoUrl}
                        bloodGroup={bloodGroupDict[petType]?.filter(({ value }) => value === bloodGroup)?.[0]?.label}
                    />
                </div>
                <div className={styles.titleWrapper}>
                    <h2 className={styles.title}>Проверьте все поля</h2>
                    <p className={styles.descr}>
                        После подтверждения изменить
                        <br />
                        критерии поиска будет нельзя
                    </p>
                </div>
                <div className={styles.bloodBox}>
                    <div className={styles.bloodBoxItem}>
                        <div className={styles.icon}>
                            <BloodVolume />
                        </div>
                        <div className={styles.text}>{bloodVolume} мл</div>
                    </div>
                    <div className={styles.bloodBoxItem}>
                        <div className={styles.icon}>
                            <Blood />
                        </div>
                        <div className={styles.bloodGroupWrapper}>
                            {desiredBloodGroups.map((group) => (
                                <div key={group} className={styles.bloodGroup}>
                                    {bloodGroupDict[petType]?.filter(({ value }) => value === group)?.[0]?.label}
                                </div>
                            ))}
                            {includeUnknownBloodGroup && <div className={styles.bloodGroup}>?</div>}
                        </div>
                    </div>
                </div>
                <div className={styles.infoItem}>
                    <div className={styles.infoItemTitle}>
                        <div className={styles.icon}>
                            <BloodComponents />
                        </div>
                        <div className={styles.text}>Компоненты крови</div>
                    </div>
                    <div className={styles.bloodComponents}>
                        {bloodComponents.map((comp) => (
                            <div key={comp} className={styles.itemText}>
                                {bloodComponentsDict.filter(({ value }) => value === comp)[0].label}
                            </div>
                        ))}
                    </div>
                </div>
                <div className={styles.infoItem}>
                    <div className={cn(styles.infoItemTitle, { [styles.noMargin]: true })}>
                        <div className={styles.icon}>
                            <Location />
                        </div>
                        <div className={styles.text}>
                            {locations
                                .map((lock) => locationsDict.filter(({ value }) => value === lock)[0].label)
                                .join(', ')}
                        </div>
                    </div>
                </div>
                <div className={styles.infoItem}>
                    <div className={cn(styles.infoItemTitle, { [styles.noMargin]: true })}>
                        <div className={cn(styles.icon, { [styles.miniPaw]: true })}>
                            <MiniPaw />
                        </div>
                        <div className={styles.text}>
                            {notifyOfSmallDonors ? 'Уведомлять небольших доноров' : 'Не уведомлять небольших доноров'}
                        </div>
                    </div>
                </div>
                {(!!description.length || !!bloodRequestPhoto) && (
                    <>
                        <Accordion icon={<Pin />} title='Дополнительная информация'>
                            <div className={cn(styles.itemText, { [styles.accordionText]: true })}>{description}</div>
                            {bloodRequestPhoto && (
                                <img
                                    alt='Фото рецепиента'
                                    className={styles.bloodRequestPhoto}
                                    src={URL.createObjectURL(bloodRequestPhoto)}
                                />
                            )}
                        </Accordion>
                    </>
                )}
                <div className={styles.buttons}>
                    <Button
                        onClick={onBackClickHandler}
                        className={cn(styles.button, styles.back)}
                        startIcon={
                            <div className={cn(styles.icon, { [styles.backArrow]: true })}>
                                <BackArrow />
                            </div>
                        }
                    >
                        Редактировать
                    </Button>
                    <Button onClick={onConfirmButtonClickHandler} className={cn(styles.button, styles.ok)}>
                        Все верно
                    </Button>
                </div>
                {isLoading && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
            </div>
        </>
    );
};

export default Check;
