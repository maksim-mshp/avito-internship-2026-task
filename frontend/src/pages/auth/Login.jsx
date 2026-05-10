import {useContext, useEffect, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Dropdown} from 'primereact/dropdown'
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

    const [role, setRole] = useState('user')
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        if (isLoggedIn) {
            navigate('/', {replace: true})
        }
    }, [isLoggedIn, navigate])

    const loginHandler = (e) => {
        e.preventDefault()
        setLoading(true)

        api.post('/dummyLogin', {role}).then(({data}) => {
            login(data)
            navigate('/assistants', {replace: true})
        }).catch((error) => {
            showError(getTranslatedError(error))
        }).finally(() => {
            setLoading(false)
        })
    }

    return (
        <section className="auth">
            <form>
                <h1>Вход</h1>
                <div className="flex flex-column gap-2 mb-3">
                    <label htmlFor="role">Роль</label>
                    <Dropdown
                        inputId="role"
                        name="role"
                        value={role}
                        options={roleOptions}
                        onChange={(e) => setRole(e.value)}
                    />
                </div>
                <Button onClick={loginHandler} loading={loading}>Войти</Button>
            </form>
        </section>
    )
}
