<template>
    <div>
        <TableSearch :query="query" :showsearch="showsearch" :options="searchOpt" :search="getData" />
        <div class="container">
            <TableCustom :columns="columns" v-model:showsearch="showsearch" :tableData="tableData" :total="page.total" :viewFunc="handleView" :delSelection="delSelection" :page-change="changePage" :editFunc="handleEdit" :handleAdd="handleAdd"  :refresh="refresh">
                <template #toolbarBtn> 
                    <!-- <el-button type="warning" :icon="CirclePlusFilled" @click="handleAdd">新增</el-button> -->
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
            <TableEdit :form-data="rowData" :options="options" :edit="isEdit" :update="updateData" />
        </el-dialog>
        <el-dialog title="查看详情" v-model="visible1" width="700px" destroy-on-close>
            <TableDetail :data="viewData">
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

<script setup lang="ts" name="system-user">
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOption, FormOptionList } from '@/types/form-option';
import { User } from '@/types/user';
import service from '@/utils/request';
import { ElMessage } from 'element-plus';
import NProgress from 'nprogress';
import { onMounted, reactive, ref, watch } from 'vue';

// 添加一个标志变量，记录组件是否已初始化
const isInitialized = ref(false);
const showsearch = ref(false);
const roleList = ref([]);
const isEdit = ref(false);
// 在 mounted 生命周期钩子中设置标志
onMounted(() => {
    isInitialized.value = true;
});
// 查询相关
const query = reactive({
    username: '',
    nickname: '',
    phone: '',
    email: '',
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '用户名：', prop: 'username' },
    { type: 'input', label: '昵称：', prop: 'nickname' },
    { type: 'input', label: '电话：', prop: 'phone' },
    { type: 'input', label: '邮箱：', prop: 'email' },
])
// 表格相关
let columns = ref([
    { type: 'selection' },
    { prop: 'id', label: 'ID', width: 55, align: 'center' },
    { prop: 'username', label: '用户名' },
    { prop: 'nickname', label: '昵称' },
    { prop: 'avatar', label: '头像',isImage:true },
    { prop: 'phone', label: '手机号' },
    { prop: 'email', label: '邮箱' },
    { prop: 'title', label: '角色' },
    { prop: 'created_at', label: '创建时间' },
    { prop: 'updated_at', label: '更新时间' },
    { prop: 'status', label: '状态' },
    // { prop: 'operator', label: '操作', width: 250 },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<User[]>([]);
const getData = async () => {
    NProgress.start();
    try {
        const params = {
        ...query,
        page: page.index,
        pagesize: page.size
    };
        const res = await service.post('/api/api/v1/admin/list',params);
        if (res.data.code === 200) {
            const list = res.data.data?.list;
    tableData.value = Array.isArray(list) ? list : [];
    page.total = res.data.data?.pageTotal || 0;
        }else{
            ElMessage.error(res.data.msg);
            // router.push('/login');
            return false;
        }
        NProgress.done();
    } catch (e) {
        ElMessage.error("未知错误");
        return false;
    }

   
};
getData();
const getRole = async () => {
    const res = await service.post('/api/api/v1/auth/role/getkv');
    if (res.data.code === 200) {
        roleList.value = res.data.data.map((item: any) => ({
            value: item.id,
            label: item.title
        }));
      
    }
}
getRole();

const refresh = async () => {
    await getData();
};
const changePage = (val: number) => {
    page.index = val;
    getData();
};
const handleChange =async (row: User) => {
    if(isInitialized.value&&row.id){
        try {
        const res = await service.post("/api/api/v1/admin/updatestatus",{"id":row.id,"status":row.status});
        if (res.data.code === 200) {
            ElMessage.success('操作成功');
            getData(); // 刷新表格
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
    }
};


// 新增/编辑弹窗相关
let options = ref<FormOption>({
    labelWidth: '100px',
    span: 12,
    list: [
        { type: 'input', label: '用户名', prop: 'username', required: true },
        { type: 'input', label: '昵称', prop: 'nickname', required: true },
        { type: 'input', label: '密码', prop: 'password'},
        { type: 'input', label: '手机号', prop: 'phone', required: true },
        { type: 'input', label: '邮箱', prop: 'email', required: true },
        { type: 'select', label: '角色', prop: 'role_id', required: true,opts:roleList.value},
        { 
            type: 'switch', 
            label: '状态', 
            prop: 'status',
            activeValue: 1,  
            inactiveValue: 0,
        }
    ]
})
watch(roleList, (newVal) => {
    const roleOpt = options.value.list.find(item => item.prop === 'role_id');
    if (roleOpt) {
        roleOpt.opts = newVal;
    }
}, { immediate: true });
const visible = ref(false);

const rowData = ref({});
const handleEdit = (row: User) => {
    rowData.value = { ...row[0] };
    const passwordOpt = options.value.list.find(item => item.prop === 'password');
    if (passwordOpt) {
        passwordOpt.required = false;
    }
    const username = options.value.list.find(item => item.prop === 'username');
    if (username) {
        username.disabled = true;
    }
    isEdit.value = true;
    visible.value = true;
};
const handleAdd = () => {
    rowData.value = { status: 1};
    const passwordOpt = options.value.list.find(item => item.prop === 'password');
    if (passwordOpt) {
        passwordOpt.required = true;
    }
    const username = options.value.list.find(item => item.prop === 'username');
    if (username) {
        username.disabled = false;
    }
    isEdit.value = false;
    visible.value = true;
};
const updateData =async (form: any)=> {
    try {
        var url="/api/api/v1/admin/add";
        if(form.id){
            url="/api/api/v1/admin/update";
        }
        const res = await service.post(url, form);
        if (res.data.code === 200) {
            ElMessage.success(res.data.msg);
            closeDialog();
            getData(); // 刷新表格
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

const handleView = (row: User) => {
    viewData.value.row = { ...row[0] }
    viewData.value.list = [
        {
            prop: 'id',
            label: 'ID',
        },
        {
            prop: 'username',
            label: '用户名',
        },
        {
            prop: 'nickname',
            label: '昵称',
        },
        {
            prop: 'phone',
            label: '电话',
        },
        {
            prop: 'email',
            label: '邮箱',
        },
        {
            prop: 'title',
            label: '角色',
        },
        {
            prop: 'created_at',
            label: '创建时间',
        },
        {
            prop: 'updated_at',
            label: '更新时间',
        },
        {
            prop: 'status',
            label: '状态',
        },
    ]
    visible1.value = true;
};

const delSelection = async (rows: any) => {
    try {
        var ids=[];
        rows.forEach((item: any) => {
            ids.push(item.id)
        })
        const res = await service.post("/api/api/v1/admin/del", {ids:ids});
        if (res.data.code === 200) {
            ElMessage.success('删除成功');
            getData(); // 刷新表格
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
}
</script>

<style scoped></style>