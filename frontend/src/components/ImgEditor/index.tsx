import cn from 'classnames';
import Cancel from 'imgs/svg/cancel';
import Edit from 'imgs/svg/edit';
import Plus from 'imgs/svg/plus';
import { ChangeEvent, FC, memo, MouseEvent, useRef, useState } from 'react';

import styles from './ImgEditor.module.less';
import PhotoViewer from './PhotoViewer';

// ═══════════════════════════════════════════════════════════════════════
// Fallback для Telegram WebView на Android:
// Telegram перехватывает <input type='file'>, открывает свою галерею,
// и выдаёт content:// URI, который уже невалиден к моменту onChange.
// arrayBuffer() падает, но браузерный декодер изображений (через <img>)
// часто имеет более широкие права и может прочитать файл.
// Поэтому: грузим в <img> → отрисовываем на <canvas> → экспортируем Blob.
// ═══════════════════════════════════════════════════════════════════════
async function readFileViaCanvas(file: File): Promise<File> {
    return new Promise((resolve, reject) => {
        const img = new Image();
        const url = URL.createObjectURL(file);

        img.onload = () => {
            URL.revokeObjectURL(url);
            const canvas = document.createElement('canvas');
            canvas.width = img.naturalWidth;
            canvas.height = img.naturalHeight;
            const ctx = canvas.getContext('2d');

            if (!ctx) {
                reject(new Error('Canvas 2D context not available'));

                return;
            }
            ctx.drawImage(img, 0, 0);

            canvas.toBlob(
                (blob) => {
                    if (blob) {
                        resolve(new File([blob], file.name, { type: file.type || 'image/jpeg' }));

                        return;
                    }

                    reject(new Error('Canvas toBlob returned null'));
                },
                file.type || 'image/jpeg',
                0.95,
            );
        };

        img.onerror = () => {
            URL.revokeObjectURL(url);
            reject(new Error('Failed to load image via blob URL'));
        };

        img.src = url;
    });
}

type Props = {
    name?: string;
    weight?: string;
    src: File | null;
    petType?: string;
    className?: string;
    showStub?: boolean;
    serverSrc?: string;
    bloodGroup?: string;
    isMiniView?: boolean;
    isEditIcon?: boolean;
    onLoad?: (photo: File | null) => void;
};

// const acceptableFormats = ['png', 'jpg', 'jpeg', 'jpe', 'webp', 'heic', 'raw'];

const ImgEditor: FC<Props> = ({
    src,
    name,
    weight,
    onLoad,
    petType,
    showStub,
    className,
    serverSrc,
    bloodGroup,
    isMiniView,
    isEditIcon,
}) => {
    const [imgSrc, setImgSrc] = useState('');
    const [file, setFile] = useState<File | null>(src);
    const [isImageLoaded, setIsImageLoaded] = useState(false);
    const [isLoadImageError, setIsLoadImageError] = useState(false);
    const [isRenderPhotoSlider, setIsRenderPhotoSlider] = useState(false);

    const fileInputRef = useRef<HTMLInputElement | null>(null);

    const onAddItemClickHandler = () => {
        fileInputRef.current?.click();
    };
    const onLoadFileHandler = async ({ currentTarget }: ChangeEvent<HTMLInputElement>) => {
        const newFile = currentTarget?.files?.[0];

        if (!newFile) {
            return;
        }

        try {
            // Способ 1: arrayBuffer — работает в iOS, Max WebView, и большинстве случаев
            const buffer = await newFile.arrayBuffer();
            const safeFile = new File([buffer], newFile.name, { type: newFile.type || 'image/jpeg' });

            setFile(safeFile);
            setIsLoadImageError(false);
            onLoad?.(safeFile);
        } catch {
            // Способ 2: Canvas fallback для Telegram WebView на Android,
            // где content:// URI уже недоступен на момент onChange
            try {
                const safeFile = await readFileViaCanvas(newFile);
                setFile(safeFile);
                setIsLoadImageError(false);
                onLoad?.(safeFile);
            } catch (canvasError) {
                console.error('[ImgEditor] all read methods failed:', canvasError);
                setIsLoadImageError(true);
            }
        }
    };

    const onDeleteClickHandler = (e: MouseEvent<HTMLDivElement>) => {
        e.stopPropagation();

        fileInputRef.current!.value = '';

        setFile(null);
        onLoad?.(null);
    };

    const onImgErrorHandler = () => {
        setIsLoadImageError(true);
        setIsImageLoaded(true);
    };

    const onImgLoadHandler = () => {
        setIsImageLoaded(true);
    };

    const onPhotoSliderToggleHandler = () => {
        setIsRenderPhotoSlider((prevState) => !prevState);
    };

    const onImgClickHandler = (e: MouseEvent<HTMLImageElement>) => {
        const newSrc = e.currentTarget?.src;

        setImgSrc(newSrc);
        onPhotoSliderToggleHandler();
    };

    const renderErrorStub = () => <div>Не удалось загрузить изображение</div>;

    const renderDeleteButton = () => (
        <div className={styles.delete} onClick={onDeleteClickHandler}>
            <Cancel />
        </div>
    );

    const renderEditButton = () => (
        <div className={styles.edit} onClick={onAddItemClickHandler}>
            <Edit />
        </div>
    );

    const renderLabels = () => (
        <>
            {!isMiniView && <span className={styles.bloodGroup}>{bloodGroup !== 'UNKNOWN' ? bloodGroup : '?'}</span>}
            {weight && (
                <div className={styles.weight}>
                    <span>{weight}</span>
                    <span className={styles.weightCaption}> кг</span>
                </div>
            )}
            <div className={styles.name}>{name?.toUpperCase()}</div>
        </>
    );

    const renderContent = () => (
        <>
            {!src && !file && !serverSrc && showStub && petType && (
                <div className={cn(styles.stub, { [styles[petType]]: true, [styles.isEdit]: serverSrc })}>
                    {isEditIcon && renderEditButton()}
                    {renderLabels()}
                </div>
            )}
            {!src && !file && !serverSrc && !showStub && (
                <div className={styles.photo} onClick={onAddItemClickHandler}>
                    <div className={styles.addPhoto}>
                        <Plus />
                    </div>
                    <span className={styles.caption}>
                        Добавьте фото
                        <br />
                        питомца
                    </span>
                </div>
            )}
            {!file && (serverSrc || src) && isLoadImageError && <div>{renderErrorStub()}</div>}
            {!file && serverSrc && !isLoadImageError && (
                <div className={styles.imgWrapper}>
                    <img
                        alt={name}
                        src={serverSrc}
                        className={styles.img}
                        onLoad={onImgLoadHandler}
                        onError={onImgErrorHandler}
                        onClick={onImgClickHandler}
                    />
                    {onLoad && !isEditIcon && renderDeleteButton()}
                    {isEditIcon && renderEditButton()}
                    {showStub && renderLabels()}
                    {showStub && <div className={styles.gradient} />}
                </div>
            )}
            {file && (
                <div className={cn(styles.imgWrapper, { [styles.imgWithLabels]: showStub })}>
                    <img className={styles.img} src={URL.createObjectURL(file)} onClick={onImgClickHandler} alt='' />
                    {onLoad && !isEditIcon && renderDeleteButton()}
                    {isEditIcon && renderEditButton()}
                    {showStub && renderLabels()}
                    {showStub && <div className={styles.gradient} />}
                </div>
            )}
            <input
                type='file'
                id='imageInput'
                accept='image/*'
                ref={fileInputRef}
                className={styles.input}
                onChange={onLoadFileHandler}
            />
        </>
    );

    return (
        <div
            className={cn(styles.container, className, {
                [styles.showStub]: showStub,
                [styles.withPhoto]: src || file || serverSrc,
                [styles.mini]: isMiniView,
            })}
        >
            {renderContent()}
            {isRenderPhotoSlider && (
                <PhotoViewer images={[{ src: imgSrc, key: imgSrc }]} onCloseClick={onPhotoSliderToggleHandler} />
            )}
        </div>
    );
};
export default memo(ImgEditor);
