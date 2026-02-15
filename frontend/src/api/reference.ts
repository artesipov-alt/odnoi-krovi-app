import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { PetGender, PetType } from './types';

export type PetTypeDict = {
    value: PetType;
    label: string;
};

export type PetGenderDict = {
    value: PetGender;
    label: string;
};

export type Dict = {
    value: string;
    label: string;
};

export type StringDict = Dict;

export type GetPetsTypesResponse = {
    data: PetTypeDict[];
};

export type GetBloodGroupsResponse = {
    data: Dict[];
};

export type GetBloodComponentsResponse = {
    data: Dict[];
};

export type GetLocationsResponse = {
    data: Dict[];
};

export type GetGendersResponse = {
    data: PetGenderDict[];
};

export type GetBreedsByTypeResponse = {
    data: Dict[];
};

export type GetLivingConditionsResponse = {
    data: Dict[];
};

export type GetHealthStatusesResponse = {
    data: Dict[];
};

export type GetReproductiveStatusesResponse = {
    data: Dict[];
};

export interface IReferenceApi {
    getPetsTypes(): AxiosPromise<GetPetsTypesResponse>;
    getBloodGroups(pet: PetType): AxiosPromise<GetBloodGroupsResponse>;
    getBloodComponents(): AxiosPromise<GetBloodComponentsResponse>;
    getLocations(): AxiosPromise<GetLocationsResponse>;
    getGenders(): AxiosPromise<GetGendersResponse>;
    getBreedsByType(pet: PetType): AxiosPromise<GetBreedsByTypeResponse>;
    getLivingConditions(): AxiosPromise<GetLivingConditionsResponse>;
    getHealthStatuses(): AxiosPromise<GetHealthStatusesResponse>;
    getReproductiveStatuses(): AxiosPromise<GetReproductiveStatusesResponse>;
}

export const REFERENCE_URL = '/v1/reference';

export const referenceApi = (): IReferenceApi => ({
    getPetsTypes() {
        return instance.get(`${REFERENCE_URL}/pet-types`);
    },
    getBloodGroups(pet) {
        return instance.get(`${REFERENCE_URL}/blood-groups/${pet}`);
    },
    getBloodComponents() {
        return instance.get(`${REFERENCE_URL}/blood-components`);
    },
    getLocations() {
        return instance.get(`${REFERENCE_URL}/locations`);
    },
    getGenders() {
        return instance.get(`${REFERENCE_URL}/genders`);
    },
    getBreedsByType(petType) {
        return instance.get(`${REFERENCE_URL}/breeds-by-type`, { params: { petType } });
    },
    getLivingConditions() {
        return instance.get(`${REFERENCE_URL}/living-conditions`);
    },
    getHealthStatuses() {
        return instance.get(`${REFERENCE_URL}/health-statuses`);
    },
    getReproductiveStatuses() {
        return instance.get(`${REFERENCE_URL}/reproductive-statuses`);
    },
});
