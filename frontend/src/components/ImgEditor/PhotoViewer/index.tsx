import { FC, useState } from 'react';
import { PhotoSlider } from 'react-photo-view';
import { DataType } from 'react-photo-view/dist/types';

type Props = {
    images: DataType[];
    selectedImg?: string;
    onCloseClick: () => void;
};
const PhotoViewer: FC<Props> = ({ selectedImg, images, onCloseClick }) => {
    const [index, setIndex] = useState(selectedImg ? images.findIndex((item) => item.src === selectedImg) || 0 : 0);

    return (
        <PhotoSlider
            visible
            index={index}
            onClose={onCloseClick}
            onIndexChange={setIndex}
            images={images.map(({ src }) => ({ src: src!, key: src! }))}
        />
    );
};
export default PhotoViewer;
