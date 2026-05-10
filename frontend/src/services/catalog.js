import {api} from './api.js'

export const getCategories = () => {
    return api.get('/categories')
}

export const getAssistants = (params) => {
    return api.get('/assistants', {params})
}
