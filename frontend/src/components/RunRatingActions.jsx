import {Button} from 'primereact/button'

export const RunRatingActions = ({rating, loading, onRate}) => {
    return (
        <div className="run-rating-actions">
            <Button
                rounded
                text
                severity={rating === 'like' ? 'success' : 'secondary'}
                icon={rating === 'like' ? 'pi pi-thumbs-up-fill' : 'pi pi-thumbs-up'}
                aria-label="Нравится"
                loading={loading === 'like'}
                disabled={Boolean(loading && loading !== 'like')}
                onClick={() => onRate('like')}
            />
            <Button
                rounded
                text
                severity={rating === 'dislike' ? 'danger' : 'secondary'}
                icon={rating === 'dislike' ? 'pi pi-thumbs-down-fill' : 'pi pi-thumbs-down'}
                aria-label="Не нравится"
                loading={loading === 'dislike'}
                disabled={Boolean(loading && loading !== 'dislike')}
                onClick={() => onRate('dislike')}
            />
        </div>
    )
}
