import {Tag} from 'primereact/tag'

export const runStatusOptions = [
    {label: 'Все статусы', value: null},
    {label: 'В ожидании', value: 'pending'},
    {label: 'Успешно', value: 'success'},
    {label: 'Ошибка', value: 'failed'},
]

const statusView = {
    pending: {icon: 'pi pi-clock', severity: 'warning', value: 'В ожидании'},
    success: {icon: 'pi pi-check', severity: 'success', value: 'Успешно'},
    failed: {icon: 'pi pi-times', severity: 'danger', value: 'Ошибка'},
}

export const RunStatusTag = ({status}) => {
    const view = statusView[status] || {icon: 'pi pi-info-circle', severity: 'secondary', value: status}

    return <Tag className="run-status-tag" icon={view.icon} severity={view.severity} value={view.value}/>
}
