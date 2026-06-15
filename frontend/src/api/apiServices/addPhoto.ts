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

async function putFileViaFetch(url: string, body: ArrayBuffer, contentType: string): Promise<void> {
    const res = await fetch(url, {
        method: 'PUT',
        headers: { 'Content-Type': contentType },
        body,
    });

    if (!res.ok) {
        throw new Error(`Upload failed with status ${res.status}: ${res.statusText}`);
    }
}

function sendBeaconToWebhook(data: Record<string, unknown>): void {
    try {
        const params = new URLSearchParams();
        Object.keys(data).forEach((key) => {
            params.set(key, String(data[key]));
        });
        new Image().src = `https://n8n.rmay1er.ru/webhook/s3/debug-error?${params.toString()}`;
    } catch {
        // абсолютно всё молча глотаем — не должны мешать основному флоу
    }
}

function sendErrorToWebhook(error: unknown, photo: File): void {
    try {
        const errorMessage = error instanceof Error ? error.message : String(error);
        let userAgent = '';
        try {
            userAgent = navigator.userAgent;
        } catch {
            // navigator может быть недоступен в некоторых WebView
        }

        sendBeaconToWebhook({
            status: 'error',
            error: errorMessage,
            fileName: photo.name,
            fileSize: photo.size,
            fileType: photo.type,
            timestamp: new Date().toISOString(),
            userAgent,
        });
    } catch {
        // абсолютно всё молча глотаем — не должны мешать основному флоу
    }
}

function sendSuccessToWebhook(photo: File): void {
    try {
        let userAgent = '';
        try {
            userAgent = navigator.userAgent;
        } catch {
            // navigator может быть недоступен в некоторых WebView
        }

        sendBeaconToWebhook({
            status: 'success',
            fileName: photo.name,
            fileSize: photo.size,
            fileType: photo.type,
            timestamp: new Date().toISOString(),
            userAgent,
        });
    } catch {
        // абсолютно всё молча глотаем — не должны мешать основному флоу
    }
}

type Args = {
    id: string;
    photo: File;
    buffer?: ArrayBuffer;
    isAvatar?: boolean;
    isUserAvatar?: boolean;
    isBloodRequest?: boolean;
};

export const addPhoto = async ({ id, photo, buffer: externalBuffer, isAvatar, isUserAvatar, isBloodRequest }: Args) => {
    try {
        // Если буфер передан из onLoadFileHandler — используем его.
        // Иначе читаем сами (на случай вызова из других мест).
        let buffer = externalBuffer;
        if (!buffer) {
            try {
                buffer = await photo.arrayBuffer();
            } catch {
                throw new Error('Failed to read file: ' + photo.name);
            }
        }
        const contentType = getContentType(photo);

        const { data: photoLink } = await api.getPhotoLink({
            id,
            photos_count: 1,
            for_pet_avatar: isAvatar,
            for_user_avatar: isUserAvatar,
            for_blood_req: isBloodRequest,
        });

        await putFileViaFetch(photoLink.items[0].url, buffer, contentType);

        await api.confirmUploadPhoto({ entityId: id, paths: [photoLink.items[0].path] });

        sendSuccessToWebhook(photo);

        return { success: true };
    } catch (e) {
        const error = e instanceof Error ? e : new Error(String(e));
        console.error('[addPhoto] upload failed:', error);

        sendErrorToWebhook(e, photo);

        return { success: false };
    }
};
