import MuiTextField from '@mui/material/TextField';
import cn from 'classnames';
import { ChangeEvent, FC, ReactNode } from 'react';

import styles from './TextField.module.less';

type Props = {
    name: string;
    value: string;
    maxLength?: number;
    disabled?: boolean;
    placeholder: string;
    inputClass?: string;
    multiline?: boolean;
    isDigitInput?: boolean;
    htmlInputClass?: string;
    endAdornment?: ReactNode;
    onBlur?: (e: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => void;
    onChange?: (e: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => void;
};

const TextField: FC<Props> = ({
    name,
    value,
    onBlur,
    onChange,
    disabled,
    maxLength,
    multiline,
    inputClass,
    placeholder,
    isDigitInput,
    endAdornment,
    htmlInputClass,
}) => (
    <MuiTextField
        fullWidth
        name={name}
        value={value}
        onBlur={onBlur}
        disabled={disabled}
        onChange={onChange}
        multiline={multiline}
        placeholder={placeholder}
        slotProps={{
            htmlInput: {
                className: cn(styles.input, htmlInputClass, { [styles.noMultiline]: !multiline }),
                maxLength,
                inputMode: isDigitInput ? 'decimal' : 'text',
            },
            input: {
                endAdornment,
                className: cn(styles.inputWrapper, inputClass, { [styles.multiline]: multiline }),
            },
        }}
        sx={{
            '& .MuiOutlinedInput-root': {
                '& .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#dee2e9',
                },
                '&.Mui-disabled .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#dee2e9',
                },
                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#dee2e9',
                    borderWidth: '1px',
                },
            },
            '& .MuiInputBase-input': {
                '&:focus': {
                    backgroundColor: '#ffffff',
                },
            },
        }}
    />
);

export default TextField;
