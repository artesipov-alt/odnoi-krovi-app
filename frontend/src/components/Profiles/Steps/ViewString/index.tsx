import cn from 'classnames';
import Info from 'imgs/svg/info';
import Lock from 'imgs/svg/lock';
import { FC, useEffect, useState } from 'react';

import styles from './ViewString.module.less';

type Props = {
    name: string;
    value: string;
    descr?: string;
    tooltip?: string;
    withLock?: boolean;
    noAlignCanter?: boolean;
};

const ViewString: FC<Props> = ({ name, value, noAlignCanter, tooltip, descr, withLock }) => {
    const [isOpenTooltip, setIsOpenTooltip] = useState(false);

    const toggleTooltip = (e: React.MouseEvent) => {
        e.stopPropagation();

        setIsOpenTooltip((prevState) => !prevState);
    };

    useEffect(() => {
        const onOutsideClickHandler = () => {
            setIsOpenTooltip(false);
        };

        window.addEventListener('click', onOutsideClickHandler);

        return () => {
            window.removeEventListener('click', onOutsideClickHandler);
        };
    }, []);

    return (
        <>
            <div className={cn(styles.container, { [styles.noAlignCanter]: noAlignCanter, [styles.isDescr]: !!descr })}>
                <div className={cn(styles.nameWrapper, { [styles.withLock]: withLock })}>
                    {withLock && (
                        <div className={styles.lockIcon}>
                            <Lock />
                        </div>
                    )}
                    <div className={styles.name}>{name}</div>
                    {!!tooltip && (
                        <div onClick={toggleTooltip} className={styles.infoIcon}>
                            <Info />
                        </div>
                    )}
                    {isOpenTooltip && <div className={styles.tooltip}>{tooltip}</div>}
                </div>
                <div className={styles.value}>{value}</div>
            </div>
            {descr && <div className={styles.descr}>{descr}</div>}
        </>
    );
};

export default ViewString;
