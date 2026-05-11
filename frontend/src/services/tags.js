export const normalizeTags = (tags = []) => {
    const seen = new Set()

    return tags.reduce((result, tag) => {
        const value = String(tag || '').trim()
        const key = value.toLowerCase()

        if (!value || seen.has(key)) {
            return result
        }

        seen.add(key)
        result.push(value)

        return result
    }, [])
}
