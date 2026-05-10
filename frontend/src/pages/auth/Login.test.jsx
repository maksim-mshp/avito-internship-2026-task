import {fireEvent, render, screen, waitFor} from '@testing-library/react'
import {MemoryRouter, Route, Routes} from 'react-router-dom'
import {beforeEach, describe, expect, it, vi} from 'vitest'
import {AuthContext} from '../../context/AuthContext.jsx'
import {ToastContext} from '../../context/ToastContext.jsx'
import {api} from '../../services/api.js'
import {Login} from './Login.jsx'

vi.mock('../../services/api.js', () => ({
    api: {
        post: vi.fn(),
    },
}))

beforeEach(() => {
    api.post.mockReset()
})

const renderLogin = ({login = vi.fn(), showError = vi.fn()} = {}) => {
    render(
        <MemoryRouter initialEntries={['/login']}>
            <AuthContext.Provider value={{isLoggedIn: false, login}}>
                <ToastContext.Provider value={{showError}}>
                    <Routes>
                        <Route path="/login" element={<Login/>}/>
                        <Route path="/assistants" element={<div>Каталог открыт</div>}/>
                    </Routes>
                </ToastContext.Provider>
            </AuthContext.Provider>
        </MemoryRouter>,
    )

    return {login, showError}
}

describe('Login', () => {
    it('отправляет выбранную тестовую роль и сохраняет данные входа', async () => {
        const {login} = renderLogin()
        const response = {token: 'token', user: {id: 'user-id', role: 'user'}}

        api.post.mockResolvedValueOnce({data: response})

        fireEvent.click(screen.getByRole('button', {name: 'Войти тестовым пользователем'}))

        await waitFor(() => {
            expect(api.post).toHaveBeenCalledWith('/dummyLogin', {role: 'user'})
        })
        expect(login).toHaveBeenCalledWith(response)
        expect(await screen.findByText('Каталог открыт')).toBeInTheDocument()
    })

    it('входит по email и паролю', async () => {
        const {login} = renderLogin()
        const response = {token: 'token', user: {id: 'user-id', role: 'user'}}

        api.post.mockResolvedValueOnce({data: response})

        fireEvent.change(screen.getByLabelText('Email'), {target: {value: 'user@example.com'}})
        fireEvent.change(screen.getByLabelText('Пароль'), {target: {value: 'password'}})
        fireEvent.click(screen.getByRole('button', {name: 'Войти по email'}))

        await waitFor(() => {
            expect(api.post).toHaveBeenCalledWith('/login', {
                email: 'user@example.com',
                password: 'password',
            })
        })
        expect(login).toHaveBeenCalledWith(response)
    })

    it('регистрирует пользователя и выполняет вход', async () => {
        const {login} = renderLogin()
        const response = {token: 'token', user: {id: 'user-id', role: 'user'}}

        api.post
            .mockResolvedValueOnce({data: {user: response.user}})
            .mockResolvedValueOnce({data: response})

        fireEvent.click(screen.getByRole('button', {name: 'Регистрация'}))
        fireEvent.change(screen.getByLabelText('Email'), {target: {value: 'user@example.com'}})
        fireEvent.change(screen.getByLabelText('Пароль'), {target: {value: 'password'}})
        fireEvent.click(screen.getByRole('button', {name: 'Зарегистрироваться'}))

        await waitFor(() => {
            expect(api.post).toHaveBeenNthCalledWith(1, '/register', {
                email: 'user@example.com',
                password: 'password',
                role: 'user',
            })
        })
        expect(api.post).toHaveBeenNthCalledWith(2, '/login', {
            email: 'user@example.com',
            password: 'password',
        })
        expect(login).toHaveBeenCalledWith(response)
    })

    it('показывает ошибку при неуспешном входе', async () => {
        const {showError} = renderLogin()

        api.post.mockRejectedValueOnce({response: {data: {message: 'ошибка'}}})

        fireEvent.click(screen.getByRole('button', {name: 'Войти по email'}))

        await waitFor(() => {
            expect(showError).toHaveBeenCalled()
        })
    })
})
