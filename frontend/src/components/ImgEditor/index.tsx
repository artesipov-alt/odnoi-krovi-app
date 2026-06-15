import cn from 'classnames';
import Cancel from 'imgs/svg/cancel';
import Edit from 'imgs/svg/edit';
import Plus from 'imgs/svg/plus';
import { ChangeEvent, FC, memo, MouseEvent, useRef, useState } from 'react';

import styles from './ImgEditor.module.less';
import PhotoViewer from './PhotoViewer';

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
    const onLoadFileHandler = ({ currentTarget }: ChangeEvent<HTMLInputElement>) => {
        const newFile = currentTarget?.files?.[0];

        if (!newFile) {
            return;
        }

        // if (!acceptableFormats.includes(newFile.type.split('/')[1])) {
        //     // не тот формат
        //     return;
        // }

        setFile(newFile);
        setIsLoadImageError(false);

        onLoad?.(newFile);
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
                accept='image/jpeg,image/png,image/webp,image/*'
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
