App.filter("stripTags", () => {
    return (text) => {
        if (!text) {
            return ""
        }
        return String(text).replace(/<[^>]+>/gm, "").replace(/&[^;]+;/gm, "")
    }
})
