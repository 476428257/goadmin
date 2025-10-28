<template>
    <div>
        <div class="container">

            <TableCustom :columns="columns" :has-pagination="false" :tableData="tableData"  :viewFunc="handleView"
                 :showsearchbtn="false" :m_d1="false" :handleAdd="handleAdd" :delSelection="delSelection" :refresh="getData"  :editFunc="handleEdit" :default-expand-all="true" >
                <template #toolbarBtn>

                </template>
                <template #status="{ rows, index }">
                    <el-switch
                        v-model="rows.status"
                        :active-value="1"
                        :inactive-value="0"
                        @change="handleChange(rows)"
                        :disabled="index === 0"
                    />
                </template>
            </TableCustom>
        </div>
        <el-dialog :title="isEdit ? '编辑' : '新增'" v-model="visible" width="700px" destroy-on-close
            :close-on-click-modal="false" @close="closeDialog">
            <TableEdit :form-data="rowData" :options="options" :edit="isEdit" :update="updateData">
            <template #pid>
                    <el-cascader v-model="rowData.pid" :options="cascaderOptions" :props="{ checkStrictly: true , emitPath: false }" @change="handleChangepid"
                         />
            </template>
            <template #permissions>
                <div style="align-items: center; margin-bottom: 15px;width:100%">
                    <el-checkbox v-model="selectAll" @change="handleSelectAll">
                        全选
                    </el-checkbox>
                    <el-checkbox v-model="expandAll" @change="handleExpandAll">
                        展开全部
                    </el-checkbox>
                </div>
                <div>
                    <el-tree
                        ref="treeRef"
                        :data="treeData"
                        node-key="id"
                        :default-expand-all="expandAll"
                        show-checkbox
                        :default-checked-keys="checkedKeys"
                        @check-change="handleTreeCheckChange"
                    />
                </div>
            </template>
            </TableEdit>
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

<script setup lang="ts" name="system-role">
import TableCustom from '@/components/table-custom.vue';
import TableDetail from '@/components/table-detail.vue';
import TableEdit from '@/components/table-edit.vue';
import { FormOption } from '@/types/form-option';
import { Role } from '@/types/role';
import service from '@/utils/request';
import { ElMessage, ElTree } from 'element-plus';
import NProgress from 'nprogress';
import { onMounted, ref } from 'vue';
const isInitialized = ref(false);

onMounted(() => {
    isInitialized.value = true;
});
// 表格相关
let columns = ref([
    { type: 'selection'},
    { prop: 'title', label: '角色名称',align: 'left',width:"600" },
    { prop: 'status', label: '状态' },
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
const tableData = ref<any[]>([]);
const cascaderOptions = ref<any[]>([]);

// 权限树相关
const treeRef = ref<InstanceType<typeof ElTree>>();
const treeData = ref<any[]>([]);
const checkedKeys = ref<string[]>([]);
const menuObj = ref<any>({});
const rulearr = ref<any>([]);

// 全选和展开控制
const selectAll = ref(true); // 默认全选
const expandAll = ref(true); // 默认展开全部
// 获取权限树数据
const getTreeData = (data: any[]) => {
    return data.map((item) => {
        const obj: any = {
            id: item.id,
            label: item.title,
        };
        if (item.children) {
            menuObj.value[item.id] = item.children.map((sub: any) => sub.id);
            obj.children = getTreeData(item.children);
        }
        return obj;
    });
};

// 初始化权限树数据
const initTreeData = async() => {


    const menuData = localStorage.getItem('allmenuData');
    if (menuData) {
        treeData.value = getTreeData(JSON.parse(menuData));
    }
};

// 检查数据
const checkData = (data: string[]) => {
    return data.filter((item) => {
        return !menuObj.value[item] || data.toString().includes(menuObj.value[item].toString());
    });
};

// 获取所有ID
const getAllIds = (treeData: any[]) => {
    let ids: any[] = [];
    treeData.forEach(item => {
        ids.push(item.id);
        if (item.children) {
            ids = ids.concat(getAllIds(item.children));
        }
    });
    return ids;
};

// 处理全选功能
const handleSelectAll = (value: boolean) => {
    if (!treeRef.value) return;
    
    if (value) {
        // 全选：设置所有节点为选中状态
        const allIds = getAllIds(treeData.value);
        treeRef.value.setCheckedKeys(allIds);
        checkedKeys.value = allIds.map(id => String(id));
    } else {
        // 取消全选：清空所有选中状态
        treeRef.value.setCheckedKeys([]);
        checkedKeys.value = [];
    }
};

// 处理展开全部功能
const handleExpandAll = (value: boolean) => {
    if (!treeRef.value) return;
    
    // 获取所有节点
    const allNodes = treeRef.value.store.nodesMap;
    
    // 遍历所有节点，设置展开/收缩状态
    Object.keys(allNodes).forEach(nodeKey => {
        const node = allNodes[nodeKey];
        if (node.childNodes && node.childNodes.length > 0) {
            node.expanded = value;
        }
    });
};

// 处理树形结构节点选中状态变化
const handleTreeCheckChange = () => {
    if (!treeRef.value) return;
    
    // 获取当前选中的节点
    const checkedNodes = treeRef.value.getCheckedKeys();
    const allIds = getAllIds(treeData.value);
    
    // 检查是否为全选状态
    selectAll.value = allIds.length > 0 && allIds.every(id => checkedNodes.includes(id));
    
    // 更新checkedKeys
    checkedKeys.value = checkedNodes.map(id => String(id));
};

function flattenTreeToMap(tree: any[], map: Record<number, string> = {}) {
  tree.forEach(node => {
    map[node.id] = node.rule
    if (node.children && node.children.length > 0) {
      flattenTreeToMap(node.children, map)
    }
  })
  return map
}

const getData = async () => {
    NProgress.start();
    try {
        const res = await service.post('/api/api/v1/auth/role/list');
        if (res.data.code === 200) {
            tableData.value = res.data.data;
            var rule=getOptions(tableData.value);
            cascaderOptions.value = rule;
            rulearr.value = flattenTreeToMap(tableData.value)
        }else{
            ElMessage.error(res.data.msg);
            return false;
        }
        NProgress.done();
    } catch (e) {
        ElMessage.error("未知错误");
        return false;
    }
};
getData();


// 新增/编辑弹窗相关
const options = ref<FormOption>({
    labelWidth: '100px',
    span: 24,
    list: [
        { type: 'input', label: '角色名称', prop: 'title', required: true },
        { type: 'parent', label: '父级', prop: 'pid',required: true },
        { 
            type: 'switch', 
            label: '状态', 
            prop: 'status',
            activeValue: 1,  
            inactiveValue: 0,
        },
        { type: 'custom', label: '权限管理', prop: 'permissions' }
    ]
})
const visible = ref(false);
const isEdit = ref(false);
const rowData =ref<any>({});
const handleEdit = (row: Role) => {
    if(row[0].id==1){
        ElMessage.error("超级管理员不允许修改");
        return false;
    }
    rowData.value = { ...row[0] };
    isEdit.value = true;
    visible.value = true;
    
    // 初始化权限树
    initTreeData();
    
    // 设置已选中的权限
    if (row[0].rule) {
        const permissions = row[0].rule.split(',').map((item: string) => item.trim());
        if (permissions.includes("*")) {
            // 如果包含"*"，表示全部权限
            checkedKeys.value = getAllIds(treeData.value);
            selectAll.value = true;
        } else {
            checkedKeys.value = checkData(permissions);
            // 检查是否为全选状态
            const allIds = getAllIds(treeData.value);
            selectAll.value = allIds.length > 0 && allIds.every(id => checkedKeys.value.includes(String(id)));
        }
    } else {
        checkedKeys.value = [];
        selectAll.value = false;
    }
    
    // 默认展开全部
    expandAll.value = true;
};
const handleAdd = () => {
    rowData.value = { status: 1, pid: 1 };
    isEdit.value = false;
    visible.value = true;
    
    // 初始化权限树
    initTreeData();
    
    // 默认全选和展开全部
    selectAll.value = true;
    expandAll.value = true;
    checkedKeys.value = getAllIds(treeData.value);
};
const updateData = async (form: any) => {
    // 获取选中的权限
    if (treeRef.value) {
        const keys = [...treeRef.value.getCheckedKeys(false), ...treeRef.value.getHalfCheckedKeys()] as number[];
        rowData.value.rule = keys.join(',');
    }
    try {
        var url="/api/api/v1/auth/role/add";
        if(form.id){
            url="/api/api/v1/auth/role/update";
        }
        form.rule=rowData.value.rule
        form.pid=rowData.value.pid
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
    rowData.value = {};
    // 重置状态
    selectAll.value = true;
    expandAll.value = true;
    checkedKeys.value = [];
};

const handleChangepid =(value: any) => {    
    // 更新选中的权限ID数组
    if (rulearr.value[value] === "*") {
        // 如果是全部权限
        checkedKeys.value = getAllIds(treeData.value);
        selectAll.value = true;
    } else {
        // 确保是字符串类型且不为空
        const ruleStr = String(rulearr.value[value]).trim();
        if (ruleStr && ruleStr !== 'undefined' && ruleStr !== 'null') {
            checkedKeys.value = checkData(ruleStr.split(',').map((item: string) => item.trim()));
            // 检查是否为全选状态
            const allIds = getAllIds(treeData.value);
            selectAll.value = allIds.length > 0 && allIds.every(id => checkedKeys.value.includes(String(id)));
        } else {
            checkedKeys.value = [];
            selectAll.value = false;
        }
    }
    // 同步更新树形组件的选中状态
    if (treeRef.value) {
        treeRef.value.setCheckedKeys(checkedKeys.value);
    } 
};

const handleChange =async (row: Role) => {
    if(isInitialized.value&&row.id){
        if(row.id==1){
            ElMessage.error("超级管理员不能修改");
            return false;
        }
        try {
        const res = await service.post("/api/api/v1/auth/role/updatestatus",{"id":row.id,"status":row.status});
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

// 查看详情弹窗相关
const visible1 = ref(false);
const viewData = ref({
    row: {},
    list: [],
    column: 1
});
const handleView = (row: Role) => {
    viewData.value.row = { ...row[0] }
    viewData.value.list = [
        {
            prop: 'id',
            label: '角色ID',
        },
        {
            prop: 'title',
            label: '角色名称',
        },
        {
            prop: 'status',
            label: '角色状态',
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
        const res = await service.post("/api/api/v1/auth/role/del", {ids:ids});
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