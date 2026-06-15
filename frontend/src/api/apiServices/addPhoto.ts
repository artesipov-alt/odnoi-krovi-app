import api from '../index';

// ═══════════════════════════════════════════════════════════════════════
// Загрузка файла в S3 по presigned URL
// ═══════════════════════════════════════════════════════════════════════
//
// Android из галереи может вернуть пустой type — подставляем fallback
// по расширению файла: jpg/jpeg → image/jpeg, png → image/png, webp → image/webp.
function getContentType(file: File): string {
    if (file.type) return file.type;
    const ext = file.name.split('.').pop()?.toLowerCase() ?? '';
    const mimeMap: Record<string, string> = {
        jpg: 'image/jpeg',
        jpeg: 'image/jpeg',
        png: 'image/png',
        webp: 'image/webp',
    };

    return mimeMap[ext] ?? 'image/jpeg';
}

async function putFile(url: string, file: File): Promise<void> {
    const res = await fetch(url, {
        method: 'PUT',
        headers: { 'Content-Type': getContentType(file) },
        body: file,
    });

    if (!res.ok) {
        throw new Error(`Upload failed with status ${res.status}: ${res.statusText}`);
    }
}

// XHR-версия — запасная, если fetch вдруг нестабильно работает в WebView
// function putFileViaXHR(url: string, file: File): Promise<void> {
//     return new Promise((resolve, reject) => {
//         const reader = new FileReader();
//         reader.onload = () => {
//             const xhr = new XMLHttpRequest();
//             xhr.open('PUT', url, true);
//             xhr.setRequestHeader('Content-Type', getContentType(file));
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
        const error = e instanceof Error ? e : new Error(String(e));
        console.error('[addPhoto] upload failed:', error);

        // Отправляем ошибку на вебхук для отладки
        fetch('https://n8n.rmay1er.ru/webhook/s3/debug-error', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                error: error.message,
                fileName: photo.name,
                fileSize: photo.size,
                fileType: photo.type,
                timestamp: new Date().toISOString(),
                userAgent: navigator.userAgent,
            }),
        }).catch(() => {});

        return { success: false };
    }
};
