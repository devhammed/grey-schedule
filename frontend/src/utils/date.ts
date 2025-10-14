
export function toLocalInputValue(dt: Date|string) {
    if (typeof dt === 'string') {
        dt = new Date(dt);
    }

    const pad = (n: number) => n.toString().padStart(2, '0')

    const year = dt.getFullYear();
    const month = pad(dt.getMonth() + 1);
    const day = pad(dt.getDate());
    const hours = pad(dt.getHours());
    const minutes = pad(dt.getMinutes());

    return `${year}-${month}-${day}T${hours}:${minutes}`;
}

export function fromLocalInputValue(v: Date|string): string {
    if (typeof v === 'string') {
        v = new Date(v);
    }

    return v.toISOString();
}