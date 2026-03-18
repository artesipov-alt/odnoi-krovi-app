import { Button } from '@mui/material';
import cn from 'classnames';
import recListOnboarding1 from 'imgs/recListOnboarding1.png';
import recListOnboarding2 from 'imgs/recListOnboarding2.png';
import MainLogo from 'imgs/svg/mainLogo';
import { FC, useCallback, useState } from 'react';
import { toast } from 'react-toastify';

import { updateUser } from 'api/apiServices/updateUser';
import { Onboarding } from 'api/user';
import Layout from 'components/Layout';

import styles from './RecipientsListOnboarding.module.less';

type Props = {
    id: string;
    refetch: () => void;
    onBoarding?: Onboarding[];
};

const RecipientsListOnboarding: FC<Props> = ({ id, refetch, onBoarding }) => {
    const [step, setStep] = useState<number>(1);

    const showToast = useCallback((text: string) => {
        toast.warn(text);
    }, []);

    const onConfirmButtonClickHandler = async () => {
        if (step === 1) {
            setStep(2);

            window.scrollTo(0, 0);

            return;
        }

        const { error } = await updateUser({
            id,
            onBoarding: onBoarding ? [...onBoarding, Onboarding.RECIPIENT_LIST] : [Onboarding.RECIPIENT_LIST],
        });

        if (error) {
            showToast('Не удалось сохранить процесс прохождения инструкции, попробуйте ещё раз.');

            return;
        }

        refetch();
    };

    return (
        <Layout>
            <div className={styles.wrapper}>
                <div className={styles.header}>
                    <div className={cn(styles.tab, { [styles.checked]: true })} />
                    <div className={cn(styles.tab, { [styles.checked]: step === 2 })} />
                </div>
                <>
                    {step === 2 && (
                        <div className={styles.logo}>
                            <MainLogo />
                        </div>
                    )}
                    <h1 className={cn(styles.title, { [styles.second]: step === 2 })}>
                        {step === 1 ? (
                            <>Выбирайте питомцев в беде или клиники, сдавайте кровь и спасайте жизни!</>
                        ) : (
                            'Бережный просмотр'
                        )}
                    </h1>
                    {step === 2 && (
                        <div className={styles.description}>
                            Часто помощь нужна питомцам в критическом состоянии, поэтому все фото скрыты по умолчанию
                        </div>
                    )}
                    <img
                        alt='search1'
                        src={step === 1 ? recListOnboarding1 : recListOnboarding2}
                        className={cn(styles.img, { [styles.second]: step === 2 })}
                    />
                    <div className={styles.button}>
                        <Button fullWidth onClick={onConfirmButtonClickHandler} className={styles.confirm}>
                            Далее
                        </Button>
                    </div>
                </>
            </div>
        </Layout>
    );
};

export default RecipientsListOnboarding;
