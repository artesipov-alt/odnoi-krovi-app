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

export type Health = {
    healthStatus: string;
    medications?: string;
    transfused?: boolean;
    lastDonation?: Date;
    surgicalInterventions?: string;
};

export type Treatments = {
    dewormingDate?: Date;
    rabiesVaccinationDate?: Date;
    infectionVaccinationDate?: Date;
    ectoparasiteTreatmentDate?: Date;
};

export type AnalysesItem = {
    analysisDate: Date;
    analysisName: string;
    analysisType: string;
};

export type Analyses = {
    leukemia?: AnalysesItem[];
    babesiosis?: AnalysesItem[];
    dirofilaria?: AnalysesItem[];
    anaplasmosis?: AnalysesItem[];
    ehrlichiosis?: AnalysesItem[];
    bartonellosis?: AnalysesItem[];
    hemoplasmosis?: AnalysesItem[];
    immunodeficiency?: AnalysesItem[];
};

export type StopFactors = {
    code: string;
    description: string;
};

export type WarnFactors = {
    code: string;
    description: string;
    subDescription: string;
};

export type DonorRestrictions = {
    stopFactors?: StopFactors[];
    warnFactors?: WarnFactors[];
};

export type Pet = {
    id: string;
    name: string;
    type: PetType;
    health?: Health;
    petStatus: Role;
    breedId?: string;
    weightKg: number;
    birthDate?: Date;
    bonuses?: Bonuses;
    ageYears?: number;
    gender?: PetGender;
    ageMonths?: number;
    bloodGroup: string;
    analyses?: Analyses;
    chipNumber?: string;
    photoUrls?: string[];
    recoveryDays?: number;
    treatments?: Treatments;
    livingCondition?: string;
    reproductiveStatus?: string;
    donorRestrictions?: DonorRestrictions;
};

export type GetPetsResponse = {
    pets: Pet[];
    totalPets: number;
    totalPlannedDonations?: number;
    totalCompletedDonations?: number;
};

export type CreatePetRequest = Omit<Pet, 'id'> & { userId: string };

export type CreatePetResponse = Pet;

export interface IPetsApi {
    getPets(id: string): AxiosPromise<GetPetsResponse>;
    createPet(data: CreatePetRequest): AxiosPromise<CreatePetResponse>;
    deletePetById(id: string): AxiosPromise<void>;
    updatePet(data: Pet): AxiosPromise<Pet>;
}

export const PETS_URL = '/v1/pet';

export const petsApi = (): IPetsApi => ({
    getPets(id) {
        return instance.get(`${PETS_URL}/user/${id}?with_all=true`);
    },
    createPet({ userId, ...params }) {
        return instance.post(`${PETS_URL}/user/${userId}`, params);
    },
    deletePetById(id) {
        return instance.delete(`${PETS_URL}/${id}`);
    },
    updatePet({ id, ...params }) {
        return instance.put(`${PETS_URL}/${id}`, params);
    },
});
