import api from '../index';

type Args = {
    id: string;
    photo: File;
    isAvatar?: boolean;
    isBloodRequest?: boolean;
};

export const addPhoto = async ({ id, photo, isAvatar, isBloodRequest }: Args) => {
    try {
        const { data: photoLink } = await api.getPhotoLink({
            id,
            photos_count: 1,
            for_pet_avatar: isAvatar,
            for_blood_req: isBloodRequest,
        });

        await fetch(photoLink.items[0].url, {
            method: 'PUT',
            headers: {
                'Content-Type': photo.type,
            },
            body: photo,
        });

        await api.confirmUploadPhoto({ entityId: id, paths: [photoLink.items[0].path] });

        return { success: true };
    } catch (e) {
        return { success: false };
    }
};
