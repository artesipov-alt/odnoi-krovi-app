import { Button } from '@mui/material';
import cn from 'classnames';
import { FC, useCallback, useState } from 'react';
import { toast } from 'react-toastify';

import { updatePoolRequest } from 'api/apiServices/updatePoolRequest';
import { Onboardings } from 'api/bloodRequest';
import Layout from 'components/Layout';

import styles from './SearchOnboarding.module.less';

export enum View {
    CARD = 'card',
    SEARCH = 'search',
}

type Props = {
    view: View;
    id?: string;
    onSucess: () => void;
    onBoarding?: Onboardings[];
};

const SearchOnboarding: FC<Props> = ({ view, id, onBoarding, onSucess }) => {
    const [step, setStep] = useState<number>(1);

    const showToast = useCallback((text: string) => {
        toast.warn(text);
    }, []);

    const onConfirmButtonClickHandler = async () => {
        if (step === 1) {
            setStep(2);

            return;
        }

        if (!id) {
            return;
        }

        const { success } = await updatePoolRequest({
            id,
            onBoarding: [...(onBoarding || []), view === View.SEARCH ? Onboardings.SEARCH : Onboardings.BLOOD_CARD],
        });

        if (!success) {
            showToast('Не удалось сохранить процесс прохождения инструкции, попробуйте ещё раз.');

            return;
        }

        onSucess();
    };

    return (
        <Layout className={cn(styles.wrapper, { [styles[view]]: view, [styles.second]: step === 2 })}>
            <div>
                <div className={styles.header}>
                    <div className={cn(styles.tab, { [styles.checked]: true })} />
                    <div className={cn(styles.tab, { [styles.checked]: step === 2 })} />
                </div>
                {view === View.SEARCH && (
                    <>
                        <h1 className={cn(styles.title, { [styles.second]: step === 2 })}>
                            {step === 1 ? (
                                <>
                                    Находите доноров или
                                    <br />
                                    бронируйте пакеты крови
                                </>
                            ) : (
                                'Детали поиска'
                            )}
                        </h1>
                        {step === 2 && (
                            <div className={styles.description}>
                                Получите найденные предложения или отмените их для
                                <br />
                                продолжения поиска.
                            </div>
                        )}
                    </>
                )}
                {view === View.CARD && (
                    <>
                        <h1 className={cn(styles.title, { [styles.card]: true })}>
                            {step === 1 ? (
                                <>
                                    Ищите кровь
                                    <br />
                                    пока не заполнится лимит
                                </>
                            ) : (
                                'Получите найденную кровь!'
                            )}
                        </h1>
                        <div className={cn(styles.description, { [styles.card]: true })}>
                            {step === 1
                                ? 'Отмените неподходящие донации, чтобы освободить шкалу и найти больше предложений'
                                : 'Чаты с донорами будут в "Найденном". Если не договоритесь - отмените предложение и продолжайте поиск.'}
                        </div>
                    </>
                )}
            </div>
            <div className={styles.button}>
                <Button fullWidth onClick={onConfirmButtonClickHandler} className={styles.confirm}>
                    Далее
                </Button>
            </div>
        </Layout>
    );
};

export default SearchOnboarding;
