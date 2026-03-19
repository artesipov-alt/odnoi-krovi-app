import { createSvgIcon } from '@mui/material';

export default createSvgIcon(
    <svg viewBox='0 0 18 18' xmlns='http://www.w3.org/2000/svg'>
        <rect width='18' height='18' rx='9' fill='white' fillOpacity='0.5' />
        <rect x='6' y='5' width='2' height='8' fill='#979797' />
        <rect x='10' y='5' width='2' height='8' fill='#979797' />
    </svg>,
    'pause',
);
