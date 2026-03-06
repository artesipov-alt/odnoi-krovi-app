import cn from 'classnames';
import Info from 'imgs/svg/info';
import { FC, ReactNode, useEffect, useState } from 'react';

import styles from './FormItem.module.less';

type Props = {
    title: string;
    tooltip?: string;
    className?: string;
    subtitle?: ReactNode;
    children?: ReactNode;
};

const FormItem: FC<Props> = ({ title, children, subtitle, className, tooltip }) => {
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
        <div className={cn(styles.formItem, className)}>
            <div className={styles.labelWrapper}>
                <div className={styles.label}>
                    {title}
                    {!!tooltip && (
                        <div onClick={toggleTooltip} className={styles.infoIcon}>
                            <Info />
                        </div>
                    )}
                </div>
                <span className={styles.subLabel}>{subtitle}</span>
                {isOpenTooltip && <div className={styles.tooltip}>{tooltip}</div>}
            </div>
            <div>{children}</div>
        </div>
    );
};

export default FormItem;
