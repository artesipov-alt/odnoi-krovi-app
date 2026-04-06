import Button from '@mui/material/Button';
import cn from 'classnames';
import { CSSProperties, FC, MouseEvent, ReactNode } from 'react';

import styles from './Curtain.module.less';

type Props = {
    title?: ReactNode;
    subTitle?: ReactNode;
    onClose?: () => void;
    onCancel?: () => void;
    onConfirm?: () => void;
    backgroundImage?: string;
    contentBorderRadius?: CSSProperties['borderRadius'];
    contentOverflow?: CSSProperties['overflow'];
    columnOfButtons?: boolean;
    noRednerButtons?: boolean;
    cancelButtonTitle?: string;
    subTitleClassName?: string;
    confirmButtonClassName?: string;
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
    backgroundImage,
    contentBorderRadius,
    contentOverflow,
    noRednerButtons,
    columnOfButtons,
    subTitleClassName,
    cancelButtonTitle,
    confirmButtonTitle,
    isDisableCancelButton,
    confirmButtonClassName,
    shouldCloseByWrapperClick,
}) => {
    const onWrapperClickHandler = (e: MouseEvent<HTMLDivElement>) => {
        if (shouldCloseByWrapperClick && e.target === e.currentTarget && onClose) {
            onClose();
        }
    };

    const contentStyle: CSSProperties | undefined =
        backgroundImage || contentBorderRadius !== undefined
            ? ({
                  ...(backgroundImage
                      ? ({
                            ['--curtain-bg-image' as any]: `url(${backgroundImage})`,
                        } satisfies CSSProperties)
                      : {}),
                  ...(contentBorderRadius !== undefined ? { borderRadius: contentBorderRadius } : {}),
                  ...(contentOverflow !== undefined ? { overflow: contentOverflow } : {}),
              } satisfies CSSProperties)
            : undefined;

    return (
        <div className={styles.wrapper} onClick={onWrapperClickHandler}>
            <div className={styles.content} style={contentStyle}>
                {!!title && <h1 className={styles.title}>{title}</h1>}
                {subTitle && <div className={cn(styles.subTitle, subTitleClassName)}>{subTitle}</div>}
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
                            <Button
                                onClick={onConfirm}
                                variant='contained'
                                className={cn(styles.button, confirmButtonClassName)}
                            >
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
