import { FC, useEffect, useState } from 'react';

type Props = {
    updatedAt: string;
    hoursToAdd: number;
    className?: string;
    onTimeEnd: () => void;
    digitClassName?: string;
    separatorClassName?: string;
};

const calculateTimeLeft = (updatedAt: string, hoursToAdd: number) => {
    const updatedTime = new Date(updatedAt).getTime();
    const endTime = updatedTime + hoursToAdd * 60 * 60 * 1000;
    const now = new Date().getTime();
    const difference = endTime - now;

    if (difference <= 0) {
        return 0;
    }

    return Math.floor(difference / 1000);
};

const Timer: FC<Props> = ({ updatedAt, hoursToAdd, onTimeEnd, className, digitClassName, separatorClassName }) => {
    const [timeLeft, setTimeLeft] = useState(() => calculateTimeLeft(updatedAt, hoursToAdd));
    const [isEnded, setIsEnded] = useState(false);

    useEffect(() => {
        const timer = setInterval(() => {
            const secondsLeft = calculateTimeLeft(updatedAt, hoursToAdd);
            setTimeLeft(secondsLeft);

            if (secondsLeft <= 0 && !isEnded) {
                setIsEnded(true);

                onTimeEnd();
            }
        }, 1000);

        return () => clearInterval(timer);
    }, [updatedAt, hoursToAdd, onTimeEnd, isEnded]);

    const formatTime = (seconds: number) => {
        const hrs = Math.floor(seconds / 3600);
        const mins = Math.floor((seconds % 3600) / 60);
        const secs = seconds % 60;

        if (hrs === 0) {
            // Показываем только минуты и секунды
            return (
                <>
                    <span className={digitClassName}>{mins.toString().padStart(2, '0')}</span>
                    <span className={separatorClassName}>:</span>
                    <span className={digitClassName}>{secs.toString().padStart(2, '0')}</span>
                </>
            );
        }

        // Показываем часы, минуты и секунды
        return (
            <>
                <span className={digitClassName}>{hrs.toString().padStart(2, '0')}</span>
                <span className={separatorClassName}>:</span>
                <span className={digitClassName}>{mins.toString().padStart(2, '0')}</span>
                <span className={separatorClassName}>:</span>
                <span className={digitClassName}>{secs.toString().padStart(2, '0')}</span>
            </>
        );
    };

    return <div className={className}>{timeLeft > 0 ? formatTime(timeLeft) : 'Время истекло'}</div>;
};

export default Timer;
