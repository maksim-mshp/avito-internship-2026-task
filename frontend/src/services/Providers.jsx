import {PrimeReactProvider} from 'primereact/api'

export const Providers = ({children}) => {
    return (
        <PrimeReactProvider value={{locale: 'ru'}}>
            {children}
        </PrimeReactProvider>
    )
}
