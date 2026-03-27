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
    taxiCompensation: boolean;
    recoveryPeriodMonths: number;
    preferredLocationIds: string[];
    compensationType: CompensationType;
    notificationFrequency: NotificationFrequency;
};

export type Identities = {
    refUrl: string;
    providerId: number;
    providerName: string;
};

export type GetUserResponse = {
    id: string;
    role?: Role;
    phone?: string;
    email?: string;
    photoUrls?: string[];
    message?: string;
    fullName: string;
    allowGeo?: boolean;
    createdAt?: string;
    consentPd?: boolean;
    locationId?: number;
    telegramId?: number;
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
    >;
};

export type UpdateUserResponse = {
    data?: string;
    message?: string;
    error?: string;
};

export type GetUserIdentitiesResponse = GetUserResponse & {
    identities: Identities[];
};

export interface IUserApi {
    getUser(id: string): AxiosPromise<GetUserResponse>;
    getUserByTelegramId(id: number): AxiosPromise<GetUserResponse>;
    updateUser(params: UpdateUserRequest): AxiosPromise<UpdateUserResponse>;
    getUserIdentities(id: string): AxiosPromise<GetUserIdentitiesResponse>;
}

export const USER_URL = '/v1/user';

export const userApi = (): IUserApi => ({
    getUser(id) {
        return instance.get(`${USER_URL}/${id}?with_donor_preference=true`);
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
});
