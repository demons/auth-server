# Список api
## `public POST /token` получение нового токена
`?grant_type=password` аутентификация по паролю

Тело запроса:
* *email*
* *password*

`?grant_type=code` аутентификация с помощью соц. сети

Тело запроса:
* *provider* (google, facebook, vk)
* *code* (код, который возвращает соц. сеть)

`?grant_type=refresh` смена токена доступа, с помощью *refresh* токена

Тело запроса:
* *refresh*

## `public POST /account/signup` регистрация нового пользователя

Тело запроса:
* *email*
* *password*

## `public POST /account/password/reset` сброс пароля
При вызове этого api, будет отправлен email на указанный адрес, со ссылкой для изменения пароля

Тело запроса:
* *email*

## `private POST /account/password/change`  изменение пароля
Заголовоки:
* *Authorization* (Bearer access_token)

Тело запроса:
* *oldPassword*
* *newPassword*