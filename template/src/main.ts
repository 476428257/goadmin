import * as ElementPlusIconsVue from '@element-plus/icons-vue';
import 'element-plus/dist/index.css';
import { createPinia } from 'pinia';
import { createApp } from 'vue';
import App from './App.vue';
import './assets/css/icon.css';
import router from './router';
import { usePermissStore } from './store/permiss';

// 页面标题：从缓存应用并监听更新
function applyTitleFromConfig(raw?: string) {
    try {
        const cfg = raw ? JSON.parse(raw || '{}') : JSON.parse(localStorage.getItem('config') || '{}');
        if (cfg && typeof cfg.title === 'string' && cfg.title) {
            document.title = cfg.title;
        }
    } catch { /* ignore */ }
}
// 首次加载应用缓存标题
applyTitleFromConfig();
// 监听来自应用内部的配置更新事件
window.addEventListener('configUpdated', (evt: any) => {
    if (evt?.detail?.config) {
        const title = evt.detail.config.title;
        if (typeof title === 'string') document.title = title || document.title;
    }
});
// 监听 storage 事件（跨标签页）
window.addEventListener('storage', (e: StorageEvent) => {
    if (e.key === 'config') {
        applyTitleFromConfig(e.newValue || undefined);
    }
});

const app = createApp(App);
app.config.warnHandler = () => {}
app.config.performance = false // 如果开了性能监控，也可以关掉
app.use(createPinia());
app.use(router);

// 注册elementplus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
    app.component(key, component);
}
// 自定义权限指令
const permiss = usePermissStore();
app.directive('permiss', {
    mounted(el, binding) {
        if (binding.value && !permiss.key.includes(String(binding.value)) && !permiss.key.includes("*")) {
            el['hidden'] = true;
        }
    },
});

app.mount('#app');
