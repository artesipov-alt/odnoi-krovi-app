import cn from 'classnames';
import Caution from 'imgs/svg/caution';
import { FC, ReactNode } from 'react';

import styles from './Alert.module.less';

export enum View {
    INFO = 'info',
    WARNING = 'warning',
}

type Props = {
    text: ReactNode;
    view?: View;
    className?: string;
};

const Alert: FC<Props> = ({ text, className, view = View.WARNING }) => (
    <div className={cn(styles.wrapper, className, { [styles.info]: view === View.INFO })}>
        <div className={cn(styles.logo, { [styles.info]: view === View.INFO })}>
            <Caution />
        </div>
        {text}
    </div>
);

export default Alert;
