export enum BirthDate {
    EXACT_DATE = 'exactDate',
    APPROXIMATE_DATE = 'approximateDate',
}

export type AnalizItem = {
    name: string;
    value: Date | null;
    defaultValue?: Date | null;
};

export type Analiz = {
    type: string;
    name: string;
    items: AnalizItem[];
};
