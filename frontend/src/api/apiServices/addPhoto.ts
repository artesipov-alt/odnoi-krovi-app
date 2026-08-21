import api from '../index';

// ═══════════════════════════════════════════════════════════════════════
// Загрузка файла в S3 по presigned URL
// ═══════════════════════════════════════════════════════════════════════
//
// Android из галереи может вернуть пустой type — подставляем fallback
// по расширению файла: jpg/jpeg → image/jpeg, png → image/png, webp → image/webp.
// function getContentType(file: File): string {
//     if (file.type) return file.type;
//     const ext = file.name.split('.').pop()?.toLowerCase() ?? '';
//     const mimeMap: Record<string, string> = {
//         jpg: 'image/jpeg',
//         jpeg: 'image/jpeg',
//         png: 'image/png',
//         webp: 'image/webp',
//     };

//     return mimeMap[ext] ?? 'image/jpeg';
// }

async function putFile(url: string, file: File): Promise<void> {
    // Читаем файл как ArrayBuffer, чтобы:
    // 1. Не зависеть от content:// URI (Android)
    // 2. Не слать Content-Type — presigned URL его не проверяет (ContentType закомментирован в Go),
    //    а отсутствие кастомных заголовков снижает требования к CORS preflight
    const buffer = await file.arrayBuffer();

    const res = await fetch(url, {
        method: 'PUT',
        body: buffer,
    });

    if (!res.ok) {
        throw new Error(`Upload failed with status ${res.status}: ${res.statusText}`);
    }
}

// XHR fallback — без кастомных заголовков, чтобы не провоцировать preflight
// Используется когда fetch PUT не работает в WebView (например, Telegram Android)
// function putFileViaXHR(url: string, file: File): Promise<void> {
//     return new Promise((resolve, reject) => {
//         const reader = new FileReader();
//         reader.onload = () => {
//             const xhr = new XMLHttpRequest();
//             xhr.open('PUT', url, true);
//             // НЕ ставим кастомные заголовки — они вызывают preflight OPTIONS,
//             // который Telegram Android WebView может заблокировать
//             xhr.onload = () => {
//                 if (xhr.status >= 200 && xhr.status < 300) {
//                     resolve();
//                 } else {
//                     reject(new Error(`Upload failed with status ${xhr.status}: ${xhr.statusText}`));
//                 }
//             };
//             xhr.onerror = () => reject(new Error('Network error during file upload'));
//             xhr.onabort = () => reject(new Error('Upload aborted'));
//             xhr.send(reader.result as ArrayBuffer);
//         };
//         reader.onerror = () => reject(new Error('Failed to read file'));
//         reader.readAsArrayBuffer(file);
//     });
// }

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
        return { success: false };
    }
};
