import { AxiosPromise } from 'axios';

import { instance } from './instance';

export type AddToPoolRequest = {
    petId: string;
    regions: number[];
    description?: string;
    bloodGroupNames: string[];
    bloodVolumeNeeded: number;
    bloodComponentIds: number[];
    smallPetsNotifyAllowed: boolean;
};

export type AddToPoolResponse = {
    id: string;
    petId: string;
    status: string;
};

export interface IBloodRequestApi {
    addToPool(params: AddToPoolRequest): AxiosPromise<AddToPoolResponse>;
}

export const BLOOD_REQUEST_URL = '/v1/blood-request';

export const bloodRequestApi = (): IBloodRequestApi => ({
    addToPool(params) {
        return instance.post(`${BLOOD_REQUEST_URL}/pool`, params);
    },
});
