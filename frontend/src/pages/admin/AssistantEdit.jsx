import {useContext, useEffect, useMemo, useState} from 'react'
import {useNavigate, useParams} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {Checkbox} from 'primereact/checkbox'
import {Chips} from 'primereact/chips'
import {Dropdown} from 'primereact/dropdown'
import {InputText} from 'primereact/inputtext'
import {InputTextarea} from 'primereact/inputtextarea'
import {Message} from 'primereact/message'
import {ProgressSpinner} from 'primereact/progressspinner'
import {ToastContext} from '../../context/ToastContext.jsx'
import {getAssistant, getCategories, updateAssistant} from '../../services/catalog.js'
import {getTranslatedError} from '../../services/errors.js'
import {normalizeTags} from '../../services/tags.js'
import '../../styles/Admin.css'

const initialForm = {
    categoryId: null,
    name: '',
    description: '',
    model: '',
    systemPrompt: '',
    exampleUserPrompt: '',
    tags: [],
    isActive: true,
}

export const AssistantEdit = () => {
    const {assistantId} = useParams()
    const navigate = useNavigate()
    const {showSuccess} = useContext(ToastContext)
    const [form, setForm] = useState(initialForm)
    const [categories, setCategories] = useState([])
    const [error, setError] = useState(null)
    const [loading, setLoading] = useState(true)
    const [saving, setSaving] = useState(false)

    const categoryOptions = useMemo(() => categories.map(category => ({
        label: category.name,
        value: category.id,
    })), [categories])

    useEffect(() => {
        let cancelled = false

        Promise.all([
            getCategories(),
            getAssistant(assistantId),
        ]).then(([categoriesResponse, assistantResponse]) => {
            if (!cancelled) {
                const assistant = assistantResponse.data
                setCategories(categoriesResponse.data.categories || [])
                setForm({
                    categoryId: assistant.categoryId,
                    name: assistant.name || '',
                    description: assistant.description || '',
                    model: assistant.model || '',
                    systemPrompt: assistant.systemPrompt || '',
                    exampleUserPrompt: assistant.exampleUserPrompt || '',
                    tags: assistant.tags || [],
                    isActive: assistant.isActive,
                })
            }
        }).catch((err) => {
            if (!cancelled) {
                setError(getTranslatedError(err))
            }
        }).finally(() => {
            if (!cancelled) {
                setLoading(false)
            }
        })

        return () => {
            cancelled = true
        }
    }, [assistantId])

    const setField = (field, value) => {
        setForm(prev => ({...prev, [field]: value}))
    }

    const submitHandler = (event) => {
        event.preventDefault()

        if (!form.categoryId) {
            setError('Выберите категорию')
            return
        }

        if (!form.name.trim() || !form.description.trim() || !form.model.trim() || !form.systemPrompt.trim()) {
            setError('Заполните обязательные поля')
            return
        }

        setSaving(true)
        setError(null)

        updateAssistant(assistantId, {
            categoryId: form.categoryId,
            name: form.name.trim(),
            description: form.description.trim(),
            model: form.model.trim(),
            systemPrompt: form.systemPrompt.trim(),
            exampleUserPrompt: form.exampleUserPrompt.trim() || null,
            tags: normalizeTags(form.tags),
            isActive: form.isActive,
        }).then(() => {
            showSuccess('Ассистент обновлён')
            navigate(`/assistants/${assistantId}`)
        }).catch((err) => {
            setError(getTranslatedError(err))
        }).finally(() => {
            setSaving(false)
        })
    }

    if (loading) {
        return (
            <div className="admin-state">
                <ProgressSpinner/>
            </div>
        )
    }

    return (
        <section className="admin-page admin-page-wide">
            <Button outlined icon="pi pi-arrow-left" label="К карточке" onClick={() => navigate(`/assistants/${assistantId}`)}/>

            <Card title="Редактирование ассистента" className="admin-card">
                <form className="admin-form" onSubmit={submitHandler}>
                    {error && <Message severity="error" text={error}/>}

                    <div className="admin-form-grid">
                        <div className="admin-field">
                            <label htmlFor="assistantEditName">Название</label>
                            <InputText
                                id="assistantEditName"
                                name="name"
                                value={form.name}
                                onChange={(event) => setField('name', event.target.value)}
                            />
                        </div>

                        <div className="admin-field">
                            <label htmlFor="assistantEditCategory">Категория</label>
                            <Dropdown
                                inputId="assistantEditCategory"
                                name="categoryId"
                                value={form.categoryId}
                                options={categoryOptions}
                                onChange={(event) => setField('categoryId', event.value)}
                                placeholder="Категория"
                            />
                        </div>

                        <div className="admin-field">
                            <label htmlFor="assistantEditModel">Модель</label>
                            <InputText
                                id="assistantEditModel"
                                name="model"
                                value={form.model}
                                onChange={(event) => setField('model', event.target.value)}
                            />
                        </div>

                        <label className="admin-checkbox" htmlFor="assistantEditIsActive">
                            <Checkbox
                                inputId="assistantEditIsActive"
                                name="isActive"
                                checked={form.isActive}
                                onChange={(event) => setField('isActive', event.checked)}
                            />
                            <span>Активен</span>
                        </label>
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantEditTags">Теги</label>
                        <Chips
                            inputId="assistantEditTags"
                            name="tags"
                            value={form.tags}
                            onChange={(event) => setField('tags', event.value || [])}
                            separator=","
                            placeholder="Тег"
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantEditDescription">Описание</label>
                        <InputTextarea
                            id="assistantEditDescription"
                            name="description"
                            value={form.description}
                            onChange={(event) => setField('description', event.target.value)}
                            rows={4}
                            autoResize
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantEditSystemPrompt">Системный промпт</label>
                        <InputTextarea
                            id="assistantEditSystemPrompt"
                            name="systemPrompt"
                            value={form.systemPrompt}
                            onChange={(event) => setField('systemPrompt', event.target.value)}
                            rows={5}
                            autoResize
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantEditExampleUserPrompt">Пример пользовательского контекста</label>
                        <InputTextarea
                            id="assistantEditExampleUserPrompt"
                            name="exampleUserPrompt"
                            value={form.exampleUserPrompt}
                            onChange={(event) => setField('exampleUserPrompt', event.target.value)}
                            rows={4}
                            autoResize
                        />
                    </div>

                    <div className="admin-actions">
                        <Button type="submit" icon="pi pi-check" label="Сохранить" loading={saving}/>
                        <Button type="button" outlined label="Отмена" onClick={() => navigate(`/assistants/${assistantId}`)}/>
                    </div>
                </form>
            </Card>
        </section>
    )
}
