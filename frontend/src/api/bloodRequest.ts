import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { Analyses, DonorRestrictions, Health, Pet, Treatments, WarnFactors } from './pets';
import { PetGender, PetType } from './types';
import { CompensationType, Role } from './user';

export enum PoolRequestStatus {
    DRAFT = 'draft',
    CLOSED = 'closed',
    ACTIVE = 'active',
    RESERVED_FULL = 'reserved_full',
}

export enum Onboardings {
    SEARCH = 'SEARCH',
    BLOOD_CARD = 'BLOOD_CARD',
}

export type AddToPoolRequest = {
    petId: string;
    regions: number[];
    description?: string;
    prioritySearch: boolean;
    bloodGroupNames: string[];
    bloodVolumeNeeded: number;
    bloodComponentIds: number[];
    smallPetsNotifyAllowed: boolean;
    includeUnknownBloodGroup: boolean;
};

export type AddToPoolResponse = {
    id: string;
    petId: string;
    status: string;
};

export enum RespondingDonorStatus {
    FAILED = 'failed',
    PENDING = 'pending',
    ACCEPTED = 'accepted',
    REJECTED = 'rejected',
    CANCELED = 'cancelled',
    COMPLETED = 'completed',
}

export type RespondingDonor = {
    id: string;
    amount: number;
    donorId: string;
    requestId: string;
    donorName: string;
    createdAt: string;
    updatedAt: string;
    isConfirmed: boolean;
    donorPhotos: string[];
    donorBloodGroup: string;
    taxiCompensation: boolean;
    warnFactors: WarnFactors[];
    status: RespondingDonorStatus;
    compensationType: CompensationType;
};

export type GetPoolRequestResponse = {
    id: string;
    petId: string;
    regions: string[];
    createdAt?: string;
    photoUrls?: string[];
    description?: string;
    suitableDonors: number;
    prioritySearch?: boolean;
    bloodGroupNames: string[];
    bloodVolumeNeeded: number;
    bloodVolumeDonated: number;
    status?: PoolRequestStatus;
    onBoarding?: Onboardings[];
    bloodComponentIds: string[];
    responses?: RespondingDonor[];
    bloodVolumeReserved?: number;
    smallPetsNotifyAllowed: true;
    acceptedDonors?: RespondingDonor[];
    includeUnknownBloodGroup?: boolean;
    completedDonations?: RespondingDonor[];
};

export type UpdatePoolRequestRequest = {
    id: string;
    onBoarding: GetPoolRequestResponse['onBoarding'];
};

export type UpdatePoolRequestResponse = {
    id: string;
    updatedAt: string;
};

export type GetDonorInfoResponse = Pet & {
    taxi: boolean;
    ownerId: string;
    ownerName: string;
    ownerPhone: string;
    availableBloodAmount: number;
    compensationType: CompensationType;
};

export type ApplyDonorRespondResponse = {
    message: string;
};

export type Application = {
    id: string;
    amount: number;
    taxiCompensation: boolean;
    status: RespondingDonorStatus;
    compensationType: CompensationType;
};

export type DonationForRecipientDonorData = {
    id: string;
    name: string;
    type: PetType;
    health: Health;
    breedId: string;
    ownerId: string;
    petStatus: Role;
    birthDate: Date;
    weightKg: number;
    bonuses: string[];
    createdAt: string;
    deletedAt: string;
    gender: PetGender;
    ownerName: string;
    updatedAt: string;
    bloodGroup: string;
    chipNumber: string;
    photoUrls: string[];
    phoneNumber: string;
    analyses?: Analyses;
    treatments: Treatments;
    livingCondition: string;
    application: Application;
    reproductiveStatus?: string;
    availableBloodAmount: number;
    donorRestrictions?: DonorRestrictions;
};

export type DonationForRecipientRecipientData = {
    petName: string;
    petType: PetType;
    photoUrls: string[];
    bloodVolumeNeeded: number;
    bloodVolumeDonated: number;
    bloodVolumeReserved: number;
};

export type GetDonationForRecipientByIdResponse = {
    donorData: DonationForRecipientDonorData;
    recipientData: DonationForRecipientRecipientData;
};

export type ConfirmDonationRequest = {
    id: string;
    amount: number;
};

export interface IBloodRequestApi {
    addToPool(params: AddToPoolRequest): AxiosPromise<AddToPoolResponse>;
    getPoolRequest(id: string): AxiosPromise<GetPoolRequestResponse>;
    updatePoolRequest(params: UpdatePoolRequestRequest): AxiosPromise<UpdatePoolRequestResponse>;
    getDonorInfo(id: string): AxiosPromise<GetDonorInfoResponse>;
    applyDonorRespond(id: string): AxiosPromise<ApplyDonorRespondResponse>;
    getDonationForRecipientById(id: string): AxiosPromise<GetDonationForRecipientByIdResponse>;
    confirmDonation(params: ConfirmDonationRequest): AxiosPromise<void>;
    rejectDonation(id: string): AxiosPromise<void>;
    closeSearch(id: string): AxiosPromise<void>;
}

export const BLOOD_REQUEST_URL = '/v1/blood-request';

export const bloodRequestApi = (): IBloodRequestApi => ({
    addToPool(params) {
        return instance.post(`${BLOOD_REQUEST_URL}/pool`, params);
    },
    getPoolRequest(id) {
        return instance.get(`${BLOOD_REQUEST_URL}/pet/${id}`);
    },
    updatePoolRequest({ id, ...params }) {
        return instance.patch(`${BLOOD_REQUEST_URL}/${id}`, params);
    },
    getDonorInfo(id) {
        return instance.get(`${BLOOD_REQUEST_URL}/donor/${id}`);
    },
    applyDonorRespond(id) {
        return instance.post(`${BLOOD_REQUEST_URL}/apply-response/${id}`);
    },
    getDonationForRecipientById(id) {
        return instance.get(`${BLOOD_REQUEST_URL}/donation/${id}`);
    },
    confirmDonation({ id, ...params }) {
        return instance.post(`${BLOOD_REQUEST_URL}/donation/${id}/confirm`, params);
    },
    rejectDonation(id) {
        return instance.post(`${BLOOD_REQUEST_URL}/donation/${id}/reject`);
    },
    closeSearch(id) {
        return instance.post(`${BLOOD_REQUEST_URL}/close/${id}`);
    },
});
