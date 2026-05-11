import {createContext, useEffect, useMemo, useState} from 'react'

const STORAGE_KEY = 'theme'
const DARK_THEME = 'dark'
const LIGHT_THEME = 'light'

const getInitialTheme = () => {
    const storedTheme = localStorage.getItem(STORAGE_KEY)
    return storedTheme === DARK_THEME ? DARK_THEME : LIGHT_THEME
}

export const ThemeContext = createContext({
    theme: LIGHT_THEME,
    isDark: false,
    toggleTheme: () => {},
})

export const ThemeProvider = ({children}) => {
    const [theme, setTheme] = useState(getInitialTheme)

    useEffect(() => {
        document.documentElement.dataset.theme = theme
        localStorage.setItem(STORAGE_KEY, theme)
    }, [theme])

    const value = useMemo(() => ({
        theme,
        isDark: theme === DARK_THEME,
        toggleTheme: () => setTheme(current => current === DARK_THEME ? LIGHT_THEME : DARK_THEME),
    }), [theme])

    return (
        <ThemeContext.Provider value={value}>
            {children}
        </ThemeContext.Provider>
    )
}
