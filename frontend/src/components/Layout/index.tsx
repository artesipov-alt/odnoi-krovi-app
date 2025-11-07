import { FC, ReactNode } from 'react';

import styles from './Layout.module.less';

type Props = {
    children: ReactNode;
};

const Layout: FC<Props> = ({ children }) => <div className={styles.wrapper}>{children}</div>;

export default Layout;
