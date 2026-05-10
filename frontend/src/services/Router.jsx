import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'
import {Login} from '../pages/auth/Login.jsx'
import {AssistantList} from '../pages/assistants/AssistantList.jsx'
import {PrivateRoute} from './PrivateRoute.jsx'

export const Router = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<Login/>}/>
                <Route path="/assistants" element={<PrivateRoute><AssistantList/></PrivateRoute>}/>
                <Route path="*" element={<Navigate to="/assistants" replace/>}/>
            </Routes>
        </BrowserRouter>
    )
}
