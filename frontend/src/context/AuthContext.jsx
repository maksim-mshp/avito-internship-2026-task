import {createContext, useMemo, useState} from 'react'

export const AuthContext = createContext(null)

const getStoredUser = () => {
    const storedUser = localStorage.getItem('user')

    return storedUser ? JSON.parse(storedUser) : null
}

export const AuthProvider = ({children}) => {
    const [token, setToken] = useState(() => localStorage.getItem('token'))
    const [user, setUser] = useState(getStoredUser)

    const login = ({token, user}) => {
        setToken(token)
        setUser(user)
        localStorage.setItem('token', token)
        localStorage.setItem('user', JSON.stringify(user))
    }

    const logout = () => {
        setToken(null)
        setUser(null)
        localStorage.removeItem('token')
        localStorage.removeItem('user')
    }

    const value = useMemo(() => ({
        token,
        user,
        isLoggedIn: Boolean(token),
        login,
        logout,
    }), [token, user])

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}
