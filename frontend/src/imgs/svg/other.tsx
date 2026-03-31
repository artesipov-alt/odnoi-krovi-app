import { createSvgIcon } from '@mui/material';

export default createSvgIcon(
    <svg width='26' height='26' viewBox='0 0 26 26' fill='none' xmlns='http://www.w3.org/2000/svg'>
        <rect width='26' height='26' rx='13' fill='url(#paint0_linear_4968_6730)' />
        <rect
            x='0.777799'
            y='0.777799'
            width='24.4444'
            height='24.4444'
            rx='12.2222'
            stroke='white'
            strokeOpacity='0.3'
            strokeWidth='1.5556'
        />
        <g filter='url(#filter0_d_4968_6730)'>
            <circle cx='7.56916' cy='12.9998' r='1.71564' fill='white' />
            <circle cx='13.0008' cy='12.9998' r='1.71564' fill='white' />
            <circle cx='18.4324' cy='12.9998' r='1.71564' fill='white' />
        </g>
        <defs>
            <filter
                id='filter0_d_4968_6730'
                x='5.33498'
                y='11.2842'
                width='15.332'
                height='4.46822'
                filterUnits='userSpaceOnUse'
                colorInterpolationFilters='sRGB'
            >
                <feFlood floodOpacity='0' result='BackgroundImageFix' />
                <feColorMatrix
                    in='SourceAlpha'
                    type='matrix'
                    values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 127 0'
                    result='hardAlpha'
                />
                <feOffset dy='0.518533' />
                <feGaussianBlur stdDeviation='0.259266' />
                <feComposite in2='hardAlpha' operator='out' />
                <feColorMatrix type='matrix' values='0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0.15 0' />
                <feBlend mode='normal' in2='BackgroundImageFix' result='effect1_dropShadow_4968_6730' />
                <feBlend mode='normal' in='SourceGraphic' in2='effect1_dropShadow_4968_6730' result='shape' />
            </filter>
            <linearGradient id='paint0_linear_4968_6730' x1='13' y1='0' x2='13' y2='26' gradientUnits='userSpaceOnUse'>
                <stop stopColor='#E7CDE8' />
                <stop offset='1' stopColor='#C790CA' />
            </linearGradient>
        </defs>
    </svg>,
    'other',
);
