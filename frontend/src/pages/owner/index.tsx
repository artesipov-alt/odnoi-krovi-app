import Bell from 'imgs/svg/bell';
import Paw from 'imgs/svg/paw';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { TelegramUser } from 'types';

import { getPets } from 'api/apiServices/getPets';
import { Pet } from 'api/pets';
import Layout from 'components/Layout';

import styles from './Owner.module.less';

type Props = {
    user: TelegramUser;
};

const Owner: FC<Props> = ({ user }) => {
    const navigate = useNavigate();

    const [pets, setPets] = useState<Pet[]>([]);

    const fetchPets = useCallback(async () => {
        const items = await getPets(user.id);

        if (items?.length) {
            setPets(items);
        }
    }, [user.id]);

    const onAddPetClickHandler = () => {
        navigate('/adding');
    };

    useEffect(() => {
        fetchPets();
    }, [fetchPets]);

    return (
        <Layout>
            {pets.length ? (
                <div>est</div>
            ) : (
                <div className={styles.noPetsWrapper}>
                    <div className={styles.header}>
                        <h2 className={styles.title}>Мои питомцы</h2>
                        <div className={styles.services}>
                            <div className={styles.bellIcon}>
                                <Bell />
                            </div>
                            <div className={styles.avatar}>{user.fullName.charAt(0).toUpperCase()}</div>
                        </div>
                    </div>
                    <div className={styles.button} onClick={onAddPetClickHandler}>
                        <div className={styles.pawIcon}>
                            <Paw />
                        </div>
                        <p className={styles.buttonText}>Добавить питомца</p>
                    </div>
                </div>
            )}
        </Layout>
    );
};

export default Owner;
