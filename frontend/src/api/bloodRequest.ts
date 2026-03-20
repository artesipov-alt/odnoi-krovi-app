import { AxiosPromise } from 'axios';

import { instance } from './instance';
import { CompensationType } from './user';

enum PoolRequest {
    DRAFT = 'draft',
    ACTIVE = 'active',
    CLOSED = 'closed',
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

export type RespondingDonor = {
    id: string;
    amount: number;
    status: string;
    donorId: string;
    requestId: string;
    donorName: string;
    createdAt: string;
    updatedAt: string;
    warnFactors: string[];
    donorPhotos: string[];
    donorBloodGroup: string;
    taxiCompensation: boolean;
    compensationType: CompensationType;
};

export type GetPoolRequestResponse = {
    id: string;
    petId: string;
    regions: string[];
    createdAt?: string;
    status?: PoolRequest;
    photoUrls?: string[];
    description?: string;
    suitableDonors: number;
    bloodGroupNames: string[];
    bloodVolumeNeeded: number;
    onBoarding?: Onboardings[];
    bloodComponentIds: string[];
    responses?: RespondingDonor[];
    bloodVolumeReserved?: number;
    smallPetsNotifyAllowed: true;
};

export type UpdatePoolRequestRequest = {
    id: string;
    onBoarding: GetPoolRequestResponse['onBoarding'];
};

export type UpdatePoolRequestResponse = {
    id: string;
    updatedAt: string;
};

export interface IBloodRequestApi {
    addToPool(params: AddToPoolRequest): AxiosPromise<AddToPoolResponse>;
    getPoolRequest(id: string): AxiosPromise<GetPoolRequestResponse>;
    updatePoolRequest(params: UpdatePoolRequestRequest): AxiosPromise<UpdatePoolRequestResponse>;
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
});
