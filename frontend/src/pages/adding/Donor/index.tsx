import Header from 'pages/adding/common/Header';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { createDonor } from 'api/apiServices/createDonor';
import { getBloodGroups } from 'api/apiServices/getBloodGroups';
import { getBreedsByType } from 'api/apiServices/getBreedsByType';
import { getGenders } from 'api/apiServices/getGenders';
import { getHealthStatuses } from 'api/apiServices/getHealthStatuses';
import { getLivingConditions } from 'api/apiServices/getLivingConditions';
import { getPetsTypes } from 'api/apiServices/getPetsTypes';
import { getReproductiveStatuses } from 'api/apiServices/getReproductiveStatuses';
import { Analyses, AnalysesItem } from 'api/pets';
import { Dict, PetGenderDict, PetTypeDict, StringDict } from 'api/reference';
import { AnalysesMapping, PetGender, PetType } from 'api/types';
import { Role } from 'api/user';

import styles from './Donor.module.less';
import Onboarding from './Onboarding';
import Check from './Steps/Check';
import Fifth from './Steps/Fifth';
import Final from './Steps/Final';
import First from './Steps/First';
import Fourth from './Steps/Fourth';
import Second from './Steps/Second';
import Third from './Steps/Third';
import { Analiz, BirthDate } from './types';

type GroupsDicts = Record<PetType, Dict[]>;

type Props = {
    userId: string;
    onBackToStart: () => void;
};

const captions = ['О питомце', 'Параметры', 'Здоровье', 'Обработки', 'Анализы'];

const Donor: FC<Props> = ({ userId, onBackToStart }) => {
    const navigate = useNavigate();

    const [step, setStep] = useState(1);
    const [isLoading, setIsLoading] = useState(false);
    const [isOnboardingFinish, setIsOnboardingFinish] = useState(false);

    // 1 step
    const [name, setName] = useState('');
    const [petType, setPetType] = useState<string>('');
    const [photo, setPhoto] = useState<File | null>(null);
    const [petGender, setPetGender] = useState<string>('');
    const [chipNumber, setChipNumber] = useState<string>(''); // 'none' - значит отсутствует
    const [exactDate, setExactDate] = useState<Date | null>(null);
    const [birthDateType, setBirthDateType] = useState<BirthDate | null>(null);
    const [approximateDateYear, setApproximateDateYear] = useState<string>('');
    const [approximateDateMonth, setApproximateDateMonth] = useState<string>('');

    // 2 step
    const [weight, setWeight] = useState<string>('');
    const [breed, setBreed] = useState<Dict | null>(null);
    const [bloodGroup, setBloodGroup] = useState<string>('');
    const [livingCondition, setLivingCondition] = useState<string>('');
    const [reproductiveStatus, setReproductiveStatus] = useState<string>('');

    // 3 step
    const [healthStatus, setHealthStatus] = useState<string>('');
    const [surgicalList, setSurgicalList] = useState<string>('');
    const [medicationsList, setMedicationsList] = useState<string>('');
    const [lastDonation, setLastDonation] = useState<Date | null>(null);
    const [isNoLastDonation, setIsNoLastDonation] = useState<boolean>(false);
    const [wasBloodTransfusion, setWasBloodTransfusion] = useState<boolean | null>(null);
    const [isTakingMedications, setIsTakingMedications] = useState<boolean | null>(null);
    const [wasSurgicalInterventions, setWasSurgicalInterventions] = useState<boolean | null>(null);

    // 4 step
    const [isNoDeworming, setIsNoDeworming] = useState<boolean>(false);
    const [dewormingDate, setDewormingDate] = useState<Date | null>(null);
    const [isNoRabiesVaccination, setIsNoRabiesVaccination] = useState<boolean>(false);
    const [rabiesVaccinationDate, setRabiesVaccinationDate] = useState<Date | null>(null);
    const [isNoInfectionsVaccination, setIsNoInfectionsVaccination] = useState<boolean>(false);
    const [isNoEctoparasitesTreatment, setIsNoEctoparasitesTreatment] = useState<boolean>(false);
    const [infectionsVaccinationDate, setInfectionsVaccinationDate] = useState<Date | null>(null);
    const [ectoparasitesTreatmentDate, setEctoparasitesTreatmentDate] = useState<Date | null>(null);

    // 5 step
    const [leukemia, setLeukemia] = useState<Analiz>({
        type: 'leukemia',
        name: 'Анализ на лейкоз (ВЛК)',
        items: [
            { name: 'ПЦР', value: null },
            { name: 'ИФА', value: null },
            { name: 'ИХА', value: null },
        ],
    });
    const [immunodeficiency, setImmunodeficiency] = useState<Analiz>({
        type: 'immunodeficiency',
        name: 'Анализ на иммунодефицит (ВИК)',
        items: [
            { name: 'ПЦР', value: null },
            { name: 'ИФА', value: null },
            { name: 'ИХА', value: null },
        ],
    });
    const [hemoplasmosis, setHemoplasmosis] = useState<Analiz>({
        type: 'hemoplasmosis',
        name: 'Анализ на гемоплазмоз',
        items: [{ name: 'ПЦР', value: null }],
    });
    const [bartonellosis, setBartonellosis] = useState<Analiz>({
        type: 'bartonellosis',
        name: 'Анализ на бартонеллез',
        items: [{ name: 'ПЦР', value: null }],
    });
    const [babesiosis, setBabesiosis] = useState<Analiz>({
        type: 'babesiosis',
        name: 'Анализ на бабезиоз',
        items: [
            { name: 'ПЦР', value: null },
            { name: 'Микроскопия мазка', value: null },
        ],
    });
    const [dirofilaria, setDirofilaria] = useState<Analiz>({
        type: 'dirofilaria',
        name: 'Анализ на дирофиляриоз',
        items: [
            { name: 'ПЦР', value: null },
            { name: 'Микроскопия мазка', value: null },
        ],
    });
    const [ehrlichiosis, setEhrlichiosis] = useState<Analiz>({
        type: 'ehrlichiosis',
        name: 'Анализ на эрлихиоз',
        items: [
            { name: 'ПЦР', value: null },
            { name: 'Экспресс-тест', value: null },
        ],
    });
    const [anaplasmosis, setAnaplasmosis] = useState<Analiz>({
        type: 'anaplasmosis',
        name: 'Анализ на анаплазмоз',
        items: [
            { name: 'ПЦР', value: null },
            { name: 'Экспресс-тест', value: null },
        ],
    });

    // dicts
    const [petTypesDict, setPetTypesDict] = useState<PetTypeDict[]>([]);
    const [petGendersDict, setPetGendersDict] = useState<PetGenderDict[]>([]);
    const [breedsDict, setBreedsDict] = useState<GroupsDicts>({} as GroupsDicts);
    const [healthStatusesDict, setHealthStatusesDict] = useState<StringDict[]>([]);
    const [livingConditionsDict, setLivingConditionsDict] = useState<StringDict[]>([]);
    const [bloodGroupDict, setBloodGroupDict] = useState<GroupsDicts>({} as GroupsDicts);
    const [reproductiveStatusesDict, setReproductiveStatusesDict] = useState<StringDict[]>([]);

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
            const dict: GroupsDicts = {} as GroupsDicts;

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

    const fetchBreedsTypes = useCallback(
        async (pets: PetTypeDict[]) => {
            const dict: GroupsDicts = {} as GroupsDicts;

            await Promise.allSettled(
                pets.map(async ({ value }) => {
                    const response = await getBreedsByType(value);

                    if (response) {
                        dict[value] = response;
                    }

                    return { [value]: response };
                }),
            ).then((result) => {
                if (result.some(({ status }) => status === 'rejected')) {
                    showToast('Не удалось загрузить словарь пород, попробуйте еще раз');
                }
            });

            setBreedsDict(dict);
        },
        [showToast],
    );

    const fetchPetTypes = useCallback(async () => {
        const response = await getPetsTypes();

        if (!response) {
            showToast('Не удалось загрузить словарь типов животных, попробуйте еще раз');

            return;
        }

        fetchBloodTypes(response);
        fetchBreedsTypes(response);

        setPetTypesDict(response);
    }, [fetchBloodTypes, fetchBreedsTypes, showToast]);

    const fetchGenders = useCallback(async () => {
        const response = await getGenders();

        if (!response) {
            showToast('Не удалось загрузить словарь полов, попробуйте еще раз');

            return;
        }

        setPetGendersDict(response);
    }, [showToast]);

    const fetchLivingConditions = useCallback(async () => {
        const response = await getLivingConditions();

        if (!response) {
            showToast('Не удалось загрузить словарь условий проживания, попробуйте еще раз');

            return;
        }

        setLivingConditionsDict(response);
    }, [showToast]);

    const fetchHealthStatuses = useCallback(async () => {
        const response = await getHealthStatuses();

        if (!response) {
            showToast('Не удалось загрузить словарь статусов здоровья, попробуйте еще раз');

            return;
        }

        setHealthStatusesDict(response);
    }, [showToast]);

    const fetchReproductiveStatuses = useCallback(async () => {
        const response = await getReproductiveStatuses();

        if (!response) {
            showToast('Не удалось загрузить словарь репродуктивных состояний, попробуйте еще раз');

            return;
        }

        setReproductiveStatusesDict(response);
    }, [showToast]);

    const onFinishOnboardingHandler = () => {
        setIsOnboardingFinish(true);
    };

    const onBackClickHandler = () => {
        if (step > 1) {
            setStep((prevState) => prevState - 1);

            return;
        }

        onBackToStart();
    };

    const resetDependentFields = () => {
        setWeight('');
        setBreed(null);
        setBloodGroup('');

        setLeukemia((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setImmunodeficiency((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setHemoplasmosis((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setBartonellosis((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setBabesiosis((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setDirofilaria((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setEhrlichiosis((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
        setAnaplasmosis((prevState) => ({
            ...prevState,
            items: prevState.items.map((item) => ({ ...item, value: null })),
        }));
    };

    const onLoadPhotoHandler = useCallback((newPhoto: File | null) => {
        setPhoto(newPhoto);
    }, []);

    const onChangeNameHandler = (newName: string) => {
        setName(newName);
    };

    const onChangePetTypeHandler = (newType: string) => {
        setPetType(newType);

        resetDependentFields();
    };

    const onChangePetGenderHandler = (newGender: string) => {
        setPetGender(newGender);
        setReproductiveStatus('');
    };

    const onChangeChipNumberHandler = (newValue: string) => {
        setChipNumber(newValue);
    };

    const onChangeBirthDateTypeHandler = (newType: BirthDate) => {
        setBirthDateType(newType);

        if (exactDate) {
            setExactDate(null);
        }

        if (approximateDateYear) {
            setApproximateDateYear('');
        }

        if (approximateDateMonth) {
            setApproximateDateMonth('');
        }
    };

    const onChangeExactDateHandler = (newDate: Date | null) => {
        setExactDate(newDate);
    };

    const onChangeApproximateDateYearHandler = (newYear: string) => {
        setApproximateDateYear(newYear);
    };

    const onChangeApproximateDateMonthHandler = (newYear: string) => {
        setApproximateDateMonth(newYear);
    };

    const onChangeBloodGroupHandler = (newBloodGroup: string) => {
        setBloodGroup(newBloodGroup);
    };

    const onChangeWeightHandler = (newWeight: string) => {
        setWeight(newWeight);
    };

    const onChangeBreedHandler = (newBreed: Dict | null) => {
        setBreed(newBreed);
    };

    const onChangeLivingConditionHandler = (newCondition: string) => {
        setLivingCondition(newCondition);
    };

    const onChangeHealthStatusHandler = (newStatus: string) => {
        setHealthStatus(newStatus);
    };

    const onChangeReproductiveStatusHandler = (newStatus: string) => {
        setReproductiveStatus(newStatus);
    };

    const onChangeLastDonationDateHandler = (newDate: Date | null) => {
        setLastDonation(newDate);
        setIsNoLastDonation(false);
    };

    const onChangeIsNoLastDonationHandler = (newValue: boolean) => {
        setLastDonation(null);
        setIsNoLastDonation(newValue);
    };

    const onChangeWasBloodTransfusionHandler = (newValue: boolean) => {
        setWasBloodTransfusion(newValue);
    };

    const onChangeIsTakingMedicationsHandler = (newValue: boolean) => {
        setIsTakingMedications(newValue);
    };

    const onChangeWasSurgicalInterventionsHandler = (newValue: boolean) => {
        setWasSurgicalInterventions(newValue);
    };

    const onChangeMedicationsListHandler = (newValue: string) => {
        setMedicationsList(newValue);
    };

    const onChangeSurgicalListHandler = (newValue: string) => {
        setSurgicalList(newValue);
    };

    const onChangeRabiesVaccinationDateHandler = (newDate: Date | null) => {
        setRabiesVaccinationDate(newDate);
        setIsNoRabiesVaccination(false);
    };

    const onChangeIsNoRabiesVaccinationHandler = (newValue: boolean) => {
        setIsNoRabiesVaccination(newValue);
        setRabiesVaccinationDate(null);
    };

    const onChangeInfectionsVaccinationDateHandler = (newDate: Date | null) => {
        setInfectionsVaccinationDate(newDate);
        setIsNoInfectionsVaccination(false);
    };

    const onChangeIsNoInfectionsVaccinationHandler = (newValue: boolean) => {
        setIsNoInfectionsVaccination(newValue);
        setInfectionsVaccinationDate(null);
    };

    const onChangeEctoparasitesTreatmentDateHandler = (newDate: Date | null) => {
        setEctoparasitesTreatmentDate(newDate);
        setIsNoEctoparasitesTreatment(false);
    };

    const onChangeIsNoEctoparasitesTreatmentHandler = (newValue: boolean) => {
        setIsNoEctoparasitesTreatment(newValue);
        setEctoparasitesTreatmentDate(null);
    };

    const onChangeDewormingDateHandler = (newDate: Date | null) => {
        setDewormingDate(newDate);
        setIsNoDeworming(false);
    };

    const onChangeIsNoDewormingHandler = (newValue: boolean) => {
        setIsNoDeworming(newValue);
        setDewormingDate(null);
    };

    const onChangeAnaliz = (type: string, analizName: string, value: Date | null) => {
        switch (type) {
            case 'leukemia': {
                setLeukemia((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'immunodeficiency': {
                setImmunodeficiency((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'hemoplasmosis': {
                setHemoplasmosis((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'bartonellosis': {
                setBartonellosis((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'babesiosis': {
                setBabesiosis((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'dirofilaria': {
                setDirofilaria((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'ehrlichiosis': {
                setEhrlichiosis((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));

                return;
            }
            case 'anaplasmosis': {
                setAnaplasmosis((prevState) => ({
                    ...prevState,
                    items: prevState.items.map((item) => {
                        if (item.name === analizName) {
                            return { ...item, value };
                        }

                        return item;
                    }),
                }));
            }
            default: {
                break;
            }
        }
    };

    const getAnalizesToRequest = () => {
        const result =
            petType === PetType.DOG
                ? [babesiosis, dirofilaria, hemoplasmosis, bartonellosis, ehrlichiosis, anaplasmosis]
                : [leukemia, immunodeficiency, hemoplasmosis, bartonellosis].reduce((acc, analiz) => {
                      acc[analiz.type] = analiz.items.reduce(
                          (res, item) => {
                              if (item.value) {
                                  res.push({
                                      analysisDate: item.value,
                                      analysisName: analiz.type,
                                      analysisType: AnalysesMapping[item.name],
                                  });
                              }

                              return res;
                          },
                          [] as unknown as AnalysesItem[],
                      );

                      return acc;
                  }, {} as Analyses);

        const filteredResult = Object.keys(result).reduce((res, key) => {
            if (result[key].length) {
                res[key] = result[key];
            }

            return res;
        }, {} as Analyses);

        if (Object.keys(filteredResult).length) {
            return filteredResult;
        }

        return undefined;
    };

    const fetchCreateDonor = async (confirmedStep: number) => {
        setIsLoading(true);

        const { success } = await createDonor({
            name,
            userId,
            photo,
            breedId: breed?.value,
            petStatus: Role.DONOR,
            weightKg: Number(weight),
            type: petType as PetType,
            gender: petGender as PetGender,
            birthDate: exactDate || undefined,
            ageYears: Number(approximateDateYear) || undefined,
            ageMonths: Number(approximateDateMonth) || undefined,
            chipNumber: chipNumber === 'none' ? undefined : chipNumber,
            bloodGroup: bloodGroupDict[petType].find((item) => item.value === bloodGroup)?.label,
            livingCondition: livingConditionsDict.find((item) => item.value === livingCondition)?.value,
            health: {
                healthStatus: healthStatusesDict.find((item) => item.value === healthStatus)?.value!,
                lastDonation: lastDonation || undefined,
                transfused: wasBloodTransfusion || undefined,
                medications: isTakingMedications ? medicationsList : undefined,
                surgicalInterventions: wasSurgicalInterventions ? surgicalList : undefined,
                reproductiveStatus: reproductiveStatusesDict.find((item) => item.value === reproductiveStatus)?.value,
            },
            treatments: {
                dewormingDate: dewormingDate || undefined,
                rabiesVaccinationDate: rabiesVaccinationDate || undefined,
                infectionVaccinationDate: infectionsVaccinationDate || undefined,
                ectoparasiteTreatmentDate: ectoparasitesTreatmentDate || undefined,
            },
            analyses: getAnalizesToRequest(),
        });

        if (success) {
            setStep(confirmedStep + 1);
            setIsLoading(false);
        } else {
            showToast('Не удалось сохранить питомца, попробуйте еще раз');
        }
    };

    const onConfirmButtonClickHandler = (confirmedStep: number) => {
        if (confirmedStep !== 0 && step === 6) {
            fetchCreateDonor(confirmedStep);

            return;
        }

        setStep(confirmedStep + 1);
    };

    useEffect(() => {
        fetchPetTypes();
        fetchGenders();
        fetchHealthStatuses();
        fetchLivingConditions();
        fetchReproductiveStatuses();
    }, [fetchGenders, fetchHealthStatuses, fetchLivingConditions, fetchPetTypes, fetchReproductiveStatuses]);

    useEffect(() => {
        document.documentElement.classList.add('useWhiteBg1');
    }, []);

    if (!isOnboardingFinish) {
        return <Onboarding onFinish={onFinishOnboardingHandler} onBackToStart={onBackToStart} />;
    }

    return (
        <>
            {step < 6 && (
                <>
                    <Header
                        step={step}
                        stepsCount={5}
                        caption={captions[step - 1]}
                        onBackClickHandler={onBackClickHandler}
                    />
                    <div className={styles.form}>
                        {step === 1 && (
                            <First
                                name={name}
                                photo={photo}
                                petType={petType}
                                exactDate={exactDate}
                                petGender={petGender}
                                petTypes={petTypesDict}
                                chipNumber={chipNumber}
                                petGenders={petGendersDict}
                                birthDateType={birthDateType}
                                onLoadPhoto={onLoadPhotoHandler}
                                onChangeName={onChangeNameHandler}
                                onChangePetType={onChangePetTypeHandler}
                                approximateDateYear={approximateDateYear}
                                approximateDateMonth={approximateDateMonth}
                                onChangeExactDate={onChangeExactDateHandler}
                                onChangePetGender={onChangePetGenderHandler}
                                onChangeChipNumber={onChangeChipNumberHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                onChangeBirthDateType={onChangeBirthDateTypeHandler}
                                onChangeApproximateDateYear={onChangeApproximateDateYearHandler}
                                onChangeApproximateDateMonth={onChangeApproximateDateMonthHandler}
                            />
                        )}
                        {step === 2 && (
                            <Second
                                breed={breed}
                                weight={weight}
                                petType={petType}
                                petGender={petGender}
                                bloodGroup={bloodGroup}
                                breedsDict={breedsDict}
                                bloodGroupDict={bloodGroupDict}
                                livingCondition={livingCondition}
                                onChangeBreed={onChangeBreedHandler}
                                onChangeWeight={onChangeWeightHandler}
                                reproductiveStatus={reproductiveStatus}
                                livingConditionsDict={livingConditionsDict}
                                onChangeBloodGroup={onChangeBloodGroupHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                reproductiveStatusesDict={reproductiveStatusesDict}
                                onChangeLivingCondition={onChangeLivingConditionHandler}
                                onChangeReproductiveStatus={onChangeReproductiveStatusHandler}
                            />
                        )}
                        {step === 3 && (
                            <Third
                                healthStatus={healthStatus}
                                lastDonation={lastDonation}
                                surgicalList={surgicalList}
                                medicationsList={medicationsList}
                                isNoLastDonation={isNoLastDonation}
                                healthStatusesDict={healthStatusesDict}
                                wasBloodTransfusion={wasBloodTransfusion}
                                isTakingMedications={isTakingMedications}
                                onChangeHealthStatus={onChangeHealthStatusHandler}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                onChangeSurgicalList={onChangeSurgicalListHandler}
                                wasSurgicalInterventions={wasSurgicalInterventions}
                                onChangeMedicationsList={onChangeMedicationsListHandler}
                                onChangeLastDonationDate={onChangeLastDonationDateHandler}
                                onChangeIsNoLastDonation={onChangeIsNoLastDonationHandler}
                                onChangeIsTakingMedications={onChangeIsTakingMedicationsHandler}
                                onChangeWasBloodTransfusion={onChangeWasBloodTransfusionHandler}
                                onChangeWasSurgicalInterventions={onChangeWasSurgicalInterventionsHandler}
                            />
                        )}
                        {step === 4 && (
                            <Fourth
                                isNoDeworming={isNoDeworming}
                                dewormingDate={dewormingDate}
                                rabiesVaccinationDate={rabiesVaccinationDate}
                                isNoRabiesVaccination={isNoRabiesVaccination}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                                onChangeIsNoDeworming={onChangeIsNoDewormingHandler}
                                onChangeDewormingDate={onChangeDewormingDateHandler}
                                infectionsVaccinationDate={infectionsVaccinationDate}
                                isNoInfectionsVaccination={isNoInfectionsVaccination}
                                ectoparasitesTreatmentDate={ectoparasitesTreatmentDate}
                                isNoEctoparasitesTreatment={isNoEctoparasitesTreatment}
                                onChangeRabiesVaccinationDate={onChangeRabiesVaccinationDateHandler}
                                onChangeIsNoRabiesVaccination={onChangeIsNoRabiesVaccinationHandler}
                                onChangeIsNoInfectionsVaccination={onChangeIsNoInfectionsVaccinationHandler}
                                onChangeInfectionsVaccinationDate={onChangeInfectionsVaccinationDateHandler}
                                onChangeIsNoEctoparasitesTreatment={onChangeIsNoEctoparasitesTreatmentHandler}
                                onChangeEctoparasitesTreatmentDate={onChangeEctoparasitesTreatmentDateHandler}
                            />
                        )}
                        {step === 5 && (
                            <Fifth
                                petType={petType}
                                leukemia={leukemia}
                                babesiosis={babesiosis}
                                dirofilaria={dirofilaria}
                                anaplasmosis={anaplasmosis}
                                ehrlichiosis={ehrlichiosis}
                                bartonellosis={bartonellosis}
                                hemoplasmosis={hemoplasmosis}
                                onChangeAnaliz={onChangeAnaliz}
                                immunodeficiency={immunodeficiency}
                                onConfirmButtonClick={onConfirmButtonClickHandler}
                            />
                        )}
                    </div>
                </>
            )}
            {step === 6 && (
                <Check
                    name={name}
                    photo={photo}
                    weight={weight}
                    leukemia={leukemia}
                    isLoading={isLoading}
                    exactDate={exactDate}
                    petTypeCode={petType}
                    chipNumber={chipNumber}
                    babesiosis={babesiosis}
                    dirofilaria={dirofilaria}
                    anaplasmosis={anaplasmosis}
                    ehrlichiosis={ehrlichiosis}
                    surgicalList={surgicalList}
                    lastDonation={lastDonation}
                    dewormingDate={dewormingDate}
                    bartonellosis={bartonellosis}
                    hemoplasmosis={hemoplasmosis}
                    medicationsList={medicationsList}
                    immunodeficiency={immunodeficiency}
                    wasBloodTransfusion={wasBloodTransfusion}
                    approximateDateYear={approximateDateYear}
                    approximateDateMonth={approximateDateMonth}
                    rabiesVaccinationDate={rabiesVaccinationDate}
                    onConfirmButtonClick={onConfirmButtonClickHandler}
                    infectionsVaccinationDate={infectionsVaccinationDate}
                    ectoparasitesTreatmentDate={ectoparasitesTreatmentDate}
                    petType={petTypesDict.filter(({ value }) => value === petType)[0].label}
                    petGender={petGendersDict.filter(({ value }) => value === petGender)[0].label}
                    breed={breedsDict[petType].filter(({ value }) => value === breed?.value)[0].label}
                    bloodGroup={bloodGroupDict[petType].filter(({ value }) => value === bloodGroup)[0]?.label}
                    livingCondition={livingConditionsDict.filter(({ value }) => value === livingCondition)[0].label}
                    healthStatus={healthStatusesDict.filter(({ value }) => value === healthStatus)[0].label}
                    reproductiveStatus={
                        reproductiveStatusesDict.filter(({ value }) => value === reproductiveStatus)[0]?.label
                    }
                />
            )}
            {step === 7 && <Final photo={photo} onBackToStart={onBackToStart} />}
        </>
    );
};

export default Donor;
