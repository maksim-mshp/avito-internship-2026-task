import {createContext, useRef} from 'react'
import {Toast} from 'primereact/toast'

export const ToastContext = createContext(null)

export const ToastProvider = ({children}) => {
    const toast = useRef(null)

    const showError = (error) => {
        toast.current.show({
            severity: 'error',
            detail: error,
            life: 3000,
        })
    }

    const showWarning = (error) => {
        toast.current.show({
            severity: 'warn',
            detail: error,
            life: 3000,
        })
    }

    const showSuccess = (message) => {
        toast.current.show({
            severity: 'success',
            detail: message,
            life: 3000,
        })
    }

    return (
        <ToastContext.Provider value={{toast, showError, showWarning, showSuccess}}>
            <Toast ref={toast} position="top-center"/>
            {children}
        </ToastContext.Provider>
    )
}
