import {useContext, useEffect, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Dropdown} from 'primereact/dropdown'
import {InputText} from 'primereact/inputtext'
import {Password} from 'primereact/password'
import {AuthContext} from '../../context/AuthContext.jsx'
import {ToastContext} from '../../context/ToastContext.jsx'
import {api} from '../../services/api.js'
import {getTranslatedError} from '../../services/errors.js'
import '../../styles/Auth.css'

const roleOptions = [
    {label: 'admin', value: 'admin'},
    {label: 'user', value: 'user'},
]

export const Login = () => {
    const navigate = useNavigate()
    const {isLoggedIn, login} = useContext(AuthContext)
    const {showError} = useContext(ToastContext)

    const [mode, setMode] = useState('login')
    const [role, setRole] = useState('user')
    const [dummyRole, setDummyRole] = useState('user')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        if (isLoggedIn) {
            navigate('/', {replace: true})
        }
    }, [isLoggedIn, navigate])

    const finishLogin = (data) => {
        login(data)
        navigate('/assistants', {replace: true})
    }

    const dummyLoginHandler = (e) => {
        e.preventDefault()
        setLoading(true)

        api.post('/dummyLogin', {role: dummyRole}).then(({data}) => {
            finishLogin(data)
        }).catch((error) => {
            showError(getTranslatedError(error))
        }).finally(() => {
            setLoading(false)
        })
    }

    const passwordLoginHandler = (e) => {
        e.preventDefault()
        setLoading(true)

        const payload = {
            email,
            password,
        }

        const request = mode === 'register'
            ? api.post('/register', {...payload, role}).then(() => api.post('/login', payload))
            : api.post('/login', payload)

        request.then(({data}) => {
            finishLogin(data)
        }).catch((error) => {
            showError(getTranslatedError(error))
        }).finally(() => {
            setLoading(false)
        })
    }

    return (
        <section className="auth">
            <form onSubmit={passwordLoginHandler}>
                <h1>Вход</h1>

                <div className="auth-mode">
                    <Button
                        type="button"
                        label="Вход по email"
                        severity={mode === 'login' ? 'primary' : 'secondary'}
                        outlined={mode !== 'login'}
                        onClick={() => setMode('login')}
                    />
                    <Button
                        type="button"
                        label="Регистрация"
                        severity={mode === 'register' ? 'primary' : 'secondary'}
                        outlined={mode !== 'register'}
                        onClick={() => setMode('register')}
                    />
                </div>

                <div className="flex flex-column gap-2 mb-3">
                    <label htmlFor="email">Email</label>
                    <InputText
                        id="email"
                        name="email"
                        type="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        autoComplete="email"
                    />
                </div>

                <div className="flex flex-column gap-2 mb-3">
                    <label htmlFor="password">Пароль</label>
                    <Password
                        inputId="password"
                        name="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        feedback={false}
                        toggleMask
                        autoComplete={mode === 'register' ? 'new-password' : 'current-password'}
                    />
                </div>

                {mode === 'register' &&
                    <div className="flex flex-column gap-2 mb-3">
                        <label htmlFor="registerRole">Роль</label>
                        <Dropdown
                            inputId="registerRole"
                            name="role"
                            value={role}
                            options={roleOptions}
                            onChange={(e) => setRole(e.value)}
                        />
                    </div>
                }

                <Button type="submit" loading={loading}>
                    {mode === 'register' ? 'Зарегистрироваться' : 'Войти по email'}
                </Button>

                <div className="auth-divider">или</div>

                <div className="flex flex-column gap-2 mb-3">
                    <label htmlFor="dummyRole">Тестовая роль</label>
                    <Dropdown
                        inputId="dummyRole"
                        name="dummyRole"
                        value={dummyRole}
                        options={roleOptions}
                        onChange={(e) => setDummyRole(e.value)}
                    />
                </div>

                <Button type="button" outlined onClick={dummyLoginHandler} loading={loading}>
                    Войти тестовым пользователем
                </Button>
            </form>
        </section>
    )
}
