App.filter("stripTags", () => text => {
    if (!text) {
        return ""
    }
    return String(text).replace(/<[^>]+>/gm, "").replace(/&[^;]+;/gm, "")
})
