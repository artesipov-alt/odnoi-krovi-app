import CircularProgress from '@mui/material/CircularProgress';
import { CSSProperties, FC } from 'react';

import styles from './Loading.module.less';

type Props = {
    size?: number;
    thickness?: number;
    className?: string;
};

const defaultSize = 98;
const defaultThickness = 7;

const Loading: FC<Props> = ({ size = defaultSize, thickness = defaultThickness, className = styles.loader }) => (
    <div className={styles.wrapper}>
        <CircularProgress
            size={size}
            color='inherit'
            thickness={thickness}
            className={className}
            sx={{
                '& .MuiCircularProgress-circle': {
                    strokeLinecap: 'round',
                },
            }}
            style={{ '--loading-size': `${size}px` } as CSSProperties}
        />
    </div>
);
export default Loading;
