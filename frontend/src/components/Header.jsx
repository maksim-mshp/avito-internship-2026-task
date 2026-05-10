import {useContext} from 'react'
import {useLocation, useNavigate} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Tag} from 'primereact/tag'
import {AuthContext} from '../context/AuthContext.jsx'
import '../styles/Header.css'

export const Header = () => {
    const location = useLocation()
    const navigate = useNavigate()
    const {logout, user} = useContext(AuthContext)

    const logoutHandler = () => {
        logout()
        navigate('/login')
    }

    return (
        <header className="app-header flex justify-content-between align-items-center p-3 pl-4 pr-4">
            <div className="menu">
                <Button
                    label="Ассистенты"
                    severity="secondary"
                    text
                    className={location.pathname.startsWith('/assistants') ? 'active' : ''}
                    onClick={() => navigate('/assistants')}
                />
                <Button
                    label="Мои запуски"
                    severity="secondary"
                    text
                    className={location.pathname === '/runs/my' ? 'active' : ''}
                    onClick={() => navigate('/runs/my')}
                />
            </div>
            <div className="flex align-items-center gap-2">
                {user?.role && <Tag value={user.role}/>}
                <Button className="p-button-danger" outlined icon="pi pi-sign-out" onClick={logoutHandler}/>
            </div>
        </header>
    )
}
