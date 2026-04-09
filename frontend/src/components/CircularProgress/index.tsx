import cn from 'classnames';
import { FC, useMemo } from 'react';

import styles from './CircularProgress.module.less';

type Props = {
    size: number;
    color?: string;
    current: number; // текущее значение
    total: number; // максимальное значение (100%)
    showDot?: boolean;
    className?: string;
    strokeWidth?: number;
};

export const CircularProgress: FC<Props> = ({
    size,
    strokeWidth = 5,
    current,
    total,
    color = '#ff4d4f',
    showDot = false,
    className,
}) => {
    // Рассчитываем прогресс в процентах: от 0 до 100
    const progress = total > 0 ? (current / total) * 100 : 0;

    const radius = (size - strokeWidth) / 2;
    const center = size / 2;
    const circumference = 2 * Math.PI * radius;
    const strokeDashoffset = useMemo(() => circumference * (1 - progress / 100), [circumference, progress]);

    const dotPosition = useMemo(() => {
        const angle = (progress / 100) * 2 * Math.PI - Math.PI / 2; // начинаем с верхней точки

        return {
            x: center + radius * Math.cos(angle),
            y: center + radius * Math.sin(angle),
        };
    }, [center, radius, progress]);

    const getDotRadius = () => {
        switch (`${current}`.length) {
            case 3: {
                return 16;
            }
            case 4: {
                return 18;
            }
            default: {
                return 14;
            }
        }
    };

    return (
        <svg
            width={size}
            height={size}
            viewBox={`0 0 ${size} ${size}`}
            className={cn(styles.progressRing, className)}
            style={{ '--progress-color': color } as React.CSSProperties}
        >
            {/* Фоновое кольцо */}
            <circle
                cx={center}
                cy={center}
                r={radius}
                fill='none'
                stroke='#e6e6e6'
                strokeWidth={strokeWidth}
                strokeLinecap='round'
            />
            {/* Прогресс-кольцо */}
            <circle
                cx={center}
                cy={center}
                r={radius}
                fill='none'
                stroke='currentColor'
                strokeWidth={strokeWidth}
                strokeDasharray={circumference}
                strokeDashoffset={strokeDashoffset}
                transform={`rotate(-90 ${center} ${center})`}
                strokeLinecap='round'
                className={styles.progressCircle}
            />
            {/* Круглая точка 28x28 с текущим значением внутри */}
            {showDot && (
                <g transform={`translate(${dotPosition.x}, ${dotPosition.y})`}>
                    {/* Фон точки */}
                    <circle cx={0} cy={0} fill='currentColor' className={styles.progressDot} r={getDotRadius()} />
                    {/* Текст: текущее число (без %) */}
                    <text
                        x={0.5}
                        y={6.5}
                        textAnchor='middle'
                        fontSize='10'
                        fill='#ffffff'
                        fontWeight='bold'
                        className={styles.progressDotText}
                    >
                        {Math.round(current)}
                    </text>
                </g>
            )}
        </svg>
    );
};
