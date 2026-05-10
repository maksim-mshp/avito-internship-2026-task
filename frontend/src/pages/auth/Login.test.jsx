import {render, screen, waitFor} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import {MemoryRouter, Route, Routes} from 'react-router-dom'
import {describe, expect, it, vi} from 'vitest'
import {AuthContext} from '../../context/AuthContext.jsx'
import {ToastContext} from '../../context/ToastContext.jsx'
import {api} from '../../services/api.js'
import {Login} from './Login.jsx'

vi.mock('../../services/api.js', () => ({
    api: {
        post: vi.fn(),
    },
}))

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
    it('отправляет выбранную роль и сохраняет данные входа', async () => {
        const user = userEvent.setup()
        const {login} = renderLogin()
        const response = {token: 'token', user: {id: 'user-id', role: 'user'}}

        api.post.mockResolvedValueOnce({data: response})

        await user.click(screen.getByRole('button', {name: 'Войти'}))

        await waitFor(() => {
            expect(api.post).toHaveBeenCalledWith('/dummyLogin', {role: 'user'})
        })
        expect(login).toHaveBeenCalledWith(response)
        expect(await screen.findByText('Каталог открыт')).toBeInTheDocument()
    })

    it('показывает ошибку при неуспешном входе', async () => {
        const user = userEvent.setup()
        const {showError} = renderLogin()

        api.post.mockRejectedValueOnce({response: {data: {message: 'ошибка'}}})

        await user.click(screen.getByRole('button', {name: 'Войти'}))

        await waitFor(() => {
            expect(showError).toHaveBeenCalled()
        })
    })
})
