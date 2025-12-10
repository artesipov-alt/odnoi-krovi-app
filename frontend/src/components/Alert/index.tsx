import cn from 'classnames';
import Caution from 'imgs/svg/caution';
import { FC } from 'react';

import styles from './Alert.module.less';

type Props = {
    text: string;
    className?: string;
};

const Alert: FC<Props> = ({ text, className }) => (
    <div className={cn(styles.wrapper, className)}>
        <div className={styles.logo}>
            <Caution />
        </div>
        {text}
    </div>
);

export default Alert;
