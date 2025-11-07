import { styled } from '@mui/material/styles';
import UISwitch, { SwitchProps } from '@mui/material/Switch';

import styles from './Switch.module.less';

const classes = {
    root: styles.root,
};

const Switch = styled((props: SwitchProps) => (
    <UISwitch focusVisibleClassName='.Mui-focusVisible' disableRipple classes={classes} {...props} />
))(({ theme }) => ({
    width: 58,
    height: 42,
    padding: 0,
    '& .MuiSwitch-switchBase': {
        padding: 0,
        margin: 5,
        transitionDuration: '300ms',
        '&.Mui-checked': {
            transform: 'translateX(16px)',
            color: '#fff',
            '& + .MuiSwitch-track': {
                backgroundColor: '#B27B89',
                opacity: 1,
                border: 0,
                ...theme.applyStyles('dark', {
                    backgroundColor: '#B27B89',
                }),
            },
            '&.Mui-disabled + .MuiSwitch-track': {
                opacity: 0.5,
            },
        },
        '&.Mui-focusVisible .MuiSwitch-thumb': {
            color: '#B27B89',
            border: '6px solid #fff',
        },
        '&.Mui-disabled .MuiSwitch-thumb': {
            color: theme.palette.grey[100],
            ...theme.applyStyles('dark', {
                color: theme.palette.grey[600],
            }),
        },
        '&.Mui-disabled + .MuiSwitch-track': {
            opacity: 0.7,
            ...theme.applyStyles('dark', {
                opacity: 0.3,
            }),
        },
    },
    '& .MuiSwitch-thumb': {
        boxSizing: 'border-box',
        width: 32,
        height: 32,
    },
    '& .MuiSwitch-track': {
        borderRadius: 86 / 2,
        backgroundColor: 'rgba(188, 195, 208, 0.50)',
        opacity: 1,
        transition: theme.transitions.create(['background-color'], {
            duration: 500,
        }),
        ...theme.applyStyles('dark', {
            backgroundColor: '#39393D',
        }),
    },
}));

export default Switch;
