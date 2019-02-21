App.filter("stripTags", () => text => !text ? "" : String(text).replace(/<[^>]+>/gm, "").replace(/&[^;]+;/gm, ""))
