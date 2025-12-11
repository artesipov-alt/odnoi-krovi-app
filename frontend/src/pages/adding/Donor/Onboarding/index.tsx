import { Button } from '@mui/material';
import cn from 'classnames';
import Caution from 'imgs/svg/caution';
import ExclamationSquare from 'imgs/svg/exclamationSquare';
import MainLogo from 'imgs/svg/mainLogo';
import { FC, useState } from 'react';

import Alert from 'components/Alert';

import styles from './Onboarding.module.less';

type Props = {
    onFinish: () => void;
    onBackToStart: () => void;
};

const items = {
    1: [
        'Возраст от 1 года до 8 лет',
        'Вакцинирован от инфекций и бешенства (после вакцинации прошло не более года и не менее месяца)',
        'Обработан от гельминтов и эктопаразитов',
        'Клинически здоров (не имеет хронических, инфекционных, аутоиммунных и онкологических заболеваний)',
        'Не принимает препараты в рамках назначенного курса лечения',
        'Не восстанавливается после хирургического вмешательства',
        'Вне периодов беременности, лактации и течки',
        'Не переливалась кровь от других животных',
        'С последней донации прошло более 2-х месяцев',
    ],
    2: [
        'Портал борется с черным донорством - фото помогает убедиться, что питомец не участвует в донациях слишком часто',
        'Если донор не похож на фото, в донации могут отказать',
    ],
};

const Onboarding: FC<Props> = ({ onFinish, onBackToStart }) => {
    const [step, setStep] = useState(1);

    const isFirstStep = step === 1;

    const onConfirmClickHandler = () => {
        if (step === 1) {
            setStep(2);

            return;
        }

        onFinish();
    };

    return (
        <div className={cn(styles.wrapper, { [styles.secondStep]: !isFirstStep })}>
            <div className={cn(styles.header, { [styles.secondStep]: !isFirstStep })}>
                <div className={styles.logo}>
                    <MainLogo />
                </div>
                <h1 className={styles.title}>
                    {isFirstStep ? (
                        <>
                            Как питомцу
                            <br />
                            стать донором?
                        </>
                    ) : (
                        'Обязательное фото'
                    )}
                </h1>
            </div>
            <div className={cn(styles.items, { [styles.secondStep]: !isFirstStep })}>
                {items[step].map((item, i) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <div key={i} className={styles.item}>
                        <div className={cn(styles.logo, { [styles.item]: true })}>
                            {isFirstStep ? <ExclamationSquare /> : <Caution />}
                        </div>
                        <p className={styles.itemText}>{item}</p>
                    </div>
                ))}
            </div>
            {isFirstStep && (
                <Alert
                    className={styles.alert}
                    text='Будьте готовы предъявить ветеринарный паспорт для подтверждения внесенных сведений'
                />
            )}
            <div className={styles.buttons}>
                <Button fullWidth onClick={onConfirmClickHandler} className={styles.confirm}>
                    {isFirstStep ? 'Хорошо' : 'Далее'}
                </Button>
                <p className={styles.link} onClick={onBackToStart}>
                    Вернуться
                </p>
            </div>
            {!isFirstStep && <img className={styles.photos} alt='pets' src='src/imgs/onboarding.png' />}
        </div>
    );
};

export default Onboarding;
