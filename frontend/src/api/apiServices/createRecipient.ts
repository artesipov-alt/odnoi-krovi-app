import { AddToPoolRequest } from '../bloodRequest';
import api from '../index';
import { CreatePetRequest } from '../pets';
import { addPhoto } from './addPhoto';

type Args = CreatePetRequest & {
    photo: File | null;
    bloodRequestPhoto: File | null;
    poolInfo: Omit<AddToPoolRequest, 'petId'>;
};

export const createRecipient = async ({ photo, poolInfo, bloodRequestPhoto, ...params }: Args) => {
    let petId;
    let poolRequestId;

    try {
        const { data } = await api.createPet(params);

        petId = data.id;
    } catch (e) {
        return { success: false, error: 'Не удалось создать профиль питомца' };
    }

    try {
        const { data } = await api.addToPool({ ...poolInfo, petId });

        poolRequestId = data.id;
    } catch (e) {
        return { success: false, error: 'Профиль питомца создано, но не удалось создать заявку поиск крови' };
    }

    if (photo) {
        try {
            const { success } = await addPhoto({ photo, id: petId, isAvatar: true });

            if (!success) {
                return { success: false, error: 'Не удалось сохранить аватарку питомца' };
            }
        } catch (e) {
            return { success: false, error: 'Не удалось сохранить аватарку питомца' };
        }
    }

    if (bloodRequestPhoto) {
        try {
            const { success } = await addPhoto({ photo: bloodRequestPhoto, id: poolRequestId, isBloodRequest: true });

            if (!success) {
                return { success: false, error: 'Не удалось прикрепить фотографию к заявке на поиск крови' };
            }
        } catch (e) {
            return { success: false, error: 'Не удалось прикрепить фотографию к заявке на поиск крови' };
        }
    }

    return { success: true };
};
