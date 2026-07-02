import Button from '@mui/material/Button';
import cn from 'classnames';
import { ChangeEvent, FC, useState } from 'react';
import { useNavigate } from 'react-router';

import { queryClient } from 'api/queryClient';

import TextField from '../TextField';
import styles from './RejectedForm.module.less';

export enum RejectView {
    FINAL = 'final',
    CANCEL = 'cancel',
    NOT_CONFIRM = 'notConfirm',
}

type Props = {
    userId: string;
    view: RejectView;
    isDonor?: boolean;
    onBack: () => void;
    onSubmit: (reason: string) => void;
};

const OTHER = 'Другое';
const MAX_LETTERS = 150;

const items = {
    [RejectView.NOT_CONFIRM]: [
        'Донация пока не проведена',
        'Указан неверный объем донации',
        'Указана неверная группа крови донора',
        'В донации участвовал другой донор',
        OTHER,
    ],
    [RejectView.CANCEL]: [
        'Не удалось связаться с хозяином',
        'Не подошли условия',
        'Похоже на мошенничество',
        'Изменились планы',
        OTHER,
    ],
};

const RejectedForm: FC<Props> = ({ userId, onSubmit, onBack, view, isDonor = false }) => {
    const navigate = useNavigate();

    const [otherValue, setOtherValue] = useState<string>('');
    const [activeItem, setActiveItem] = useState<string>(view === RejectView.FINAL ? '' : items[view][0]);

    const onItemClickHandler = (item: string) => () => {
        setActiveItem(item);
        setOtherValue('');
    };

    const onSubmitClickHandler = () => {
        onSubmit(activeItem === OTHER ? otherValue : activeItem);
    };

    const onOtherValueChangeHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setOtherValue(value);
    };

    const onToPetsClickHandler = async () => {
        await queryClient.invalidateQueries({ queryKey: ['pets', userId] });

        navigate(isDonor ? '/' : '/owner#recipient');
    };

    const getTitle = () => {
        switch (view) {
            case RejectView.CANCEL: {
                return (
                    <>
                        Очень жаль!
                        <br />
                        Почему отменилась донация?
                    </>
                );
            }
            case RejectView.NOT_CONFIRM: {
                return (
                    <>
                        Почему не подтверждаете
                        <br />
                        донацию?
                    </>
                );
            }
            default: {
                return (
                    <>
                        Спасибо за
                        <br />
                        информацию!
                    </>
                );
            }
        }
    };

    return (
        <div className={styles.wrapper}>
            <div className={styles.content}>
                <h1 className={cn(styles.title, { [styles.isFinal]: view === RejectView.FINAL })}>{getTitle()}</h1>
                {view === RejectView.FINAL ? (
                    <>
                        <p className={styles.finalDescr}>
                            {isDonor ? 'можете запланировать новую донацию' : 'можете найти другого донора'}
                        </p>
                        <Button variant='contained' className={styles.button} onClick={onToPetsClickHandler}>
                            К питомцам
                        </Button>
                    </>
                ) : (
                    <>
                        <div>
                            {items[view].map((item) => (
                                <div className={styles.item} key={item} onClick={onItemClickHandler(item)}>
                                    <div className={cn(styles.radio, { [styles.checked]: activeItem === item })} />
                                    <p className={styles.descr}>{item}</p>
                                </div>
                            ))}
                            {activeItem === OTHER && (
                                <div className={styles.textareaWrapper}>
                                    <TextField
                                        multiline
                                        name='other'
                                        value={otherValue}
                                        maxLength={MAX_LETTERS}
                                        placeholder='Укажите причину'
                                        htmlInputClass={styles.textarea}
                                        onChange={onOtherValueChangeHandler}
                                    />
                                    <div
                                        className={cn(styles.counter, {
                                            [styles.bigText]: otherValue.length === MAX_LETTERS,
                                        })}
                                    >
                                        {otherValue.length}/{MAX_LETTERS}
                                    </div>
                                </div>
                            )}
                        </div>
                        <Button
                            variant='contained'
                            onClick={onSubmitClickHandler}
                            className={cn(styles.button, { [styles.disabled]: activeItem === OTHER && !otherValue })}
                        >
                            {view === RejectView.CANCEL ? 'Отправить' : 'Отправить и вернуться к донации'}
                        </Button>
                        <p className={styles.backLink} onClick={onBack}>
                            Вернуться
                        </p>
                    </>
                )}
            </div>
        </div>
    );
};

export default RejectedForm;
