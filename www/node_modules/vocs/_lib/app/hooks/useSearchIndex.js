import { getSearchIndex } from 'virtual:searchIndex';
import MiniSearch from 'minisearch';
import { useEffect, useState } from 'react';
let promise;
export function useSearchIndex() {
    const [searchIndex, setSearchIndex] = useState();
    useEffect(() => {
        ;
        (async () => {
            if (!promise)
                promise = getSearchIndex();
            const json = await promise;
            const searchIndex = MiniSearch.loadJSON(json, {
                fields: ['title', 'titles', 'text'],
                searchOptions: {
                    boost: { title: 4, text: 2, titles: 1 },
                    fuzzy: 0.2,
                    prefix: true,
                    // ...(theme.value.search?.provider === 'local' &&
                    //   theme.value.search.options?.miniSearch?.searchOptions),
                },
                storeFields: ['href', 'html', 'isPage', 'text', 'title', 'titles'],
                // ...(theme.value.search?.provider === 'local' &&
                //   theme.value.search.options?.miniSearch?.options),
            });
            setSearchIndex(searchIndex);
        })();
    }, []);
    useEffect(() => {
        if (!import.meta.hot)
            return;
        // TODO: Update index
        import.meta.hot.accept('virtual:searchIndex', (m) => {
            if (m) {
                console.log('update', m);
            }
        });
    }, []);
    return searchIndex;
}
//# sourceMappingURL=useSearchIndex.js.map