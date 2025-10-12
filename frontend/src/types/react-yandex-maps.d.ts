// frontend/src/react-yandex-maps.d.ts
declare module '@pbe/react-yandex-maps' {
  import { FC } from 'react';

  interface MapProps {
    defaultState: {
      center: [number, number];
      zoom: number;
    };
    width?: string;
    height?: string;
  }

  export const YMaps: FC<any>;
  export const Map: FC<MapProps>;
}