import {useContext, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import {Button} from 'primereact/button'
import {Card} from 'primereact/card'
import {InputText} from 'primereact/inputtext'
import {InputTextarea} from 'primereact/inputtextarea'
import {Message} from 'primereact/message'
import {ToastContext} from '../../context/ToastContext.jsx'
import {createCategory} from '../../services/catalog.js'
import {getTranslatedError} from '../../services/errors.js'
import '../../styles/Admin.css'

export const CategoryCreate = () => {
    const navigate = useNavigate()
    const {showSuccess} = useContext(ToastContext)
    const [name, setName] = useState('')
    const [description, setDescription] = useState('')
    const [error, setError] = useState(null)
    const [loading, setLoading] = useState(false)

    const submitHandler = (event) => {
        event.preventDefault()

        if (!name.trim()) {
            setError('Введите название категории')
            return
        }

        setLoading(true)
        setError(null)

        createCategory({
            name: name.trim(),
            description: description.trim() || null,
        }).then(() => {
            showSuccess('Категория создана')
            navigate('/assistants')
        }).catch((err) => {
            setError(getTranslatedError(err))
        }).finally(() => {
            setLoading(false)
        })
    }

    return (
        <section className="admin-page">
            <Button outlined icon="pi pi-arrow-left" label="К каталогу" onClick={() => navigate('/assistants')}/>

            <Card title="Новая категория" className="admin-card">
                <form className="admin-form" onSubmit={submitHandler}>
                    {error && <Message severity="error" text={error}/>}

                    <div className="admin-field">
                        <label htmlFor="categoryName">Название</label>
                        <InputText
                            id="categoryName"
                            name="name"
                            value={name}
                            onChange={(event) => setName(event.target.value)}
                        />
                    </div>

                    <div className="admin-field">
                        <label htmlFor="categoryDescription">Описание</label>
                        <InputTextarea
                            id="categoryDescription"
                            name="description"
                            value={description}
                            onChange={(event) => setDescription(event.target.value)}
                            rows={5}
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
