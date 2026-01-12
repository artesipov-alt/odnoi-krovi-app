import { AxiosPromise } from 'axios';

import { instance } from './instance';

enum Role {
  USER = 'user',
  ADMIN = 'admin',
  donor = 'donor',
  CLINIC = 'clinic',
}

export type GetUserResponse = {
  id: string;
  role?: Role;
  phone?: string;
  email?: string;
  message?: string;
  fullName: string;
  allowGeo?: boolean;
  createdAt?: string;
  consentPd?: boolean;
  locationId?: number;
  telegramId?: number;
  onBoarding?: boolean;
  organizationName?: string;
};

export type UpdateUserRequest = {
  id: string;
  phone?: string;
  email?: string;
  fullName: string;
  allowGeo?: boolean;
  locationId?: number;
  onBoarding?: boolean;
};

export type UpdateUserResponse = {
  data?: string;
  message?: string;
  error?: string;
};

export interface IUserApi {
  getUser(id: number): AxiosPromise<GetUserResponse>;
  getUserByTelegramId(id: number): AxiosPromise<GetUserResponse>;
  updateUser(params: UpdateUserRequest): AxiosPromise<UpdateUserResponse>;
}

export const USER_URL = '/v1/user';

export const userApi = (): IUserApi => ({
  getUser(id) {
    return instance.get(`${USER_URL}${id}`);
  },
  getUserByTelegramId(id) {
    return instance.get(`${USER_URL}/telegram?telegram_id=${id}`);
  },
  updateUser({ id, ...params }) {
    return instance.put(`${USER_URL}/${id}`, params);
  },
});
