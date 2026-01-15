import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { DesktopDatePicker as MuiDatePicker } from '@mui/x-date-pickers/DesktopDatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { ru } from 'date-fns/locale/ru';
import Cancel from 'imgs/svg/cancel';
import { FC, useState } from 'react';

import styles from './DatePicker.module.less';

type Props = {
    value: Date | null;
    backgroundColor?: string;
    onChange: (date: Date | null) => void;
};

const DatePicker: FC<Props> = ({ value, onChange, backgroundColor }) => {
    const [isOpen, setIsOpen] = useState<boolean>(false);

    const minDate = new Date();
    minDate.setFullYear(minDate.getFullYear() - 40);

    const onChangeHandler = (newValue: Date | null) => {
        onChange(newValue);
    };

    const onDeleteClickHandler = () => {
        onChange(null);
    };

    return (
        <div className={styles.picker}>
            <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ru}>
                <MuiDatePicker
                    value={value}
                    open={isOpen}
                    minDate={minDate}
                    maxDate={new Date()}
                    onChange={onChangeHandler}
                    onOpen={() => setIsOpen(true)}
                    onClose={() => setIsOpen(false)}
                    format={value ? 'dd.MM.yyyy' : 'дд.мм.гггг'}
                    slotProps={{
                        textField: {
                            fullWidth: true,
                            placeholder: 'Выберите дату',
                            onClick: () => setIsOpen(true), // Открытие при клике на поле
                            readOnly: true, // Предотвращает появление экранной клавиатуры на мобильных
                            sx: {
                                // Настройка основного фона поля
                                backgroundColor: backgroundColor || 'white',
                                borderRadius: '16px',
                                border: '1px solid #dee2e9',
                                userSelect: 'none', // Отключаем выделение текста
                                '& .MuiPickersOutlinedInput-root': {
                                    color: value ? '#8B7069' : '#ACACAC',
                                    height: '48px',
                                },
                                '& .MuiPickersOutlinedInput-notchedOutline': {
                                    border: 'none',
                                },
                                '&.Mui-focused': {
                                    backgroundColor: '#ffffff',
                                },
                                '&.Mui-focused .MuiPickersOutlinedInput-notchedOutline': {
                                    borderWidth: '0 !important',
                                },
                                // Настройка бордера через обертку поля (notchedOutline)
                                '& .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#dee2e9',
                                    borderWidth: '1px',
                                },
                                // Стили при наведении
                                '&:hover .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#dee2e9',
                                },
                            },
                        },
                        // Скрываем заголовок "Select date"
                        toolbar: {
                            hidden: true,
                        },
                        // Убираем кнопки (Ок, Отмена и т.д.)
                        actionBar: {
                            actions: [],
                        },
                    }}
                    slots={{
                        openPickerButton: () => null,
                    }}
                />
            </LocalizationProvider>
            {value && (
                <div className={styles.delete} onClick={onDeleteClickHandler}>
                    <Cancel />
                </div>
            )}
        </div>
    );
};

export default DatePicker;
