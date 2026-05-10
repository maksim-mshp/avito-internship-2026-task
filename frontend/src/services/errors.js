const TRANSLATED_ERRORS = {
    INVALID_REQUEST: 'Некорректный запрос',
    UNAUTHORIZED: 'Необходим вход',
    FORBIDDEN: 'Доступ запрещён',
    NOT_FOUND: 'Не найдено',
    EMAIL_TAKEN: 'Email уже занят',
    INVALID_CREDENTIALS: 'Неверный email или пароль',
    INTERNAL_ERROR: 'Внутренняя ошибка сервера',
    'Network Error': 'Сетевая ошибка',
}

export const getTranslatedError = (error) => {
    const responseError = error.response?.data?.error
    const message = responseError?.code || responseError?.message || error.message

    return TRANSLATED_ERRORS[message] || responseError?.message || message
}
