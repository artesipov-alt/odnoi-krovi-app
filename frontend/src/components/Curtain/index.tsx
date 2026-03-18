import Button from '@mui/material/Button';
import cn from 'classnames';
import { FC, MouseEvent, ReactNode } from 'react';

import styles from './Curtain.module.less';

type Props = {
    title: ReactNode;
    subTitle?: ReactNode;
    onClose?: () => void;
    onCancel?: () => void;
    onConfirm?: () => void;
    columnOfButtons?: boolean;
    cancelButtonTitle?: string;
    noRednerButtons?: boolean;
    confirmButtonTitle?: string;
    isDisableCancelButton?: boolean;
    shouldCloseByWrapperClick?: boolean;
    children?: ReactNode;
};

const Curtain: FC<Props> = ({
    title,
    onClose,
    subTitle,
    onCancel,
    children,
    onConfirm,
    noRednerButtons,
    columnOfButtons,
    cancelButtonTitle,
    confirmButtonTitle,
    isDisableCancelButton,
    shouldCloseByWrapperClick,
}) => {
    const onWrapperClickHandler = (e: MouseEvent<HTMLDivElement>) => {
        if (shouldCloseByWrapperClick && e.target === e.currentTarget && onClose) {
            onClose();
        }
    };

    return (
        <div className={styles.wrapper} onClick={onWrapperClickHandler}>
            <div className={styles.content}>
                <h1 className={styles.title}>{title}</h1>
                {subTitle && <div className={styles.subTitle}>{subTitle}</div>}
                {children}
                {!noRednerButtons && (
                    <div className={cn(styles.buttons, { [styles.column]: columnOfButtons })}>
                        <Button
                            onClick={onCancel}
                            fullWidth={columnOfButtons}
                            variant='contained'
                            className={cn(styles.button, { [styles.disabled]: isDisableCancelButton })}
                        >
                            {cancelButtonTitle}
                        </Button>
                        {columnOfButtons ? (
                            <p className={styles.link} onClick={onConfirm}>
                                {confirmButtonTitle || 'Вернуться'}
                            </p>
                        ) : (
                            <Button onClick={onConfirm} className={styles.button} variant='contained'>
                                {confirmButtonTitle}
                            </Button>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default Curtain;
