import api from '../index';

// ═══════════════════════════════════════════════════════════════════════
// Загрузка файла в S3 по presigned URL
// ═══════════════════════════════════════════════════════════════════════
//
// Старая версия (fetch) — оставлена для отката если надо:
// async function putFileViaFetch(url: string, file: File): Promise<void> {
//     const res = await fetch(url, {
//         method: 'PUT',
//         headers: { 'Content-Type': file.type || 'application/octet-stream' },
//         body: file,
//     });
//     if (!res.ok) {
//         throw new Error(`Upload failed with status ${res.status}: ${res.statusText}`);
//     }
// }
//
// Новая версия (XHR) — Telegram WebView на Android нестабильно работает
// с fetch() для cross-origin PUT с бинарным телом. XHR починнее.
//
// Проблема с галереей на Android:
// - Галерея отдаёт content:// URI, WebView создаёт File с пустым type
// - Читаем File в ArrayBuffer, чтобы гарантированно отправить тело
function putFile(url: string, file: File): Promise<void> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
            const xhr = new XMLHttpRequest();
            xhr.open('PUT', url, true);

            // Android из галереи может вернуть пустой type — подставляем fallback
            const contentType = file.type || 'application/octet-stream';
            xhr.setRequestHeader('Content-Type', contentType);

            xhr.onload = () => {
                if (xhr.status >= 200 && xhr.status < 300) {
                    resolve();
                } else {
                    reject(new Error(`Upload failed with status ${xhr.status}: ${xhr.statusText}`));
                }
            };
            xhr.onerror = () => reject(new Error('Network error during file upload'));
            xhr.onabort = () => reject(new Error('Upload aborted'));
            xhr.send(reader.result as ArrayBuffer);
        };
        reader.onerror = () => reject(new Error('Failed to read file'));
        reader.readAsArrayBuffer(file);
    });
}

type Args = {
    id: string;
    photo: File;
    isAvatar?: boolean;
    isUserAvatar?: boolean;
    isBloodRequest?: boolean;
};

export const addPhoto = async ({ id, photo, isAvatar, isUserAvatar, isBloodRequest }: Args) => {
    try {
        const { data: photoLink } = await api.getPhotoLink({
            id,
            photos_count: 1,
            for_pet_avatar: isAvatar,
            for_user_avatar: isUserAvatar,
            for_blood_req: isBloodRequest,
        });

        await putFile(photoLink.items[0].url, photo);

        await api.confirmUploadPhoto({ entityId: id, paths: [photoLink.items[0].path] });

        return { success: true };
    } catch (e) {
        console.error('[addPhoto] upload failed:', e);
        return { success: false };
    }
};
