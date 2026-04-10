# TODO: Удалить проверки BloodGroupName после установки значений по умолчанию

После того, как убедимся, что у всех питомцев в базе данных установлена группа крови по умолчанию (так что BloodGroupName никогда не nil), удалить следующие проверки nil, чтобы упростить код. Эти проверки были добавлены как хотфиксы, чтобы предотвратить паники, когда группы крови отсутствовали.

## Файлы для обновления:

1. **odnoi-krovi-app/backend/internal/infra/presistance/domainmapper/bloodreq_mapper.go**
   - Удалить проверку `if resp.Edges.Donor.Edges.BloodGroupRef != nil` около строки 72.
   - Вернуть прямое присваивание: `app.DonorBloodGroup = resp.Edges.Donor.Edges.BloodGroupRef.BloodGroup`
   - Причина: С установленными по умолчанию, BloodGroupRef всегда будет существовать.

2. **odnoi-krovi-app/backend/internal/application/bloodsearch/cmd/accept_response.go**
   - Удалить переменные `donorBloodGroup` и `recipientBloodGroup` и их проверки nil (около строк 118-127).
   - Вернуть прямое разыменование: `BloodGroup: *donorPet.BloodGroupName,` и `BloodGroup: *recipientPet.BloodGroupName,`
   - Причина: BloodGroupName всегда будет не-nil после применения по умолчанию.

3. **odnoi-krovi-app/backend/internal/application/bloodsearch/cmd/confirm_donation.go**
   - Удалить переменную `donorBloodGroup` и её проверку nil (около строк 113-117).
   - Вернуть прямое разыменование: `BloodGroup: *donorPet.BloodGroupName,`
   - Причина: BloodGroupName всегда будет не-nil после применения по умолчанию.

## Шаги для завершения:
- Запустить миграцию базы данных/скрипт для установки групп крови по умолчанию для всех питомцев, где blood_group_id NULL.
- Обновить логику создания/обновления питомцев, чтобы blood_group_id всегда устанавливался.
- Удалить проверки, перечисленные выше.
- Протестировать приложение, чтобы убедиться в отсутствии регрессий.
- Удалить этот файл TODO.

## Другие TODO:
- [ ] Добавить здесь будущие задачи по мере необходимости.