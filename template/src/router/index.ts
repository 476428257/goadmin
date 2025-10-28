import axios from 'axios';
import { ElMessage } from 'element-plus';
import NProgress from 'nprogress';
import 'nprogress/nprogress.css';
import { ref } from 'vue';
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { usePermissStore } from '../store/permiss';
import Home from '../views/home.vue';

const routes: RouteRecordRaw[] = [
    {
        path: '/',
        name: 'Home',
        component: Home,
        redirect: '/dashboard',
        children: [
            {
                path: 'ucenter',
                meta: {
                    title: '个人中心',
                },
                component: () => import( '../views/pages/ucenter.vue'),
            },
            {
                path: 'theme',
                meta: {
                    title: '主题',
                    noAuth: true,
                },
                component: () => import( '../views/pages/theme.vue'),
            },
        ]
    },
    {
        path: '/login',
        meta: {
            title: '登录',
            noAuth: true,
        },
        component: () => import( '../views/pages/login.vue'),
    },
    {
        path: '/403',
        meta: {
            title: '没有权限',
            noAuth: true,
        },
        component: () => import( '../views/pages/403.vue'),
    },
    {
        path: '/404',
        meta: {
            title: '找不到页面',
            noAuth: true,
        },
        component: () => import( '../views/pages/404.vue'),
    },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});
const modules = import.meta.glob('../views/**/*.vue');
const routepage = ref<any[]>([]);
let routed = false;
let configed = false;
let isBootstrapped = false;
let bootstrapPromise: Promise<void> | null = null;

const bootstrap = async () => {
    if (isBootstrapped) return;
    if (bootstrapPromise) return bootstrapPromise;

    bootstrapPromise = (async () => {
        try {
            const storedData = localStorage.getItem('routepage');
            routepage.value = storedData ? JSON.parse(storedData) : [];
            if (routepage.value.length === 0) {
                const res = await axios.post('/api/api/v1/auth/rule/getmenu');
                if (res.data.code === 200) {
                    routepage.value = res.data.data;
                    localStorage.setItem('routepage', JSON.stringify(routepage.value));
                } else {
                    ElMessage.error(res.data.msg);
                    return; // 失败时直接返回，保持状态为未完成
                }
            }

            const config = localStorage.getItem('config');
            if (!config) {
                const res = await axios.post('/api/api/v1/config/getconfig');
                if (res.data.code === 200) {
                    localStorage.setItem('config', JSON.stringify(res.data.data));
                } else {
                    ElMessage.error(res.data.msg);
                    return;
                }
            }

            // 注入/刷新动态路由
            routepage.value.forEach(route => {
                if (router.hasRoute(route.route)) {
                    router.removeRoute(route.route);
                }
            });
            routepage.value.forEach(menu => {
                const componentPath = `../views/${menu.pagepath}.vue`;
                const component = modules[componentPath];
                if (!component) {
                    console.warn(`页面路径未找到: ${componentPath}`);
                    return;
                }
                router.addRoute('Home', {
                    path: menu.route,
                    name: menu.route,
                    meta: { title: menu.title, permiss: menu.id },
                    component,
                });
            });
            router.addRoute({
                path: "/:path(.*)",
                redirect: '/404',
            });

            routed = true;
            configed = true;
            isBootstrapped = true;
        } catch (e: any) {
            ElMessage.error(e.message || '初始化失败');
        } finally {
            bootstrapPromise = null;
        }
    })();

    return bootstrapPromise;
};

router.beforeEach(async (to, from, next) => {
    NProgress.start();
    let didBootstrap = false;
    if (!isBootstrapped) {
        await bootstrap();
        didBootstrap = true;
    }

    if (didBootstrap) {
        // 第一次动态注入后，使用 replace 重新匹配一次路由
        next({ ...to, replace: true });
        return;
    }

    const role = localStorage.getItem('username');
    const permiss = usePermissStore();
    if (!role && to.meta.noAuth !== true) {
        next('/login');
    } else if (typeof to.meta.permiss == 'string' && !permiss.key.includes(to.meta.permiss) && !permiss.key.includes("*")) {
        next('/403');
    } else {
        next();
    }
});

router.afterEach(() => {
    NProgress.done();
});

export default router;
