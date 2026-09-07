import cn from 'classnames';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Caution from 'imgs/svg/caution';
import { FC, ReactNode, useState } from 'react';

import styles from './Alert.module.less';

export enum View {
    WARNING = 'warning',
    INFO_WITHOUT_ICON = 'infoWithoutIcon',
}

type Props = {
    view?: View;
    text: ReactNode;
    reason?: string;
    className?: string;
};

const Alert: FC<Props> = ({ text, className, reason, view = View.WARNING }) => {
    const [isReasonOpen, setIsReasonOpen] = useState(false);

    const toggleReason = () => {
        setIsReasonOpen((prevState) => !prevState);
    };

    return (
        <div className={cn(styles.alert, className, { [styles.infoWithoutIcon]: view === View.INFO_WITHOUT_ICON })}>
            <div className={cn(styles.wrapper)}>
                <div className={cn(styles.logo, { [styles.infoWithoutIcon]: view === View.INFO_WITHOUT_ICON })}>
                    <Caution />
                </div>
                {text}
            </div>
            {reason && (
                <div onClick={toggleReason} className={styles.reason}>
                    {isReasonOpen ? (
                        `Причина: ${reason}`
                    ) : (
                        <div className={styles.start}>
                            Узнать причину
                            <div className={styles.icon}>
                                <BackAngularArrow />
                            </div>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default Alert;
