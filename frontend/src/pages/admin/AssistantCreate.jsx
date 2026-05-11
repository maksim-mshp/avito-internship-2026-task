import {useContext, useEffect, useMemo, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {Checkbox} from 'primereact/checkbox'
import {Chips} from 'primereact/chips'
import {Dropdown} from 'primereact/dropdown'
import {InputText} from 'primereact/inputtext'
import {InputTextarea} from 'primereact/inputtextarea'
import {Message} from 'primereact/message'
import {ToastContext} from '../../context/ToastContext.jsx'
import {createAssistant, getCategories} from '../../services/catalog.js'
import {getTranslatedError} from '../../services/errors.js'
import {normalizeTags} from '../../services/tags.js'
import '../../styles/Admin.css'

const initialForm = {
    categoryId: null,
    name: '',
    description: '',
    model: 'mock-smart',
    systemPrompt: '',
    exampleUserPrompt: '',
    tags: [],
    isActive: true,
}

export const AssistantCreate = () => {
    const navigate = useNavigate()
    const {showSuccess} = useContext(ToastContext)
    const [form, setForm] = useState(initialForm)
    const [categories, setCategories] = useState([])
    const [error, setError] = useState(null)
    const [loading, setLoading] = useState(false)

    const categoryOptions = useMemo(() => categories.map(category => ({
        label: category.name,
        value: category.id,
    })), [categories])

    useEffect(() => {
        let cancelled = false

        getCategories().then(({data}) => {
            if (!cancelled) {
                setCategories(data.categories || [])
            }
        }).catch((err) => {
            if (!cancelled) {
                setError(getTranslatedError(err))
            }
        })

        return () => {
            cancelled = true
        }
    }, [])

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

        setLoading(true)
        setError(null)

        createAssistant({
            categoryId: form.categoryId,
            name: form.name.trim(),
            description: form.description.trim(),
            model: form.model.trim(),
            systemPrompt: form.systemPrompt.trim(),
            exampleUserPrompt: form.exampleUserPrompt.trim() || null,
            tags: normalizeTags(form.tags),
            isActive: form.isActive,
        }).then(({data}) => {
            showSuccess('Ассистент создан')
            navigate(`/assistants/${data.id}`)
        }).catch((err) => {
            setError(getTranslatedError(err))
        }).finally(() => {
            setLoading(false)
        })
    }

    return (
        <section className="admin-page admin-page-wide">
            <Button outlined icon="pi pi-arrow-left" label="К каталогу" onClick={() => navigate('/assistants')}/>

            <Card title="Новый ассистент" className="admin-card">
                <form className="admin-form" onSubmit={submitHandler}>
                    {error && <Message severity="error" text={error}/>}

                    <div className="admin-form-grid">
                        <div className="admin-field">
                            <label htmlFor="assistantName">Название</label>
                            <InputText
                                id="assistantName"
                                name="name"
                                value={form.name}
                                onChange={(event) => setField('name', event.target.value)}
                            />
                        </div>

                        <div className="admin-field">
                            <label htmlFor="assistantCategory">Категория</label>
                            <Dropdown
                                inputId="assistantCategory"
                                name="categoryId"
                                value={form.categoryId}
                                options={categoryOptions}
                                onChange={(event) => setField('categoryId', event.value)}
                                placeholder="Категория"
                            />
                        </div>

                        <div className="admin-field">
                            <label htmlFor="assistantModel">Модель</label>
                            <InputText
                                id="assistantModel"
                                name="model"
                                value={form.model}
                                onChange={(event) => setField('model', event.target.value)}
                            />
                        </div>

                        <label className="admin-checkbox" htmlFor="assistantIsActive">
                            <Checkbox
                                inputId="assistantIsActive"
                                name="isActive"
                                checked={form.isActive}
                                onChange={(event) => setField('isActive', event.checked)}
                            />
                            <span>Активен</span>
                        </label>
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantTags">Теги</label>
                        <Chips
                            inputId="assistantTags"
                            name="tags"
                            value={form.tags}
                            onChange={(event) => setField('tags', event.value || [])}
                            separator=","
                            placeholder="Тег"
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantDescription">Описание</label>
                        <InputTextarea
                            id="assistantDescription"
                            name="description"
                            value={form.description}
                            onChange={(event) => setField('description', event.target.value)}
                            rows={4}
                            autoResize
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantSystemPrompt">Системный промпт</label>
                        <InputTextarea
                            id="assistantSystemPrompt"
                            name="systemPrompt"
                            value={form.systemPrompt}
                            onChange={(event) => setField('systemPrompt', event.target.value)}
                            rows={5}
                            autoResize
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="assistantExampleUserPrompt">Пример пользовательского контекста</label>
                        <InputTextarea
                            id="assistantExampleUserPrompt"
                            name="exampleUserPrompt"
                            value={form.exampleUserPrompt}
                            onChange={(event) => setField('exampleUserPrompt', event.target.value)}
                            rows={4}
                            autoResize
                        />
                    </div>

                    <div className="admin-actions">
                        <Button type="submit" icon="pi pi-check" label="Создать" loading={loading}/>
                        <Button type="button" outlined label="Отмена" onClick={() => navigate('/assistants')}/>
                    </div>
                </form>
            </Card>
        </section>
    )
}
