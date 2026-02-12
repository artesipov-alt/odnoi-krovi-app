import cn from 'classnames';
import BackArrow from 'imgs/svg/backArrow';
import MainLogo from 'imgs/svg/mainLogo';
import { FC } from 'react';

import styles from './Header.module.less';

type Props = {
    step: number;
    caption: string;
    stepsCount: number;
    onBackClickHandler: () => void;
};

const Header: FC<Props> = ({ step, stepsCount, onBackClickHandler, caption }) => (
    <>
        <div className={styles.header}>
            <div className={styles.title}>
                <div className={styles.logo}>
                    <MainLogo />
                </div>
                <div className={styles.descr}>
                    <h4 className={styles.step}>
                        Шаг {step} из {stepsCount}
                    </h4>
                    <span className={styles.caption}>{caption}</span>
                </div>
            </div>
            <div className={styles.back} onClick={onBackClickHandler}>
                <BackArrow />
            </div>
        </div>
        <div className={styles.progressWrapper}>
            <div
                className={cn(styles.progress, {
                    [styles.twoSteps]: stepsCount === 2,
                    [styles.fiveSteps]: stepsCount === 5,
                    [styles.one]: step === 1,
                    [styles.two]: step === 2,
                    [styles.three]: step === 3,
                    [styles.four]: step === 4,
                    [styles.five]: step === 5,
                })}
            />
        </div>
    </>
);

export default Header;
