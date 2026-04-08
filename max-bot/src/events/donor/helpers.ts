export const generateMessage = (params: {
  recipientPetName: string;
  recipientPetNeededVolume: number;
  recipientPetSearchingBloodGroup: string[];
  donorName: string;
  donorBloodGroup: string;
}): string => {
  return `На Ваш поиск (${params.recipientPetName}, ${params.recipientPetNeededVolume} мл, группа ${params.recipientPetSearchingBloodGroup.join(", ")}) откликнулся донор ${params.donorName} (группа ${params.donorBloodGroup}).
Ознакомьтесь с информацией о доноре в приложении.`;
};
