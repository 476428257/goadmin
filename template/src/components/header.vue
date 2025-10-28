<template>
    <div class="header">
        <!-- 折叠按钮 -->
        <div class="header-left">
            <img class="logo" :src="loginConfig('logo')" alt="" />
            <div class="web-title">{{ loginConfig('title') }}</div>
            <div class="collapse-btn" @click="collapseChage">
                <el-icon v-if="sidebar.collapse">
                    <Expand />
                </el-icon>
                <el-icon v-else>
                    <Fold />
                </el-icon>
            </div>
        </div>
        <div class="header-right">
            <div class="header-user-con">
                <div class="btn-icon" @click="router.push('/theme')">
                    <el-tooltip effect="dark" content="设置主题" placement="bottom">
                        <i class="el-icon-lx-skin"></i>
                    </el-tooltip>
                </div>
                <!-- 用户头像 -->
                <el-avatar class="user-avator" :size="30" :src="currentAvatar" />
                <!-- 用户名下拉菜单 -->
                <el-dropdown class="user-name" trigger="click" @command="handleCommand">
                    <span class="el-dropdown-link">
                        {{ username }}
                        <el-icon class="el-icon--right">
                            <arrow-down />
                        </el-icon>
                    </span>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item command="user">个人中心</el-dropdown-item>
                            <el-dropdown-item divided command="loginout">退出登录</el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>
            </div>
        </div>
    </div>
</template>
<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { useSidebarStore } from '../store/sidebar';

const username: string | null = localStorage.getItem('username');
// 当前头像URL（包含时间戳防止缓存）
const avatarImg = ref("");
avatarImg.value = localStorage.getItem('avatar') || '';

const currentAvatar = computed(() => {
    if (!avatarImg.value) return '';
    // 添加时间戳参数防止浏览器缓存
    const separator = avatarImg.value.includes('?') ? '&' : '?';
    return `${avatarImg.value}${separator}t=${Date.now()}`;
});

// 配置缓存（响应式）
const configState = ref<Record<string, any>>({});
const getStoredConfig = () => {
    const config = localStorage.getItem('config');
    if (!config) return {} as Record<string, any>;
    try {
        return JSON.parse(config) || {};
    } catch {
        return {};
    }
};
configState.value = getStoredConfig();

const loginConfig = (field: string) => {
    return (configState.value as any)?.[field] ?? '';
};

// 监听头像更新事件
const handleAvatarUpdate = (event: CustomEvent) => {
    if (event.detail && event.detail.avatar) {
        avatarImg.value = event.detail.avatar;
    }
};

// 监听配置更新事件（来自 config.vue 保存成功后）
const handleConfigUpdated = (event: CustomEvent) => {
    if (event.detail && event.detail.config) {
        configState.value = event.detail.config;
    }
};

const sidebar = useSidebarStore();
// 侧边栏折叠
const collapseChage = () => {
    sidebar.handleCollapse();
};

onMounted(() => {
    if (document.body.clientWidth < 1500) {
        collapseChage();
    }
    // 监听头像更新事件
    window.addEventListener('avatarUpdated', handleAvatarUpdate as EventListener);
    // 监听配置更新事件
    window.addEventListener('configUpdated', handleConfigUpdated as EventListener);

});

onUnmounted(() => {
    // 清理事件监听
    window.removeEventListener('avatarUpdated', handleAvatarUpdate as EventListener);
    window.removeEventListener('configUpdated', handleConfigUpdated as EventListener);
});

// 用户名下拉菜单选择事件
const router = useRouter();
const handleCommand = (command: string) => {
    if (command == 'loginout') {
        localStorage.clear()
        router.push('/login');
    } else if (command == 'user') {
        router.push('/ucenter');
    }
};

</script>
<style scoped>
.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-sizing: border-box;
    width: 100%;
    height: 70px;
    color: var(--header-text-color);
    background-color: var(--header-bg-color);
    border-bottom: 1px solid #ddd;
}

.header-left {
    display: flex;
    align-items: center;
    padding-left: 20px;
    height: 100%;
}

.logo {
    width: 35px;
}

.web-title {
    margin: 0 40px 0 10px;
    font-size: 22px;
}

.collapse-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
    padding: 0 10px;
    cursor: pointer;
    opacity: 0.8;
    font-size: 22px;
}

.collapse-btn:hover {
    opacity: 1;
}

.header-right {
    float: right;
    padding-right: 50px;
}

.header-user-con {
    display: flex;
    height: 70px;
    align-items: center;
}

.btn-fullscreen {
    transform: rotate(45deg);
    margin-right: 5px;
    font-size: 24px;
}

.btn-icon {
    position: relative;
    width: 30px;
    height: 30px;
    text-align: center;
    cursor: pointer;
    display: flex;
    align-items: center;
    color: var(--header-text-color);
    margin: 0 5px;
    font-size: 20px;
}

.btn-bell-badge {
    position: absolute;
    right: 4px;
    top: 0px;
    width: 8px;
    height: 8px;
    border-radius: 4px;
    background: #f56c6c;
    color: var(--header-text-color);
}

.user-avator {
    margin: 0 10px 0 20px;
}

.el-dropdown-link {
    color: var(--header-text-color);
    cursor: pointer;
    display: flex;
    align-items: center;
}

.el-dropdown-menu__item {
    text-align: center;
}
</style>
