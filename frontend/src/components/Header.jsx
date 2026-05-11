import {useContext} from 'react'
import {useLocation, useNavigate} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Tag} from 'primereact/tag'
import {AuthContext} from '../context/AuthContext.jsx'
import {ThemeContext} from '../context/ThemeContext.jsx'
import '../styles/Header.css'

export const Header = () => {
    const location = useLocation()
    const navigate = useNavigate()
    const {logout, user} = useContext(AuthContext)
    const {isDark, toggleTheme} = useContext(ThemeContext)

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
                {user?.role === 'admin' &&
                    <>
                        <Button
                            label="Все запуски"
                            severity="secondary"
                            text
                            className={location.pathname === '/admin/runs' ? 'active' : ''}
                            onClick={() => navigate('/admin/runs')}
                        />
                        <Button
                            label="Новый ассистент"
                            severity="secondary"
                            text
                            className={location.pathname === '/admin/assistants/new' ? 'active' : ''}
                            onClick={() => navigate('/admin/assistants/new')}
                        />
                        <Button
                            label="Новая категория"
                            severity="secondary"
                            text
                            className={location.pathname === '/admin/categories/new' ? 'active' : ''}
                            onClick={() => navigate('/admin/categories/new')}
                        />
                    </>
                }
            </div>
            <div className="flex align-items-center gap-2">
                <Button
                    className="theme-toggle"
                    outlined
                    rounded
                    icon={isDark ? 'pi pi-sun' : 'pi pi-moon'}
                    aria-label={isDark ? 'Включить светлую тему' : 'Включить тёмную тему'}
                    onClick={toggleTheme}
                />
                {user?.role && <Tag value={user.role}/>}
                <Button className="p-button-danger" outlined icon="pi pi-sign-out" onClick={logoutHandler}/>
            </div>
        </header>
    )
}
