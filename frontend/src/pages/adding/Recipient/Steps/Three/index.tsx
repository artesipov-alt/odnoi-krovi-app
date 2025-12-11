import { Button } from '@mui/material';
import TextField from '@mui/material/TextField';
import { ChangeEvent, FC, useState } from 'react';

import Alert from 'components/Alert';
import ImgEditor from 'components/ImgEditor';

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
            <div className={styles.formItem}>
                <div className={styles.labelWrapper}>
                    <p className={styles.label}>Добавьте фото</p>
                    <span className={styles.subLabel}>Необязательно</span>
                </div>
                <ImgEditor src={photo} onLoad={onLoadPhoto} className={styles.photo} />
            </div>
            <div className={styles.formItem}>
                <div className={styles.labelWrapper}>
                    <p className={styles.label}>Опишите ситуацию</p>
                    <span className={styles.subLabel}>Необязательно</span>
                </div>
                <div className={styles.textFieldWrapper}>
                    <TextField
                        fullWidth
                        multiline
                        name='description'
                        value={descrValue}
                        onBlur={onDescriptionBlurHandler}
                        placeholder='Почему вы ищете помощь?'
                        onChange={onDescriptionChangeHandler}
                        slotProps={{
                            htmlInput: { className: styles.input, maxLength: MAX_LETTERS },
                            input: { className: styles.inputWrapper },
                        }}
                        sx={{
                            '& .MuiOutlinedInput-root': {
                                '& .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#dee2e9',
                                },
                                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#dee2e9',
                                    borderWidth: '1px',
                                },
                            },
                        }}
                    />
                    <div className={styles.counter}>
                        {descrValue.length}/{MAX_LETTERS}
                    </div>
                </div>
            </div>
            <Button fullWidth className={styles.confirm} onClick={onConfirmButtonClickHandler}>
                Далее
            </Button>
        </>
    );
};

export default Three;
