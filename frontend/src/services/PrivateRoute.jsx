import {useContext} from 'react'
import {Navigate} from 'react-router-dom'
import {AuthContext} from '../context/AuthContext.jsx'
import {Header} from '../components/Header.jsx'

export const PrivateRoute = ({children, role}) => {
    const {isLoggedIn, user} = useContext(AuthContext)

    if (!isLoggedIn) {
        return <Navigate to="/login" replace/>
    }

    if (role && user?.role !== role) {
        return <Navigate to="/assistants" replace/>
    }

    return (
        <>
            <Header/>
            <main>
                {children}
            </main>
        </>
    )
}
