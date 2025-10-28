<template>
    <div>
        <TableSearch :query="query" :showsearch="showsearch" :options="searchOpt" :search="getData" />
        <div class="container">
            <TableCustom :columns="columns" v-model:showsearch="showsearch" :tableData="tableData" :total="page.total" :viewFunc="handleView" :delSelection="delSelection" :page-change="changePage" :editFunc="handleEdit" :handleAdd="handleAdd"  :refresh="refresh">
                <template #toolbarBtn> 
                    <!-- <el-button type="warning" :icon="CirclePlusFilled" @click="handleAdd">新增</el-button> -->
                </template>
                <template #is_hot="{ rows }">
                    <el-switch
                        v-model="rows.is_hot"
                        :active-value="1"
                        :inactive-value="0"
                        @change="handleChange(rows)"
                    />
                </template>
                <template #is_top="{ rows }">
                    <el-switch
                        v-model="rows.is_top"
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
            <TableEdit :form-data="rowData" :options="options" :edit="isEdit" :update="updateData" >
   
            </TableEdit>
        </el-dialog>
        <el-dialog title="查看详情" v-model="visible1" width="700px" destroy-on-close>
            <TableEdit :form-data="rowData" :options="options" :edit="isEdit" >
                <template #status="{ rows }">
                    <el-switch
                        v-model="rows.status"
                        :active-value="1"
                        :inactive-value="0"
                        disabled
                    />
                </template>
            </TableEdit>
        </el-dialog>
    </div>
</template>

<script setup lang="ts" name="system-user">
import TableCustom from '@/components/table-custom.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOption, FormOptionList } from '@/types/form-option';
import service from '@/utils/request';
import { ElMessage } from 'element-plus';
import NProgress from 'nprogress';
import { onMounted, reactive, ref } from 'vue';

// 添加一个标志变量，记录组件是否已初始化
const isInitialized = ref(false);
const showsearch = ref(false);
const isEdit = ref(false);
// 在 mounted 生命周期钩子中设置标志
onMounted(() => {
    isInitialized.value = true;
});
// 查询相关
const query = reactive({
    title: '',
    is_hot: '',
    is_top: '',
    status: '',
    created_at: [],
    updated_at: [],
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '标题：', prop: 'title' },
    { type: 'select', label: '热门：', prop: 'is_hot',opts:[{"label":"请选择", "value":""},{"label":"是", "value":1},{"label":"否", "value":0}] },
    { type: 'select', label: '置顶：', prop: 'is_top',opts:[{"label":"请选择", "value":""},{"label":"是", "value":1},{"label":"否", "value":0}] },
    { type: 'select', label: '状态：', prop: 'status',opts:[{"label":"请选择", "value":""},{"label":"启用", "value":1},{"label":"禁用", "value":0}]},
    { type: 'daterange', label: '创建时间：', prop: 'created_at', format: 'YYYY-MM-DD HH:mm:ss' },
    { type: 'daterange', label: '更新时间：', prop: 'updated_at', format: 'YYYY-MM-DD HH:mm:ss' },
])
// 表格相关
let columns = ref([
    { type: 'selection' },
    { prop: 'id', label: 'ID', width: 55, align: 'center' },
    { prop: 'title', label: '标题',width: '300', ellipsisClick: true },
    { prop: 'image', label: '图片',isImage:true },
    { prop: 'is_hot', label: '热门' },
    { prop: 'is_top', label: '置顶' },
    { prop: 'status', label: '状态' },
    { prop: 'created_at', label: '创建时间' },
    { prop: 'updated_at', label: '更新时间' },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<any>([]);
const getData = async () => {
    NProgress.start();
    try {
        const params = {
            ...query,
            page: page.index,
            pagesize: page.size
        };
        
        // 判断并移除空值字段
        if (params.is_hot === '' ) {
            delete params.is_hot;
        }
        if (params.is_top === '' ) {
            delete params.is_top;
        }
        if (params.status === '' ) {
            delete params.status;
        }
        const res = await service.post('/api/api/v1/article/list',params);
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
const refresh = async () => {
    await getData();
};
const changePage = (val: number) => {
    page.index = val;
    getData();
};
const handleChange =async (row:any) => {
    if(isInitialized.value&&row.id){
        try {
        const res = await service.post("/api/api/v1/article/updatestatus",{"id":row.id,"status":row.status,"is_hot":row.is_hot,"is_top":row.is_top});
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
    span: 24, // 修改为24，让富文本编辑器可以占满宽度
    list: [
        { type: 'input', label: '标题', prop: 'title', required: true},
        { type: 'upload', label: '图片', prop: 'image',multiple:true, required: true},
        { type: 'switch', label: '热门', prop: 'is_hot',activeValue: 1,  
            inactiveValue: 0},
        { type: 'switch', label: '置顶', prop: 'is_top',activeValue: 1,  
            inactiveValue: 0},
        { 
            type: 'switch', 
            label: '状态', 
            prop: 'status',
            activeValue: 1,  
            inactiveValue: 0,
        },
        { type: 'editor', label: '内容', prop: 'content', required: true, height: '400px',},       
    ]
})

const visible = ref(false);

const rowData = ref({});
const handleEdit = async (row: any) => {
    try {
        const res = await service.post("/api/api/v1/article/info", {id:row[0].id});
        if (res.data.code === 200) {
            rowData.value =res.data.data;
            isEdit.value = true;
            visible.value = true;
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
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
const updateData = async (form: any) => {
    try {
        var url = "/api/api/v1/article/add";
        if (form.id) {
            url = "/api/api/v1/article/update";
        }
        form.image = form.image.join(',');
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

const handleView = async (row: any) => {
    try {
        const res = await service.post("/api/api/v1/article/info", {id:row[0].id});
        if (res.data.code === 200) {
            rowData.value =res.data.data;
            isEdit.value = true;
            visible1.value = true;
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
};

const delSelection = async (rows: any) => {
    try {
        var ids=[];
        rows.forEach((item: any) => {
            ids.push(item.id)
        })
        const res = await service.post("/api/api/v1/article/del", {ids:ids});
        if (res.data.code === 200) {
            ElMessage.success(res.data.msg);
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