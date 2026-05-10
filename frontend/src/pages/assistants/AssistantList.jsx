import {useContext, useEffect, useMemo, useReducer, useState} from 'react'
import {useSearchParams} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {Checkbox} from 'primereact/checkbox'
import {Dropdown} from 'primereact/dropdown'
import {IconField} from 'primereact/iconfield'
import {InputIcon} from 'primereact/inputicon'
import {InputText} from 'primereact/inputtext'
import {Message} from 'primereact/message'
import {Paginator} from 'primereact/paginator'
import {ProgressSpinner} from 'primereact/progressspinner'
import {Tag} from 'primereact/tag'
import {AuthContext} from '../../context/AuthContext.jsx'
import {getAssistants, getCategories} from '../../services/catalog.js'
import {getTranslatedError} from '../../services/errors.js'
import '../../styles/Assistants.css'

const PAGE_SIZE = 9
const initialState = {
    assistants: [],
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
                assistants: action.assistants,
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

export const AssistantList = () => {
    const {user} = useContext(AuthContext)
    const [searchParams, setSearchParams] = useSearchParams()
    const [{assistants, pagination, loading, error}, dispatch] = useReducer(reducer, initialState)
    const [categories, setCategories] = useState([])
    const [searchDraft, setSearchDraft] = useState(searchParams.get('q') || '')

    const page = Number(searchParams.get('page') || 1)
    const q = searchParams.get('q') || ''
    const categoryId = searchParams.get('categoryId') || null
    const includeInactive = searchParams.get('includeInactive') === 'true'
    const isAdmin = user?.role === 'admin'

    const categoryOptions = useMemo(() => [
        {label: 'Все категории', value: null},
        ...categories.map(category => ({label: category.name, value: category.id})),
    ], [categories])

    useEffect(() => {
        getCategories().then(({data}) => {
            setCategories(data.categories || [])
        }).catch((err) => {
            dispatch({type: 'error', error: getTranslatedError(err)})
        })
    }, [])

    useEffect(() => {
        let cancelled = false
        dispatch({type: 'loading'})

        getAssistants({
            page,
            pageSize: PAGE_SIZE,
            q: q || undefined,
            categoryId: categoryId || undefined,
            includeInactive: isAdmin && includeInactive ? true : undefined,
        }).then(({data}) => {
            if (!cancelled) {
                dispatch({
                    type: 'success',
                    assistants: data.assistants || [],
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
    }, [categoryId, includeInactive, isAdmin, page, q])

    const updateParams = (changes) => {
        const next = new URLSearchParams(searchParams)

        Object.entries(changes).forEach(([key, value]) => {
            if (value === null || value === undefined || value === '' || value === false) {
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

    const applySearch = (event) => {
        event.preventDefault()
        updateParams({q: searchDraft.trim()})
    }

    const clearFilters = () => {
        setSearchDraft('')
        setSearchParams({})
    }

    const renderContent = () => {
        if (loading) {
            return (
                <div className="assistants-state">
                    <ProgressSpinner/>
                </div>
            )
        }

        if (error) {
            return <Message severity="error" text={error}/>
        }

        if (assistants.length === 0) {
            return <div className="assistants-state">Ассистенты не найдены</div>
        }

        return (
            <>
                <div className="assistants-grid">
                    {assistants.map(assistant => (
                        <Card key={assistant.id} className="assistant-card" title={assistant.name}>
                            <div className="assistant-card-meta">
                                {assistant.categoryName && <Tag value={assistant.categoryName}/>}
                                <Tag value={assistant.model} severity="secondary"/>
                                {isAdmin && !assistant.isActive && <Tag value="выключен" severity="danger"/>}
                            </div>
                            <p>{assistant.description}</p>
                            {assistant.exampleUserPrompt &&
                                <small>Пример контекста: {assistant.exampleUserPrompt}</small>
                            }
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
        <section className="assistants-page">
            <div className="assistants-header">
                <div>
                    <h1>Каталог ассистентов</h1>
                    <p>Выберите ассистента по категории или найдите по названию и описанию.</p>
                </div>
            </div>

            <form className="assistants-filters" onSubmit={applySearch}>
                <IconField iconPosition="left" className="assistants-search">
                    <InputIcon className="pi pi-search"/>
                    <InputText
                        id="assistants-search"
                        name="q"
                        value={searchDraft}
                        onChange={(event) => setSearchDraft(event.target.value)}
                        placeholder="Поиск"
                    />
                </IconField>

                <Dropdown
                    inputId="categoryId"
                    name="categoryId"
                    value={categoryId}
                    options={categoryOptions}
                    onChange={(event) => updateParams({categoryId: event.value})}
                    placeholder="Категория"
                />

                {isAdmin &&
                    <label className="assistants-checkbox">
                        <Checkbox
                            inputId="includeInactive"
                            name="includeInactive"
                            checked={includeInactive}
                            onChange={(event) => updateParams({includeInactive: event.checked})}
                        />
                        <span>Показать неактивных</span>
                    </label>
                }

                <Button type="submit" icon="pi pi-search" label="Найти"/>
                <Button type="button" outlined icon="pi pi-filter-slash" label="Сбросить" onClick={clearFilters}/>
            </form>

            {renderContent()}
        </section>
    )
}
