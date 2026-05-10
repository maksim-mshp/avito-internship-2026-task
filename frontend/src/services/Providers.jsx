import {PrimeReactProvider} from 'primereact/api'
import {AuthProvider} from '../context/AuthContext.jsx'
import {ToastProvider} from '../context/ToastContext.jsx'

export const Providers = ({children}) => {
    return (
        <PrimeReactProvider value={{locale: 'ru'}}>
            <AuthProvider>
                <ToastProvider>
                    {children}
                </ToastProvider>
            </AuthProvider>
        </PrimeReactProvider>
    )
}
