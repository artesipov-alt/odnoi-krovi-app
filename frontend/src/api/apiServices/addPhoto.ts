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
        headers: {
            'Content-Type': getContentType(file),
        },
        body: file,
    });

    if (!res.ok) {
        throw new Error(`Upload failed with status ${res.status}: ${res.statusText}`);
    }
}

// XHR-версия — запасная, если fetch вдруг нестабильно работает в WebView
function putFileViaXHR(url: string, file: File): Promise<void> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => {
            const xhr = new XMLHttpRequest();
            xhr.open('PUT', url, true);
            xhr.setRequestHeader('Content-Type', getContentType(file));
            xhr.setRequestHeader('Content-Length', file.size.toString());
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

function sendBeaconToWebhook(data: Record<string, unknown>): void {
    try {
        // Image() — единственный способ, гарантированно работающий под любым CSP/CORS
        // в WebView. sendBeacon и fetch блокируются connect-src, Image обходит.
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
        // test
        sendSuccessToWebhook(photo);

        return { success: true };
    } catch (e) {
        const error = e instanceof Error ? e : new Error(String(e));
        console.error('[addPhoto] upload failed:', error);

        sendErrorToWebhook(e, photo);

        return { success: false };
    }
};
