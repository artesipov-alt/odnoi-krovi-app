import { Button } from '@mui/material';
import onboarding from 'imgs/recipientOnboarding.png';
import { FC } from 'react';

import { updateUser } from 'api/apiServices/updateUser';
import { Onboarding } from 'api/user';
import Layout from 'components/Layout';

import styles from './RecipientOnboarding.module.less';

type Props = {
    id: string;
    onboardings?: Onboarding[];
    onConfirmButtonClick: () => void;
};

const RecipientOnboarding: FC<Props> = ({ id, onboardings, onConfirmButtonClick }) => {
    const onConfirmButtonClickHandler = async () => {
        await updateUser({
            id,
            onBoarding: onboardings ? [...onboardings, Onboarding.FIND_BLOOD] : [Onboarding.FIND_BLOOD],
        });

        onConfirmButtonClick();
    };

    return (
        <Layout>
            <div className={styles.wrapper}>
                <h1 className={styles.title}>
                    Ищите кровь и следите
                    <br />
                    за статусом поиска
                </h1>
                <img className={styles.photos} alt='map' src={onboarding} />
                <Button fullWidth onClick={onConfirmButtonClickHandler} className={styles.confirm}>
                    Далее
                </Button>
            </div>
        </Layout>
    );
};

export default RecipientOnboarding;
