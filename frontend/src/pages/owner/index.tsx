import cn from 'classnames';
import DonorButton from 'imgs/svg/donorButton';
import Paw from 'imgs/svg/paw';
import RecipientButton from 'imgs/svg/recipientButton';
import { FC, useCallback, useEffect, useLayoutEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { TelegramUser } from 'types';

import { getPets } from 'api/apiServices/getPets';
import { Pet } from 'api/pets';
import Layout from 'components/Layout';

import styles from './Owner.module.less';

type Props = {
    user: TelegramUser;
};

type View = 'donor' | 'recipient';

const Owner: FC<Props> = ({ user }) => {
    const navigate = useNavigate();

    const [pets, setPets] = useState<Pet[]>([]);
    const [view, setView] = useState<View>('recipient');

    const fetchPets = useCallback(async () => {
        const items = await getPets(user.id);

        if (items?.length) {
            setPets(items);
        }
    }, [user.id]);

    const onButtonClickHandler = (newView: View) => () => {
        setView(newView);
    };

    const onAddPetClickHandler = () => {
        navigate('/adding');
    };

    useEffect(() => {
        fetchPets();
    }, [fetchPets]);

    useLayoutEffect(() => {
        if (window.location.hash === '#recipient') {
            setView('recipient');
        }

        if (window.location.hash === '#donor') {
            setView('donor');
        }
    }, []);

    return (
        <Layout>
            <div className={styles.wrapper}>
                <div className={styles.header}>
                    <div className={styles.avatar}>{user.fullName.charAt(0).toUpperCase()}</div>
                </div>
                {!pets.length ? (
                    <div className={styles.button} onClick={onAddPetClickHandler}>
                        <div className={styles.pawIcon}>
                            <Paw />
                        </div>
                        <p className={styles.buttonText}>Добавить питомца</p>
                    </div>
                ) : (
                    <>
                        <div className={styles.showcase}>
                            {pets.map((pet) => (
                                <div key={pet.id} className={cn(styles.pet, { [styles[pet.type]]: true })}>
                                    <div className={styles.photo}>
                                        {!!pet.photoUrls?.[0] && (
                                            <img className={styles.img} src={pet.photoUrls?.[0]} alt={pet.name} />
                                        )}
                                        <p className={styles.name}>{pet.name.toUpperCase()}</p>
                                        <div className={styles.gradient} />
                                    </div>
                                </div>
                            ))}
                        </div>
                        <div className={styles.footer}>
                            <div
                                onClick={onButtonClickHandler('donor')}
                                className={cn(styles.viewButton, { [styles.checked]: view === 'donor' })}
                            >
                                <div className={cn(styles.buttonIcon, { [styles.checked]: view === 'donor' })}>
                                    <DonorButton />
                                </div>
                                <p className={styles.buttonText}>Стать донором</p>
                            </div>
                            <div onClick={onAddPetClickHandler} className={styles.pawButton}>
                                <Paw />
                            </div>
                            <div
                                onClick={onButtonClickHandler('recipient')}
                                className={cn(styles.viewButton, { [styles.checked]: view === 'recipient' })}
                            >
                                <div className={cn(styles.buttonIcon, { [styles.checked]: view === 'recipient' })}>
                                    <RecipientButton />
                                </div>
                                <p className={styles.buttonText}>Найти кровь</p>
                            </div>
                        </div>
                    </>
                )}
            </div>
        </Layout>
    );
};

export default Owner;
