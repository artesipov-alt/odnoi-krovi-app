import cn from 'classnames';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import { FC, ReactNode } from 'react';

import styles from './Header.module.less';

type Props = {
    title: string;
    icon: ReactNode;
    onClose: () => void;
    isEditMode: boolean;
};

const Header: FC<Props> = ({ onClose, isEditMode, icon, title }) => (
    <div className={cn(styles.header, { [styles.isEdit]: isEditMode })}>
        <div className={styles.back} onClick={onClose}>
            <BackAngularArrow />
        </div>
        <h2 className={styles.title}>{title}</h2>
        <div className={cn(styles.titleIcon)}>{icon}</div>
    </div>
);

export default Header;
