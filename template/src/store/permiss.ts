import { defineStore } from 'pinia';

export const usePermissStore = defineStore('permiss', {
    state: () => {
        const rulearr = localStorage.getItem('rulearr');
        return {
            key: rulearr,
        };
    },
    actions: {
        handleSet(val: string[]) {
            this.key = val;
        },
    },
});
