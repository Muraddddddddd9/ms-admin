package constants

const (
	ErrServerError = "Ошибка сервера"

	ErrInvalidInput       = "Данные введены неверно"
	ErrDataInCollection   = "Данные находятся в коллекции %s"
	ErrDataAlreadyExists  = "Данные уже существуют"
	ErrFieldCannotEmpty   = "Поле '%s' не может быть пустым"
	ErrLoadEnv            = "Не удалось загрузить env"
	ErrUserNotFound       = "Пользователь не найден"
	ErrCollectionNotFound = "Не верная коллекция"

	ErrAdminConfig       = "Адрес электронной почты администратора/пароль не заданы в конфигурации"
	ErrHashPassword      = "Не удалось хэшировать пароль: %v"
	ErrCreateAdminStatus = "Не удалось создать статус 'админ': %v"
	ErrAdminNotFound     = "Не удалось найти статус 'админ': %v"
	ErrCreateAdmin       = "Не удалось создать администратора: %v"
	ErrCheckAdmin        = "Не удалось проверить наличие администратора: %v"
	ErrCreateStatus      = "Не удалось создать статус: %v. Ошибка: %v\n"

	ErrDeleteData  = "Ошибка в удалении"
	ErrUpdateData  = "Обновление данных провалилось"
	ErrDataLogging = "Ошибка в логировании данных"
	ErrGetData     = "Ошибка в получении данных"

	ErrTeacherNotFound = "Учитель не найден"
	ErrObjectNotFound  = "Предмет не найден"
	ErrGroupNotFound   = "Группа не найдена"
	ErrStatusNotFound  = "Статус не найден"
	ErrSessionNotFound = "Сессия не найдена"

	ErrInvalidDataStudent        = "Неверные данные студента"
	ErrInvalidDataTeacher        = "Неверные данные учителя"
	ErrInvalidDataStatus         = "Неверные данные статуса"
	ErrInvalidDataGroup          = "Неверные данные группы"
	ErrInvalidDataObject         = "Неверные данные предмета"
	ErrInvalidDataObjectForGroup = "Неверные данные предмета для группы"
)

const (
	SuccConnectMongo       = "Подключение к MONGODB - успешно"
	SuccConnectRedis       = "Подключение к REDIS - успешно"
	SuccDataAdd            = "Данные добавлены с ID: %v"
	SuccDataDelete         = "Было удалено %v из %v"
	SuccDataUpdate         = "Данные обновлены"
	SuccCreateAdmin        = "Администратор создан"
	SuccAdminAlreadyExists = "Администратор уже существует"
)
