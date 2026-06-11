import cn from 'classnames';
import { ChangeEvent, FC, KeyboardEvent, useEffect, useRef, useState } from 'react';

import styles from './SMSInput.module.less';

type Props = {
    className: string;
    isCodeNotValid: boolean;
    onFill: (code: string) => void;
};

const inputLength = 4;

const SMSInput: FC<Props> = ({ className, onFill, isCodeNotValid }) => {
    const [code, setCode] = useState<string[]>(Array(inputLength).fill(''));

    const inputsRef = useRef<(HTMLInputElement | null)[]>([]);

    const onChangeHandler = (i: number) => (e: ChangeEvent<HTMLInputElement>) => {
        const { value } = e.target;

        if (/^\d?$/.test(value)) {
            const newCode = [...code];

            newCode[i] = value;

            setCode(newCode);

            // Перейти к следующему полю
            if (value && i < inputLength - 1) {
                inputsRef.current[i + 1]?.focus();
            }

            if (value === '' && code.every((digit) => digit !== '')) {
                onFill('');

                return;
            }

            if (value && newCode.every((digit) => digit !== '')) {
                onFill(newCode.join(''));
            }
        }
    };

    const onKeyDownHandler = (i: number) => (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Backspace' && !code[i] && i > 0) {
            inputsRef.current[i - 1]?.focus();
        }
    };

    const onClickHandler = (i: number) => () => {
        inputsRef.current[i]?.focus();
    };

    useEffect(() => {
        inputsRef.current[0]?.focus();
    }, []);

    return (
        <div className={cn(styles.container, className)}>
            {code.map((digit, i) => (
                <input
                    /* eslint-disable-next-line react/no-array-index-key */
                    key={i}
                    type='text'
                    maxLength={1}
                    value={digit}
                    inputMode='numeric'
                    name={`digit-${i}`}
                    onClick={onClickHandler(i)}
                    onChange={onChangeHandler(i)}
                    onKeyDown={onKeyDownHandler(i)}
                    className={cn(styles.input, { [styles.fill]: digit, [styles.error]: isCodeNotValid })}
                    ref={(el) => {
                        inputsRef.current[i] = el;
                    }}
                />
            ))}
        </div>
    );
};

export default SMSInput;
