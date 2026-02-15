import { AxiosPromise } from 'axios';

import { instance } from './instance';

export type GetPhotoLinkRequest = {
    id: string;
    photos_count: number;
    for_pet_avatar?: boolean;
    for_user_avatar?: boolean;
    for_blood_req?: boolean;
};

export type UploadItem = {
    url: 'string';
    path: 'string';
};

export type GetPhotoLinkResponse = {
    items: UploadItem[];
};

export type ConfirmUploadPhotoRequest = {
    entityId: string;
    paths: string[];
};

export interface IPhotoApi {
    getPhotoLink(params: GetPhotoLinkRequest): AxiosPromise<GetPhotoLinkResponse>;
    confirmUploadPhoto(params: ConfirmUploadPhotoRequest): AxiosPromise<void>;
}

export const PHOTO_URL = '/v1';

export const photoApi = (): IPhotoApi => ({
    getPhotoLink({ id, photos_count, for_pet_avatar, for_user_avatar, for_blood_req }) {
        let params = '';

        if (for_pet_avatar) {
            params = `&for_pet_avatar=${for_pet_avatar}`;
        } else if (for_user_avatar) {
            params = `&for_user_avatar=${for_user_avatar}`;
        } else if (for_blood_req) {
            params = `&for_blood_req=${for_blood_req}`;
        }

        return instance.post(`${PHOTO_URL}/uploads/presign/${id}?photos_count=${photos_count}${params}`);
    },
    confirmUploadPhoto(body) {
        return instance.post(`${PHOTO_URL}/uploads/confirm`, body);
    },
});
