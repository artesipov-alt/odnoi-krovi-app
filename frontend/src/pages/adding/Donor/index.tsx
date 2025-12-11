import { Button } from '@mui/material';
import { FC, useState } from 'react';

import Onboarding from './Onboarding';
import styles from './Onboarding/Onboarding.module.less';

type Props = {
    onBackToStart: () => void;
};

const Donor: FC<Props> = ({ onBackToStart }) => {
    const [step, setStep] = useState(1);
    const [isOnboardingFinish, setIsOnboardingFinish] = useState(false);

    const onFinishOnboardingHandler = () => {
        setIsOnboardingFinish(true);
    };

    return (
        <>
            {!isOnboardingFinish && <Onboarding onFinish={onFinishOnboardingHandler} onBackToStart={onBackToStart} />}
            {isOnboardingFinish && (
                <Button fullWidth onClick={onBackToStart} className={styles.confirm}>
                    Назад
                </Button>
            )}
        </>
    );
};

export default Donor;
