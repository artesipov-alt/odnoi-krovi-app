declare module '*.less' {
    const classes: { [key: string]: string };
    export default classes;
}

declare module '*.svg' {
    import React from 'react';

    const content: React.FunctionComponent<React.SVGProps<SVGElement>>;
    export default content;
}

declare module '*.jpg' {
    const image: string;
    export default image;
}

declare module '*.png' {
    const image: string;
    export default image;
}
