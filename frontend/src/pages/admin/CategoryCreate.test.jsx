import {render, screen, waitFor} from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import {MemoryRouter, Route, Routes} from 'react-router-dom'
import {describe, expect, it, vi} from 'vitest'
import {ToastContext} from '../../context/ToastContext.jsx'
import {createCategory} from '../../services/catalog.js'
import {CategoryCreate} from './CategoryCreate.jsx'

vi.mock('../../services/catalog.js', () => ({
    createCategory: vi.fn(),
}))

const renderCategoryCreate = ({showSuccess = vi.fn()} = {}) => {
    render(
        <MemoryRouter initialEntries={['/admin/categories/new']}>
            <ToastContext.Provider value={{showSuccess}}>
                <Routes>
                    <Route path="/admin/categories/new" element={<CategoryCreate/>}/>
                    <Route path="/assistants" element={<div>Каталог открыт</div>}/>
                </Routes>
            </ToastContext.Provider>
        </MemoryRouter>,
    )

    return {showSuccess}
}

describe('CategoryCreate', () => {
    it('не отправляет пустое название категории', async () => {
        const user = userEvent.setup()

        renderCategoryCreate()

        await user.click(screen.getByRole('button', {name: 'Создать'}))

        expect(screen.getByText('Введите название категории')).toBeInTheDocument()
        expect(createCategory).not.toHaveBeenCalled()
    })

    it('создает категорию и возвращает к каталогу', async () => {
        const user = userEvent.setup()
        const {showSuccess} = renderCategoryCreate()

        createCategory.mockResolvedValueOnce({data: {id: 'category-id'}})

        await user.type(screen.getByLabelText('Название'), 'Новая категория')
        await user.type(screen.getByLabelText('Описание'), 'Описание категории')
        await user.click(screen.getByRole('button', {name: 'Создать'}))

        await waitFor(() => {
            expect(createCategory).toHaveBeenCalledWith({
                name: 'Новая категория',
                description: 'Описание категории',
            })
        })
        expect(showSuccess).toHaveBeenCalledWith('Категория создана')
        expect(await screen.findByText('Каталог открыт')).toBeInTheDocument()
    })
})
