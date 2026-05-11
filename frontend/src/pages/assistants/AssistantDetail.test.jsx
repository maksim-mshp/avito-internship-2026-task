import {fireEvent, render, screen, waitFor} from '@testing-library/react'
import {MemoryRouter, Route, Routes} from 'react-router-dom'
import {beforeEach, describe, expect, it, vi} from 'vitest'
import {AuthContext} from '../../context/AuthContext.jsx'
import {addFavoriteAssistant, getAssistant, removeFavoriteAssistant, runAssistant} from '../../services/catalog.js'
import {AssistantDetail} from './AssistantDetail.jsx'

vi.mock('../../services/catalog.js', () => ({
    addFavoriteAssistant: vi.fn(),
    getAssistant: vi.fn(),
    removeFavoriteAssistant: vi.fn(),
    runAssistant: vi.fn(),
}))

beforeEach(() => {
    addFavoriteAssistant.mockReset()
    getAssistant.mockReset()
    removeFavoriteAssistant.mockReset()
    runAssistant.mockReset()
})

const assistant = {
    id: 'assistant-id',
    categoryName: 'Разработка',
    name: 'Проверочный ассистент',
    description: 'Помогает проверять запуск',
    model: 'mock-smart',
    systemPrompt: 'секретный системный промпт',
    exampleUserPrompt: 'пример контекста',
    isFavorite: false,
    isActive: true,
}

const renderDetail = () => {
    render(
        <MemoryRouter initialEntries={['/assistants/assistant-id']}>
            <AuthContext.Provider value={{user: {role: 'user'}}}>
                <Routes>
                    <Route path="/assistants/:assistantId" element={<AssistantDetail/>}/>
                </Routes>
            </AuthContext.Provider>
        </MemoryRouter>,
    )
}

describe('AssistantDetail', () => {
    it('загружает ассистента без показа системного промпта', async () => {
        getAssistant.mockResolvedValueOnce({data: assistant})

        renderDetail()

        expect(await screen.findByText('Проверочный ассистент')).toBeInTheDocument()
        expect(screen.getByDisplayValue('пример контекста')).toBeInTheDocument()
        expect(screen.queryByText('секретный системный промпт')).not.toBeInTheDocument()
    })

    it('запускает ассистента только с пользовательским контекстом', async () => {
        getAssistant.mockResolvedValueOnce({data: assistant})
        runAssistant.mockResolvedValueOnce({
            data: {
                id: 'run-id',
                status: 'success',
                output: 'готовый ответ',
            },
        })

        renderDetail()

        const textarea = await screen.findByLabelText('Пользовательский контекст')
        fireEvent.change(textarea, {target: {value: 'новый контекст'}})
        fireEvent.click(screen.getByRole('button', {name: 'Запустить'}))

        await waitFor(() => {
            expect(runAssistant).toHaveBeenCalledWith('assistant-id', 'новый контекст')
        })
        expect(await screen.findByText('готовый ответ')).toBeInTheDocument()
    })

    it('добавляет ассистента в избранное', async () => {
        getAssistant.mockResolvedValueOnce({data: assistant})
        addFavoriteAssistant.mockResolvedValueOnce({})

        renderDetail()

        fireEvent.click(await screen.findByRole('button', {name: 'В избранное'}))

        await waitFor(() => {
            expect(addFavoriteAssistant).toHaveBeenCalledWith('assistant-id')
        })
        expect(await screen.findByRole('button', {name: 'В избранном'})).toBeInTheDocument()
    })
})
