import cn from 'classnames';
import {
    BloodAndBreedGroupsDict,
    useBloodComponentsQuery,
    useLocationsQuery,
    usePetTypesAndBloodGroupsQuery,
} from 'hooks/useDicts';
import { usePetsQuery } from 'hooks/usePetsQuery';
import Header from 'pages/adding/common/Header';
import { FC, useCallback, useEffect, useMemo, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { createRecipient } from 'api/apiServices/createRecipient';
import { updatePet } from 'api/apiServices/updatePet';
import { Pet } from 'api/pets';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import { Role } from 'api/user';
import Alert from 'components/Alert';
import Curtain from 'components/Curtain';
import Loading from 'components/Loading';

import Check from '../Steps/Check';
import Final from '../Steps/Final';
import Second from '../Steps/Second';
import Three from '../Steps/Three';
import styles from './Search.module.less';

type Props = {
    petId?: string;
    userId: string;
};

const captions = ['Критерии поиска', 'Дополнительно'];

const Search: FC<Props> = ({ petId, userId }) => {
    const navigate = useNavigate();

    const [step, setStep] = useState<number>(1);
    const [isLoading, setIsLoading] = useState(false);

    const [locations, setLocations] = useState<string[]>([]);
    const [bloodGroup, setBloodGroup] = useState<string>('');
    const [bloodVolume, setBloodVolume] = useState<string>('');
    const [description, setDescription] = useState<string>('');
    const [bloodComponents, setBloodComponents] = useState<string[]>([]);
    const [desiredBloodGroups, setDesiredBloodGroups] = useState<string[]>([]);
    const [notifyOfSmallDonors, setNotifyOfSmallDonors] = useState(false);
    const [bloodRequestPhoto, setBloodRequestPhoto] = useState<File | null>(null);

    const { data: pets = [], isLoading: isPetsLoading, refetch } = usePetsQuery(userId);
    const { data: locationsDict = [], isError: isErrorLocations } = useLocationsQuery();
    const { data: bloodComponentsDict = [], isError: isErrorBloodComponents } = useBloodComponentsQuery();
    const { data: { bloodGroupDict = {} } = {}, isError: isErrorPetTypesAndBloodGroups } =
        usePetTypesAndBloodGroupsQuery();

    const selectedPet = useMemo(() => pets?.find((pet) => pet.id === petId), [petId, pets]);

    const showToast = useCallback(
        (text: string) => {
            toast.warn(text, {
                onClose: () => {
                    navigate('/owner');
                },
            });
        },
        [navigate],
    );

    const onBackClickHandler = () => {
        if (step > 1) {
            setStep((prevState) => prevState - 1);

            return;
        }

        navigate('/owner#recipient');
    };

    const onChangeBloodGroupHandler = (newBloodGroup: string) => () => {
        setBloodGroup(newBloodGroup);
        setDesiredBloodGroups([newBloodGroup]);
    };

    const onCancelClickHandler = () => {
        navigate('/owner#recipient');
    };

    const onChangeDesiredBloodGroupHandler = (newBloodGroups: string[]) => {
        setDesiredBloodGroups(newBloodGroups);
    };

    const onChangeBloodComponentsHandler = (newComponents: string[]) => {
        setBloodComponents(newComponents);
    };

    const onChangeLocationsHandler = (newLocations: string[]) => {
        setLocations(newLocations);
    };

    const onChangeBloodVolumeHandler = (newVolume: string) => {
        setBloodVolume(newVolume);
    };

    const onChangeNotifyOfSmallDonors = (isChecked: boolean) => {
        setNotifyOfSmallDonors(isChecked);
    };

    const onDescriptionChangeHandler = (newDescr: string) => {
        setDescription(newDescr);
    };

    const onLoadBloodRequestPhotoHandler = useCallback((newPhoto: File | null) => {
        setBloodRequestPhoto(newPhoto);
    }, []);

    const fetchCreateRecipient = async (confirmedStep: number) => {
        setIsLoading(true);

        const { success, error } = await createRecipient({
            petId,
            name: selectedPet?.name!,
            photo: null,
            userId,
            bloodRequestPhoto,
            petStatus: Role.NONE,
            type: selectedPet?.type as PetType,
            weightKg: Number(selectedPet?.weightKg),
            bloodGroup: bloodGroupDict[selectedPet?.type || ''].find((item) => item.value === bloodGroup)?.label,
            poolInfo: {
                description,
                bloodVolumeNeeded: Number(bloodVolume),
                regions: locations as unknown as number[],
                smallPetsNotifyAllowed: notifyOfSmallDonors,
                bloodComponentIds: bloodComponents as unknown as number[],
                bloodGroupNames: bloodGroupDict[selectedPet?.type || ''].reduce((res, item) => {
                    if (desiredBloodGroups.includes(item.value)) {
                        res.push(item.label);
                    }

                    return res;
                }, [] as string[]),
            },
        });

        if (success) {
            setStep(confirmedStep);
            setIsLoading(false);

            await queryClient.invalidateQueries({ queryKey: ['pets', userId] });
        } else {
            showToast(error || '');
        }
    };

    const onConfirmButtonClickHandler = (confirmedStep: number) => {
        if (confirmedStep !== 0 && step === 3) {
            fetchCreateRecipient(confirmedStep);

            return;
        }

        if (confirmedStep === 0) {
            setStep(1);

            return;
        }

        setStep(confirmedStep);
    };

    const onCurtainConfirmClickHandler = async () => {
        if (!petId) {
            return;
        }

        const { success } = await updatePet({ id: petId, bloodGroup } as Pet);

        if (success) {
            await refetch();
        } else {
            showToast('Не удалось сохранить выбранную группу крови, пожалуйста, попробуйте еще раз');
        }
    };

    useEffect(() => {
        if (selectedPet?.bloodGroup) {
            setDesiredBloodGroups([
                bloodGroupDict[selectedPet?.type || '']?.filter(({ label }) => label === selectedPet?.bloodGroup)[0]
                    ?.value,
            ]);
        }
    }, [bloodGroupDict, selectedPet]);

    useEffect(() => {
        if (isErrorBloodComponents) {
            showToast('Не удалось загрузить словарь компонентов крови, попробуйте перезагрузить приложение');
        }
    }, [isErrorBloodComponents, showToast]);

    useEffect(() => {
        if (isErrorLocations) {
            showToast('Не удалось загрузить словарь регионов, попробуйте перезагрузить приложение');
        }
    }, [isErrorLocations, showToast]);

    useEffect(() => {
        if (isErrorPetTypesAndBloodGroups) {
            showToast('Не удалось загрузить словарь типов животных и групп крови, попробуйте перезагрузить приложение');
        }
    }, [isErrorPetTypesAndBloodGroups, showToast]);

    if (!pets?.length) {
        return null;
    }

    return (
        <div className={styles.wrapper}>
            {step < 3 && (
                <>
                    <Header
                        step={step}
                        stepsCount={2}
                        caption={captions[step - 1]}
                        onBackClickHandler={onBackClickHandler}
                    />
                    <div className={styles.form}>
                        {step === 1 && (
                            <Second
                                locations={locations}
                                bloodVolume={bloodVolume}
                                locationsDict={locationsDict}
                                bloodComponents={bloodComponents}
                                petType={selectedPet?.type || ''}
                                desiredBloodGroups={desiredBloodGroups}
                                weight={`${selectedPet?.weightKg}` || ''}
                                notifyOfSmallDonors={notifyOfSmallDonors}
                                bloodComponentsDict={bloodComponentsDict}
                                onChangeLocations={onChangeLocationsHandler}
                                onChangeBloodVolume={onChangeBloodVolumeHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                onChangeBloodComponents={onChangeBloodComponentsHandler}
                                onChangeNotifyOfSmallDonors={onChangeNotifyOfSmallDonors}
                                bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
                                onChangeDesiredBloodGroups={onChangeDesiredBloodGroupHandler}
                                bloodGroup={
                                    bloodGroupDict[selectedPet?.type || '']?.filter(
                                        ({ label }) => label === selectedPet?.bloodGroup,
                                    )[0]?.value || selectedPet?.bloodGroup
                                }
                            />
                        )}
                        {step === 2 && (
                            <Three
                                photo={bloodRequestPhoto}
                                description={description}
                                onLoadPhoto={onLoadBloodRequestPhotoHandler}
                                onDescriptionChange={onDescriptionChangeHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                            />
                        )}
                    </div>
                </>
            )}
            {step === 3 && (
                <Check
                    photo={null}
                    isLoading={isLoading}
                    locations={locations}
                    bloodVolume={bloodVolume}
                    description={description}
                    locationsDict={locationsDict}
                    name={selectedPet?.name || ''}
                    petType={selectedPet?.type || ''}
                    bloodComponents={bloodComponents}
                    bloodRequestPhoto={bloodRequestPhoto}
                    photoUrl={selectedPet?.photoUrls?.[0]}
                    desiredBloodGroups={desiredBloodGroups}
                    weight={`${selectedPet?.weightKg}` || ''}
                    bloodComponentsDict={bloodComponentsDict}
                    notifyOfSmallDonors={notifyOfSmallDonors}
                    onConfirmButtonClick={onConfirmButtonClickHandler}
                    bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
                    bloodGroup={
                        bloodGroupDict[selectedPet?.type || '']?.filter(
                            ({ label }) => label === selectedPet?.bloodGroup,
                        )[0]?.value || selectedPet?.bloodGroup
                    }
                />
            )}
            {step === 4 && <Final fromSearch onBackToStart={onCancelClickHandler} />}
            {!selectedPet?.bloodGroup && (
                <Curtain
                    columnOfButtons
                    cancelButtonTitle='Далее'
                    title='Выберите группу крови'
                    confirmButtonTitle='Вернуться'
                    onConfirm={onCancelClickHandler}
                    onCancel={onCurtainConfirmClickHandler}
                    isDisableCancelButton={bloodGroup === ''}
                    subTitle={
                        <div className={styles.curtainSubtitle}>
                            Для поиска крови необходимо знать родную
                            <br />
                            группу крови питомца
                        </div>
                    }
                >
                    {!!bloodGroupDict && !!selectedPet?.type && (
                        <div className={styles.blood}>
                            <div
                                className={cn(styles.bloodGroups, {
                                    [styles.dogGroup]: selectedPet.type === PetType.DOG,
                                })}
                            >
                                {bloodGroupDict[selectedPet.type]?.map(({ label, value }) => (
                                    <div
                                        key={value}
                                        onClick={onChangeBloodGroupHandler(label)}
                                        className={cn(styles.bloodItem, { [styles.checked]: bloodGroup === label })}
                                    >
                                        {label}
                                    </div>
                                ))}
                            </div>
                            {selectedPet.type === PetType.DOG && (
                                <Alert
                                    className={styles.alert}
                                    text='Сведения вносятся по системе DEA.&nbsp;Не используйте данные из других систем (KAI, DAL)'
                                />
                            )}
                        </div>
                    )}
                </Curtain>
            )}
            {isPetsLoading && (
                <div className={styles.loading}>
                    <Loading size={90} thickness={4} />
                </div>
            )}
        </div>
    );
};

export default Search;
