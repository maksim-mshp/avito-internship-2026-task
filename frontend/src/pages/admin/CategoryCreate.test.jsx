import {fireEvent, render, screen, waitFor} from '@testing-library/react'
import {MemoryRouter, Route, Routes} from 'react-router-dom'
import {beforeEach, describe, expect, it, vi} from 'vitest'
import {ToastContext} from '../../context/ToastContext.jsx'
import {createCategory} from '../../services/catalog.js'
import {CategoryCreate} from './CategoryCreate.jsx'

vi.mock('../../services/catalog.js', () => ({
    createCategory: vi.fn(),
}))

beforeEach(() => {
    createCategory.mockReset()
})

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
        renderCategoryCreate()

        fireEvent.click(screen.getByRole('button', {name: 'Создать'}))

        expect(screen.getByText('Введите название категории')).toBeInTheDocument()
        expect(createCategory).not.toHaveBeenCalled()
    })

    it('создает категорию и возвращает к каталогу', async () => {
        const {showSuccess} = renderCategoryCreate()

        createCategory.mockResolvedValueOnce({data: {id: 'category-id'}})

        fireEvent.change(screen.getByLabelText('Название'), {target: {value: 'Новая категория'}})
        fireEvent.change(screen.getByLabelText('Описание'), {target: {value: 'Описание категории'}})
        fireEvent.click(screen.getByRole('button', {name: 'Создать'}))

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
