import {api} from './api.js'

export const getMyRuns = (params) => {
    return api.get('/runs/my', {params})
}
