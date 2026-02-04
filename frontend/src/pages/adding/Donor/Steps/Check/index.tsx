import { Button } from '@mui/material';
import cn from 'classnames';
import Analizes from 'imgs/svg/analizes';
import BackArrow from 'imgs/svg/backArrow';
import Health from 'imgs/svg/health';
import NoSuccess from 'imgs/svg/noSuccess';
import Params from 'imgs/svg/params';
import Processing from 'imgs/svg/processing';
import Success from 'imgs/svg/success';
import Accordion from 'pages/adding/common/Accordion';
import { FC, useRef, useState } from 'react';
import { getCorrectDeclension, getDateFormat, Variants } from 'utils/utils';

import { PetType } from 'api/types';
import ImgEditor from 'components/ImgEditor';
import Loading from 'components/Loading';

import { Analiz } from '../../types';
import styles from './Check.module.less';

type Props = {
    name: string;
    breed: string;
    weight: string;
    petType: string;
    leukemia: Analiz;
    petGender: string;
    isLoading: boolean;
    babesiosis: Analiz;
    chipNumber: string;
    photo: File | null;
    bloodGroup: string;
    petTypeCode: string;
    dirofilaria: Analiz;
    anaplasmosis: Analiz;
    ehrlichiosis: Analiz;
    surgicalList: string;
    healthStatus: string;
    bartonellosis: Analiz;
    hemoplasmosis: Analiz;
    exactDate: Date | null;
    medicationsList: string;
    livingCondition: string;
    immunodeficiency: Analiz;
    lastDonation: Date | null;
    dewormingDate: Date | null;
    reproductiveStatus: string;
    approximateDateYear: string;
    approximateDateMonth: string;
    rabiesVaccinationDate: Date | null;
    wasBloodTransfusion: boolean | null;
    infectionsVaccinationDate: Date | null;
    ectoparasitesTreatmentDate: Date | null;
    onConfirmButtonClick: (step: number) => void;
};

type Param = {
    name: string;
    analiz?: boolean;
    value: string | null;
    vaccination?: boolean;
};

const Check: FC<Props> = ({
    name,
    photo,
    breed,
    weight,
    petType,
    leukemia,
    isLoading,
    petGender,
    exactDate,
    babesiosis,
    bloodGroup,
    chipNumber,
    petTypeCode,
    dirofilaria,
    anaplasmosis,
    ehrlichiosis,
    surgicalList,
    lastDonation,
    healthStatus,
    dewormingDate,
    hemoplasmosis,
    bartonellosis,
    livingCondition,
    medicationsList,
    immunodeficiency,
    reproductiveStatus,
    wasBloodTransfusion,
    approximateDateYear,
    approximateDateMonth,
    onConfirmButtonClick,
    rabiesVaccinationDate,
    infectionsVaccinationDate,
    ectoparasitesTreatmentDate,
}) => {
    const [isOpenView, setIsOpenView] = useState<boolean>(false);

    const accordionsRef = useRef<Record<number, boolean>>({ 1: false, 2: false, 3: false, 4: false });

    const onEditClickHandler = () => {
        onConfirmButtonClick(0);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(6);
    };

    const onAccordionToggleHandler = (id: number) => (isOpen: boolean) => {
        accordionsRef.current[id] = isOpen;

        if (Object.values(accordionsRef.current).some((item) => item)) {
            setIsOpenView(true);
        }

        if (Object.values(accordionsRef.current).every((item) => !item)) {
            setIsOpenView(false);
        }
    };

    const getAge = () => {
        if (!approximateDateYear && !approximateDateMonth) {
            return null;
        }

        const years = approximateDateYear
            ? `${approximateDateYear} ${getCorrectDeclension(Variants.YEARS, +approximateDateYear)}`
            : '';
        const months = approximateDateMonth
            ? `${approximateDateMonth} ${getCorrectDeclension(Variants.MONTHS, +approximateDateMonth)}`
            : '';

        return `${years} ${months}`;
    };

    const getAnalizesValues = (analiz: Analiz) => {
        const result = analiz.items.reduce((acc, item) => {
            if (item.value) {
                acc.push(`${item.name} ${getDateFormat(item.value)}`);
            }

            return acc;
        }, [] as string[]);

        if (!result.length) {
            return ' ';
        }

        return result.join(',');
    };

    const renderAnalizString = (value: string | null) => {
        if (!value) {
            return '';
        }

        return value.split(',').map((item) => <div className={styles.analizItem}>{item}</div>);
    };

    const renderString = (params: Param[]) => (
        <div className={styles.accordionContent}>
            {params
                .filter(({ value }) => !!value)
                .map((param) => (
                    <div className={cn(styles.params, { [styles.analiz]: param.analiz })} key={param.name}>
                        <div className={cn(styles.paramName, { [styles.withIcon]: param.vaccination })}>
                            {param.vaccination && (
                                <div className={cn(styles.icon, { [styles.accStringIcon]: true })}>
                                    {param.value === ' ' ? <NoSuccess /> : <Success />}
                                </div>
                            )}
                            {param.name}
                        </div>
                        <div
                            className={cn(styles.paramValue, {
                                [styles.vaccination]: param.vaccination,
                                [styles.analiz]: param.analiz,
                            })}
                        >
                            {param.analiz ? renderAnalizString(param.value) : param.value}
                        </div>
                    </div>
                ))}
        </div>
    );

    return (
        <>
            <div className={cn(styles.header, { [styles.isOpen]: isOpenView })}>
                <h2 className={cn(styles.title, { [styles.isOpen]: isOpenView })}>Проверьте все поля</h2>
            </div>
            <div className={cn(styles.content, { [styles.isOpen]: isOpenView })}>
                <div className={styles.photo}>
                    <ImgEditor showStub src={photo} name={name} bloodGroup={bloodGroup} />
                </div>
                <div className={cn(styles.accordions, { [styles.isOpen]: isOpenView })}>
                    <Accordion title='Параметры' icon={<Params />} onToggle={onAccordionToggleHandler(1)}>
                        {renderString([
                            { name: 'Вид', value: petType },
                            { name: 'Пол', value: petGender },
                            { name: 'Вес', value: `${weight} кг` },
                            { name: '№ чипа', value: chipNumber === 'none' ? 'Отсутствует' : chipNumber },
                            {
                                name: 'Дата рождения',
                                value: exactDate ? getDateFormat(exactDate) : null,
                            },
                            { name: 'Возраст', value: getAge() },
                            { name: 'Порода', value: breed },
                            { name: 'Условия содержания', value: livingCondition },
                            { name: 'Состояние питомца', value: reproductiveStatus || null },
                        ])}
                    </Accordion>
                    <Accordion title='Здоровье' icon={<Health />} onToggle={onAccordionToggleHandler(2)}>
                        {renderString([
                            { name: 'Состояние здоровья', value: healthStatus },
                            {
                                name: 'Последняя донация',
                                value: lastDonation ? getDateFormat(lastDonation) : 'Не был донором',
                            },
                            {
                                name: 'Переливания питомцу',
                                value: wasBloodTransfusion ? 'Проводились' : 'Не проводились',
                            },
                            {
                                name: 'Прием препаратов',
                                value: medicationsList || 'Не принимает',
                            },
                            {
                                name: 'Хирургические вмешательства',
                                value: surgicalList || 'Отсутствуют',
                            },
                        ])}
                    </Accordion>
                    <Accordion title='Обработки' icon={<Processing />} onToggle={onAccordionToggleHandler(3)}>
                        {renderString([
                            {
                                vaccination: true,
                                name: 'Вакцинация от бешенства',
                                value: rabiesVaccinationDate ? getDateFormat(rabiesVaccinationDate) : ' ',
                            },
                            {
                                vaccination: true,
                                name: 'Вакцинация от инфекций',
                                value: infectionsVaccinationDate ? getDateFormat(infectionsVaccinationDate) : ' ',
                            },
                            {
                                vaccination: true,
                                name: 'Дегельминтизация',
                                value: dewormingDate ? getDateFormat(dewormingDate) : ' ',
                            },
                            {
                                vaccination: true,
                                name: 'Обработка от эктопаразитов',
                                value: ectoparasitesTreatmentDate ? getDateFormat(ectoparasitesTreatmentDate) : ' ',
                            },
                        ])}
                    </Accordion>
                    <Accordion title='Анализы' icon={<Analizes />} onToggle={onAccordionToggleHandler(4)}>
                        {renderString(
                            petTypeCode === PetType.DOG
                                ? [
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Бабезиоз',
                                          value: getAnalizesValues(babesiosis),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Дирофиляриоз',
                                          value: getAnalizesValues(dirofilaria),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Гемоплазмоз',
                                          value: getAnalizesValues(hemoplasmosis),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Бартонеллез',
                                          value: getAnalizesValues(bartonellosis),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Эрлихиоз',
                                          value: getAnalizesValues(ehrlichiosis),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Анаплазмоз',
                                          value: getAnalizesValues(anaplasmosis),
                                      },
                                  ]
                                : [
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Лейкоз',
                                          value: getAnalizesValues(leukemia),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Иммунодефицит',
                                          value: getAnalizesValues(immunodeficiency),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Гемоплазмоз',
                                          value: getAnalizesValues(hemoplasmosis),
                                      },
                                      {
                                          analiz: true,
                                          vaccination: true,
                                          name: 'Бартонеллез',
                                          value: getAnalizesValues(bartonellosis),
                                      },
                                  ],
                        )}
                    </Accordion>
                </div>
            </div>
            <div className={styles.buttons}>
                <Button
                    onClick={onEditClickHandler}
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
        </>
    );
};

export default Check;
