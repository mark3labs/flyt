export function debounce(fn, delay) {
    let invoked = false;
    return () => {
        invoked = true;
        setTimeout(() => {
            if (invoked)
                fn();
            invoked = false;
        }, delay);
    };
}
//# sourceMappingURL=debounce.js.map