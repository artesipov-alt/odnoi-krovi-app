import { Button } from '@mui/material';
import cn from 'classnames';
import Success from 'imgs/svg/success';
import { FC } from 'react';

import { PetType } from 'api/types';
import DatePicker from 'components/DatePicker';

import { Analiz } from '../../types';
import styles from './Fifth.module.less';

type Props = {
    petType: string;
    leukemia: Analiz;
    babesiosis: Analiz;
    dirofilaria: Analiz;
    ehrlichiosis: Analiz;
    anaplasmosis: Analiz;
    hemoplasmosis: Analiz;
    bartonellosis: Analiz;
    immunodeficiency: Analiz;
    onConfirmButtonClick: (step: number) => void;
    onChangeAnaliz: (type: string, analizName: string, value: Date | null) => void;
};

const Fifth: FC<Props> = ({
    petType,
    leukemia,
    babesiosis,
    dirofilaria,
    ehrlichiosis,
    anaplasmosis,
    bartonellosis,
    hemoplasmosis,
    onChangeAnaliz,
    immunodeficiency,
    onConfirmButtonClick,
}) => {
    const isDog = petType === PetType.DOG;

    const onDateChangeHandler = (type: string, name: string) => (value: Date | null) => {
        onChangeAnaliz(type, name, value);
    };

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(5);
    };

    return (
        <>
            <h2 className={cn(styles.title, { [styles.isDog]: isDog })}>
                Укажите даты последних
                <br />
                проведенных анализов
            </h2>
            {(isDog
                ? [babesiosis, dirofilaria, hemoplasmosis, bartonellosis, ehrlichiosis, anaplasmosis]
                : [leukemia, immunodeficiency, hemoplasmosis, bartonellosis]
            ).map((analiz) => (
                <div key={analiz.name} className={styles.group}>
                    <p className={styles.analiz}>{analiz.name}</p>
                    {analiz.items.map((item) => (
                        <div key={analiz.type} className={styles.item}>
                            <p className={cn(styles.name, { [styles.isFilled]: !!item.value })}>{item.name}</p>
                            <div className={styles.picker}>
                                <DatePicker value={item.value} onChange={onDateChangeHandler(analiz.type, item.name)} />
                            </div>
                            {item.value && (
                                <div className={styles.icon}>
                                    <Success />
                                </div>
                            )}
                        </div>
                    ))}
                </div>
            ))}
            <Button fullWidth onClick={onConfirmButtonClickHandler} className={styles.confirm}>
                Далее
            </Button>
        </>
    );
};

export default Fifth;
