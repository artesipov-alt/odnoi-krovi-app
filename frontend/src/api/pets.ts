import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { PetGender, PetType } from './types';
import { Role } from './user';

type Bonuses = {
    isArtist: boolean;
    isGuideDog: boolean;
    isTherapist: boolean;
    isFormerDonor: boolean;
};

type Health = {
    healthStatus: string;
    medications?: string;
    transfused?: boolean;
    lastDonation?: string;
    reproductiveStatus?: string;
    surgicalInterventions?: string;
};

type Treatments = {
    dewormingDate?: string;
    rabiesVaccinationDate?: string;
    ectoparasiteTreatmentDate?: string;
    infectionVaccinationDate?: string;
};

export type Pet = {
    id: string;
    name: string;
    type: PetType;
    analyses?: any; // TODO добавить тип
    health?: Health;
    petStatus: Role;
    breedId?: number;
    weightKg: number;
    bonuses?: Bonuses;
    photoUrl?: string;
    ageYears?: number;
    gender?: PetGender;
    ageMonths?: number;
    birthDate?: string;
    bloodGroup: string;
    chipNumber?: string;
    treatments?: Treatments;
    livingCondition?: string;
};

export type GetPetsResponse = Pet[];

export type CreatePetRequest = {
    name: string;
    type: string;
    userId: string;
    petStatus: Role;
    weightKg?: number;
    bloodGroup: string;
};

export type CreatePetResponse = Pet;

export type GetPhotoLinkResponse = {
    url: 'string';
    path: 'string';
};

export type ConfirmUploadPhotoResponse = {
    publicUrl: string;
};

export interface IPetsApi {
    getPets(id: string): AxiosPromise<GetPetsResponse>;
    createPet(data: CreatePetRequest): AxiosPromise<CreatePetResponse>;
    deletePetById(id: string): AxiosPromise<void>;
    getPhotoLink(id: string): AxiosPromise<GetPhotoLinkResponse>;
    confirmUploadPhoto(path: string): AxiosPromise<ConfirmUploadPhotoResponse>;
}

export const PETS_URL = '/v1/pet';

export const petsApi = (): IPetsApi => ({
    getPets(id) {
        return instance.get(`${PETS_URL}/user/${id}`);
    },
    createPet({ userId, ...params }) {
        return instance.post(`${PETS_URL}/user/${userId}`, params);
    },
    deletePetById(id) {
        return instance.delete(`${PETS_URL}/user/${id}`);
    },
    getPhotoLink(id) {
        return instance.get(`${PETS_URL}/upload/avatar/${id}`);
    },
    confirmUploadPhoto(path) {
        return instance.post(`${PETS_URL}/upload/avatar/confirm/${path}`);
    },
});
