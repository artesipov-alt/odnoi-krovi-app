import cn from 'classnames';
import BloodSearch from 'imgs/svg/bloodSearch';
import DonorButton from 'imgs/svg/donorButton';
import Paw from 'imgs/svg/paw';
import RecipientButton from 'imgs/svg/recipientButton';
import { FC, MouseEvent, useCallback, useEffect, useLayoutEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { TelegramUser } from 'types';

import { getPets } from 'api/apiServices/getPets';
import { Pet } from 'api/pets';
import { Role } from 'api/user';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import styles from './Owner.module.less';
import PetProfile from './Profiles/Pet';

type Props = {
    user: TelegramUser;
};

type View = 'donor' | 'recipient';

const Owner: FC<Props> = ({ user }) => {
    const navigate = useNavigate();

    const [pets, setPets] = useState<Pet[]>([]);
    const [view, setView] = useState<View>('recipient');
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [selectedPet, setSelectedPet] = useState<Pet | null>(null);
    const [isPetProfileOpen, setIsPetProfileOpen] = useState<boolean>(false);

    const fetchPets = useCallback(async () => {
        const items = await getPets(user.id);

        if (!items) {
            setPets([]);
            setSelectedPet(null);
            setIsLoading(false);

            return;
        }

        if (items.length) {
            setPets(items);
        }

        setSelectedPet((prevState) => {
            if (prevState) {
                return items?.find((pet) => pet.id === prevState.id) || null;
            }

            return prevState;
        });

        setIsLoading(false);
    }, [user]);

    const onButtonClickHandler = (newView: View) => () => {
        window.location.hash = `#${newView}`;

        setView(newView);
    };

    const onPetProfileToggleHandler = (petData: Pet | null) => () => {
        setSelectedPet(petData);

        setIsPetProfileOpen((prevState) => !prevState);
    };

    const onAddPetClickHandler = () => {
        navigate('/adding');
    };

    const onLabelClickHandler =
        (label: 'startSearch' | 'activeSearch', petId: string) => (e: MouseEvent<HTMLDivElement>) => {
            e.stopPropagation();

            if (label === 'startSearch') {
                navigate(`/adding#startSearch_${petId}`);
            }
        };

    const renderLabel = (petsStatus: Role, petId: string) => {
        switch (true) {
            case petsStatus === Role.RECIPIENT: {
                return (
                    <div
                        onClick={onLabelClickHandler('activeSearch', petId)}
                        className={cn(styles.label, { [styles.activeSearch]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <BloodSearch />
                        </div>
                        <div className={styles.labelText}>
                            Идет поиск <p className={styles.labelArrow}>⟶</p>
                        </div>
                    </div>
                );
            }
            case petsStatus === Role.NONE || petsStatus === Role.DONOR: {
                return (
                    <div
                        onClick={onLabelClickHandler('startSearch', petId)}
                        className={cn(styles.label, { [styles.startSearch]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <RecipientButton />
                        </div>
                        Начать поиск
                    </div>
                );
            }
            default: {
                return null;
            }
        }
    };

    useEffect(() => {
        setIsLoading(true);

        fetchPets();
    }, [fetchPets]);

    useLayoutEffect(() => {
        if (window.location.hash === '#recipient') {
            setView('recipient');

            return;
        }

        if (window.location.hash === '#donor') {
            setView('donor');

            return;
        }

        window.location.hash = '#recipient';
    }, []);

    if (isPetProfileOpen && selectedPet) {
        return (
            <Layout>
                <PetProfile updatePets={fetchPets} onClose={onPetProfileToggleHandler(null)} {...selectedPet} />
            </Layout>
        );
    }

    return (
        <Layout>
            <div className={cn(styles.wrapper, { [styles.isPets]: !!pets.length })}>
                <div className={styles.header}>
                    <div className={styles.avatar}>{user.fullName.charAt(0).toUpperCase()}</div>
                </div>
                {isLoading && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
                {!isLoading && !pets.length && (
                    <div className={styles.button} onClick={onAddPetClickHandler}>
                        <div className={styles.pawIcon}>
                            <Paw />
                        </div>
                        <p className={styles.pawButtonText}>Добавить питомца</p>
                    </div>
                )}
                {!isLoading && !!pets.length && (
                    <>
                        <div className={styles.showcase}>
                            {pets.map((pet) => (
                                <div key={`${pet.id}`} className={cn(styles.pet, { [styles[pet.type]]: true })}>
                                    <div className={styles.photo} onClick={onPetProfileToggleHandler(pet)}>
                                        {!!pet.photoUrls?.[0] && (
                                            <img className={styles.img} src={pet.photoUrls?.[0]} alt={pet.name} />
                                        )}
                                        <div className={styles.bloodGroup}>{pet.bloodGroup || '?'}</div>
                                        <div className={styles.photoFooter}>
                                            <p className={styles.name}>{pet.name.toUpperCase()}</p>
                                            {view === 'recipient' && renderLabel(pet.petStatus, pet.id)}
                                        </div>
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
