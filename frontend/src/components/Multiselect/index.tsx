import { Box, Chip, MenuItem, OutlinedInput, Select, SelectChangeEvent } from '@mui/material';
import cn from 'classnames';
import ArrowDown from 'imgs/svg/arrowDown';
import { FC, ReactElement, useState } from 'react';

import styles from './Multiselect.module.less';

type Dict = {
    value: string;
    label: string;
}[];

type Props = {
    dict: Dict;
    selectValue: string[];
    onChange: (event: SelectChangeEvent<string[]>) => void;
};

const ITEM_HEIGHT = 48;
const ITEM_PADDING_TOP = 8;
const MenuProps = {
    PaperProps: {
        style: {
            maxHeight: ITEM_HEIGHT * 4.5 + ITEM_PADDING_TOP,
            width: 250,
        },
    },
};

const Multiselect: FC<Props> = ({ selectValue, onChange, dict }) => {
    const [isOpen, setIsOpen] = useState(false);

    const onIsOpenClickToggle = () => {
        setIsOpen((prevState) => !prevState);
    };

    return (
        <Select
            multiple
            fullWidth
            displayEmpty
            open={isOpen}
            value={selectValue}
            onOpen={onIsOpenClickToggle}
            onClose={onIsOpenClickToggle}
            classes={{
                root: styles.rootSelect,
                select: styles.select,
                icon: styles.selectIcon,
            }}
            endAdornment={
                <div className={cn(styles.selectArrow, { [styles.isOpen]: isOpen })}>
                    <ArrowDown />
                </div>
            }
            onChange={onChange}
            sx={{
                '& .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#dee2e9',
                },
                '&:hover .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#dee2e9',
                },
                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                    borderColor: '#dee2e9',
                    borderWidth: '1px',
                },
            }}
            input={
                (
                    <OutlinedInput
                        sx={{
                            '&.MuiNotchedOutlined-root-MuiOutlinedInput-notchedOutline': {
                                borderColor: 'red',
                            },
                        }}
                    />
                ) as ReactElement
            }
            renderValue={(selected) => {
                if (!selected.length) {
                    return <span className={styles.placeholder}>Выберите из списка</span>;
                }

                return (
                    <Box sx={{ borderColor: 'primary.main', display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
                        {selected.map((value) => (
                            <Chip
                                key={value}
                                classes={{
                                    root: styles.chip,
                                }}
                                label={dict.filter((el) => el.value === value)[0].label}
                            />
                        ))}
                    </Box>
                );
            }}
            MenuProps={MenuProps}
        >
            {dict.map(({ label, value }) => (
                <MenuItem key={label} value={value}>
                    {label}
                </MenuItem>
            ))}
        </Select>
    );
};

export default Multiselect;
