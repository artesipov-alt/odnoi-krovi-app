import { AxiosPromise } from 'axios';

import { PoolRequestStatus } from './bloodRequest';
import { instance } from './instance';
import { PetType } from './types';
import { CompensationType } from './user';

export enum RecipientStatus {
    DRAFT = 'draft',
    ACTIVE = 'active',
    CLOSED = 'closed',
}

export enum DonorStatus {
    FAILED = 'failed',
    PENDING = 'pending',
    ACCEPTED = 'accepted',
    REJECTED = 'rejected',
    CANCELLED = 'cancelled',
    COMPLETED = 'completed',
}

export enum BonusType {
    LOCK = 'lock',
    FOOD = 'food',
    OTHER = 'other',
    PREPARATION = 'preparation',
}

export enum PlatformName {
    WB = 'wb',
    OZON = 'озон',
    FOUR_PAWS = 'четыре лапы',
    TAILY_PLATFORM = 'taily platform',
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
    prioritySearch?: boolean;
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

export type Bonus = {
    type: BonusType;
    partner: string;
    description: string;
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
    prioritySearch?: boolean;
    bloodVolumeNeeded: number;
    availableBonuses?: Bonus[];
    advancedInfo?: AdvancedInfo;
    defaultPrefs?: DefaultPrefs;
    bloodVolumeReserved: number;
    searchingBloodNames: string[];
    matchingDonors?: MatchingDonor[];
    includeUnknownBloodGroup?: boolean;
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

export type PlannedDonationDonorInfo = {
    // id донации
    id: string;
    amount: number;
    petName: string;
    createdAt: string;
    bonuses?: Bonus[];
    updatedAt: string;
    donorPetID: string;
    photoUrls: string[];
    status: DonorStatus;
    isConfirmed: boolean;
    rejectedReason: string;
    donorBloodGroup: string;
    taxiCompensation: boolean;
    compensationType: CompensationType;
};

export type PlannedDonationRecipientInfo = {
    // id питомца
    id: string;
    ownerID: string;
    petName: string;
    petType: PetType;
    regions: string[];
    updatedAt: string;
    createdAt: string;
    deletedAt: string;
    ownerName: string;
    bloodGroup: string;
    photoUrls?: string[];
    status: PoolRequestStatus;
    bloodVolumeNeeded: number;
    bloodVolumeDonated: number;
    advancedInfo: AdvancedInfo;
    bloodVolumeReserved: number;
    searchingBloodNames: string[];
    includeUnknownBloodGroup?: boolean;
};

export type PlannedDonation = {
    applicationData: PlannedDonationDonorInfo;
    recipientData: PlannedDonationRecipientInfo;
};

export type GetPlannedDonationsResponse = {
    items: PlannedDonation[];
};

export type CompleteDonationRequest = {
    id: string;
    amount: number;
};

export type GetCompletedDonationsResponse = {
    total: number;
    items: PlannedDonation[];
    totalDonatedVolume: number;
    totalCompletedDonations: number;
};

export type BonusCategory = {
    id: string;
    stage: string;
    target: string;
    userId: string;
    expiresAt: string;
    deletedAt: string;
    updatedAt: string;
    createdAt: string;
    recipient: string;
    promoCode: string;
    assignedAt: string;
    platformUrl: string;
    category: BonusType;
    description: string;
    partnerName: string;
    subcategory: string;
    platformName: string;
};

export type GetAllBonusesResponse = {
    priority: number;
    [BonusType.FOOD]: BonusCategory[];
    [BonusType.OTHER]: BonusCategory[];
    [BonusType.PREPARATION]: BonusCategory[];
};

export interface IDonorApi {
    getRecipientsList(id: string, status?: RecipientStatus): AxiosPromise<GetRecipientsListResponse>;
    getRecipientDetails(id: string): AxiosPromise<GetRecipientDetailsResponse>;
    bloodSearchApply(params: BloodSearchApplyRequest): AxiosPromise<BloodSearchApplyResponse>;
    getPlannedDonations(id: string): AxiosPromise<GetPlannedDonationsResponse>;
    cancelDonation(id: string): AxiosPromise<void>;
    completeDonation(params: CompleteDonationRequest): AxiosPromise<void>;
    getCompletedDonations(id: string): AxiosPromise<GetCompletedDonationsResponse>;
    getAllBonuses(id: string): AxiosPromise<GetAllBonusesResponse>;
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
    getPlannedDonations(id) {
        return instance.get(`${DONOR_URL}/planned-donations/${id}`);
    },
    cancelDonation(id) {
        return instance.post(`${DONOR_URL}/donation/${id}/cancel`);
    },
    completeDonation({ id, ...params }) {
        return instance.post(`${DONOR_URL}/donation/${id}/complete`, params);
    },
    getCompletedDonations(id) {
        return instance.get(`${DONOR_URL}/completed-donations/${id}`);
    },
    getAllBonuses(id) {
        return instance.get(`${DONOR_URL}/bonuses/${id}`);
    },
});
