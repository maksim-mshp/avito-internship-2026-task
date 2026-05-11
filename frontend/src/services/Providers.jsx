import {PrimeReactProvider} from 'primereact/api'
import {AuthProvider} from '../context/AuthContext.jsx'
import {ThemeProvider} from '../context/ThemeContext.jsx'
import {ToastProvider} from '../context/ToastContext.jsx'

export const Providers = ({children}) => {
    return (
        <PrimeReactProvider value={{locale: 'ru'}}>
            <ThemeProvider>
                <AuthProvider>
                    <ToastProvider>
                        {children}
                    </ToastProvider>
                </AuthProvider>
            </ThemeProvider>
        </PrimeReactProvider>
    )
}
