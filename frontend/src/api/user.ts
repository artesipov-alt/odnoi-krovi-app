import { AxiosPromise } from 'axios';

import { instance } from './instance';

export enum Role {
    USER = 'user',
    NONE = 'none',
    ADMIN = 'admin',
    DONOR = 'donor',
    CLINIC = 'clinic',
    RECIPIENT = 'recipient',
    RECOVERING = 'recovering',
    BLOOD_FOUND = 'blood_found',
    PLANNED_DONATION = 'planned_donation',
}

export enum Onboarding {
    START = 'START',
    FIND_BLOOD = 'FIND_BLOOD',
    RECIPIENT_LIST = 'RECIPIENT_LIST',
}

export enum CompensationType {
    FREE = 'free',
    PAID = 'paid',
    FOOD = 'food',
}

export enum NotificationFrequency {
    NEVER = 'never',
    DAILY = 'daily',
    WEEKLY = 'weekly',
    IMMEDIATELY = 'immediately',
}

export type DonorPreference = {
    id: string;
    userId: string;
    createdAt: string;
    deletedAt?: string;
    updatedAt?: string;
    openForContact: boolean;
    taxiCompensation: boolean;
    recoveryPeriodMonths: number;
    preferredLocationIds: string[];
    compensationType: CompensationType;
    notificationFrequency: NotificationFrequency;
};

export type Identities = {
    refUrl?: string;
    providerId: number | string;
    providerName: string;
};

export type GetUserResponse = {
    id: string;
    role?: Role;
    phone?: string;
    email?: string;
    message?: string;
    fullName: string;
    verified: boolean;
    allowGeo?: boolean;
    createdAt?: string;
    consentPd?: boolean;
    locationId?: number;
    telegramId?: number;
    photoUrls?: string[];
    identities: Identities[];
    onBoarding?: Onboarding[];
    organizationName?: string;
    donorPreference?: DonorPreference;
};

export type UpdateUserRequest = {
    id: string;
    phone?: string;
    email?: string;
    fullName?: string;
    allowGeo?: boolean;
    locationId?: number;
    onBoarding?: Onboarding[];
    donorPreference?: Pick<
        DonorPreference,
        | 'compensationType'
        | 'notificationFrequency'
        | 'preferredLocationIds'
        | 'recoveryPeriodMonths'
        | 'taxiCompensation'
        | 'openForContact'
    >;
};

export type UpdateUserResponse = {
    data?: string;
    message?: string;
    error?: string;
};

export type GetUserIdentitiesResponse = GetUserResponse;

export type GetUserContactsRequest = {
    id: string;
    provider: string;
};

export type DefaultResponse = {
    message: string;
};

export type UpdatePhoneRequest = {
    id: string;
    phone: string;
};

export type VerifyPhoneRequest = {
    id: string;
    code: string;
};

export interface IUserApi {
    getUser(id: string): AxiosPromise<GetUserResponse>;
    getUserByTelegramId(id: number): AxiosPromise<GetUserResponse>;
    updateUser(params: UpdateUserRequest): AxiosPromise<UpdateUserResponse>;
    getUserIdentities(id: string): AxiosPromise<GetUserIdentitiesResponse>;
    getUserContacts(params: GetUserContactsRequest): AxiosPromise<DefaultResponse>;
    updatePhone(params: UpdatePhoneRequest): AxiosPromise<DefaultResponse>;
    verifyPhone(params: VerifyPhoneRequest): AxiosPromise<DefaultResponse>;
}

export const USER_URL = '/v1/user';

export const userApi = (): IUserApi => ({
    getUser(id) {
        return instance.get(`${USER_URL}/${id}?with_donor_preference=true&with_identities=true`);
    },
    getUserByTelegramId(id) {
        return instance.get(`${USER_URL}/telegram/${id}`);
    },
    updateUser({ id, ...params }) {
        return instance.put(`${USER_URL}/${id}`, params);
    },
    getUserIdentities(id) {
        return instance.get(`${USER_URL}/${id}?with_identities=true`);
    },
    getUserContacts({ id, provider }) {
        return instance.get(`${USER_URL}/${id}/contact?provider=${provider}`);
    },
    updatePhone({ id, ...params }) {
        return instance.post(`${USER_URL}/${id}/phone`, params);
    },
    verifyPhone({ id, ...params }) {
        return instance.post(`${USER_URL}/${id}/phone/verify`, params);
    },
});
