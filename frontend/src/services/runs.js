import {api} from './api.js'

export const getMyRuns = (params) => {
    return api.get('/runs/my', {params})
}

export const getAdminRuns = (params) => {
    return api.get('/admin/runs', {params})
}
