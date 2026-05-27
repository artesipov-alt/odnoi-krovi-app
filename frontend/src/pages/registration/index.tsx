import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import cn from 'classnames';
import useBodyScrollLock from 'hooks/useBodyScrollLock';
import { ChangeEvent, FC, useState } from 'react';
import InputMask from 'react-input-mask';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { updateUser } from 'api/apiServices/updateUser';
import { queryClient } from 'api/queryClient';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

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

    const [phone, setPhone] = useState<Input>({ value: '' });
    const [email, setEmail] = useState<Input>({ value: '' });
    const [name, setName] = useState<Input>({ value: fullName });

    const [isLoading, setIsLoading] = useState(false);

    useBodyScrollLock(isLoading);

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

        setIsLoading(true);

        const { data, error } = await updateUser({
            id: userId,
            phone: phone.value,
            email: email.value,
            fullName: name.value,
        });

        if (!data || error) {
            toast.warn(error); // TODO ?

            setIsLoading(false);

            return;
        }

        await queryClient.invalidateQueries({ queryKey: ['userById', userId] });
        await initialize();

        navigate('/owner');
    };

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
                            className={styles.agreementLabelLink}
                            href='https://однойкрови.рф/docs#n-a9dea2ae-b2a2-4bc6-b0ea-1eca71588ab0'
                            target='_blank'
                            rel='noreferrer'
                        >
                            Пользовательское соглашение
                        </a>{' '}
                        и{' '}
                        <a
                            className={styles.agreementLabelLink}
                            href='https://однойкрови.рф/docs#n-80ae6549-954a-4c8c-bc54-357a93ec3dee'
                            target='_blank'
                            rel='noreferrer'
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
            </div>
        </Layout>
    );
};

export default Registration;
