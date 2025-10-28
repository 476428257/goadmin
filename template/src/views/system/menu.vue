<template>
    <div>
        <div class="container">
            <TableCustom ref="tableCustomRef" :columns="columns" :tableData="menuData" row-key="id" :has-pagination="false" :showsearchbtn="false" :m_d1="false"
                :viewFunc="handleView" :delSelection="delSelection" :editFunc="handleEdit" :refresh="refresh" :handleAdd="handleAdd">
                <template #toolbarBtn>
                    <el-button 
                        type="primary" 
                        :icon="expandAllState ? Fold : Expand"
                        @click="toggleExpandAll"
                    >
                        {{ expandAllState ? '收缩全部' : '展开全部' }}
                    </el-button>
                </template>
                <template #icon="{ rows }">
                    <el-icon>
                        <component :is="rows.icon"></component>
                    </el-icon>
                </template>
                <template #is_menu="{ rows }">
                    <el-switch
                        v-model="rows.is_menu"
                        :active-value="1"
                        :inactive-value="0"
                        @change="handleChange(rows)"
                    />
                </template>
                <template #status="{ rows }">
                    <el-switch
                        v-model="rows.status"
                        :active-value="1"
                        :inactive-value="0"
                        @change="handleChange(rows)"
                    />
                </template>
            </TableCustom>
        </div>
        <el-dialog :title="isEdit ? '编辑' : '新增'" v-model="visible" width="700px" destroy-on-close
            :close-on-click-modal="false" @close="closeDialog">
            <TableEdit :form-data="rowData" :options="options" :edit="isEdit" :update="updateData">
                <template #pid>
                    <el-cascader v-model="rowData.pid" :options="cascaderOptions" :props="{ checkStrictly: true , emitPath: false }"
                         />
                </template>
            </TableEdit>
        </el-dialog>
        <el-dialog title="查看详情" v-model="visible1" width="700px" destroy-on-close>
            <TableDetail :data="viewData">
                <template #icon="{ rows }">
                    <el-icon>
                        <component :is="rows.icon"></component>
                    </el-icon>
                </template>
                <template #is_menu="{ rows }">
                    <el-switch
                        v-model="rows.is_menu"
                        :active-value="1"
                        :inactive-value="0"
                        disabled
                    />
                </template>
                <template #status="{ rows }">
                    <el-switch
                        v-model="rows.status"
                        :active-value="1"
                        :inactive-value="0"
                        disabled
                    />
                </template>
            </TableDetail>
        </el-dialog>
    </div>
</template>

<script setup lang="ts" name="system-menu">
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import { FormOption } from '@/types/form-option';
import { Menus } from '@/types/menu';
import service from '@/utils/request';
import { Expand, Fold } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import NProgress from 'nprogress';
import { onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

// 添加一个标志变量，记录组件是否已初始化
const isInitialized = ref(false);

// 在 mounted 生命周期钩子中设置标志
onMounted(() => {
    isInitialized.value = true;
});
const router = useRouter();

// 展开/收缩状态管理
const expandAllState = ref(false); // false表示未展开，true表示已展开
const tableCustomRef = ref();

// 表格相关
let columns = ref([
    { type: 'selection' },
    { prop: 'title', label: '菜单名称', align: 'left',width:"320" },
    { prop: 'icon', label: '图标' },
    { prop: 'route', label: '路由路径',width:"310"  },
    { prop: 'is_menu', label: '是否菜单' },
    { prop: 'pagepath', label: '页面路径' },
    { prop: 'status', label: '状态' },
    { prop: 'weigh', label: '权重' },
])

const getOptions = (data: any) => {
    return data.map(item => {
        const a: any = {
            label: item.title,
            value: item.id,
        }
        if (item.children) {
            a.children = getOptions(item.children)
        }
        return a
    })
}
const refresh = async () => {
    NProgress.start();
    await fetchMenuData();
    NProgress.done();
};

const menuData = ref<Menus[]>([]);
const cascaderOptions = ref<any[]>([]);

const fetchMenuData = async () => {
    try {
        const res = await service.post('/api/api/v1/auth/rule/list',{is_menu:0});
        if (res.data.code === 200) {
            menuData.value = res.data.data.rules;
            var rule=getOptions(menuData.value);
            // 正确添加顶级节点，确保value为数字0而不是字符串
            cascaderOptions.value = [{label:"顶级节点", value:"顶"}, ...rule];
            if(ischange.value){
                localStorage.removeItem('menuData')
                localStorage.removeItem('routepage')
                localStorage.removeItem('allmenuData')
                setTimeout(() => {
                router.go(0)  // 刷新当前路由
                }, 1000)
            }
        }else{
            ElMessage.error(res.data.msg);
            // router.push('/login');
            return false;
        }
    } catch (e) {
        ElMessage.error("未知错误");
        return false;
    }
};

fetchMenuData();


// 新增/编辑弹窗相关
let options = ref<FormOption>({
    labelWidth: '100px',
    span: 12,
    list: [
        { type: 'input', label: '菜单名称', prop: 'title', required: true },
        { type: 'input', label: '路由路径', prop: 'route', required: true },
        { type: 'input', label: '页面路径', prop: 'pagepath' },
        { type: 'input', label: '图标', prop: 'icon',suffixLink: 'https://element-plus.org/zh-CN/component/icon.html#icon-collection', 
        suffixTooltip: '选择图标'},
        { type: 'parent', label: '父菜单', prop: 'pid' },
        { type: 'number', label: '权重', prop: 'weigh' },
        { 
            type: 'switch', 
            label: '是否菜单', 
            prop: 'is_menu',
            activeValue: 1,  
            inactiveValue: 0,
        },
        { 
            type: 'switch', 
            label: '状态', 
            prop: 'status',
            activeValue: 1,  
            inactiveValue: 0,
        }
    ]
})
const visible = ref(false);
const isEdit = ref(false);
const ischange = ref(false);
const rowData = ref<any>({});
const handleEdit = (row: any) => {
    if(row[0].pid==0){
        row[0].pid="顶"
    }
    rowData.value = { ...row[0] };
    isEdit.value = true;
    visible.value = true;
};
const handleAdd = () => {
    rowData.value = {is_menu: 0, status: 1,pid:"顶"};
    isEdit.value = false;
    visible.value = true;
};

const updateData = async (form: any) => {
    try {
        var url="/api/api/v1/auth/rule/add";
        if(form.id){
            url="/api/api/v1/auth/rule/update";
        }
        form.pid=rowData.value.pid
        if(form.pid =="顶"){
            form.pid = 0;
        }
        const res = await service.post(url, form);
        if (res.data.code === 200) {
            ElMessage.success(res.data.msg);
            closeDialog();
            ischange.value=true
            fetchMenuData(); // 刷新表格
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
};

const closeDialog = () => {
    visible.value = false;
    isEdit.value = false;
};

// 查看详情弹窗相关
const visible1 = ref(false);
const viewData = ref({
    row: {},
    list: []
});
const handleView = (row: Menus) => {
    viewData.value.row = { ...row[0] }
    viewData.value.list = [
        {
            prop: 'id',
            label: '菜单ID',
        },
        {
            prop: 'pid',
            label: '父菜单ID',
        },
        {
            prop: 'title',
            label: '菜单名称',
        },
        {
            prop: 'icon',
            label: '图标',
        },
        {
            prop: 'route',
            label: '路由路径',
        },
        {
            prop: 'weigh',
            label: '权重',
        },
        {
            prop: 'is_menu',
            label: '是否菜单',
        },
        {
            prop: 'status',
            label: '状态',
        },
    ]
    visible1.value = true;
};

// 删除相关
const delSelection =async (row: any) => {
    try {
        var ids=[];
        row.forEach((item: any) => {
            ids.push(item.id)
        })
        const res = await service.post("/api/api/v1/auth/rule/del", {ids:ids});
        if (res.data.code === 200) {
            ElMessage.success('删除成功');
            closeDialog();
            ischange.value=true
            fetchMenuData(); // 刷新表格
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
}



const changedata= async (data)=>{
    try {
        const res = await service.post("/api/api/v1/auth/rule/updatestatus", data);
        if (res.data.code === 200) {
            ElMessage.success('操作成功');
            ischange.value=true
            fetchMenuData(); // 刷新表格
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
}

// 在 script setup 中添加事件处理函数
const handleChange = (row: Menus) => {
    if(isInitialized.value&&row.id){
        changedata({"is_menu":row.is_menu,"id":row.id,"status":row.status})
    }
};

// 展开/收缩所有树形表格节点
const toggleExpandAll = () => {
    if (tableCustomRef.value) {
        expandAllState.value = !expandAllState.value;
        tableCustomRef.value.toggleExpandAll(expandAllState.value);
    }
};

</script>

<style scoped></style>