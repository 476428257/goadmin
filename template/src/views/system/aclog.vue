<template>
    <div>
        <TableSearch :query="query" :showsearch="showsearch" :options="searchOpt" :search="getData" />
        <div class="container">
            <TableCustom :columns="columns" v-model:showsearch="showsearch" :tableData="tableData" :total="page.total" :currentPage="page.index" :pageSize="page.size" :hasadd="false" :hasedit="false" :hasdel="false" :hasview="false" :changePage="changePage"  :refresh="refresh">
                <template #toolbarBtn> 
                    <!-- <el-button type="warning" :icon="CirclePlusFilled" @click="handleAdd">新增</el-button> -->
                </template>
            </TableCustom>

        </div>
    </div>
</template>

<script setup lang="ts" name="system-user">
import TableCustom from '@/components/table-custom.vue';
import TableSearch from '@/components/table-search.vue';
import { FormOptionList } from '@/types/form-option';
import service from '@/utils/request';
import { ElMessage } from 'element-plus';
import NProgress from 'nprogress';
import { reactive, ref } from 'vue';

const showsearch = ref(false);
// 查询相关
const query = reactive({
    username: '',
    path: '',
    request_data: '',
    created_at: [],
});
const searchOpt = ref<FormOptionList[]>([
    { type: 'input', label: '用户名：', prop: 'username' },
    { type: 'input', label: '接口路径：', prop: 'path' },
    { type: 'input', label: '数据内容：', prop: 'request_data' },
    { type: 'daterange', label: '操作时间：', prop: 'created_at',format: 'YYYY-MM-DD HH:mm:ss'  },
])
// 表格相关
let columns = ref([
    { type: 'selection' },
    { prop: 'id', label: 'ID'},
    { prop: 'username', label: '用户名',width:"150" },
    { prop: 'path', label: '接口路径：' },
    { prop: 'request_data', label: '数据内容',width: '300', ellipsisClick: true },
    { prop: 'created_at', label: '操作时间' },
])
const page = reactive({
    index: 1,
    size: 10,
    total: 0,
})
const tableData = ref<any[]>([]);
const getData = async () => {
    NProgress.start();
    try {
        const params = {
        ...query,
        page: page.index,
        pagesize: page.size
    };
        const res = await service.post('/api/api/v1/system/aclog',params);
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
const changePage = async(val: number) => {
    page.index = val;
   await  getData();
};
const refresh = async () => {
    await getData();
};

</script>

<style scoped>
</style>