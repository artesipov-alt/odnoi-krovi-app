import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { createPet } from 'api/apiServices/createPet';
import { getBloodComponents } from 'api/apiServices/getBloodComponents';
import { getBloodGroups } from 'api/apiServices/getBloodGroups';
import { getLocations } from 'api/apiServices/getLocations';
import { getPetsTypes } from 'api/apiServices/getPetsTypes';
import { Dict, PetTypeDict } from 'api/reference';
import { PetType } from 'api/types';
import { Role } from 'api/user';

import Header from '../common/Header';
import styles from './Recipient.module.less';
import Check from './Steps/Check';
import Final from './Steps/Final';
import First from './Steps/First';
import Second from './Steps/Second';
import Three from './Steps/Three';
import fs from 'fs';

type Props = {
    userId: string;
    onBackToStart: () => void;
};

type BloodGroupsDicts = Record<PetType, Dict[]>;

const captions = ['О питомце', 'Критерии поиска', 'Дополнительно'];

const Recipient: FC<Props> = ({ userId, onBackToStart }) => {
    const navigate = useNavigate();

    const [name, setName] = useState('');
    const [step, setStep] = useState<number>(1);
    const [weight, setWeight] = useState<string>('');
    const [petType, setPetType] = useState<string>('');
    const [photo, setPhoto] = useState<File | null>(null);
    const [bloodGroup, setBloodGroup] = useState<string>('');
    const [locations, setLocations] = useState<string[]>([]);
    const [bloodVolume, setBloodVolume] = useState<string>('');
    const [description, setDescription] = useState<string>('');
    const [bloodComponents, setBloodComponents] = useState<string[]>([]);
    const [notifyOfSmallDonors, setNotifyOfSmallDonors] = useState(false);
    const [desiredBloodGroups, setDesiredBloodGroups] = useState<string[]>([]);

    const [isLoading, setIsLoading] = useState(false);

    const [locationsDict, setLocationsDict] = useState<Dict[]>([]);
    const [petTypesDict, setPetTypesDict] = useState<PetTypeDict[]>([]);
    const [bloodComponentsDict, setBloodComponentsDict] = useState<Dict[]>([]);
    const [bloodGroupDict, setBloodGroupDict] = useState<BloodGroupsDicts>({} as BloodGroupsDicts);

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

    const fetchBloodTypes = useCallback(
        async (pets: PetTypeDict[]) => {
            const dict: BloodGroupsDicts = {} as BloodGroupsDicts;

            await Promise.allSettled(
                pets.map(async ({ value }) => {
                    const response = await getBloodGroups(value);

                    if (response) {
                        dict[value] = response;
                    }

                    return { [value]: response };
                }),
            ).then((result) => {
                if (result.some(({ status }) => status === 'rejected')) {
                    showToast('Не удалось загрузить словарь групп крови, попробуйте еще раз');
                }
            });

            setBloodGroupDict(dict);
        },
        [showToast],
    );

    const fetchBloodComponents = useCallback(async () => {
        const response = await getBloodComponents();

        if (!response) {
            showToast('Не удалось загрузить словарь компонентов крови, попробуйте еще раз');

            return;
        }

        setBloodComponentsDict(response);
    }, [showToast]);

    const fetchLocations = useCallback(async () => {
        const response = await getLocations();

        if (!response) {
            showToast('Не удалось загрузить словарь регионов, попробуйте еще раз');

            return;
        }

        setLocationsDict(response);
    }, [showToast]);

    const fetchPetTypes = useCallback(async () => {
        const response = await getPetsTypes();

        if (!response) {
            showToast('Не удалось загрузить словарь типов животных, попробуйте еще раз');

            return;
        }

        fetchBloodTypes(response);

        setPetTypesDict(response);
    }, [fetchBloodTypes, showToast]);

    const fetchCreateRecipient = async (confirmedStep: number) => {
        // const { success } = await createPet({
        //     name,
        //     photo,
        //     userId,
        //     type: petType,
        //     weightKg: Number(weight),
        //     petStatus: Role.RECIPIENT,
        //     bloodGroup: `${bloodGroup}`,
        //     poolInfo: {
        //         description,
        //         bloodVolumeNeeded: Number(bloodVolume),
        //         regions: locations as unknown as number[],
        //         smallPetsNotifyAllowed: notifyOfSmallDonors,
        //         bloodComponentIds: bloodComponents as unknown as number[],
        //         bloodGroupIds: desiredBloodGroups.map((item) => String(item)),
        //     },
        // });

        // if (success) {
        if (true) {
            setStep(confirmedStep + 1);
        } else {
            showToast('Не удалось сохранить питомца, попробуйте еще раз');
        }
    };

    const onBackClickHandler = () => {
        if (step > 1) {
            setStep((prevState) => prevState - 1);

            return;
        }

        onBackToStart();
    };

    const onLoadPhotoHandler = useCallback((newPhoto: File | null) => {
        setPhoto(newPhoto);
    }, []);

    const onChangeNameHandler = (newName: string) => {
        setName(newName);
    };

    const onChangePetTypeHandler = (newType: string) => {
        setPetType(newType);
        setBloodVolume('');
        setDesiredBloodGroups([]);
    };

    const onChangeWeightHandler = (newWeight: string) => {
        setWeight(newWeight);
        setBloodVolume('');
    };

    const onChangeBloodGroupHandler = (newBloodGroup: string) => {
        setBloodGroup(newBloodGroup);
        setDesiredBloodGroups([newBloodGroup]);
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

    const onConfirmButtonClickHandler = (confirmedStep: number) => {
        if (confirmedStep !== 0 && step === 4) {
            fetchCreateRecipient(confirmedStep);

            return;
        }

        setStep(confirmedStep + 1);
    };

    useEffect(() => {
        fetchPetTypes();
        fetchLocations();
        fetchBloodComponents();
    }, [fetchBloodComponents, fetchLocations, fetchPetTypes]);

    useEffect(() => {
        document.documentElement.classList.add('useWhiteBg1');
    }, []);

    return (
        <>
            {step < 4 && (
                <>
                    <Header
                        step={step}
                        stepsCount={3}
                        caption={captions[step - 1]}
                        onBackClickHandler={onBackClickHandler}
                    />
                    <div className={styles.form}>
                        {step === 1 && (
                            <First
                                name={name}
                                photo={photo}
                                weight={weight}
                                petType={petType}
                                petTypes={petTypesDict}
                                bloodGroup={bloodGroup}
                                bloodGroupDict={bloodGroupDict}
                                onLoadPhoto={onLoadPhotoHandler}
                                onChangeName={onChangeNameHandler}
                                onChangeWeight={onChangeWeightHandler}
                                onChangePetType={onChangePetTypeHandler}
                                onChangeBloodGroup={onChangeBloodGroupHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                            />
                        )}
                        {step === 2 && (
                            <Second
                                weight={weight}
                                petType={petType}
                                locations={locations}
                                bloodGroup={bloodGroup}
                                bloodVolume={bloodVolume}
                                locationsDict={locationsDict}
                                bloodGroupDict={bloodGroupDict}
                                bloodComponents={bloodComponents}
                                desiredBloodGroups={desiredBloodGroups}
                                notifyOfSmallDonors={notifyOfSmallDonors}
                                bloodComponentsDict={bloodComponentsDict}
                                onChangeLocations={onChangeLocationsHandler}
                                onChangeBloodVolume={onChangeBloodVolumeHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                onChangeBloodComponents={onChangeBloodComponentsHandler}
                                onChangeNotifyOfSmallDonors={onChangeNotifyOfSmallDonors}
                                onChangeDesiredBloodGroups={onChangeDesiredBloodGroupHandler}
                            />
                        )}
                        {step === 3 && (
                            <Three
                                photo={photo}
                                description={description}
                                onLoadPhoto={onLoadPhotoHandler}
                                onDescriptionChange={onDescriptionChangeHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                            />
                        )}
                    </div>
                </>
            )}
            {step === 4 && (
                <Check
                    name={name}
                    photo={photo}
                    weight={weight}
                    petType={petType}
                    isLoading={isLoading}
                    locations={locations}
                    bloodGroup={bloodGroup}
                    bloodVolume={bloodVolume}
                    description={description}
                    locationsDict={locationsDict}
                    bloodGroupDict={bloodGroupDict}
                    bloodComponents={bloodComponents}
                    desiredBloodGroups={desiredBloodGroups}
                    bloodComponentsDict={bloodComponentsDict}
                    notifyOfSmallDonors={notifyOfSmallDonors}
                    onConfirmButtonClick={onConfirmButtonClickHandler}
                />
            )}
            {step === 5 && <Final onBackToStart={onBackToStart} />}
        </>
    );
};

export default Recipient;
