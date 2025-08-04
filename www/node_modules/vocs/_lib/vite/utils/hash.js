import { createHash } from 'node:crypto';
export function hash(text, length) {
    let hash = createHash('sha256').update(text).digest('hex');
    if (length)
        hash = hash.substring(0, length);
    return hash;
}
//# sourceMappingURL=hash.js.map