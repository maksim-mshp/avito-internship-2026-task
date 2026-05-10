import {api} from './api.js'

export const getCategories = () => {
    return api.get('/categories')
}

export const createCategory = (data) => {
    return api.post('/categories', data)
}

export const getAssistants = (params) => {
    return api.get('/assistants', {params})
}

export const createAssistant = (data) => {
    return api.post('/assistants', data)
}

export const updateAssistant = (assistantId, data) => {
    return api.put(`/assistants/${assistantId}`, data)
}

export const getAssistant = (assistantId) => {
    return api.get(`/assistants/${assistantId}`)
}

export const runAssistant = (assistantId, userPrompt) => {
    return api.post(`/assistants/${assistantId}/run`, {userPrompt})
}
