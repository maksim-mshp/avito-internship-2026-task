import {useEffect, useReducer, useState} from 'react'
import {useSearchParams} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {Dialog} from 'primereact/dialog'
import {Dropdown} from 'primereact/dropdown'
import {Message} from 'primereact/message'
import {Paginator} from 'primereact/paginator'
import {ProgressSpinner} from 'primereact/progressspinner'
import {RunRatingActions} from '../../components/RunRatingActions.jsx'
import {RunStatusTag, runStatusOptions} from '../../components/RunStatusTag.jsx'
import {getTranslatedError} from '../../services/errors.js'
import {getMyRuns, setRunRating} from '../../services/runs.js'
import '../../styles/Runs.css'

const PAGE_SIZE = 10

const initialState = {
    runs: [],
    pagination: {page: 1, pageSize: PAGE_SIZE, total: 0},
    loading: true,
    error: null,
}

const reducer = (state, action) => {
    switch (action.type) {
        case 'loading':
            return {...state, loading: true, error: null}
        case 'success':
            return {
                runs: action.runs,
                pagination: action.pagination,
                loading: false,
                error: null,
            }
        case 'error':
            return {...state, loading: false, error: action.error}
        case 'rating':
            return {
                ...state,
                runs: state.runs.map(run => run.id === action.run.id ? action.run : run),
            }
        default:
            return state
    }
}

const formatDate = (value) => {
    if (!value) {
        return '—'
    }

    return new Date(value).toLocaleString('ru-RU')
}

const shorten = (value, length = 140) => {
    if (!value) {
        return '—'
    }

    return value.length > length ? `${value.slice(0, length)}...` : value
}

export const MyRuns = () => {
    const [searchParams, setSearchParams] = useSearchParams()
    const [{runs, pagination, loading, error}, dispatch] = useReducer(reducer, initialState)
    const [selectedRun, setSelectedRun] = useState(null)
    const [ratingLoading, setRatingLoading] = useState(null)

    const page = Number(searchParams.get('page') || 1)
    const status = searchParams.get('status') || null

    useEffect(() => {
        let cancelled = false
        dispatch({type: 'loading'})

        getMyRuns({
            page,
            pageSize: PAGE_SIZE,
            status: status || undefined,
        }).then(({data}) => {
            if (!cancelled) {
                dispatch({
                    type: 'success',
                    runs: data.runs || [],
                    pagination: data.pagination || {page, pageSize: PAGE_SIZE, total: 0},
                })
            }
        }).catch((err) => {
            if (!cancelled) {
                dispatch({type: 'error', error: getTranslatedError(err)})
            }
        })

        return () => {
            cancelled = true
        }
    }, [page, status])

    const updateParams = (changes) => {
        const next = new URLSearchParams(searchParams)

        Object.entries(changes).forEach(([key, value]) => {
            if (value === null || value === undefined || value === '') {
                next.delete(key)
            } else {
                next.set(key, String(value))
            }
        })

        if (!Object.prototype.hasOwnProperty.call(changes, 'page')) {
            next.delete('page')
        }

        setSearchParams(next)
    }

    const rateRun = (run, rating) => {
        setRatingLoading(`${run.id}:${rating}`)

        setRunRating(run.id, rating).then(({data}) => {
            dispatch({type: 'rating', run: data})
            setSelectedRun(prev => prev?.id === data.id ? data : prev)
        }).catch((err) => {
            dispatch({type: 'error', error: getTranslatedError(err)})
        }).finally(() => {
            setRatingLoading(null)
        })
    }

    const renderContent = () => {
        if (loading) {
            return (
                <div className="runs-state">
                    <ProgressSpinner/>
                </div>
            )
        }

        if (error) {
            return <Message severity="error" text={error}/>
        }

        if (runs.length === 0) {
            return <div className="runs-state">Запусков пока нет</div>
        }

        return (
            <>
                <div className="runs-list">
                    {runs.map(run => (
                        <Card key={run.id} className="run-card">
                            <div className="run-card-header">
                                <div>
                                    <h2>{run.assistantName || 'Ассистент'}</h2>
                                    <span>{formatDate(run.createdAt)}</span>
                                </div>
                                <RunStatusTag status={run.status}/>
                            </div>
                            <div className="run-card-body">
                                <div>
                                    <span>Пользовательский контекст</span>
                                    <p>{shorten(run.userPrompt)}</p>
                                </div>
                                <div>
                                    <span>Ответ</span>
                                    <p>{shorten(run.output || run.error)}</p>
                                </div>
                            </div>
                            <div className="run-card-actions">
                                {run.status === 'success' && run.output &&
                                    <RunRatingActions
                                        rating={run.rating}
                                        loading={ratingLoading?.startsWith(`${run.id}:`) ? ratingLoading.split(':')[1] : null}
                                        onRate={(rating) => rateRun(run, rating)}
                                    />
                                }
                                <Button
                                    outlined
                                    icon="pi pi-external-link"
                                    label="Открыть полный ответ"
                                    onClick={() => setSelectedRun(run)}
                                />
                            </div>
                        </Card>
                    ))}
                </div>
                <Paginator
                    first={(pagination.page - 1) * pagination.pageSize}
                    rows={pagination.pageSize}
                    totalRecords={pagination.total}
                    onPageChange={(event) => updateParams({page: event.page + 1})}
                />
            </>
        )
    }

    return (
        <section className="runs-page">
            <div className="runs-header">
                <div>
                    <h1>Мои запуски</h1>
                    <p>История обращений к ассистентам с результатами и ошибками.</p>
                </div>
                <Dropdown
                    inputId="runStatus"
                    name="status"
                    value={status}
                    options={runStatusOptions}
                    onChange={(event) => updateParams({status: event.value})}
                    placeholder="Статус"
                />
            </div>

            {renderContent()}

            <Dialog
                header={selectedRun?.assistantName || 'Запуск ассистента'}
                visible={Boolean(selectedRun)}
                onHide={() => setSelectedRun(null)}
                className="run-dialog"
                modal
            >
                {selectedRun &&
                    <div className="run-dialog-content">
                        <RunStatusTag status={selectedRun.status}/>
                        <div>
                            <span>Пользовательский контекст</span>
                            <p>{selectedRun.userPrompt}</p>
                        </div>
                        <div>
                            <span>Ответ</span>
                            <p>{selectedRun.output || selectedRun.error || '—'}</p>
                        </div>
                        {selectedRun.status === 'success' && selectedRun.output &&
                            <RunRatingActions
                                rating={selectedRun.rating}
                                loading={
                                    ratingLoading?.startsWith(`${selectedRun.id}:`)
                                        ? ratingLoading.split(':')[1]
                                        : null
                                }
                                onRate={(rating) => rateRun(selectedRun, rating)}
                            />
                        }
                    </div>
                }
            </Dialog>
        </section>
    )
}
