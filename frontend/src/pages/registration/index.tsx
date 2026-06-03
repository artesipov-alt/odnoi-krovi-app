import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import cn from 'classnames';
import useBodyScrollLock from 'hooks/useBodyScrollLock';
import { ChangeEvent, FC, useEffect, useRef, useState } from 'react';
import InputMask from 'react-input-mask';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { updatePhone } from 'api/apiServices/updatePhone';
import { updateUser } from 'api/apiServices/updateUser';
import { verifyPhone } from 'api/apiServices/verifyPhone';
import { queryClient } from 'api/queryClient';
import Curtain from 'components/Curtain';
import Layout from 'components/Layout';
import Loading from 'components/Loading';
import SMSInput from 'components/SmsInput';

import styles from './Registration.module.less';

type Props = {
    userId: string;
    fullName: string;
    initialize: () => Promise<void>;
};

type Input = {
    value: string;
    isError?: boolean;
};

const emailRegexp = /^\w+([+.-]?\w+)*@\w+([.-]?\w+)*(\.\w+)+$/i;

const Registration: FC<Props> = ({ userId, fullName, initialize }) => {
    const navigate = useNavigate();

    const [code, setCode] = useState('');
    const [phone, setPhone] = useState<Input>({ value: '' });
    const [email, setEmail] = useState<Input>({ value: '' });
    const [name, setName] = useState<Input>({ value: fullName });
    const [timeLeft, setTimeLeft] = useState<number>(0);
    const [isCodeNotValid, setIsCodeNotValid] = useState<boolean>(false);
    const [timeOfOpenCurtain, setTimeOfOpenCurtain] = useState<number>(0);
    const [isConfirmCurtainOpen, setIsConfirmCurtainOpen] = useState<boolean>(false);

    const [isLoading, setIsLoading] = useState(false);

    const updatedParams = useRef({ name: fullName, phone: '', email: '' });

    useBodyScrollLock(isLoading || isConfirmCurtainOpen);

    const isValidEmail = () => email.value.match(emailRegexp);

    const onChangeInput =
        (inputName?: string) =>
        ({ target: { value } }: ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
            switch (true) {
                case inputName === 'phone': {
                    setPhone({ value });

                    return;
                }
                case inputName === 'email': {
                    setEmail({ value });

                    return;
                }
                default: {
                    setName({ value });
                }
            }
        };

    const onBlurEmailHandler = () => {
        if (!isValidEmail()) {
            setEmail((prevState) => ({ ...prevState, isError: true }));

            toast.warn('Неверный формат e-mail адреса');
        }
    };

    const onCloseConfirmCurtainHandler = () => {
        setIsConfirmCurtainOpen(false);
    };

    const onConfirmClickHandler = async () => {
        let isInvalidValue = false;

        if (!name.value) {
            isInvalidValue = true;

            setName({ isError: true, value: '' });
        }

        if (!phone.value) {
            isInvalidValue = true;

            setPhone({ isError: true, value: '' });
        }

        if (phone.value && /_/.test(phone.value)) {
            isInvalidValue = true;

            setPhone(({ value }) => ({ isError: true, value }));

            toast.warn('Неверный формат номера телефона');
        }

        if (!email.value) {
            isInvalidValue = true;

            setEmail({ isError: true, value: '' });
        }

        if (email.value && !isValidEmail()) {
            isInvalidValue = true;

            toast.warn('Неверный формат e-mail адреса');
        }

        if (isInvalidValue) {
            return;
        }

        setTimeLeft(60);
        setTimeOfOpenCurtain(Date.now());
        setIsConfirmCurtainOpen(true);

        if (updatedParams.current.name !== name.value || updatedParams.current.email !== email.value) {
            const updateUserResponse = await updateUser({
                id: userId,
                email: email.value,
                fullName: name.value,
            });

            updatedParams.current.name = name.value;
            updatedParams.current.email = email.value;

            if (updateUserResponse.error) {
                toast.warn(updateUserResponse.error);
            }
        }

        const updatePhoneResponse = await updatePhone({ id: userId, phone: phone.value });

        if (updatePhoneResponse.error) {
            toast.warn(updatePhoneResponse.error);

            onCloseConfirmCurtainHandler();
        }
    };

    const onFillCodeHandler = (confirmCode: string) => {
        setCode(confirmCode);
        setIsCodeNotValid(false);
    };

    const onConfirmCodeClickHandler = async () => {
        setIsLoading(true);

        const { error } = await verifyPhone({ id: userId, code });

        if (error) {
            setIsLoading(false);
            setIsCodeNotValid(true);

            toast.warn(error);

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['userById', userId] });
        await initialize();

        navigate('/owner');
    };

    const onGetCodeClickHandler = () => {
        setCode('');
        onConfirmClickHandler();
    };

    const getNewCodeBlock = () => {
        const timeNotExpired = timeLeft > 0;

        return (
            <div
                onClick={timeNotExpired ? undefined : onGetCodeClickHandler}
                className={cn(styles.newCode, { [styles.expired]: !timeNotExpired })}
            >
                {timeNotExpired ? `Запросить код заново через ${timeLeft} сек.` : 'Запросить код заново'}
            </div>
        );
    };

    useEffect(() => {
        if (timeOfOpenCurtain === 0) return;

        const endTime = timeOfOpenCurtain + 60000; // +1 минута
        const timer = setInterval(() => {
            const now = Date.now();
            const diff = Math.floor((endTime - now) / 1000);

            if (diff <= 0) {
                setTimeLeft(0);
                clearInterval(timer);
            } else {
                setTimeLeft(diff);
            }
        }, 1000);

        return () => clearInterval(timer);
    }, [timeOfOpenCurtain]);

    return (
        <Layout className={styles.wrapper}>
            <div className={styles.header}>
                <h2 className={styles.title}>Добро пожаловать,</h2>
                <h4 className={styles.subTitle}>{fullName}</h4>
            </div>
            <div className={styles.form}>
                <TextField
                    fullWidth
                    variant='standard'
                    value={name.value}
                    onChange={onChangeInput()}
                    placeholder='Как вас зовут?'
                    slotProps={{
                        htmlInput: { className: cn(styles.input, { [styles.inputError]: name.isError } as any) },
                        input: {
                            className: cn(styles.inputWrapper, { [styles.inputWrapperError]: name.isError } as any),
                        },
                    }}
                />
                <InputMask mask='+79999999999' value={phone.value} onChange={onChangeInput('phone')}>
                    {(inputProps) => (
                        <TextField
                            {...inputProps}
                            fullWidth
                            type='tel'
                            placeholder='+7'
                            variant='standard'
                            className={styles.textField}
                            slotProps={{
                                htmlInput: {
                                    className: cn(styles.input, { [styles.inputError]: phone.isError } as any),
                                },
                                input: {
                                    className: cn(styles.inputWrapper, {
                                        [styles.inputWrapperError]: phone.isError,
                                    } as any),
                                },
                            }}
                        />
                    )}
                </InputMask>
                <TextField
                    fullWidth
                    variant='standard'
                    value={email.value}
                    placeholder='E-mail'
                    onBlur={onBlurEmailHandler}
                    onChange={onChangeInput('email')}
                    slotProps={{
                        htmlInput: { className: cn(styles.input, { [styles.inputError]: email.isError } as any) },
                        input: {
                            className: cn(styles.inputWrapper, {
                                [styles.inputWrapperError]: email.isError,
                            } as any),
                        },
                    }}
                />
                <div className={styles.links}>
                    <p className={styles.agreementLabel}>
                        Продолжая, Вы принимаете{' '}
                        <a
                            target='_blank'
                            rel='noreferrer'
                            className={styles.agreementLabelLink}
                            href='https://однойкрови.рф/docs#n-a9dea2ae-b2a2-4bc6-b0ea-1eca71588ab0'
                        >
                            Пользовательское соглашение
                        </a>{' '}
                        и{' '}
                        <a
                            target='_blank'
                            rel='noreferrer'
                            className={styles.agreementLabelLink}
                            href='https://однойкрови.рф/docs#n-80ae6549-954a-4c8c-bc54-357a93ec3dee'
                        >
                            Политику конфиденциальности
                        </a>
                    </p>
                </div>
                <Button fullWidth variant='contained' className={styles.button} onClick={onConfirmClickHandler}>
                    Продолжить
                </Button>
                {isLoading && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
                {isConfirmCurtainOpen && (
                    <Curtain
                        noRednerButtons
                        title='Код подтверждения'
                        shouldCloseByWrapperClick
                        onClose={onCloseConfirmCurtainHandler}
                        subTitleClassName={styles.confirmSubtitle}
                        subTitle={`Робот позвонит на номер ${phone.value} и назовёт код`}
                    >
                        <SMSInput
                            onFill={onFillCodeHandler}
                            className={styles.codeInput}
                            isCodeNotValid={isCodeNotValid}
                        />
                        <div className={styles.buttons}>
                            <Button
                                fullWidth
                                onClick={onConfirmCodeClickHandler}
                                className={cn(styles.confirm, { [styles.enabled]: !!code && !isCodeNotValid })}
                            >
                                {isCodeNotValid ? 'Неверный код' : 'Подтвердить'}
                            </Button>
                            {getNewCodeBlock()}
                        </div>
                    </Curtain>
                )}
            </div>
        </Layout>
    );
};

export default Registration;
