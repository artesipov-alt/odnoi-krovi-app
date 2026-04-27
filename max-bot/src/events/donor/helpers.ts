export const generateMessage = (params: {
  recipientPetName: string;
  recipientPetNeededVolume: number;
  recipientPetSearchingBloodGroup: string[];
  donorName: string;
  donorBloodGroup: string;
}): string => {
  const donorName = params.donorName?.trim() || "Анонимный донор";
  const donorBloodGroup =
    params.donorBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.donorBloodGroup;

  return `На Ваш поиск (${params.recipientPetName}, ${params.recipientPetNeededVolume} мл, группа ${params.recipientPetSearchingBloodGroup.join(", ")}) откликнулся донор ${donorName} (группа ${donorBloodGroup}).
Ознакомьтесь с информацией о доноре в приложении.`;
};
