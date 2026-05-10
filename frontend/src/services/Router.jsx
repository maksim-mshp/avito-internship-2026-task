import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'
import {Login} from '../pages/auth/Login.jsx'
import {AssistantDetail} from '../pages/assistants/AssistantDetail.jsx'
import {AssistantList} from '../pages/assistants/AssistantList.jsx'
import {AssistantCreate} from '../pages/admin/AssistantCreate.jsx'
import {AssistantEdit} from '../pages/admin/AssistantEdit.jsx'
import {AdminRuns} from '../pages/admin/AdminRuns.jsx'
import {CategoryCreate} from '../pages/admin/CategoryCreate.jsx'
import {MyRuns} from '../pages/runs/MyRuns.jsx'
import {PrivateRoute} from './PrivateRoute.jsx'

export const Router = () => {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<Login/>}/>
                <Route path="/assistants" element={<PrivateRoute><AssistantList/></PrivateRoute>}/>
                <Route path="/assistants/:assistantId" element={<PrivateRoute><AssistantDetail/></PrivateRoute>}/>
                <Route path="/runs/my" element={<PrivateRoute><MyRuns/></PrivateRoute>}/>
                <Route path="/admin/runs" element={<PrivateRoute role="admin"><AdminRuns/></PrivateRoute>}/>
                <Route path="/admin/assistants/new" element={<PrivateRoute role="admin"><AssistantCreate/></PrivateRoute>}/>
                <Route
                    path="/admin/assistants/:assistantId/edit"
                    element={<PrivateRoute role="admin"><AssistantEdit/></PrivateRoute>}
                />
                <Route path="/admin/categories/new" element={<PrivateRoute role="admin"><CategoryCreate/></PrivateRoute>}/>
                <Route path="*" element={<Navigate to="/assistants" replace/>}/>
            </Routes>
        </BrowserRouter>
    )
}
