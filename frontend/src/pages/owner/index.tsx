import cn from 'classnames';
import { useGetUserById } from 'hooks/useGetUserById';
import { usePetsQuery } from 'hooks/usePetsQuery';
import BloodSearch from 'imgs/svg/bloodSearch';
import BloodFound from 'imgs/svg/bloodFound';
import DonorButton from 'imgs/svg/donorButton';
import Paw from 'imgs/svg/paw';
import RecipientButton from 'imgs/svg/recipientButton';
import { FC, MouseEvent, useEffect, useLayoutEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { TelegramUser } from 'types';

import { Pet } from 'api/pets';
import { Onboarding, Role } from 'api/user';
import Layout from 'components/Layout';
import Loading from 'components/Loading';

import RecipientOnboarding from './Onboardings/Recipient';
import styles from './Owner.module.less';
import PetProfile from './Profiles/Pet';

type Props = {
    user: TelegramUser;
};

type View = Role.DONOR | Role.RECIPIENT | Role.BLOOD_FOUND;

const Owner: FC<Props> = ({ user }) => {
    const navigate = useNavigate();

    const [view, setView] = useState<View>(Role.RECIPIENT);
    const [selectedPet, setSelectedPet] = useState<Pet | null>(null);
    const [isPetProfileOpen, setIsPetProfileOpen] = useState<boolean>(false);

    const { data: pets = [], isLoading, refetch } = usePetsQuery(user.id);
    const { data: userData, isLoading: isUserDataLoading, refetch: refetchUserData } = useGetUserById(user.id);

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
        (label: 'startSearch' | 'activeSearch' | 'bloodFound', petId: string) => (e: MouseEvent<HTMLDivElement>) => {
            e.stopPropagation();

            if (label === 'startSearch') {
                navigate(`/adding#startSearch_${petId}`);
            }

            if (label === 'activeSearch') {
                navigate(`/search/${petId}`);
            }

            if (label === 'bloodFound') {
                navigate(`/search/${petId}?bloodFound=true`);
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
            case petsStatus === Role.BLOOD_FOUND: {
                return (
                    <div
                        onClick={onLabelClickHandler('bloodFound', petId)}
                        className={cn(styles.label, { [styles.bloodFound]: true })}
                    >
                        <div className={styles.searchIcon}>
                            <BloodFound />
                        </div>
                        <div className={styles.labelText}>Нашли кровь</div>
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

    useLayoutEffect(() => {
        if (window.location.hash === '#recipient') {
            setView(Role.RECIPIENT);

            return;
        }

        if (window.location.hash === '#donor') {
            setView(Role.DONOR);

            return;
        }

        window.location.hash = '#recipient';
    }, []);

    useEffect(() => {
        if (!pets) {
            setSelectedPet(null);

            return;
        }

        setSelectedPet((prevState) => {
            if (prevState) {
                return pets?.find((pet) => pet.id === prevState.id) || null;
            }

            return prevState;
        });
    }, [pets]);

    if (
        (view === Role.RECIPIENT || view === Role.BLOOD_FOUND) &&
        (!userData?.onBoarding || !userData?.onBoarding?.includes(Onboarding.FIND_BLOOD))
    ) {
        return (
            <RecipientOnboarding
                id={user.id}
                onboardings={userData?.onBoarding}
                onConfirmButtonClick={refetchUserData}
            />
        );
    }

    if (isPetProfileOpen && selectedPet) {
        return (
            <Layout>
                <PetProfile updatePets={refetch} onClose={onPetProfileToggleHandler(null)} {...selectedPet} />
            </Layout>
        );
    }

    return (
        <Layout>
            <div className={cn(styles.wrapper, { [styles.isPets]: !!pets?.length })}>
                <div className={styles.header}>
                    <div className={styles.avatar}>{user.fullName.charAt(0).toUpperCase()}</div>
                </div>
                {(isLoading || isUserDataLoading) && (
                    <div className={styles.loading}>
                        <Loading size={90} thickness={4} />
                    </div>
                )}
                {!(isLoading || isUserDataLoading) && !pets?.length && (
                    <div className={styles.button} onClick={onAddPetClickHandler}>
                        <div className={styles.pawIcon}>
                            <Paw />
                        </div>
                        <p className={styles.pawButtonText}>Добавить питомца</p>
                    </div>
                )}
                {!(isLoading || isUserDataLoading) && !!pets?.length && (
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
                                            {(view === Role.RECIPIENT || view === Role.BLOOD_FOUND) &&
                                                renderLabel(pet.petStatus, pet.id)}
                                        </div>
                                        <div className={styles.gradient} />
                                    </div>
                                </div>
                            ))}
                        </div>
                        <div className={styles.footer}>
                            <div
                                onClick={onButtonClickHandler(Role.DONOR)}
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
                                onClick={onButtonClickHandler(Role.RECIPIENT)}
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
