import cn from 'classnames';
import { FC, ReactNode } from 'react';

import styles from './Layout.module.less';

type Props = {
    className?: string;
    children: ReactNode;
};

const Layout: FC<Props> = ({ children, className }) => <div className={cn(styles.wrapper, className)}>{children}</div>;

export default Layout;
