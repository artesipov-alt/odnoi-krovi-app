import useBodyScrollLock from 'hooks/useBodyScrollLock';
import {
    BloodAndBreedGroupsDict,
    useBloodComponentsQuery,
    useLocationsQuery,
    usePetTypesAndBloodGroupsQuery,
} from 'hooks/useDicts';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { createRecipient } from 'api/apiServices/createRecipient';
import { queryClient } from 'api/queryClient';
import { PetType } from 'api/types';
import { Role } from 'api/user';

import Header from '../common/Header';
import styles from './Recipient.module.less';
import Check from './Steps/Check';
import Final from './Steps/Final';
import First from './Steps/First';
import Second from './Steps/Second';
import Three from './Steps/Three';

type Props = {
    userId: string;
    onBackToStart: () => void;
};

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
    const [bloodRequestPhoto, setBloodRequestPhoto] = useState<File | null>(null);
    const [includeUnknownBloodGroup, setIncludeUnknownBloodGroup] = useState(false);

    const [isLoading, setIsLoading] = useState(false);

    const { data: locationsDict = [], isError: isErrorLocations } = useLocationsQuery();
    const { data: bloodComponentsDict = [], isError: isErrorBloodComponents } = useBloodComponentsQuery();
    const { data: { petTypesDict = [], bloodGroupDict = {} } = {}, isError: isErrorPetTypesAndBloodGroups } =
        usePetTypesAndBloodGroupsQuery();

    useBodyScrollLock(isLoading);

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

    const fetchCreateRecipient = async (confirmedStep: number) => {
        setIsLoading(true);

        const { success, error } = await createRecipient({
            name,
            photo,
            userId,
            bloodRequestPhoto,
            petStatus: Role.NONE,
            type: petType as PetType,
            weightKg: Number(weight.replace(',', '.')),
            bloodGroup: bloodGroupDict[petType].find((item) => item.value === bloodGroup)?.label,
            poolInfo: {
                description,
                prioritySearch: false,
                includeUnknownBloodGroup,
                regions: locations as unknown as number[],
                smallPetsNotifyAllowed: notifyOfSmallDonors,
                bloodComponentIds: bloodComponents as unknown as number[],
                bloodVolumeNeeded: Number(bloodVolume.replace(',', '.')),
                bloodGroupNames: bloodGroupDict[petType].reduce((res, item) => {
                    if (desiredBloodGroups.includes(item.value)) {
                        res.push(item.label);
                    }

                    return res;
                }, [] as string[]),
            },
        });

        if (success) {
            setStep(confirmedStep + 1);

            await queryClient.invalidateQueries({ queryKey: ['pets', userId] });
        } else {
            showToast(error || '');
        }

        setIsLoading(false);
    };

    const onBackClickHandler = () => {
        if (step > 1) {
            setStep((prevState) => prevState - 1);

            return;
        }

        onBackToStart();
    };

    const onLoadBloodRequestPhotoHandler = useCallback((newPhoto: File | null) => {
        setBloodRequestPhoto(newPhoto);
    }, []);

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

    const onChangeIncludeUnknownBloodGroupHandler = (isChecked: boolean) => {
        setIncludeUnknownBloodGroup(isChecked);
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
        document.documentElement.classList.add('useWhiteBg1');
    }, []);

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
                                onLoadPhoto={onLoadPhotoHandler}
                                onChangeName={onChangeNameHandler}
                                onChangeWeight={onChangeWeightHandler}
                                onChangePetType={onChangePetTypeHandler}
                                onChangeBloodGroup={onChangeBloodGroupHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
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
                                bloodComponents={bloodComponents}
                                desiredBloodGroups={desiredBloodGroups}
                                notifyOfSmallDonors={notifyOfSmallDonors}
                                bloodComponentsDict={bloodComponentsDict}
                                onChangeLocations={onChangeLocationsHandler}
                                onChangeBloodVolume={onChangeBloodVolumeHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                includeUnknownBloodGroup={includeUnknownBloodGroup}
                                onChangeBloodComponents={onChangeBloodComponentsHandler}
                                onChangeNotifyOfSmallDonors={onChangeNotifyOfSmallDonors}
                                bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
                                onChangeDesiredBloodGroups={onChangeDesiredBloodGroupHandler}
                                onChangeIncludeUnknownBloodGroup={onChangeIncludeUnknownBloodGroupHandler}
                            />
                        )}
                        {step === 3 && (
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
                    bloodComponents={bloodComponents}
                    bloodRequestPhoto={bloodRequestPhoto}
                    desiredBloodGroups={desiredBloodGroups}
                    bloodComponentsDict={bloodComponentsDict}
                    notifyOfSmallDonors={notifyOfSmallDonors}
                    onConfirmButtonClick={onConfirmButtonClickHandler}
                    includeUnknownBloodGroup={includeUnknownBloodGroup}
                    bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
                />
            )}
            {step === 5 && <Final onBackToStart={onBackToStart} />}
        </>
    );
};

export default Recipient;
