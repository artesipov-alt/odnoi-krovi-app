import cn from 'classnames';
import { FC, ReactNode } from 'react';

import styles from './FormItem.module.less';

type Props = {
    title: string;
    className?: string;
    subtitle?: ReactNode;
    children?: ReactNode;
};

const FormItem: FC<Props> = ({ title, children, subtitle, className }) => (
    <div className={cn(styles.formItem, className)}>
        <div className={styles.labelWrapper}>
            <p className={styles.label}>{title}</p>
            <span className={styles.subLabel}>{subtitle}</span>
        </div>
        <div>{children}</div>
    </div>
);

export default FormItem;
