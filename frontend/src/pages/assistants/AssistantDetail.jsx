import {useContext, useEffect, useReducer, useState} from 'react'
import {useNavigate, useParams} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {InputTextarea} from 'primereact/inputtextarea'
import {Message} from 'primereact/message'
import {ProgressSpinner} from 'primereact/progressspinner'
import {Tag} from 'primereact/tag'
import {RunRatingActions} from '../../components/RunRatingActions.jsx'
import {AuthContext} from '../../context/AuthContext.jsx'
import {addFavoriteAssistant, getAssistant, removeFavoriteAssistant, runAssistant} from '../../services/catalog.js'
import {getTranslatedError} from '../../services/errors.js'
import {setRunRating} from '../../services/runs.js'
import '../../styles/Assistants.css'

const initialState = {
    assistant: null,
    loading: true,
    error: null,
}

const reducer = (state, action) => {
    switch (action.type) {
        case 'loading':
            return {...state, loading: true, error: null}
        case 'success':
            return {assistant: action.assistant, loading: false, error: null}
        case 'favorite':
            return {...state, assistant: {...state.assistant, isFavorite: action.isFavorite}}
        case 'error':
            return {...state, loading: false, error: action.error}
        default:
            return state
    }
}

export const AssistantDetail = () => {
    const {assistantId} = useParams()
    const navigate = useNavigate()
    const {user} = useContext(AuthContext)
    const [{assistant, loading, error}, dispatch] = useReducer(reducer, initialState)
    const [userPrompt, setUserPrompt] = useState('')
    const [run, setRun] = useState(null)
    const [runError, setRunError] = useState(null)
    const [submitting, setSubmitting] = useState(false)
    const [favoriteLoading, setFavoriteLoading] = useState(false)
    const [ratingLoading, setRatingLoading] = useState(null)

    useEffect(() => {
        let cancelled = false
        dispatch({type: 'loading'})

        getAssistant(assistantId).then(({data}) => {
            if (!cancelled) {
                dispatch({type: 'success', assistant: data})
                setUserPrompt(data.exampleUserPrompt || '')
            }
        }).catch((err) => {
            if (!cancelled) {
                dispatch({type: 'error', error: getTranslatedError(err)})
            }
        })

        return () => {
            cancelled = true
        }
    }, [assistantId])

    const submitHandler = (event) => {
        event.preventDefault()
        const prompt = userPrompt.trim()

        if (!prompt) {
            setRunError('Введите пользовательский контекст')
            return
        }

        setSubmitting(true)
        setRun(null)
        setRunError(null)

        runAssistant(assistantId, prompt).then(({data}) => {
            setRun(data)
            if (data.status === 'failed') {
                setRunError(data.error || 'Запуск завершился ошибкой')
            }
        }).catch((err) => {
            setRunError(getTranslatedError(err))
        }).finally(() => {
            setSubmitting(false)
        })
    }

    const toggleFavorite = () => {
        setFavoriteLoading(true)

        const request = assistant.isFavorite
            ? removeFavoriteAssistant(assistant.id)
            : addFavoriteAssistant(assistant.id)

        request.then(() => {
            dispatch({type: 'favorite', isFavorite: !assistant.isFavorite})
        }).catch((err) => {
            setRunError(getTranslatedError(err))
        }).finally(() => {
            setFavoriteLoading(false)
        })
    }

    const rateRun = (rating) => {
        if (!run?.id) {
            return
        }

        setRatingLoading(rating)
        setRunError(null)

        setRunRating(run.id, rating).then(({data}) => {
            setRun(data)
        }).catch((err) => {
            setRunError(getTranslatedError(err))
        }).finally(() => {
            setRatingLoading(null)
        })
    }

    if (loading) {
        return (
            <div className="assistants-state">
                <ProgressSpinner/>
            </div>
        )
    }

    if (error) {
        return (
            <section className="assistant-detail">
                <Button outlined icon="pi pi-arrow-left" label="К каталогу" onClick={() => navigate('/assistants')}/>
                <Message severity="error" text={error}/>
            </section>
        )
    }

    return (
        <section className="assistant-detail">
            <div className="assistant-detail-actions">
                <Button outlined icon="pi pi-arrow-left" label="К каталогу" onClick={() => navigate('/assistants')}/>
                {user?.role === 'admin' &&
                    <Button
                        outlined
                        icon="pi pi-pencil"
                        label="Редактировать"
                        onClick={() => navigate(`/admin/assistants/${assistant.id}/edit`)}
                    />
                }
                <Button
                    outlined
                    severity={assistant.isFavorite ? 'danger' : 'secondary'}
                    icon={assistant.isFavorite ? 'pi pi-heart-fill' : 'pi pi-heart'}
                    label={assistant.isFavorite ? 'В избранном' : 'В избранное'}
                    loading={favoriteLoading}
                    onClick={toggleFavorite}
                />
            </div>

            <Card title={assistant.name} className="assistant-detail-card">
                <div className="assistant-card-meta">
                    {assistant.categoryName && <Tag value={assistant.categoryName}/>}
                    <Tag value={assistant.model} severity="secondary"/>
                    {!assistant.isActive && <Tag value="выключен" severity="danger"/>}
                </div>

                {assistant.tags?.length > 0 &&
                    <div className="assistant-tags">
                        {assistant.tags.map(tag => (
                            <Tag key={tag} value={`#${tag}`} severity="info"/>
                        ))}
                    </div>
                }

                <p>{assistant.description}</p>

                {assistant.exampleUserPrompt &&
                    <div className="assistant-example">
                        <span>Пример контекста</span>
                        <p>{assistant.exampleUserPrompt}</p>
                    </div>
                }
            </Card>

            <Card title="Запуск ассистента" className="assistant-detail-card">
                <form className="assistant-run-form" onSubmit={submitHandler}>
                    <label htmlFor="userPrompt">Пользовательский контекст</label>
                    <InputTextarea
                        id="userPrompt"
                        name="userPrompt"
                        value={userPrompt}
                        onChange={(event) => setUserPrompt(event.target.value)}
                        rows={7}
                        autoResize
                    />
                    <Button type="submit" icon="pi pi-send" label="Запустить" loading={submitting}/>
                </form>
            </Card>

            {runError && <Message severity="error" text={runError}/>}

            {run?.output &&
                <Card title="Ответ ассистента" className="assistant-detail-card">
                    <p className="assistant-output">{run.output}</p>
                    <RunRatingActions rating={run.rating} loading={ratingLoading} onRate={rateRun}/>
                </Card>
            }
        </section>
    )
}
