import {useEffect, useMemo, useReducer, useState} from 'react'
import {useSearchParams} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {Dialog} from 'primereact/dialog'
import {Dropdown} from 'primereact/dropdown'
import {Message} from 'primereact/message'
import {Paginator} from 'primereact/paginator'
import {ProgressSpinner} from 'primereact/progressspinner'
import {RunStatusTag, runStatusOptions} from '../../components/RunStatusTag.jsx'
import {getAssistants} from '../../services/catalog.js'
import {getTranslatedError} from '../../services/errors.js'
import {getAdminRuns} from '../../services/runs.js'
import '../../styles/Runs.css'

const PAGE_SIZE = 20

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

export const AdminRuns = () => {
    const [searchParams, setSearchParams] = useSearchParams()
    const [{runs, pagination, loading, error}, dispatch] = useReducer(reducer, initialState)
    const [assistants, setAssistants] = useState([])
    const [selectedRun, setSelectedRun] = useState(null)

    const page = Number(searchParams.get('page') || 1)
    const status = searchParams.get('status') || null
    const assistantId = searchParams.get('assistantId') || null

    const assistantOptions = useMemo(() => [
        {label: 'Все ассистенты', value: null},
        ...assistants.map(assistant => ({label: assistant.name, value: assistant.id})),
    ], [assistants])

    useEffect(() => {
        let cancelled = false

        getAssistants({
            page: 1,
            pageSize: 100,
            includeInactive: true,
        }).then(({data}) => {
            if (!cancelled) {
                setAssistants(data.assistants || [])
            }
        }).catch(() => {
            if (!cancelled) {
                setAssistants([])
            }
        })

        return () => {
            cancelled = true
        }
    }, [])

    useEffect(() => {
        let cancelled = false
        dispatch({type: 'loading'})

        getAdminRuns({
            page,
            pageSize: PAGE_SIZE,
            status: status || undefined,
            assistantId: assistantId || undefined,
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
    }, [assistantId, page, status])

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

    const clearFilters = () => {
        setSearchParams({})
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
            return <div className="runs-state">Запусков не найдено</div>
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
                            <div className="run-card-body run-card-body-admin">
                                <div>
                                    <span>Пользователь</span>
                                    <p>{run.userId}</p>
                                </div>
                                <div>
                                    <span>Пользовательский контекст</span>
                                    <p>{shorten(run.userPrompt)}</p>
                                </div>
                                <div>
                                    <span>Ответ</span>
                                    <p>{shorten(run.output || run.error)}</p>
                                </div>
                            </div>
                            <Button
                                outlined
                                icon="pi pi-external-link"
                                label="Открыть полный ответ"
                                onClick={() => setSelectedRun(run)}
                            />
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
                    <h1>Все запуски</h1>
                    <p>Административный список запусков ассистентов с фильтрами.</p>
                </div>
                <div className="runs-filters">
                    <Dropdown
                        inputId="adminRunAssistant"
                        name="assistantId"
                        value={assistantId}
                        options={assistantOptions}
                        onChange={(event) => updateParams({assistantId: event.value})}
                        placeholder="Ассистент"
                    />
                    <Dropdown
                        inputId="adminRunStatus"
                        name="status"
                        value={status}
                        options={runStatusOptions}
                        onChange={(event) => updateParams({status: event.value})}
                        placeholder="Статус"
                    />
                    <Button type="button" outlined icon="pi pi-filter-slash" label="Сбросить" onClick={clearFilters}/>
                </div>
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
                            <span>Пользователь</span>
                            <p>{selectedRun.userId}</p>
                        </div>
                        <div>
                            <span>Пользовательский контекст</span>
                            <p>{selectedRun.userPrompt}</p>
                        </div>
                        <div>
                            <span>Ответ</span>
                            <p>{selectedRun.output || selectedRun.error || '—'}</p>
                        </div>
                    </div>
                }
            </Dialog>
        </section>
    )
}
