import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { PetType } from './types';

enum Sex {
    MALE = 'male',
    FEMALE = 'female',
}

export type Pet = {
    id: string;
    gender: Sex;
    name: string;
    type: PetType;
    ownerId: string;
    bloodGroup: string; // как будет приходить?
    breed: string; // как будет приходить?
    photoUrl: string;
    ageYears?: number;
    hasChip?: boolean;
    latitude?: number;
    weightKg?: number;
    ageMonths?: number;
    longitude?: number;
    chipNumber?: number;
    isGuideDog?: boolean;
    sterilized?: boolean;
    isTherapist?: boolean;
    dewormingDate?: string;
    livingCondition?: string;
    vaccinationDate?: string;
    knowsBloodGroup?: boolean;
    ectoparasiteDate?: string;
    lastTransfusionDate?: string;
};

export type GetPetsResponse = Pet[];

export interface IPetsApi {
    getPets(id: string): AxiosPromise<GetPetsResponse>;
}

export const PETS_URL = 'pets/';

export const petsApi = (): IPetsApi => ({
    getPets(id) {
        return instance.get(`${PETS_URL}user/${id}`);
    },
});
