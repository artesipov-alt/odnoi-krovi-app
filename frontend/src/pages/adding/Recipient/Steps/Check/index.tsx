import { Button } from '@mui/material';
import cn from 'classnames';
import AccordionArrow from 'imgs/svg/accordionArrow';
import BackArrow from 'imgs/svg/backArrow';
import Blood from 'imgs/svg/blood';
import BloodComponents from 'imgs/svg/bloodComponents';
import BloodVolume from 'imgs/svg/bloodVolume';
import Location from 'imgs/svg/location';
import MiniPaw from 'imgs/svg/miniPaw';
import Pin from 'imgs/svg/pin';
import { FC, useState } from 'react';

import { Dict } from 'api/reference';
import { PetType } from 'api/types';
import ImgEditor from 'components/ImgEditor';

import styles from './Check.module.less';

type Props = {
    name: string;
    weight: string;
    petType: string;
    bloodGroup: string;
    photo: File | null;
    locations: string[];
    bloodVolume: string;
    description: string;
    locationsDict: Dict[];
    onBackClick: () => void;
    bloodComponents: string[];
    bloodComponentsDict: Dict[];
    desiredBloodGroups: string[];
    notifyOfSmallDonors: boolean;
    bloodGroupDict: Record<PetType, Dict[]>;
    onConfirmButtonClick: (step: number) => void;
};

const Check: FC<Props> = ({
    name,
    photo,
    weight,
    petType,
    locations,
    bloodGroup,
    onBackClick,
    description,
    bloodVolume,
    locationsDict,
    bloodGroupDict,
    bloodComponents,
    desiredBloodGroups,
    bloodComponentsDict,
    notifyOfSmallDonors,
    onConfirmButtonClick,
}) => {
    const [isAccordionOpen, setIsAccordionOpen] = useState(false);

    const onAccordionClickHandler = () => {
        setIsAccordionOpen((prevState) => !prevState);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(4);
    };

    return (
        <div className={styles.wrapper}>
            <div className={styles.container}>
                <ImgEditor
                    showStub
                    src={photo}
                    name={name}
                    weight={weight}
                    petType={petType}
                    className={styles.photo}
                    bloodGroup={bloodGroupDict[petType].filter(({ value }) => value === Number(bloodGroup))[0].label}
                />
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
                                    {bloodGroupDict[petType].filter(({ value }) => value === Number(group))[0].label}
                                </div>
                            ))}
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
                                {bloodComponentsDict.filter(({ value }) => value === Number(comp))[0].label}
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
                                .map((lock) => locationsDict.filter(({ value }) => value === Number(lock))[0].label)
                                .join(' ')}
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
                {!!description.length && (
                    <div className={styles.accordion}>
                        <div
                            onClick={onAccordionClickHandler}
                            className={cn(styles.infoItemTitle, { [styles.noMargin]: true })}
                        >
                            <div className={styles.icon}>
                                <Pin />
                            </div>
                            <div className={styles.text}>Дополнительная информация</div>
                            <div
                                className={cn(styles.icon, {
                                    [styles.accordionIcon]: true,
                                    [styles.isOpen]: isAccordionOpen,
                                })}
                            >
                                <AccordionArrow />
                            </div>
                        </div>
                        {isAccordionOpen && (
                            <div className={cn(styles.itemText, { [styles.accordionText]: true })}>{description}</div>
                        )}
                    </div>
                )}
                <div className={styles.buttons}>
                    <Button
                        onClick={onBackClick}
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
            </div>
        </div>
    );
};

export default Check;
