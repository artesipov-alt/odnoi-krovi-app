import { Button } from '@mui/material';
import cn from 'classnames';
import FormItem from 'pages/adding/common/FormItem';
import { ChangeEvent, FC, useState } from 'react';

import Alert from 'components/Alert';
import ImgEditor from 'components/ImgEditor';
import TextField from 'components/TextField';

import styles from './Three.module.less';

type Props = {
    photo: File | null;
    description: string;
    onLoadPhoto: (photo: File | null) => void;
    onConfirmButtonClick: (step: number) => void;
    onDescriptionChange: (newDescr: string) => void;
};

const MAX_LETTERS = 500;

const Three: FC<Props> = ({ onConfirmButtonClick, description, onDescriptionChange, photo, onLoadPhoto }) => {
    const [descrValue, setDescrValue] = useState(description);

    const onConfirmButtonClickHandler = () => {
        onConfirmButtonClick(3);
    };

    const onDescriptionChangeHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        setDescrValue(value);
    };

    const onDescriptionBlurHandler = ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
        onDescriptionChange(value);
    };

    return (
        <>
            <div className={styles.formItem}>
                <Alert
                    className={styles.alert}
                    text='Подробности могут вызвать эмоциональный отклик у хозяев доноров и увеличить шансы найти помощь'
                />
            </div>
            <FormItem title='Добавьте фото' subtitle='Необязательно'>
                <ImgEditor src={photo} onLoad={onLoadPhoto} className={styles.photo} />
            </FormItem>
            <FormItem title='Опишите ситуацию' subtitle='Необязательно'>
                <div className={styles.textFieldWrapper}>
                    <TextField
                        multiline
                        name='description'
                        value={descrValue}
                        maxLength={MAX_LETTERS}
                        htmlInputClass={styles.input}
                        inputClass={styles.inputWrapper}
                        onBlur={onDescriptionBlurHandler}
                        placeholder='Почему вы ищете помощь?'
                        onChange={onDescriptionChangeHandler}
                    />
                    <div className={cn(styles.counter, { [styles.bigText]: descrValue.length === MAX_LETTERS })}>
                        {descrValue.length}/{MAX_LETTERS}
                    </div>
                </div>
            </FormItem>
            <Button fullWidth className={styles.confirm} onClick={onConfirmButtonClickHandler}>
                Далее
            </Button>
        </>
    );
};

export default Three;
