import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { PetType } from './types';

export type PetDict = {
    value: PetType;
    label: string;
};

export type Dict = {
    value: number;
    label: string;
};

export type GetPetsTypesResponse = {
    data: PetDict[];
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

export interface IReferenceApi {
    getPetsTypes(): AxiosPromise<GetPetsTypesResponse>;
    getBloodGroups(pet: PetType): AxiosPromise<GetBloodGroupsResponse>;
    getBloodComponents(): AxiosPromise<GetBloodComponentsResponse>;
    getLocations(): AxiosPromise<GetLocationsResponse>;
}

export const REFERENCE_URL = 'reference/';

export const referenceApi = (): IReferenceApi => ({
    getPetsTypes() {
        return instance.get(`${REFERENCE_URL}pet-types/`);
    },
    getBloodGroups(pet) {
        return instance.get(`${REFERENCE_URL}blood-groups/${pet}`);
    },
    getBloodComponents() {
        return instance.get(`${REFERENCE_URL}blood-components/`);
    },
    getLocations() {
        return instance.get(`${REFERENCE_URL}locations/`);
    },
});
