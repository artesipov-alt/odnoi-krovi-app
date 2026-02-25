import cn from 'classnames';
import {
    BloodAndBreedGroupsDict,
    useGendersQuery,
    useHealthStatusesQuery,
    useLivingConditionsQuery,
    usePetTypesAndBloodGroupsQuery,
    useReproductiveStatusesQuery,
} from 'hooks/useDicts';
import AccordionArrow from 'imgs/svg/accordionArrow';
import Analizes from 'imgs/svg/analizes';
import BackAngularArrow from 'imgs/svg/backAngularArrow';
import Basket from 'imgs/svg/basket';
import Blood from 'imgs/svg/blood';
import DonorButton from 'imgs/svg/donorButton';
import Edit from 'imgs/svg/edit';
import Health from 'imgs/svg/health';
import Info from 'imgs/svg/info';
import Params from 'imgs/svg/params';
import Processing from 'imgs/svg/processing';
import RecipientButton from 'imgs/svg/recipientButton';
import { FC, useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import { toast } from 'react-toastify';

import { addPhoto } from 'api/apiServices/addPhoto';
import { deletePetById } from 'api/apiServices/deletePetById';
import { Pet } from 'api/pets';
import { PetType } from 'api/types';
import { Role } from 'api/user';
import Curtain from 'components/Curtain';
import ImgEditor from 'components/ImgEditor';

import styles from './Pet.module.less';
import AnalysesStep from './Steps/Analyses';
import HealthStep from './Steps/Health';
import ParamsStep from './Steps/Params';
import TreatmentsStep from './Steps/Treatments';

type Props = Pet & {
    onClose?: () => void;
    updatePets: () => void;
};

enum TileName {
    PARAMS = 'params',
    SEARCH = 'search',
    HEALTH = 'health',
    ANALYSES = 'analyses',
    DONATIONS = 'donations',
    TREATMENTS = 'treatments',
}

enum CurtainSteps {
    DONE = 'done',
    CONFIRMATIONS = 'confirmations',
}

type CurtainType = {
    isOpen: boolean;
    step?: CurtainSteps;
};

const tiles = [
    { name: TileName.PARAMS, title: 'Параметры', icon: <Params /> },
    { name: TileName.HEALTH, title: 'Здоровье', icon: <Health /> },
    { name: TileName.TREATMENTS, title: 'Обработки', icon: <Processing /> },
    { name: TileName.ANALYSES, title: 'Анализы', icon: <Analizes /> },
    { name: TileName.DONATIONS, title: 'История донаций', icon: <DonorButton /> },
    { name: TileName.SEARCH, title: 'История поисков', icon: <RecipientButton /> },
];

const dogAnalizesCount = 6;
const catAnalizesCount = 4;

const PetProfile: FC<Props> = ({
    id,
    name,
    type,
    health = {},
    gender,
    breedId,
    onClose,
    weightKg,
    analyses = {},
    photoUrls,
    birthDate,
    petStatus,
    updatePets,
    bloodGroup,
    chipNumber,
    treatments = {},
    livingCondition,
    reproductiveStatus,
}) => {
    const navigate = useNavigate();

    const [photo, setPhoto] = useState<File | null>(null);
    const [isOpenTooltip, setIsOpenTooltip] = useState(false);
    const [isEditMode, setIsEditMode] = useState<boolean>(false);
    const [curtain, setCurtain] = useState<CurtainType>({ isOpen: false });
    const [needUpdatePets, setNeedUpdatePets] = useState<boolean>(false);
    const [stepOnEditing, setStepOnEditing] = useState<TileName | null>(null);
    const [isPhotoWasDeleted, setIsPhotoWasDeleted] = useState<boolean>(false);

    // dicts
    const { data: petGendersDict = [], isError: isErrorGenders } = useGendersQuery();
    const { data: healthStatusesDict = [], isError: isErrorHealthStatuses } = useHealthStatusesQuery();
    const { data: livingConditionsDict = [], isError: isErrorLivingConditions } = useLivingConditionsQuery();
    const { data: reproductiveStatusesDict = [], isError: isErrorReproductiveStatusesDict } =
        useReproductiveStatusesQuery();
    const {
        data: { petTypesDict = [], bloodGroupDict = {}, breedsDict = {} } = {},
        isError: isErrorPetTypesAndBloodGroups,
    } = usePetTypesAndBloodGroupsQuery();

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

    const onCloseClickHandler = async () => {
        if (needUpdatePets) {
            await updatePets();
        }

        onClose?.();
    };

    const toggleEditMode = () => {
        setIsEditMode((prevState) => !prevState);
    };

    const toggleTooltip = (e: React.MouseEvent) => {
        e.stopPropagation();

        setIsOpenTooltip((prevState) => !prevState);
    };

    const updatePhoto = useCallback(
        async (newPhoto: File) => {
            const { success } = await addPhoto({ id, photo: newPhoto, isAvatar: true });

            if (success) {
                setNeedUpdatePets(true);
            } else {
                showToast('Не удалось обновить фотографию, попробуйте еще раз');
            }
        },
        [id, showToast],
    );

    const onLoadPhotoHandler = useCallback(
        async (newPhoto: File | null) => {
            if (newPhoto === null) {
                setIsPhotoWasDeleted(true);
            } else {
                setIsPhotoWasDeleted(false);

                updatePhoto(newPhoto);
            }

            setPhoto(newPhoto);
        },
        [updatePhoto],
    );

    const getIsTileDisable = (tileName: TileName) => {
        switch (tileName) {
            case TileName.HEALTH: {
                return !Object.keys(health || {}).length;
            }
            case TileName.TREATMENTS: {
                return !Object.keys(treatments || {}).length;
            }
            case TileName.ANALYSES: {
                return !Object.keys(analyses || {}).length;
            }
            default: {
                return false;
            }
        }
    };

    const onTileClickHandler = (tileName: TileName) => () => {
        setStepOnEditing(tileName);
    };

    const onTileBackHandler = () => {
        setStepOnEditing(null);
    };

    const onDeleteClickHandler = () => {
        setCurtain({
            isOpen: true,
            step: CurtainSteps.CONFIRMATIONS,
        });
    };

    const onConfirmClickHandler = async () => {
        if (curtain.step === CurtainSteps.CONFIRMATIONS) {
            const result = await deletePetById(id);

            if (!result) {
                showToast('Не удалось удалить питомца, попробуйте еще раз');

                setCurtain({
                    isOpen: false,
                });

                return;
            }

            setNeedUpdatePets(true);
            setCurtain({
                isOpen: true,
                step: CurtainSteps.DONE,
            });
        } else {
            navigate('/adding');
        }
    };

    const onCancelClickHandler = async () => {
        if (curtain.step === CurtainSteps.CONFIRMATIONS) {
            setCurtain({
                isOpen: false,
            });
        } else {
            onCloseClickHandler();
        }
    };

    const onErrorUpdateHandler = () => {
        showToast('Не удалось обновить параметры, попробуйте еще раз');
    };

    useEffect(() => {
        const onOutsideClickHandler = () => {
            setIsOpenTooltip(false);
        };

        window.addEventListener('click', onOutsideClickHandler);

        return () => {
            window.removeEventListener('click', onOutsideClickHandler);
        };
    }, []);

    useEffect(() => {
        if (isErrorPetTypesAndBloodGroups) {
            showToast('Не удалось загрузить словарь типов животных и групп крови, попробуйте перезагрузить приложение');
        }
    }, [isErrorPetTypesAndBloodGroups, showToast]);

    useEffect(() => {
        if (isErrorGenders) {
            showToast('Не удалось загрузить словарь полов, попробуйте перезагрузить приложение');
        }
    }, [isErrorGenders, showToast]);

    useEffect(() => {
        if (isErrorHealthStatuses) {
            showToast('Не удалось загрузить словарь статусов здоровья, попробуйте перезагрузить приложение');
        }
    }, [isErrorHealthStatuses, showToast]);

    useEffect(() => {
        if (isErrorLivingConditions) {
            showToast('Не удалось загрузить словарь условий проживания, попробуйте перезагрузить приложение');
        }
    }, [isErrorLivingConditions, showToast]);

    useEffect(() => {
        if (isErrorReproductiveStatusesDict) {
            showToast('Не удалось загрузить словарь репродуктивных состояний, попробуйте перезагрузить приложение');
        }
    }, [isErrorReproductiveStatusesDict, showToast]);

    if (stepOnEditing === TileName.PARAMS) {
        return (
            <ParamsStep
                id={id}
                name={name}
                type={type}
                gender={gender}
                breedId={breedId}
                weightKg={weightKg}
                birthDate={birthDate}
                isEditMode={isEditMode}
                chipNumber={chipNumber}
                bloodGroup={bloodGroup}
                petTypes={petTypesDict}
                petGenders={petGendersDict}
                onClose={onTileBackHandler}
                onSuccessUpdate={updatePets}
                livingCondition={livingCondition}
                onErrorUpdate={onErrorUpdateHandler}
                reproductiveStatus={reproductiveStatus}
                livingConditionsDict={livingConditionsDict}
                breedsDict={breedsDict as BloodAndBreedGroupsDict}
                reproductiveStatusesDict={reproductiveStatusesDict}
                bloodGroupDict={bloodGroupDict as BloodAndBreedGroupsDict}
            />
        );
    }

    if (stepOnEditing === TileName.HEALTH) {
        return (
            <HealthStep
                id={id}
                isEditMode={isEditMode}
                onClose={onTileBackHandler}
                onSuccessUpdate={updatePets}
                transfused={health?.transfused}
                medications={health?.medications}
                healthStatus={health?.healthStatus}
                onErrorUpdate={onErrorUpdateHandler}
                healthStatusesDict={healthStatusesDict}
                surgicalInterventions={health?.surgicalInterventions}
            />
        );
    }

    if (stepOnEditing === TileName.TREATMENTS) {
        return (
            <TreatmentsStep
                id={id}
                isEditMode={isEditMode}
                onClose={onTileBackHandler}
                onSuccessUpdate={updatePets}
                onErrorUpdate={onErrorUpdateHandler}
                dewormingDate={treatments.dewormingDate}
                rabiesVaccinationDate={treatments.rabiesVaccinationDate}
                infectionVaccinationDate={treatments.infectionVaccinationDate}
                ectoparasiteTreatmentDate={treatments.ectoparasiteTreatmentDate}
            />
        );
    }

    if (stepOnEditing === TileName.ANALYSES) {
        return (
            <AnalysesStep
                petId={id}
                petType={type}
                analyses={analyses}
                isEditMode={isEditMode}
                onClose={onTileBackHandler}
                onSuccessUpdate={updatePets}
                onErrorUpdate={onErrorUpdateHandler}
            />
        );
    }

    return (
        <div className={styles.wrapper}>
            <div className={styles.header}>
                <div className={styles.back} onClick={isEditMode ? toggleEditMode : onCloseClickHandler}>
                    <BackAngularArrow />
                </div>
                <h2 className={styles.title}>{isEditMode ? 'Редактирование питомца' : name.toUpperCase()}</h2>
                <div className={cn(styles.buttons, { [styles.isEditMode]: isEditMode })}>
                    {petStatus === Role.NONE && (
                        <div className={styles.button}>
                            <div onClick={toggleEditMode} className={cn(styles.icon, { [styles.edit]: true })}>
                                <Edit />
                            </div>
                        </div>
                    )}
                    <div onClick={onDeleteClickHandler} className={styles.button}>
                        <div className={cn(styles.icon, { [styles.basket]: true })}>
                            <Basket />
                        </div>
                    </div>
                </div>
            </div>
            <div className={styles.info}>
                <div className={styles.photoWrapper}>
                    <ImgEditor
                        showStub
                        petType={type}
                        isEditIcon={isEditMode}
                        bloodGroup={bloodGroup}
                        className={styles.photo}
                        name={isEditMode ? name : undefined}
                        src={isPhotoWasDeleted ? photo : null}
                        onLoad={isEditMode ? onLoadPhotoHandler : undefined}
                        serverSrc={isPhotoWasDeleted ? undefined : photoUrls?.[0]}
                    />
                </div>
                <div className={styles.labels}>
                    {petStatus === Role.DONOR && (
                        <div className={styles.donation}>
                            <div className={styles.labelIcon}>
                                <Blood />
                            </div>
                            <div className={styles.params}>
                                <p className={styles.labelInfoTitle}>Примерный объем донации</p>
                                <div className={styles.labelInfoValue}>
                                    <p className={styles.labelDescr}>
                                        {Number((weightKg * (type === PetType.DOG ? 17.6 : 13.2)).toFixed(2))} мл
                                    </p>
                                    <div onClick={toggleTooltip} className={styles.infoIcon}>
                                        <Info />
                                    </div>
                                </div>
                                {isOpenTooltip && (
                                    <div className={styles.tooltip}>
                                        До 20% объема циркулирующей крови - не более 17,6 мл/кг
                                    </div>
                                )}
                            </div>
                        </div>
                    )}
                </div>
            </div>
            <div className={styles.tiles}>
                {tiles.map(({ name: tileName, title, icon }) => (
                    <div
                        key={tileName}
                        onClick={onTileClickHandler(tileName)}
                        className={cn(styles.tile, {
                            [styles.disabled]: getIsTileDisable(tileName) && !isEditMode,
                            [styles.noActive]: tileName === TileName.DONATIONS || tileName === TileName.SEARCH,
                        })}
                    >
                        <div
                            className={cn(styles.tileIcon, {
                                [styles.needFill]: tileName === TileName.DONATIONS || tileName === TileName.SEARCH,
                            })}
                        >
                            {icon}
                        </div>
                        <p className={styles.tileTitle}>{title}</p>
                        {isEditMode ? (
                            <div className={styles.tileEdit}>
                                <Edit />
                            </div>
                        ) : (
                            <div className={styles.arrowTileIcon}>
                                <AccordionArrow />
                            </div>
                        )}
                        {tileName === TileName.ANALYSES && (
                            <div className={styles.analizesCount}>
                                {Object.keys(analyses).length} из{' '}
                                {type === PetType.DOG ? dogAnalizesCount : catAnalizesCount}
                            </div>
                        )}
                    </div>
                ))}
            </div>
            {curtain.isOpen && (
                <Curtain
                    title={
                        curtain.step === CurtainSteps.CONFIRMATIONS ? (
                            <>
                                Вы точно хотите удалить
                                <br />
                                профиль питомца?
                            </>
                        ) : (
                            <>
                                Профиль питомца
                                <br />
                                удален
                            </>
                        )
                    }
                    onCancel={onCancelClickHandler}
                    onConfirm={onConfirmClickHandler}
                    confirmButtonTitle={curtain.step === CurtainSteps.CONFIRMATIONS ? 'Удалить' : 'Добавить нового'}
                    cancelButtonTitle={curtain.step === CurtainSteps.CONFIRMATIONS ? 'Не удалять' : 'К питомцам'}
                    subTitle={curtain.step === CurtainSteps.CONFIRMATIONS ? 'Данные о нем будут потеряны' : undefined}
                />
            )}
        </div>
    );
};

export default PetProfile;
