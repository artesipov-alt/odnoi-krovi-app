import { useQuery } from '@tanstack/react-query';

import { getBloodComponents } from 'api/apiServices/getBloodComponents';
import { getBloodGroups } from 'api/apiServices/getBloodGroups';
import { getBreedsByType } from 'api/apiServices/getBreedsByType';
import { getGenders } from 'api/apiServices/getGenders';
import { getHealthStatuses } from 'api/apiServices/getHealthStatuses';
import { getLivingConditions } from 'api/apiServices/getLivingConditions';
import { getLocations } from 'api/apiServices/getLocations';
import { getPetsTypes } from 'api/apiServices/getPetsTypes';
import { getReproductiveStatuses } from 'api/apiServices/getReproductiveStatuses';
import { Dict, PetTypeDict } from 'api/reference';
import { PetType } from 'api/types';

export type BloodAndBreedGroupsDict = Record<PetType, Dict[]>;

type PetTypesAndBloodGroupsData = {
    petTypesDict: PetTypeDict[];
    breedsDict: BloodAndBreedGroupsDict;
    bloodGroupDict: BloodAndBreedGroupsDict;
};

const queryParams = {
    staleTime: 15 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    retryDelay: 1000,
    retry: 2,
};

export const useBloodComponentsQuery = () =>
    useQuery({
        ...queryParams,
        queryKey: ['bloodComponents'],
        queryFn: getBloodComponents,
    });

export const useLocationsQuery = () =>
    useQuery({
        ...queryParams,
        queryKey: ['locations'],
        queryFn: getLocations,
    });

export const useGendersQuery = () =>
    useQuery({
        ...queryParams,
        queryKey: ['genders'],
        queryFn: getGenders,
    });

export const useHealthStatusesQuery = () =>
    useQuery({
        ...queryParams,
        queryKey: ['healthStatuses'],
        queryFn: getHealthStatuses,
    });

export const useLivingConditionsQuery = () =>
    useQuery({
        ...queryParams,
        queryKey: ['livingConditions'],
        queryFn: getLivingConditions,
    });

export const useReproductiveStatusesQuery = () =>
    useQuery({
        ...queryParams,
        queryKey: ['reproductiveStatuses'],
        queryFn: getReproductiveStatuses,
    });

export const usePetTypesAndBloodGroupsQuery = () =>
    useQuery<PetTypesAndBloodGroupsData>({
        ...queryParams,
        queryKey: ['petTypesAndBloodGroups'],
        queryFn: async (): Promise<PetTypesAndBloodGroupsData> => {
            // Шаг 1: Получаем типы животных
            const petTypesDict = await getPetsTypes();

            if (!petTypesDict) {
                throw new Error('Не удалось загрузить типы животных');
            }

            // Шаг 2: Параллельно загружаем группы крови и породы для каждого типа
            const bloodGroupDict: BloodAndBreedGroupsDict = {} as BloodAndBreedGroupsDict;
            const breedsDict: BloodAndBreedGroupsDict = {} as BloodAndBreedGroupsDict;

            const results = await Promise.allSettled(
                petTypesDict.map(async ({ value }) => {
                    const [bloodGroups, breeds] = await Promise.all([
                        getBloodGroups(value).catch(() => null),
                        getBreedsByType(value).catch(() => null),
                    ]);

                    if (bloodGroups) {
                        bloodGroupDict[value] = bloodGroups;
                    }

                    if (breeds) {
                        breedsDict[value] = breeds;
                    }

                    return { type: value, bloodGroups, breeds };
                }),
            );

            // Если все запросы провалились — считаем ошибкой
            if (results.every((result) => result.status === 'rejected')) {
                throw new Error('Не удалось загрузить данные о группах крови или породах');
            }

            return { petTypesDict, bloodGroupDict, breedsDict };
        },
    });
