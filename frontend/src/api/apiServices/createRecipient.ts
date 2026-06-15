import { AddToPoolRequest } from '../bloodRequest';
import api from '../index';
import { CreatePetRequest } from '../pets';
import { addPhoto } from './addPhoto';

type Args = CreatePetRequest & {
    petId?: string;
    photo: File | null;
    photoBuffer?: ArrayBuffer;
    bloodRequestPhoto: File | null;
    bloodRequestPhotoBuffer?: ArrayBuffer;
    poolInfo: Omit<AddToPoolRequest, 'petId'>;
};

export const createRecipient = async ({
    petId,
    photo,
    photoBuffer,
    poolInfo,
    bloodRequestPhoto,
    bloodRequestPhotoBuffer,
    ...params
}: Args) => {
    let newPetId;
    let poolRequestId;

    if (!petId) {
        try {
            const { data } = await api.createPet(params);

            newPetId = data.id;
        } catch (e) {
            return { success: false, error: 'Не удалось создать профиль питомца' };
        }
    }

    try {
        const { data } = await api.addToPool({ ...poolInfo, petId: petId || newPetId });

        poolRequestId = data.id;
    } catch (e) {
        return { success: false, error: 'Профиль питомца создан, но не удалось создать заявку поиск крови' };
    }

    if (!petId) {
        if (photo) {
            try {
                const { success } = await addPhoto({
                    photo,
                    buffer: photoBuffer,
                    id: petId || newPetId,
                    isAvatar: true,
                });

                if (!success) {
                    return { success: false, error: 'Не удалось сохранить аватарку питомца' };
                }
            } catch (e) {
                return { success: false, error: 'Не удалось сохранить аватарку питомца' };
            }
        }
    }

    if (bloodRequestPhoto) {
        try {
            const { success } = await addPhoto({
                photo: bloodRequestPhoto,
                buffer: bloodRequestPhotoBuffer,
                id: poolRequestId,
                isBloodRequest: true,
            });

            if (!success) {
                return { success: false, error: 'Не удалось прикрепить фотографию к заявке на поиск крови' };
            }
        } catch (e) {
            return { success: false, error: 'Не удалось прикрепить фотографию к заявке на поиск крови' };
        }
    }

    return { success: true };
};
