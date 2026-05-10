import {StrictMode} from 'react'
import {createRoot} from 'react-dom/client'
import {addLocale} from 'primereact/api'
import {ru} from 'primelocale/js/ru.js'

import 'primereact/resources/themes/lara-light-indigo/theme.css'
import 'primereact/resources/primereact.min.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css'

import './index.css'
import {Providers} from './services/Providers.jsx'
import {Router} from './services/Router.jsx'

addLocale('ru', ru)

createRoot(document.getElementById('root')).render(
    <StrictMode>
        <Providers>
            <Router/>
        </Providers>
    </StrictMode>,
)
