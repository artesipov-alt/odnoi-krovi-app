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

export enum AnalysesEnum {
    LEUKEMIA = 'leukemia',
    BABESIOSIS = 'babesiosis',
    DIROFILARIA = 'dirofilaria',
    EHRLICHIOSIS = 'ehrlichiosis',
    ANAPLASMOSIS = 'anaplasmosis',
    BARTONELLOSIS = 'bartonellosis',
    HEMOPLASMOSIS = 'hemoplasmosis',
    IMMUNODEFICIENCY = 'immunodeficiency',
}

export const AnalysesNamesMapping = {
    [AnalysesEnum.BABESIOSIS]: 'Анализ на бабезиоз',
    [AnalysesEnum.LEUKEMIA]: 'Анализ на лейкоз (ВЛК)',
    [AnalysesEnum.EHRLICHIOSIS]: 'Анализ на эрлихиоз',
    [AnalysesEnum.ANAPLASMOSIS]: 'Анализ на анаплазмоз',
    [AnalysesEnum.DIROFILARIA]: 'Анализ на дирофиляриоз',
    [AnalysesEnum.BARTONELLOSIS]: 'Анализ на бартонеллез',
    [AnalysesEnum.HEMOPLASMOSIS]: 'Анализ на гемоплазмоз',
    [AnalysesEnum.IMMUNODEFICIENCY]: 'Анализ на иммунодефицит (ВИК)',
};
