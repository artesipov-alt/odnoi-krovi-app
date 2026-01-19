import cn from 'classnames';
import Caution from 'imgs/svg/caution';
import { FC, ReactNode } from 'react';

import styles from './Alert.module.less';

export enum View {
    WARNING = 'warning',
    INFO_WITHOUT_ICON = 'infoWithoutIcon',
}

type Props = {
    text: ReactNode;
    view?: View;
    className?: string;
};

const Alert: FC<Props> = ({ text, className, view = View.WARNING }) => (
    <div className={cn(styles.wrapper, className, { [styles.infoWithoutIcon]: view === View.INFO_WITHOUT_ICON })}>
        <div className={cn(styles.logo, { [styles.infoWithoutIcon]: view === View.INFO_WITHOUT_ICON })}>
            <Caution />
        </div>
        {text}
    </div>
);

export default Alert;
