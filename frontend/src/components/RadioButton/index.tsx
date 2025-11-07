import Brightness1Icon from '@mui/icons-material/Brightness1';
import Brightness1OutlinedIcon from '@mui/icons-material/Brightness1Outlined';
import { Icon, Radio, RadioProps } from '@mui/material';
import { FC } from 'react';

import styles from './radioButton.module.less';

const classes = {
    root: styles.root,
    checked: styles.checked,
    disabled: styles.disabled,
};
const RadioButton: FC<RadioProps> = ({ ...props }) => (
    <Radio
        classes={classes}
        icon={<Icon sx={{ fontSize: 31 }} component={Brightness1OutlinedIcon} />}
        checkedIcon={<Icon sx={{ fontSize: 31 }} component={Brightness1Icon} />}
        {...props}
    />
);

export default RadioButton;
