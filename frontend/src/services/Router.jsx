import {BrowserRouter, Route, Routes} from 'react-router-dom'
import App from '../App.jsx'
import {Login} from '../pages/auth/Login.jsx'

export const Router = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<Login/>}/>
                <Route path="*" element={<App/>}/>
            </Routes>
        </BrowserRouter>
    )
}
