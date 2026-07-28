import { Autocomplete, AutocompleteInputChangeReason, TextField as MuiTextField } from '@mui/material';
import cn from 'classnames';
import ArrowDown from 'imgs/svg/arrowDown';
import { FC, SyntheticEvent } from 'react';

import { Dict } from 'api/reference';

import styles from './MultiAutocomplete.module.less';

type Props = {
    dict: Dict[];
    value: Dict[];
    maxCount?: number;
    onChange: (value: Dict[]) => void;
};

const MultiAutocomplete: FC<Props> = ({ dict, value, onChange, maxCount = 3 }) => {
    const onChangeLocationsHandler = (
        _: SyntheticEvent<Element, Event>,
        newValue: string,
        reason: AutocompleteInputChangeReason,
    ) => {
        if (reason !== 'input') {
            return;
        }

        if (newValue === '') {
            onChange([]);
        }
    };

    const onAutocompleteChangeHandler = (_, newValue: Dict[] | null) => {
        if (!!newValue?.length && newValue.length > maxCount) {
            return;
        }

        onChange(newValue || []);
    };

    return (
        <Autocomplete
            multiple
            value={value}
            options={dict}
            filterSelectedOptions
            id='MultiAutocomplete'
            disableCloseOnSelect
            onChange={onAutocompleteChangeHandler}
            noOptionsText='Нет подходящих вариантов'
            onInputChange={onChangeLocationsHandler}
            getOptionLabel={(option) => option.label}
            renderInput={(params) => (
                <MuiTextField
                    {...params}
                    placeholder='Выберите из списка'
                    className={cn({ [styles.input]: value.length === 3 })}
                />
            )}
            sx={{
                '& .MuiOutlinedInput-root': {
                    padding: `${value.length ? '6px' : '0'} 9px`,
                    minHeight: '50px',
                    borderRadius: '16px',
                    maxHeight: 48 * 4.5 + 8,
                    backgroundColor: 'white',
                    color: value && '#8B7069',
                    '& .MuiOutlinedInput-notchedOutline': {
                        borderColor: '#dee2e9',
                    },
                    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                        borderColor: '#dee2e9',
                        borderWidth: '1px',
                    },
                },
                '& .MuiAutocomplete-popupIndicator': {
                    marginRight: '4px',
                    transition: 'transform 0.3s ease',
                },
            }}
            popupIcon={
                <div className={styles.selectArrow}>
                    <ArrowDown />
                </div>
            }
            slotProps={{
                popupIndicator: {
                    disableRipple: true,
                },
            }}
        />
    );
};

export default MultiAutocomplete;
