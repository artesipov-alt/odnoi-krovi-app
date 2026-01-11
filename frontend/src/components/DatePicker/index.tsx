import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { DatePicker as MuiDatePicker } from '@mui/x-date-pickers/DatePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { ru } from 'date-fns/locale/ru';
import { FC, useRef, useState } from 'react';

type Props = {
    value: Date | null;
    onChange: (date: Date | null) => void;
};

const DatePicker: FC<Props> = ({ value, onChange }) => {
    const [isOpen, setIsOpen] = useState<boolean>(false);

    const onChangeHandler = (newValue: Date | null) => {
        onChange(newValue);

        setIsOpen(false);
    };

    return (
        <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ru}>
            <MuiDatePicker
                value={value}
                open={isOpen}
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
                            backgroundColor: 'white',
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
    );
};

export default DatePicker;
