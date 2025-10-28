<template>
    <div class="sidebar">
        <el-menu
            class="sidebar-el-menu"
            :default-active="onRoutes"
            :collapse="sidebar.collapse"
            :background-color="sidebar.bgColor"
            :text-color="sidebar.textColor"
            router
        >
            <template v-for="item in menuData">
                <template v-if="item.children">
                    <el-sub-menu :index="item.route" :key="item.route" v-permiss="item.id">
                        <template #title>
                            <el-icon>
                                <component :is="item.icon"></component>
                            </el-icon>
                            <span>{{ item.title }}</span>
                        </template>
                        <template v-for="subItem in item.children">
                            <el-sub-menu
                                v-if="subItem.children"
                                :index="subItem.route"
                                :key="subItem.route"
                                v-permiss="item.id"
                            >
                                <template #title>{{ subItem.title }}</template>
                                <el-menu-item
                                    v-for="(threeItem, i) in subItem.children"
                                    :key="i"
                                    :index="threeItem.route"
                                >
                                    {{ threeItem.title }}
                                </el-menu-item>
                            </el-sub-menu>
                            <el-menu-item v-else :index="subItem.route" v-permiss="item.id">
                                {{ subItem.title }}
                            </el-menu-item>
                        </template>
                    </el-sub-menu>
                </template>
                <template v-else>
                    <el-menu-item :index="item.route" :key="item.route" v-permiss="item.id">
                        <el-icon>
                            <component :is="item.icon"></component>
                        </el-icon>
                        <template #title>{{ item.title }}</template>
                    </el-menu-item>
                </template>
            </template>
        </el-menu>
    </div>
</template>

<script setup lang="ts">
// import { menuData } from '@/components/menu';
import { Menus } from '@/types/menu';
import service from '@/utils/request';
import { ElMessage } from 'element-plus';
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useSidebarStore } from '../store/sidebar';

const route = useRoute();
const onRoutes = computed(() => {
    return route.path;
});
const sidebar = useSidebarStore();
var menuData = ref<Menus[]>([])
const fetchData = async () => {
    try {
        const storedData = localStorage.getItem('menuData');
        menuData.value = storedData ? JSON.parse(storedData) : [];
        if(menuData.value.length === 0){
            const res = await service.post('/api/api/v1/auth/rule/list',{is_menu:1});
        if (res.data.code === 200) {
            menuData.value = res.data.data.rules;
            localStorage.setItem('menuData', JSON.stringify(menuData.value));
            localStorage.setItem('allmenuData', JSON.stringify(res.data.data.allrules));
        }else{
            ElMessage.error(res.data.msg);
            // router.push('/login');
            return false;
        }
        }
        
    } catch (e) {
        ElMessage.error("未知错误1");
        return false;
    }
};

fetchData();
</script>

<style scoped>
.sidebar {
    display: block;
    position: absolute;
    left: 0;
    top: 70px;
    bottom: 0;
    overflow-y: scroll;
}

.sidebar::-webkit-scrollbar {
    width: 0;
}

.sidebar-el-menu:not(.el-menu--collapse) {
    width: 250px;
}

.sidebar-el-menu {
    min-height: 100%;
}
</style>
