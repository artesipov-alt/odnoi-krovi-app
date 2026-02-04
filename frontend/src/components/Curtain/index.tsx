import Button from '@mui/material/Button';
import { FC, ReactNode } from 'react';

import styles from './Curtain.module.less';

type Props = {
    title: ReactNode;
    subTitle?: string;
    onCancel: () => void;
    onConfirm: () => void;
    cancelButtonTitle: string;
    confirmButtonTitle: string;
};

const Curtain: FC<Props> = ({ title, subTitle, confirmButtonTitle, cancelButtonTitle, onConfirm, onCancel }) => (
    <div className={styles.wrapper}>
        <div className={styles.content}>
            <h1 className={styles.title}>{title}</h1>
            {subTitle && <p className={styles.subTitle}>{subTitle}</p>}
            <div className={styles.buttons}>
                <Button onClick={onCancel} className={styles.button} variant='contained'>
                    {cancelButtonTitle}
                </Button>
                <Button onClick={onConfirm} className={styles.button} variant='contained'>
                    {confirmButtonTitle}
                </Button>
            </div>
        </div>
    </div>
);

export default Curtain;
