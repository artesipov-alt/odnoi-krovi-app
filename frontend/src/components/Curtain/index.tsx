import Button from '@mui/material/Button';
import cn from 'classnames';
import { FC, ReactNode } from 'react';

import styles from './Curtain.module.less';

type Props = {
    title: ReactNode;
    subTitle?: ReactNode;
    onCancel: () => void;
    onConfirm: () => void;
    columnOfButtons?: boolean;
    cancelButtonTitle: string;
    confirmButtonTitle: string;
    isDisableCancelButton?: boolean;
    children?: ReactNode;
};

const Curtain: FC<Props> = ({
    title,
    subTitle,
    onCancel,
    children,
    onConfirm,
    columnOfButtons,
    cancelButtonTitle,
    confirmButtonTitle,
    isDisableCancelButton,
}) => (
    <div className={styles.wrapper}>
        <div className={styles.content}>
            <h1 className={styles.title}>{title}</h1>
            {subTitle && <div className={styles.subTitle}>{subTitle}</div>}
            {children}
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
                        Вернуться
                    </p>
                ) : (
                    <Button onClick={onConfirm} className={styles.button} variant='contained'>
                        {confirmButtonTitle}
                    </Button>
                )}
            </div>
        </div>
    </div>
);

export default Curtain;
