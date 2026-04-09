import cn from 'classnames';
import AccordionArrow from 'imgs/svg/accordionArrow';
import { FC, ReactNode, useState } from 'react';

import styles from './Accordion.module.less';

type Props = {
    title: string;
    icon: ReactNode;
    className?: string;
    children: ReactNode;
    onToggle?: (isOpen: boolean) => void;
};

const Accordion: FC<Props> = ({ children, title, icon, onToggle, className }) => {
    const [isOpen, setIsOpen] = useState(false);

    const onClickHandler = () => {
        setIsOpen((prev) => !prev);

        onToggle?.(!isOpen);
    };

    return (
        <div className={cn(styles.accordion, className)}>
            <div onClick={onClickHandler} className={cn(styles.infoItemTitle, { [styles.noMargin]: true })}>
                <div className={styles.icon}>{icon}</div>
                <div className={styles.title}>{title}</div>
                <div
                    className={cn(styles.icon, {
                        [styles.accordionIcon]: true,
                        [styles.isOpen]: isOpen,
                    })}
                >
                    <AccordionArrow />
                </div>
            </div>
            {isOpen && children}
        </div>
    );
};

export default Accordion;
