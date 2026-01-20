export enum PetType {
    DOG = 'dog',
    CAT = 'cat',
}

export enum PetGender {
    MALE = 'male',
    FEMALE = 'female',
}
export enum AnalysesTypes {
    PCR = 'PCR',
    ICA = 'ICA',
    ELISA = 'ELISA',
    EXPRESS = 'Express',
    MICROSCOPY = 'Microscopy',
}

export const AnalysesMapping = {
    ПЦР: AnalysesTypes.PCR,
    ИХА: AnalysesTypes.ICA,
    'Экспресс-тест': AnalysesTypes.EXPRESS,
    'Микроскопия мазка': AnalysesTypes.MICROSCOPY,
    ИФА: AnalysesTypes.ELISA,
};
