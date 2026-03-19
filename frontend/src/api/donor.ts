import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { PetType } from './types';
import { CompensationType } from './user';

export enum RecipientStatus {
    DRAFT = 'draft',
    ACTIVE = 'active',
    CLOSED = 'closed',
}
export type MatchingDonors = {
    petId: string;
    petName: string;
    photoUrls: string[];
    donorBloodGroup: string;
};

export type RecipientItem = {
    id: string;
    petId: string;
    petName: string;
    petType: PetType;
    photoUrls?: string[];
    bloodGroupName: string;
    status: RecipientStatus;
    prioritySearch: boolean;
    bloodVolumeRemaining: number;
    matchingDonors?: MatchingDonors[];
};

export type GetRecipientsListResponse = {
    total: number;
    items: RecipientItem[];
};

export type AdvancedInfo = {
    photoUrls: string[];
    description: string;
};

export type DefaultPrefs = {
    bonuses: string[];
    taxiCompensation: boolean;
    compensationType: CompensationType;
};

export type MatchingDonor = {
    petId: string;
    amount: number;
    petName: string;
    photoUrls: string[];
    donorBloodGroup: string;
};

export type GetRecipientDetailsResponse = {
    id: string;
    petId: string;
    petName: string;
    petType: PetType;
    regions: string[];
    ownerName: string;
    photoUrls?: string[];
    bloodGroupName: string;
    status: RecipientStatus;
    prioritySearch: boolean;
    bloodVolumeNeeded: number;
    advancedInfo?: AdvancedInfo;
    defaultPrefs?: DefaultPrefs;
    bloodVolumeReserved: number;
    searchingBloodNames: string[];
    matchingDonors?: MatchingDonor[];
};

export type BloodSearchApplyRequest = {
    id: string;
    donorId: string;
    taxiCompensation: boolean;
    compensationType: CompensationType;
};

export type BloodSearchApplyResponse = {
    id: string;
    donorId: string;
    createdAt: string;
    requestId: string;
    status: RecipientStatus;
};

export interface IDonorApi {
    getRecipientsList(id: string, status?: RecipientStatus): AxiosPromise<GetRecipientsListResponse>;
    getRecipientDetails(id: string): AxiosPromise<GetRecipientDetailsResponse>;
    bloodSearchApply(params: BloodSearchApplyRequest): AxiosPromise<BloodSearchApplyResponse>;
}

export const DONOR_URL = '/v1/donor';

export const donorApi = (): IDonorApi => ({
    getRecipientsList(id, status) {
        return instance.get(`${DONOR_URL}/recipient-list/${id}?status=${status || RecipientStatus.ACTIVE}`);
    },
    getRecipientDetails(id) {
        return instance.get(`${DONOR_URL}/recipient-details/${id}`);
    },
    bloodSearchApply({ id, ...params }) {
        return instance.post(`${DONOR_URL}/recipient/${id}/apply`, params);
    },
});
